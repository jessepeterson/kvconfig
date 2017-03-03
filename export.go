package kvconfig

import (
	"reflect"
	"strconv"

	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

type exportState struct {
	structCounter
	depth int
}

// Uses reflection to walk the structure i and set values in the key/value interface kv.
func Export(i interface{}, kv Setter) error {
	s := exportState{}
	s.structCounter = make(structCounter)
	return exportWalk(reflect.ValueOf(i), nil, kv, &s)
}

func exportWalk(v reflect.Value, sfield *structAndField, kv Setter, s *exportState) (err error) {
	s.depth += 1

	kn, knok := keyname(sfield, s.structCounter)

	switch v.Kind() {
	case reflect.Map:
		err = exportMap(v, kv, s)
	case reflect.Slice:
		err = exportSlice(v, kv, s)
	case reflect.Struct:
		err = exportStruct(v, kv, s)
	case reflect.String:
		if knok {
			kv.Set(kn, v.String())
		}
	case reflect.Int:
		if knok {
			kv.Set(kn, strconv.Itoa(int(v.Int())))
		}
	case reflect.Interface:
		if v.NumMethod() == 0 {
			err = exportWalk(v.Elem(), sfield, kv, s)
		}
	case reflect.Ptr:
		t := v.Interface()
		switch t.(type) {
		case *rsa.PrivateKey:
			if knok {
				pk := t.(*rsa.PrivateKey)
				kv.Set(kn, marshalRSAPrivateKey(pk))
			}
		case *tls.Certificate:
			if keyRaw, ct, ok := keynameRaw(sfield, s.structCounter); ok {
				exportTLSCertificate(kv, keyRaw, ct, t.(*tls.Certificate))
			}

		default:
			err = exportWalk(v.Elem(), sfield, kv, s)
		}
	}
	s.depth -= 1
	return
}

func exportStruct(v reflect.Value, kv Setter, s *exportState) (err error) {
	s.structCounter.Increment(v.Type())

	for f := 0; f < v.NumField(); f += 1 {
		sfield := structAndField{v.Type(), v.Type().Field(f)}
		err = exportWalk(v.Field(f), &sfield, kv, s)
		if err != nil {
			break
		}
	}

	return
}

func exportSlice(v reflect.Value, kv Setter, s *exportState) (err error) {
	for i := 0; i < v.Len(); i += 1 {
		err = exportWalk(v.Index(i), nil, kv, s)
		if err != nil {
			break
		}
	}
	return
}

func exportMap(v reflect.Value, kv Setter, s *exportState) (err error) {
	for _, key := range v.MapKeys() {
		err = exportWalk(v.MapIndex(key), nil, kv, s)
		if err != nil {
			break
		}
	}
	return
}

func marshalRSAPrivateKey(pk *rsa.PrivateKey) string {
	der := x509.MarshalPKCS1PrivateKey(pk)
	return base64.StdEncoding.EncodeToString(der)
}

func exportTLSCertificate(kv Setter, name string, ct int, tlsCert *tls.Certificate) {
	if tlsCert == nil {
		return
	}

	tC := *tlsCert

	for i := 0; i < len(tC.Certificate); i++ {
		if len(tC.Certificate[i]) > 0 {
			certStr := base64.StdEncoding.EncodeToString(tC.Certificate[i])
			var keyName string
			if i < 1 {
				keyName = fmt.Sprintf("%s_cert_%d", name, ct)
			} else {
				keyName = fmt.Sprintf("%s_cert%d_%d", name, i+1, ct)
			}
			kv.Set(keyName, certStr)
		}
	}

	if tC.PrivateKey != nil {
		keyBytes := x509.MarshalPKCS1PrivateKey(tC.PrivateKey.(*rsa.PrivateKey))
		keyStr := base64.StdEncoding.EncodeToString(keyBytes)
		keyName := fmt.Sprintf("%s_pk_%d", name, ct)
		kv.Set(keyName, keyStr)
	}
}
