package jsval

import (
	"errors"

	"github.com/lestrrat/go-jsref"
	"github.com/lestrrat/go-jsschema"
	"github.com/lestrrat/go-pdebug"
)

func New() *JSVal {
	return &JSVal{
		refs: make(map[string]Constraint),
	}
}

type buildctx struct {
	V *JSVal
	S *schema.Schema
	R map[string]struct{}
}

func (v *JSVal) Build(s *schema.Schema) (err error) {
	if pdebug.Enabled {
		g := pdebug.IPrintf("START JSVal.Build")
		defer func() {
			if err == nil {
				g.IRelease("END JSVal.Build (OK)")
			} else {
				g.IRelease("END JSVal.Build (FAIL): %s", err)
			}
		}()
	}

	ctx := buildctx{
		V: v,
		S: s,
		R: map[string]struct{}{}, // names of references used
	}
	c, err := buildFromSchema(&ctx, s)
	if err != nil {
		return err
	}

	// Now, resolve references that were used in the schema
	if len(ctx.R) > 0 {
		if pdebug.Enabled {
			pdebug.Printf("Checking references now")
		}
		r := jsref.New()
		for ref := range ctx.R {
			if pdebug.Enabled {
				pdebug.Printf("Building constraints for reference '%s'", ref)
			}

			if ref == "#" {
				if pdebug.Enabled {
					pdebug.Printf("'%s' resolves to the main schema", ref)
				}
				v.refs[ref] = c
				continue
			}

			thing, err := r.Resolve(s, ref)
			if err != nil {
				return err
			}

			s1, ok := thing.(*schema.Schema)
			if !ok {
				return errors.New("resolved result is not a schema")
			}

			c1, err := buildFromSchema(&ctx, s1)
			if err != nil {
				return err
			}
			v.refs[ref] = c1
		}
	}

	v.root = c
	return nil
}

func (v *JSVal) Validate(x interface{}) error {
	return v.root.Validate(x)
}

func (v *JSVal) SetRoot(c Constraint) {
	v.root = c
}

func (v *JSVal) Root() Constraint {
	return v.root
}

type unresolved struct {
	V *JSVal
	S *schema.Schema
	R *jsref.Resolver
}

func (v *JSVal) GetReference(ref string) (Constraint, error) {
	v.reflock.Lock()
	defer v.reflock.Unlock()
	c, ok := v.refs[ref]
	if !ok {
		return nil, errors.New("reference '" + ref + "' not found")
	}

	return c, nil
}

/*
	switch c.(type) {
	case Constraint:
		return c.(Constraint), nil
	case unresolved:
		u := c.(unresolved)
		thing, err := u.R.Resolve(u.S, ref)
		if err != nil {
			return nil, err
		}

		s, ok := thing.(*schema.Schema)
		if !ok {
			return nil, errors.New("resolved to something other than a schema")
		}

		c1, err := buildFromSchema(s)
		if err != nil {
			return nil, err
		}
		v.refs[ref] = c1
		return c1, nil
	default:
		return nil, errors.New("invalid reference")
	}
*/

func (v *JSVal) SetReference(ref string, c Constraint) {
	if pdebug.Enabled {
		pdebug.Printf("JSVal.SetReference %s", ref)
	}

	v.reflock.Lock()
	defer v.reflock.Unlock()
	v.refs[ref] = c
}

/*
func (ctx *buildctx) resolve(ref string) (Constraint, error) {
	c, err := ctx.V.GetReference(ref)
	if err == nil {
		return c, nil
	}

	thing, err := ctx.resolver.Resolve(ctx.schema, ref)
	if err != nil {
		return nil, err
	}

	s, ok := thing.(*schema.Schema)
	if !ok {
		return nil, errors.New("resolved to something other than a schema")
	}

	c, err = ctx.V.buildFromSchema(s)
	if err != nil {
		return nil, err
	}

	ctx.V.SetReference(ref, c)
	return c, nil
}
*/

func buildFromSchema(ctx *buildctx, s *schema.Schema) (Constraint, error) {
	if ref := s.Reference; ref != "" {
		c := Reference(ctx.V)
		if err := c.buildFromSchema(ctx, s); err != nil {
			return nil, err
		}
		ctx.R[ref] = struct{}{}
		return c, nil
	}

	ct := All()

	switch {
	case s.Not != nil:
		if pdebug.Enabled {
			pdebug.Printf("Not constraint")
		}
		ct.Add(NilConstraint)
	case len(s.AllOf) > 0:
		if pdebug.Enabled {
			pdebug.Printf("AllOf constraint")
		}
		ac := All()
		for _, s1 := range s.AllOf {
			c1, err := buildFromSchema(ctx, s1)
			if err != nil {
				return nil, err
			}
			ac.Add(c1)
		}
		ct.Add(ac.Reduce())
	case len(s.AnyOf) > 0:
		if pdebug.Enabled {
			pdebug.Printf("AnyOf constraint")
		}
		ac := Any()
		for _, s1 := range s.AnyOf {
			c1, err := buildFromSchema(ctx, s1)
			if err != nil {
				return nil, err
			}
			ac.Add(c1)
		}
		ct.Add(ac.Reduce())
	case len(s.OneOf) > 0:
		if pdebug.Enabled {
			pdebug.Printf("OneOf constraint")
		}
		ct.Add(NilConstraint)
	}

	if l := len(s.Type); l > 0 {
		tct := Any()
		for _, st := range s.Type {
			var c Constraint
			switch st {
			case schema.StringType:
				c = String()
			case schema.NumberType:
				c = Number()
			case schema.IntegerType:
				c = Integer()
			case schema.BooleanType:
				c = Boolean()
			case schema.ArrayType:
				c = Array()
			case schema.ObjectType:
				c = Object()
			default:
				return nil, errors.New("unknown type: " + st.String())
			}
			if err := c.buildFromSchema(ctx, s); err != nil {
				return nil, err
			}
			tct.Add(c)
		}
		ct.Add(tct.Reduce())
	} else {
		// No type?! deduce which constraints apply
		if len(s.Properties) > 0 || (s.AdditionalProperties != nil && s.AdditionalProperties.Schema == nil) {
			oc := Object()
			if err := oc.buildFromSchema(ctx, s); err != nil {
				return nil, err
			}
			ct.Add(oc)
		}

		// All else failed, check if we have some enumeration?
		if len(s.Enum) > 0 {
			ec := Enum(s.Enum...)
			ct.Add(ec)
		}
	}

	return ct.Reduce(), nil
}

func matchenum(v interface{}, values []interface{}) bool {
	for _, x := range values {
		if x == v {
			return true
		}
	}
	return false
}
