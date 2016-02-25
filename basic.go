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
	if rv == zeroval {
		return nil
	}

	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if rv.IsNil() {
			return nil
		}
	}
	return errors.New("value is not null")
}

func Not(c Constraint) *NotConstraint {
	return &NotConstraint{child: c}
}

func (nc NotConstraint) HasDefault() bool {
		return false
}

func (nc NotConstraint) DefaultValue() interface{} {
	return nil
}

func (nc NotConstraint) Validate(v interface{}) error {
	if err := nc.child.Validate(v); err == nil {
		return errors.New("'not' validation failed")
	}
	return nil
}
