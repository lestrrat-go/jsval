package jsval

import (
	"errors"
	"sync"

	"github.com/lestrrat/go-jsschema"
	"github.com/lestrrat/go-pdebug"
)

// ReferenceConstraint is a constraint where its actual definition
// is stored elsewhere.
type ReferenceConstraint struct {
	V         *JSVal
	lock      sync.Mutex
	resolved  Constraint
	reference string
}

func Reference(v *JSVal) *ReferenceConstraint {
	return &ReferenceConstraint{
		V: v,
	}
}

func (r *ReferenceConstraint) buildFromSchema(ctx *buildctx, s *schema.Schema) error {
	pdebug.Printf("ReferenceConstraint.buildFromSchema '%s'", s.Reference)
	if s.Reference == "" {
		return errors.New("schema does not contain a reference")
	}
	r.reference = s.Reference

	return nil
}
func (r *ReferenceConstraint) Resolved() (c Constraint, err error) {
	if pdebug.Enabled {
		g := pdebug.IPrintf("START ReferenceConstraint.Resolved '%s'", r.reference)
		defer func() {
			if err == nil {
				g.IRelease("END ReferenceConstraint.Resolved '%s' (OK)", r.reference)
			} else {
				g.IRelease("END ReferenceConstraint.Resolved '%s' (FAIL): %s", r.reference, err)
			}
		}()
	}
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.resolved != nil {
		if pdebug.Enabled {
			pdebug.Printf("Reference is already resolved")
		}
		return r.resolved, nil
	}

	c, err = r.V.GetReference(r.reference)
	if err != nil {
		return nil, err
	}
	r.resolved = c
	return c, nil
}

func (r *ReferenceConstraint) RefersTo(s string) *ReferenceConstraint {
	r.reference = s
	return r
}

func (r *ReferenceConstraint) Default(_ interface{}) {
}

func (r *ReferenceConstraint) DefaultValue() interface{} {
	c, err := r.Resolved()
	if err != nil {
		return nil
	}
	return c.DefaultValue()
}

func (r *ReferenceConstraint) HasDefault() bool {
	c, err := r.Resolved()
	if err != nil {
		return false
	}
	return c.HasDefault()
}

func (r *ReferenceConstraint) Required(_ bool) {
}

func (r *ReferenceConstraint) IsRequired() bool {
	c, err := r.Resolved()
	if err != nil {
		return false
	}
	return c.IsRequired()
}

func (r *ReferenceConstraint) Validate(v interface{}) (err error) {
	if pdebug.Enabled {
		g := pdebug.IPrintf("START ReferenceConstraint.Validate")
		defer func() {
			if err == nil {
				g.IRelease("END ReferenceConstraint.Validate (PASS)")
			} else {
				g.IRelease("END ReferenceConstraint.Validate (FAIL): %s", err)
			}
		}()
	}

	c, err := r.Resolved()
	if err != nil {
		return err
	}
	return c.Validate(v)
}