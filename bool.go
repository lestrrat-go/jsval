package jsval

import (
	"errors"
	"reflect"
)

func Boolean() *BooleanConstraint {
	return &BooleanConstraint{}
}

func (bc *BooleanConstraint) Default(v interface{}) *BooleanConstraint {
	bc.defaultValue.initialized = true
	bc.defaultValue.value = v
	return bc
}

func (b *BooleanConstraint) Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Bool:
	default:
		return errors.New("value is not a boolean")
	}
	return nil
}
