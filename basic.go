package jsval

import "github.com/lestrrat/go-jsschema"

func (dv defaultValue) HasDefault() bool {
	return dv.initialized
}

func (dv defaultValue) DefaultValue() interface{} {
	return dv.value
}

func (nc nilConstraint) Validate(_ interface{}) error {
	return nil
}

func (nc nilConstraint) FromSchema(_ *schema.Schema) error {
	return nil
}

func (nc nilConstraint) HasDefault() bool {
	return false
}

func (nc nilConstraint) DefaultValue() interface{} {
	return nil
}

func (nc nilConstraint) IsRequired() bool {
	return false
}

func (nc nilConstraint) Required(_ bool) {
}

func (r required) IsRequired() bool {
	return bool(r)
}

func (r *required) Required(v bool) {
	*r = required(v)
}
