package jsval

import (
	"errors"

	"github.com/lestrrat/go-pdebug"
)

func New() *JSVal {
	return &JSVal{
		refs: make(map[string]Constraint),
	}
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

func (v *JSVal) GetReference(ref string) (Constraint, error) {
	v.reflock.Lock()
	defer v.reflock.Unlock()
	c, ok := v.refs[ref]
	if !ok {
		return nil, errors.New("reference '" + ref + "' not found")
	}

	return c, nil
}

func (v *JSVal) SetReference(ref string, c Constraint) {
	if pdebug.Enabled {
		pdebug.Printf("JSVal.SetReference %s", ref)
	}

	v.reflock.Lock()
	defer v.reflock.Unlock()
	v.refs[ref] = c
}
