// Package jsval implements an input validator, based on JSON Schema.
// The main purpose is to validate JSON Schemas (see
// https://github.com/lestrrat/go-jsschema), and to automatically
// generate validators from schemas, but jsval can be used independently
// of JSON Schema.
package jsval

// New creates a new JSVal instance.
func New() *JSVal {
	return &JSVal{
		ConstraintMap: &ConstraintMap{},
	}
}

// Validate validates the input, and return an error
// if any of the validations fail
func (v *JSVal) Validate(x interface{}) error {
	return v.root.Validate(x)
}

// SetRoot sets the root Constraint object.
func (v *JSVal) SetRoot(c Constraint) *JSVal {
	v.root = c
	return v
}

// Root returns the root Constraint object.
func (v *JSVal) Root() Constraint {
	return v.root
}

// SetConstraintMap allows you to set the map that is referred to
// when resolving JSON references within constraints. By setting
// this to a common map, for example, you can share the same references
// to save you some memory space and sanity. See an example in the
// `generated_validator_test.go` file.
func (v *JSVal) SetConstraintMap(cm *ConstraintMap) *JSVal {
	v.ConstraintMap = cm
	return v
}
