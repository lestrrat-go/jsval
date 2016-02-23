package jsval

import (
	"errors"
	"reflect"

	"github.com/lestrrat/go-jsschema"
)

func (c *ArrayConstraint) FromSchema(s *schema.Schema) error {
	var err error
	if !s.Items.TupleMode {
		c.itemspec, err = FromSchema(s.Items.Schemas[0])
		if err != nil {
			return err
		}
	} else {
		c.positionalItems = make([]Constraint, len(s.Items.Schemas))
		for i, espec := range s.Items.Schemas {
			c.positionalItems[i], err = FromSchema(espec)
			if err != nil {
				return err
			}
		}

		if aitems := s.AdditionalItems; aitems != nil {
			if as := aitems.Schema; as != nil {
				c.additionalItems, err = FromSchema(as)
				if err != nil {
					return err
				}
			} else {
				c.additionalItems = NilConstraint
			}
		}
	}

	return nil
}

func Array() *ArrayConstraint {
	return &ArrayConstraint{}
}

func (c *ArrayConstraint) Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice:
	default:
		return errors.New("value must be a slice")
	}

	l := rv.Len()
	if celem := c.itemspec; celem != nil {
		// if this is set, then all items must fulfill this.
		// additional items are ignored
		for i := 0; i < l; i++ {
			if err := celem.Validate(rv.Index(i).Interface()); err != nil {
				return err
			}
		}
	} else {
		// otherwise, check the positional specs, and apply the
		// additionalItems constraint
		for i, cpos := range c.positionalItems {
			ev := rv.Index(i)
			// We don't do defaults and stuff here, because it's virtually
			// impossible to tell if this is "uninitialized"
			if err := cpos.Validate(ev.Interface()); err != nil {
				return err
			}
		}

		lp := len(c.positionalItems)
		if l >= lp { // we got more than positional schemas
			cadd := c.additionalItems
			if cadd == nil { // you can't have additionalItems!
				return errors.New("additional elements found in array")
			}
			for i := lp - 1; i < l; i++ {
				if err := cadd.Validate(rv.Index(i).Interface()); err != nil {
					return err
				}
			}
		}
	}
	return nil
}