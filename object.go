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
	for pname, pdef := range s.Properties {
		cprop, err := buildFromSchema(ctx, pdef)
		if err != nil {
			return err
		}

		if s.IsPropRequired(pname) {
			cprop.Required(true)
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
	return nil
}

func Object() *ObjectConstraint {
	return &ObjectConstraint{
		properties:           map[string]Constraint{},
		additionalProperties: nil,
	}
}

func (o *ObjectConstraint) AdditionalProperties(c Constraint) *ObjectConstraint {
	o.additionalProperties = c
	return o
}

func (o *ObjectConstraint) AddProp(name string, c Constraint) *ObjectConstraint {
	o.lock.Lock()
	defer o.lock.Unlock()

	o.properties[name] = c
	return o
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
	props := map[string]struct{}{}
	for _, k := range fields {
		props[k] = struct{}{}
	}

	// Now, for all known constraints, validate the prop
	// create a copy of properties so that we don't have to keep the lock
	propdefs := make(map[string]Constraint)
	o.lock.Lock()
	for pname, c := range o.properties {
		propdefs[pname] = c
	}
	o.lock.Unlock()

	for pname, c := range propdefs {
		if pdebug.Enabled {
			pdebug.Printf("Validating property '%s'", pname)
		}
		pval := rv.MapIndex(reflect.ValueOf(pname))
		if pval == zeroval {
			if pdebug.Enabled {
				pdebug.Printf("Property '%s' does not exist", pname)
			}
			if c.IsRequired() { // required, and not present.
				return errors.New("object property '" + pname + "' is required")
			}
			if c.HasDefault() { // check default...
				dv := c.DefaultValue()
				pval = reflect.ValueOf(dv)
			}

			// tricky! this field must be deleted from the props map before
			// going into the next iteration
			delete(props, pname)
			continue
		}
		delete(props, pname)

		if err := c.Validate(pval.Interface()); err != nil {
			return errors.New("object property '" + pname + "' validation failed: " + err.Error())
		}
	}

	if len(props) > 0 {
		c := o.additionalProperties
		if c == nil {
			return errors.New("additional items are not allowed")
		}

		for pname := range props {
			pval := rv.MapIndex(reflect.ValueOf(pname))
			if err := c.Validate(pval.Interface()); err != nil {
				return errors.New("object property for '" + pname + "' validation failed: " + err.Error())
			}
		}
	}
	return nil
}
