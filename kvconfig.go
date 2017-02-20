// Package kvconfig implements a system of mapping Go structures to and from key/value stores.
//
// Key names end in an underscore and integer (e.g. "_2").
// This is to facilitate e.g. arrays of structures and the multiple values they hold.
// When parsing CLI arguments or envvars names may be transformed to conform.
// When specified on structures the field tag is "config" followed by the key name.
package kvconfig

import (
	"fmt"
	"reflect"
)

type Setter interface {
	Set(string, string)
}

type Getter interface {
	Get(string) string
	Exists(string) bool
}

type structCounter map[reflect.Type]int

func (s structCounter) Increment(t reflect.Type) {
	if _, ok := s[t]; !ok {
		s[t] = 0
	}
	s[t] += 1
}

func (s structCounter) Current(t reflect.Type) int {
	if _, ok := s[t]; !ok {
		s[t] = 0
	}
	return s[t]
}

type structAndField struct {
	structType reflect.Type
	field      reflect.StructField
}

// Tries to derive a numeric-ending key name from the number of times we've seen a structure
func keyname(sfield *structAndField, c structCounter) (string, bool) {
	if sfield == nil || sfield.structType == nil {
		return "", false
	}
	if _, ok := sfield.field.Tag.Lookup("config"); !ok {
		return "", false
	}
	return fmt.Sprintf("%s_%d", sfield.field.Tag.Get("config"), c.Current(sfield.structType)-1), true
}
