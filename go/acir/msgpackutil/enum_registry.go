package msgpackutil

import (
	"fmt"
	"reflect"
)

// EnumVariant is the constraint each variant type must satisfy
type EnumVariant interface {
	SerdeName() string
}

// EnumRegistry pairs an enum type's [Schema] with a per-tag constructor table,
// derived from a single ordered exemplar list.
type EnumRegistry[V EnumVariant] struct {
	Schema       Schema
	constructors []func() V
}

// NewEnumRegistry builds a registry from an ordered list of exemplar values
// (typically zero-valued pointers).
func NewEnumRegistry[V EnumVariant](exemplars []V) EnumRegistry[V] {
	nameToTag := make(map[string]int, len(exemplars))
	constructors := make([]func() V, len(exemplars))
	for i, ex := range exemplars {
		nameToTag[ex.SerdeName()] = i
		t := reflect.TypeOf(ex).Elem()
		constructors[i] = func() V { return reflect.New(t).Interface().(V) }
	}
	return EnumRegistry[V]{Schema: NewSchema(nameToTag), constructors: constructors}
}

// New constructs the variant for tag, or errors if tag is out of range.
func (r EnumRegistry[V]) New(tag int) (V, error) {
	var zero V
	if tag < 0 || tag >= len(r.constructors) {
		return zero, fmt.Errorf("unknown variant tag %d", tag)
	}
	return r.constructors[tag](), nil
}
