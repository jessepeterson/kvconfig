package kvconfig

import (
	"reflect"
	"strconv"

	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"

	"crypto/tls"
	"fmt"
)

type importState struct {
	structCounter
	depth int
}

// Uses reflection to walk the structure i and create or set new elements from the key/value interface kv.
func Import(kv Getter, i interface{}) error {
	s := importState{}
	s.structCounter = make(structCounter)
	return importWalk(kv, reflect.ValueOf(i), nil, &s)
}

func importWalk(kv Getter, v reflect.Value, sfield *structAndField, s *importState) (err error) {
	if (v.Kind() == reflect.Interface && v.NumMethod() == 0) || (v.Kind() == reflect.Ptr && v.Elem().Kind() != reflect.Invalid) {
		v = v.Elem()
	}

	kn, knok := keyname(sfield, s.structCounter)

	s.depth += 1
	switch v.Kind() {
	case reflect.Struct:
		err = importStruct(kv, v, s)
	case reflect.Slice:
		err = importSlice(kv, v, s)
	case reflect.Int:
		if knok {
			i, _ := strconv.Atoi(kv.Get(kn))
			v.SetInt(int64(i))
		}
	case reflect.String:
		if knok {
			v.SetString(kv.Get(kn))
		}
	case reflect.Ptr:
		if knok {
			t := v.Interface()
			switch t.(type) {
			case *rsa.PrivateKey:
				cert := unmarshalRSAPrivateKey(kv.Get(kn))
				v.Set(reflect.ValueOf(cert))
			case *tls.Certificate:
				n, ct, ok := keynameRaw(sfield, s.structCounter)
				if ok {
					tlsCert, ok := importTLSCertificate(kv, n, ct)
					if ok {
						s.structCounter.Increment(v.Type())
						v.Set(reflect.ValueOf(tlsCert))
					}
				}
			default:
				switch v.Type().Elem().Kind() {
				case reflect.Int, reflect.String:
					v.Set(reflect.New(v.Type().Elem()))
					err = importWalk(kv, v, sfield, s)
				}
			}
		}
	}
	s.depth -= 1

	return
}

func importSlice(kv Getter, v reflect.Value, s *importState) (err error) {
	for i := 0; i < v.Len(); i += 1 {
		err = importWalk(kv, v.Index(i), nil, s)
		if err != nil {
			break
		}
	}

	sliceType := v.Type().Elem()
	structCand := sliceType

	if sliceType.Kind() == reflect.Ptr {
		structCand = sliceType.Elem()
	}

	if structCand.Kind() == reflect.Struct {
		for newStruct, ok := importNewStruct(kv, structCand, s); ok; newStruct, ok = importNewStruct(kv, structCand, s) {

			// Grow the slice if necessary.
			// Borrowed from https://golang.org/src/encoding/xml/read.go
			n := v.Len()
			if n >= v.Cap() {
				ncap := 2 * n
				if ncap < 4 {
					ncap = 4
				}
				new := reflect.MakeSlice(v.Type(), n, ncap)
				reflect.Copy(new, v)
				v.Set(new)
			}

			v.SetLen(n + 1)
			v.Index(n).Set(newStruct)

			err = importStruct(kv, newStruct.Elem(), s)
		}
	}

	return
}

func importStruct(kv Getter, v reflect.Value, s *importState) (err error) {
	s.structCounter.Increment(v.Type())

	for f := 0; f < v.NumField(); f += 1 {
		sfield := structAndField{v.Type(), v.Type().Field(f)}
		err = importWalk(kv, v.Field(f), &sfield, s)
		if err != nil {
			break
		}
	}
	return
}

func importNewStruct(kv Getter, t reflect.Type, s *importState) (reflect.Value, bool) {
	if t.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}
	var newStruct reflect.Value
	var newStructPtr reflect.Value

	for f := 0; f < t.NumField(); f += 1 {
		field := t.Field(f)
		kn, knok := keyname(&structAndField{t, field}, s.structCounter)
		if _, ok := kv.Lookup(kn); knok && ok {
			if !newStruct.IsValid() {
				newStructPtr = reflect.New(t)
				newStruct = newStructPtr.Elem()
				break
			}
		}
	}

	return newStructPtr, reflect.Value{} != newStruct
}

func unmarshalRSAPrivateKey(s string) *rsa.PrivateKey {
	x509bytes, _ := base64.StdEncoding.DecodeString(s)
	cert, _ := x509.ParsePKCS1PrivateKey(x509bytes)
	return cert
}

func importTLSCertificate(kv Getter, name string, ct int) (*tls.Certificate, bool) {
	_, certOk := kv.Lookup(fmt.Sprintf("%s_cert_%d", name, ct))
	keyStr, keyOk := kv.Lookup(fmt.Sprintf("%s_pk_%d", name, ct))

	if !keyOk || !certOk {
		return nil, false
	}

	tC := tls.Certificate{}

	for i := 0; ; i++ {
		var keyName string
		if i < 1 {
			keyName = fmt.Sprintf("%s_cert_%d", name, ct)
		} else {
			keyName = fmt.Sprintf("%s_cert%d_%d", name, i+1, ct)
		}

		if certStr, ok := kv.Lookup(keyName); ok {
			certBytes, _ := base64.StdEncoding.DecodeString(certStr)
			tC.Certificate = append(tC.Certificate, certBytes)
		} else {
			break
		}
	}

	keyBytes, _ := base64.StdEncoding.DecodeString(keyStr)
	tC.PrivateKey, _ = x509.ParsePKCS1PrivateKey(keyBytes)

	return &tC, true
}
