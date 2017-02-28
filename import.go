package kvconfig

import (
	"reflect"
	"strconv"

	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
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
	if v.Kind() == reflect.Interface && v.NumMethod() == 0 || v.Kind() == reflect.Ptr {
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

	s.structCounter.Increment(t)

	var newStruct reflect.Value
	var newStructPtr reflect.Value

	for f := 0; f < t.NumField(); f += 1 {
		field := t.Field(f)
		kn, knok := keyname(&structAndField{t, field}, s.structCounter)
		if _, ok := kv.Lookup(kn); knok && ok {
			if !newStruct.IsValid() {
				newStructPtr = reflect.New(t)
				newStruct = newStructPtr.Elem()
			}

			newField := newStruct.Field(f)
			newType := newField.Type()

			switch newType.Kind() {
			case reflect.Int:
				i, _ := strconv.Atoi(kv.Get(kn))
				newField.SetInt(int64(i))
			case reflect.String:
				newField.SetString(kv.Get(kn))
			case reflect.Ptr:
				t := newField.Interface()
				switch t.(type) {
				case *rsa.PrivateKey:
					cert := unmarshalRSAPrivateKey(kv.Get(kn))
					newField.Set(reflect.ValueOf(cert))
				}
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
