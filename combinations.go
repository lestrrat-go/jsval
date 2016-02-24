package jsval

import "errors"

func (c *comboconstraint) Add(v Constraint) {
	c.constraints = append(c.constraints, v)
}

func (c *comboconstraint) Constraints() []Constraint {
	return c.constraints
}

func reduceCombined(cc CombinedConstraint) Constraint {
	l := cc.Constraints()
	if len(l) == 1 {
		return l[0]
	}
	return cc
}

func Any() *AnyConstraint {
	return &AnyConstraint{}
}

func (c *AnyConstraint) Reduce() Constraint {
	return reduceCombined(c)
}

func (c *AnyConstraint) Add(c2 Constraint) *AnyConstraint {
	c.comboconstraint.Add(c2)
	return c
}

func (c *AnyConstraint) Validate(v interface{}) error {
	for _, celem := range c.constraints {
		if err := celem.Validate(v); err == nil {
			return nil
		}
	}
	return errors.New("could not validate against any of the constraints")
}

func All() *AllConstraint {
	return &AllConstraint{}
}

func (c *AllConstraint) Reduce() Constraint {
	return reduceCombined(c)
}

func (c *AllConstraint) Add(c2 Constraint) *AllConstraint {
	c.comboconstraint.Add(c2)
	return c
}

func (c *AllConstraint) Validate(v interface{}) error {
	for _, celem := range c.constraints {
		if err := celem.Validate(v); err != nil {
			return err
		}
	}
	return nil
}

