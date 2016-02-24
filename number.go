package jsval

import (
	"errors"
	"reflect"

	"github.com/lestrrat/go-jsschema"
	"github.com/lestrrat/go-pdebug"
)

func (nc *NumberConstraint) buildFromSchema(ctx *buildctx, s *schema.Schema) error {
	if !s.Type.Contains(schema.NumberType) && !s.Type.Contains(schema.IntegerType) {
		return errors.New("schema is not for number")
	}

	if s.Minimum.Initialized {
		nc.Minimum(s.Minimum.Val)
		if s.ExclusiveMinimum.Initialized {
			nc.exclusiveMinimum = s.ExclusiveMinimum.Val
		}
	}

	if s.Maximum.Initialized {
		nc.Maximum(s.Maximum.Val)
		if s.ExclusiveMaximum.Initialized {
			nc.exclusiveMaximum = s.ExclusiveMaximum.Val
		}
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
	if nc.enums == nil {
		nc.enums = Enum()
	}
	nc.enums.Enum(l)
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

func (nc *NumberConstraint) ExclusiveMinimum(b bool) *NumberConstraint {
	nc.exclusiveMinimum = b
	return nc
}

func (nc *NumberConstraint) ExclusiveMaximum(b bool) *NumberConstraint {
	nc.exclusiveMaximum = b
	return nc
}

func (nc *NumberConstraint) Validate(v interface{}) (err error) {
	if pdebug.Enabled {
		g := pdebug.IPrintf("START NumberConstraint.Validate")
		defer func() {
			if err == nil {
				g.IRelease("END NumberConstraint.Validate (PASS)")
			} else {
				g.IRelease("END NumberConstraint.Validate (FAIL): %s", err)
			}
		}()
	}

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
		if pdebug.Enabled {
			pdebug.Printf("Checking Minimum (%f)", nc.minimum)
		}
		if nc.minimum > f {
			return errors.New("numeric value less than minimum")
		}
	}

	if nc.applyMaximum {
		if pdebug.Enabled {
			pdebug.Printf("Checking Maximum (%f)", nc.maximum)
		}
		if nc.maximum < f {
			return errors.New("numeric value greater than maximum")
		}
	}

	if enum := nc.enums; enum != nil {
		if err := enum.Validate(f); err != nil {
			return err
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

func (ic *IntegerConstraint) Maximum(n float64) *IntegerConstraint {
	ic.applyMaximum = true
	ic.maximum = n
	return ic
}

func (ic *IntegerConstraint) Minimum(n float64) *IntegerConstraint {
	ic.applyMinimum = true
	ic.minimum = n
	return ic
}

func (ic *IntegerConstraint) ExclusiveMinimum(b bool) *IntegerConstraint {
	ic.NumberConstraint.ExclusiveMinimum(b)
	return ic
}

func (ic *IntegerConstraint) ExclusiveMaximum(b bool) *IntegerConstraint {
	ic.NumberConstraint.ExclusiveMaximum(b)
	return ic
}

func (ic *IntegerConstraint) Validate(v interface{}) (err error) {
	if pdebug.Enabled {
		g := pdebug.IPrintf("START IntegerConstraint.Validate")
		defer func() {
			if err == nil {
				g.IRelease("END IntegerConstraint.Validate (PASS)")
			} else {
				g.IRelease("END IntegerConstraint.Validate (FAIL): %s", err)
			}
		}()
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Interface, reflect.Ptr:
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ic.NumberConstraint.Validate(float64(rv.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return ic.NumberConstraint.Validate(float64(rv.Uint()))
	default:
		return errors.New("value is not an int/uint")
	}
}
