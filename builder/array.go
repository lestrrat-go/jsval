package builder

import (
	"github.com/lestrrat/go-jsschema"
	"github.com/lestrrat/go-jsval"
	"github.com/lestrrat/go-pdebug"
)

func buildArrayConstraint(ctx *buildctx, c *jsval.ArrayConstraint, s *schema.Schema) (err error) {
	if pdebug.Enabled {
		g := pdebug.IPrintf("START buildArrayConstraint")
		defer func() {
			if err == nil {
				g.IRelease("END buildArrayConstraint (PASS)")
			} else {
				g.IRelease("END buildArrayConstraint (FAIL): %s", err)
			}
		}()
	}

	if items := s.Items; items != nil {
		if !items.TupleMode {
			specs, err := buildFromSchema(ctx, items.Schemas[0])
			if err != nil {
				return err
			}
			c.Items(specs)
		} else {
			specs := make([]jsval.Constraint, len(items.Schemas))
			for i, espec := range items.Schemas {
				item, err := buildFromSchema(ctx, espec)
				if err != nil {
					return err
				}
				specs[i] = item
			}
			c.PositionalItems(specs)

			if aitems := s.AdditionalItems; aitems != nil {
				if as := aitems.Schema; as != nil {
					spec, err := buildFromSchema(ctx, as)
					if err != nil {
						return err
					}
					c.AdditionalItems(spec)
				} else {
					c.AdditionalItems(jsval.NilConstraint)
				}
			}
		}
	}

	if s.MinItems.Initialized {
		c.MinItems(s.MinItems.Val)
	}

	if s.MaxItems.Initialized {
		c.MaxItems(s.MaxItems.Val)
	}

	if s.UniqueItems.Initialized {
		c.UniqueItems(s.UniqueItems.Val)
	}

	return nil
}