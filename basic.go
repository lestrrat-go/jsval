package jsval

import (
	"errors"
	"reflect"
)

func (dv defaultValue) HasDefault() bool {
	return dv.initialized
}

func (dv defaultValue) DefaultValue() interface{} {
	return dv.value
}

func (nc emptyConstraint) Validate(_ interface{}) error {
	return nil
}

func (nc emptyConstraint) HasDefault() bool {
	return false
}

func (nc emptyConstraint) DefaultValue() interface{} {
	return nil
}

func (nc nullConstraint) HasDefault() bool {
	return false
}

func (nc nullConstraint) DefaultValue() interface{} {
	return nil
}

func (nc nullConstraint) Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv == zeroval || rv.IsNil() {
		return nil
	}
	return errors.New("value is not null")
}
