package jsval

import (
	"errors"
	"fmt"
	"reflect"
)

// Boolean creates a new BooleanConsraint
func Boolean() *BooleanConstraint {
	return &BooleanConstraint{}
}

// Default specifies the default value to apply
func (bc *BooleanConstraint) Default(v interface{}) *BooleanConstraint {
	bc.defaultValue.initialized = true
	bc.defaultValue.value = v
	return bc
}

// ExpectValue allows you to set the value to expect
func (bc *BooleanConstraint) ExpectValue(value bool) *BooleanConstraint {
	b := new(bool)
	*b = value
	bc.expected = b
	return bc
}

// Validate vaidates the value against the given value
func (bc *BooleanConstraint) Validate(v interface{}) error {
	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Bool:
	default:
		return errors.New("value is not a boolean")
	}

	if bc.expected != nil {
		if v := rv.Bool(); *bc.expected != v {
			return fmt.Errorf(
				"expected value to be %v, but got %v",
				*bc.expected, v)
		}
	}

	return nil
}
