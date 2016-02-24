package jsval

import (
	"errors"
	"reflect"

	"github.com/lestrrat/go-jsschema"
	"github.com/lestrrat/go-pdebug"
)

func (c *ArrayConstraint) buildFromSchema(ctx *buildctx, s *schema.Schema) error {
	var err error
	if items := s.Items; items != nil {
		if !items.TupleMode {
			c.itemspec, err = buildFromSchema(ctx, items.Schemas[0])
			if err != nil {
				return err
			}
		} else {
			c.positionalItems = make([]Constraint, len(items.Schemas))
			for i, espec := range items.Schemas {
				c.positionalItems[i], err = buildFromSchema(ctx, espec)
				if err != nil {
					return err
				}
			}

			if aitems := s.AdditionalItems; aitems != nil {
				if as := aitems.Schema; as != nil {
					c.additionalItems, err = buildFromSchema(ctx, as)
					if err != nil {
						return err
					}
				} else {
					c.additionalItems = NilConstraint
				}
			}
		}
	}

	if s.MinItems.Initialized {
		c.minItems = s.MinItems.Val
	}

	if s.MaxItems.Initialized {
		c.maxItems = s.MaxItems.Val
	}

	if s.UniqueItems.Initialized {
		c.uniqueItems = s.UniqueItems.Val
	}

	return nil
}

func Array() *ArrayConstraint {
	return &ArrayConstraint{
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
	if celem := c.itemspec; celem != nil {
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

func (c *ArrayConstraint) MinItems(i int) *ArrayConstraint {
	c.minItems = i
	return c
}

func (c *ArrayConstraint) MaxItems(i int) *ArrayConstraint {
	c.maxItems = i
	return c
}

func (c *ArrayConstraint) UniqueItems() *ArrayConstraint {
	c.uniqueItems = true
	return c
}