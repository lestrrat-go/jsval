package jsval

import (
	"errors"
	"reflect"

	"github.com/lestrrat/go-pdebug"
)

func Array() *ArrayConstraint {
	return &ArrayConstraint{
		additionalItems: NilConstraint,
		minItems: -1,
	}
}

func (c *ArrayConstraint) Validate(v interface{}) (err error) {
	if pdebug.Enabled {
		g := pdebug.IPrintf("START ArrayConstraint.Validate")
		defer func() {
			if err == nil {
				g.IRelease("END ArrayConstraint.Validate (PASS)")
			} else {
				g.IRelease("END ArrayConstraint.Validate (FAIL): %s", err)
			}
		}()
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice:
	default:
		return errors.New("value must be a slice")
	}

	l := rv.Len()
	if celem := c.items; celem != nil {
		if pdebug.Enabled {
			pdebug.Printf("Checking if all items match a spec")
		}
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
		lv := rv.Len()
		for i, cpos := range c.positionalItems {
			if lv <= i {
				break
			}
			if pdebug.Enabled {
				pdebug.Printf("Checking positional item at '%d'", i)
			}
			ev := rv.Index(i)
			// We don't do defaults and stuff here, because it's virtually
			// impossible to tell if this is "uninitialized"
			if err := cpos.Validate(ev.Interface()); err != nil {
				return err
			}
		}

		lp := len(c.positionalItems)
		if lp > 0 && l > lp { // we got more than positional schemas
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

func (c *ArrayConstraint) AdditionalItems(ac Constraint) *ArrayConstraint {
	c.additionalItems = ac
	return c
}

func (c *ArrayConstraint) Items(ac Constraint) *ArrayConstraint {
	c.items = ac
	return c
}

func (c *ArrayConstraint) MinItems(i int) *ArrayConstraint {
	c.minItems = i
	return c
}

func (c *ArrayConstraint) MaxItems(i int) *ArrayConstraint {
	c.maxItems = i
	return c
}

func (c *ArrayConstraint) PositionalItems(ac []Constraint) *ArrayConstraint {
	c.positionalItems = ac
	return c
}

func (c *ArrayConstraint) UniqueItems(b bool) *ArrayConstraint {
	c.uniqueItems = b
	return c
}