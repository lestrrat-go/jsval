package jsval

import (
	"errors"
	"reflect"

	"github.com/lestrrat/go-jsschema"
	"github.com/lestrrat/go-pdebug"
)

func (c *ObjectConstraint) buildFromSchema(ctx *buildctx, s *schema.Schema) error {
	if pdebug.Enabled {
		g := pdebug.IPrintf("START ObjectConstraint.FromSchema")
		defer g.IRelease("END ObjectConstraint.FromSchema")
	}

	if l := s.Required; len(l) > 0 {
		c.Required(l...)
	}

	for pname, pdef := range s.Properties {
		cprop, err := buildFromSchema(ctx, pdef)
		if err != nil {
			return err
		}

		c.AddProp(pname, cprop)
	}

	if aprops := s.AdditionalProperties; aprops != nil {
		if sc := aprops.Schema; sc != nil {
			aitem, err := buildFromSchema(ctx, sc)
			if err != nil {
				return err
			}
			c.AdditionalProperties(aitem)
		} else {
			c.AdditionalProperties(NilConstraint)
		}
	}

	for from, to := range s.Dependencies.Names {
		c.PropDependency(from, to...)
	}
	/* TODO: unimplemented */
	/*
		for from, to := range s.Dependencies.Schemas {
			c.SchemaDependency(from, to)
		}
	*/

	return nil
}

func Object() *ObjectConstraint {
	return &ObjectConstraint{
		additionalProperties: nil,
		properties:           make(map[string]Constraint),
		propdeps:             make(map[string][]string),
		required:             make(map[string]struct{}),
	}
}

func (o *ObjectConstraint) Required(l ...string) *ObjectConstraint {
	o.reqlock.Lock()
	defer o.reqlock.Unlock()

	for _, pname := range l {
		o.required[pname] = struct{}{}
	}
	return o
}

func (o *ObjectConstraint) IsPropRequired(s string) bool {
	o.reqlock.Lock()
	defer o.reqlock.Unlock()

	_, ok := o.required[s]
	return ok
}

func (o *ObjectConstraint) AdditionalProperties(c Constraint) *ObjectConstraint {
	o.additionalProperties = c
	return o
}

func (o *ObjectConstraint) AddProp(name string, c Constraint) *ObjectConstraint {
	o.proplock.Lock()
	defer o.proplock.Unlock()

	o.properties[name] = c
	return o
}

func (o *ObjectConstraint) PropDependency(from string, to ...string) *ObjectConstraint {
	o.deplock.Lock()
	defer o.deplock.Unlock()

	l := o.propdeps[from]
	l = append(l, to...)
	o.propdeps[from] = l
	return o
}

/* TODO: Properly implement this */
/*
func (o *ObjectConstraint) SchemaDependency(from string, s *schema.Schema) *ObjectConstraint {
	o.deplock.Lock()
	defer o.deplock.Unlock()

	o.schemadeps[from] = s
	return o
}
*/

func (o *ObjectConstraint) GetPropDependencies(from string) []string {
	o.deplock.Lock()
	defer o.deplock.Unlock()

	l, ok := o.propdeps[from]
	if !ok {
		return nil
	}

	return l
}

func (o *ObjectConstraint) Validate(v interface{}) (err error) {
	if pdebug.Enabled {
		g := pdebug.IPrintf("START ObjectConstraint.Validate")
		defer func() {
			if err == nil {
				g.IRelease("END ObjectConstraint.Validate (PASS)")
			} else {
				g.IRelease("END ObjectConstraint.Validate (FAIL): %s", err)
			}
		}()
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface:
		rv = rv.Elem()
	}

	var fields []string
	switch rv.Kind() {
	case reflect.Struct:
		if pdebug.Enabled {
			pdebug.Printf("Validation target is a struct")
		}
		fields = o.FieldNamesFromStruct(rv)
	case reflect.Map:
		if pdebug.Enabled {
			pdebug.Printf("Validation target is a map")
		}
		if rv.Type().Key().Kind() != reflect.String {
			return errors.New("only maps with string keys can be handled")
		}
		for _, v := range rv.MapKeys() {
			fields = append(fields, v.String())
		}
	default:
		return errors.New("value is not map/object")
	}

	// Find the list of field names that were passed to us
	// "premain" shows extra props, if any.
	// "pseen" shows props that we have already seen
	premain := map[string]struct{}{}
	pseen := map[string]struct{}{}
	for _, k := range fields {
		premain[k] = struct{}{}
	}

	// Now, for all known constraints, validate the prop
	// create a copy of properties so that we don't have to keep the lock
	propdefs := make(map[string]Constraint)
	o.proplock.Lock()
	for pname, c := range o.properties {
		propdefs[pname] = c
	}
	o.proplock.Unlock()

	for pname, c := range propdefs {
		if pdebug.Enabled {
			pdebug.Printf("Validating property '%s'", pname)
		}

		pval := rv.MapIndex(reflect.ValueOf(pname))
		if pval == zeroval {
			if pdebug.Enabled {
				pdebug.Printf("Property '%s' does not exist", pname)
			}
			if o.IsPropRequired(pname) { // required, and not present.
				return errors.New("object property '" + pname + "' is required")
			}

			// At this point we know that the property was not present
			// and that this field was indeed not required.
			if !c.HasDefault() {
				// We have no default. we can safely continue
				continue
			}

			// We have default
			dv := c.DefaultValue()
			pval = reflect.ValueOf(dv)
		}

		// delete from remaining props
		delete(premain, pname)
		// ...and add to props that we have seen
		pseen[pname] = struct{}{}

		if err := c.Validate(pval.Interface()); err != nil {
			return errors.New("object property '" + pname + "' validation failed: " + err.Error())
		}
	}

	if len(premain) > 0 {
		c := o.additionalProperties
		if c == nil {
			return errors.New("additional items are not allowed")
		}

		for pname := range premain {
			pval := rv.MapIndex(reflect.ValueOf(pname))
			if err := c.Validate(pval.Interface()); err != nil {
				return errors.New("object property for '" + pname + "' validation failed: " + err.Error())
			}
		}
	}

	for pname := range pseen {
		if deps := o.GetPropDependencies(pname); len(deps) > 0 {
			for _, dep := range deps {
				if _, ok := pseen[dep]; !ok {
					return errors.New("required dependency '" + dep + "' is mising")
				}
			}

			// can't, and shouldn't do object validation after checking prop deps
			continue
		}

		/* Not implemented yet. do we want to? */
		/*
			if depschema := o.GetSchemaDependency(pname); depschema != nil {

			}
		*/
	}

	return nil
}
