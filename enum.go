package jsval

import (
	"errors"

	"github.com/lestrrat/go-pdebug"
)

// Enum creates a new EnumConstraint
func Enum(v ...interface{}) *EnumConstraint {
	return &EnumConstraint{enums: v}
}

// Enum method sets the possible enumerations
func (c *EnumConstraint) Enum(v []interface{}) *EnumConstraint {
	c.enums = v
	return c
}

// Validate validates the value against the list of enumerations
func (c *EnumConstraint) Validate(v interface{}) (err error) {
	if pdebug.Enabled {
		g := pdebug.IPrintf("START EnumConstraint.Validate")
		defer func() {
			if err == nil {
				g.IRelease("END EnumConstraint.Validate (PASS)")
			} else {
				g.IRelease("END EnumConstraint.Validate (FAIL): %s", err)
			}
		}()
	}
	for _, e := range c.enums {
		if e == v {
			return nil
		}
	}
	return errors.New("value is not in enumeration")
}
