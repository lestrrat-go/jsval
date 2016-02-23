package jsval

import (
	"errors"

	"github.com/lestrrat/go-jsschema"
)

func FromSchema(s *schema.Schema) (Constraint, error) {
	clist := make([]Constraint, len(s.Type))
	for i, st := range s.Type {
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
		if err := c.FromSchema(s); err != nil {
			return nil, err
		}
		clist[i] = c
	}

	switch len(clist) {
	case 0:
		return NilConstraint, nil
	case 1:
		return clist[0], nil
	default:
		return Any(clist), nil
	}
}

func matchenum(v interface{}, values []interface{}) bool {
	for _, x := range values {
		if x == v {
			return true
		}
	}
	return false
}
