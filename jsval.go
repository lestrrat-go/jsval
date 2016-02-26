package jsval

// Package jsval implements an input validator, based on JSON Schema.
// The main purpose is to validate JSON Schemas (see
// https://github.com/lestrrat/go-jsschema), and to automatically
// generate validators from schemas, but jsval can be used independently
// of JSON Schema.

import (
	"errors"

	"github.com/lestrrat/go-pdebug"
)

// New creates a new JSVal instance.
func New() *JSVal {
	return &JSVal{
		refs: make(map[string]Constraint),
	}
}

// Validate validates the input, and return an error
// if any of the validations fail
func (v *JSVal) Validate(x interface{}) error {
	return v.root.Validate(x)
}

// SetRoot sets the root Constraint object.
func (v *JSVal) SetRoot(c Constraint) {
	v.root = c
}

// Root returns the root Constraint object.
func (v *JSVal) Root() Constraint {
	return v.root
}

// GetReference returns the Constraint object pointed at by `ref`.
// It will return an error if a matching constraint has not already
// been registered
func (v *JSVal) GetReference(ref string) (Constraint, error) {
	v.reflock.Lock()
	defer v.reflock.Unlock()
	c, ok := v.refs[ref]
	if !ok {
		return nil, errors.New("reference '" + ref + "' not found")
	}

	return c, nil
}

// SetReference sets a Constraint object to be associated with `ref`.
func (v *JSVal) SetReference(ref string, c Constraint) {
	if pdebug.Enabled {
		pdebug.Printf("JSVal.SetReference %s", ref)
	}

	v.reflock.Lock()
	defer v.reflock.Unlock()
	v.refs[ref] = c
}
