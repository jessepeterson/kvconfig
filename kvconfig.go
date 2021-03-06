// Package kvconfig implements a system of mapping Go structures to and from key/value stores.
//
// Key names end in an underscore and integer (e.g. "_2").
// This is to facilitate e.g. arrays of structures and the multiple values they hold.
// When parsing CLI arguments or envvars names may be transformed to conform.
// When specified on structures the field tag is "kvconfig" followed by the key name.
package kvconfig

import (
	"fmt"
	"reflect"
)

const structTagName = "kvconfig"

type Setter interface {
	Set(string, string)
}

type Getter interface {
	Get(string) string
	Lookup(string) (string, bool)
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
	name, ct, ok := keynameRaw(sfield, c)

	if !ok {
		return "", false
	}

	return fmt.Sprintf("%s_%d", name, ct), true
}

func keynameRaw(sfield *structAndField, c structCounter) (string, int, bool) {
	if sfield == nil || sfield.structType == nil {
		return "", 0, false
	}
	lTagName, ok := sfield.field.Tag.Lookup(structTagName)
	if !ok {
		return "", 0, false
	}
	ct := c.Current(sfield.structType)
	if ct >= 1 {
		ct = ct - 1
	} else {
		ct = 0
	}
	return lTagName, ct, true
}
