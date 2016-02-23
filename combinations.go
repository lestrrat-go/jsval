package jsval

import "errors"

func Any(clist []Constraint) *AnyConstraint {
	l := make([]Constraint, len(clist))
	copy(l, clist)
	return &AnyConstraint{
		constraints: l,
	}
}

func (c *AnyConstraint) Validate(v interface{}) error {
	for _, celem := range c.constraints {
		if err := celem.Validate(v); err == nil {
			return nil
		}
	}
	return errors.New("could not validate against any of the constraints")
}
