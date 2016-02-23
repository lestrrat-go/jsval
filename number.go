package jsval

import (
	"errors"
	"reflect"

	"github.com/lestrrat/go-jsschema"
)

func (nc *NumberConstraint) FromSchema(s *schema.Schema) error {
	if !s.Type.Contains(schema.NumberType) && !s.Type.Contains(schema.IntegerType) {
		return errors.New("schema is not for number")
	}

	if s.Minimum.Initialized {
		nc.Minimum(s.Minimum.Val)
	}
	if s.Maximum.Initialized {
		nc.Maximum(s.Maximum.Val)
	}

	if lst := s.Enum; len(lst) > 0 {
		nc.Enum(lst)
	}

	if v := s.Default; v != nil {
		nc.Default(v)
	}

	return nil
}

func (nc *NumberConstraint) Enum(l []interface{}) *NumberConstraint {
	nc.enum = l
	return nc
}

func (nc *NumberConstraint) Default(v interface{}) *NumberConstraint {
	nc.defaultValue.initialized = true
	nc.defaultValue.value = v
	return nc
}

func (nc *NumberConstraint) Maximum(n float64) *NumberConstraint {
	nc.applyMaximum = true
	nc.maximum = n
	return nc
}

func (nc *NumberConstraint) Minimum(n float64) *NumberConstraint {
	nc.applyMinimum = true
	nc.minimum = n
	return nc
}

func (nc *NumberConstraint) Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface:
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Float32, reflect.Float64:
	default:
		return errors.New("value is not a float")
	}

	f := rv.Float()
	if nc.applyMinimum {
		if nc.minimum > f {
			return errors.New("numeric value less than minimum")
		}
	}

	if nc.applyMaximum {
		if nc.maximum < f {
			return errors.New("numeric value greater than maximum")
		}
	}

	if enum := nc.enum; enum != nil {
		if !matchenum(f, enum) {
			return errors.New("value not in enumeration")
		}
	}

	return nil
}

func Number() *NumberConstraint {
	return &NumberConstraint{
		applyMinimum: false,
		applyMaximum: false,
	}
}

func Integer() *IntegerConstraint {
	c := &IntegerConstraint{}
	c.applyMinimum = false
	c.applyMaximum = false
	return c
}

func (ic *IntegerConstraint) Validate(v interface{}) error {
	rv := reflect.ValueOf(v).Elem()
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ic.NumberConstraint.Validate(float64(rv.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return ic.NumberConstraint.Validate(float64(rv.Uint()))
	default:
		return errors.New("value is not an int/uint")
	}
}
