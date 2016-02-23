package jsval

import (
	"errors"
	"reflect"
	"regexp"

	"github.com/lestrrat/go-jsschema"
)

func (sc *StringConstraint) Default(v interface{}) *StringConstraint {
	sc.defaultValue.initialized = true
	sc.defaultValue.value = v
	return sc
}

func (c *StringConstraint) FromSchema(s *schema.Schema) error {
	if !s.Type.Contains(schema.StringType) {
		return errors.New("schema is not for string")
	}

	if s.MaxLength.Initialized {
		c.MaxLength(s.MaxLength.Val)
	}

	if s.MinLength.Initialized {
		c.MinLength(s.MinLength.Val)
	}

	if pat := s.Pattern; pat != nil {
		c.Regexp(pat)
	}

	if lst := s.Enum; len(lst) > 0 {
		c.Enum(lst)
	}

	if v := s.Default; v != nil {
		c.Default(v)
	}

	return nil
}

// Note that StringConstraint does not apply default values to the
// incoming string value, because the Zero value for string ("")
// can be a perfectly reasonable value.
//
// The caller is the only person who can determine if a string
// value is "unavailable"
func (s *StringConstraint) Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface:
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.String:
	default:
		return errors.New("value is not a string")
	}

	str := rv.String()
	ls := len(str)
	if s.maxLength > 0 && ls > s.maxLength {
		return errors.New("string longer than maxLength")
	}

	if s.minLength > -1 && ls < s.minLength {
		return errors.New("string shorter than minLength")
	}

	if rx := s.regexp; rx != nil {
		if !rx.MatchString(str) {
			return errors.New("string does not match regular expression")
		}
	}

	if enum := s.enum; enum != nil {
		if !matchenum(str, enum) {
			return errors.New("value not in enumeration")
		}
	}

	return nil
}

func (sc *StringConstraint) Enum(l []interface{}) *StringConstraint {
	sc.enum = l
	return sc
}

func (sc *StringConstraint) MaxLength(l int) *StringConstraint {
	sc.maxLength = l
	return sc
}

func (sc *StringConstraint) MinLength(l int) *StringConstraint {
	sc.minLength = l
	return sc
}

func (sc *StringConstraint) RegexpString(pat string) *StringConstraint {
	return sc.Regexp(regexp.MustCompile(pat))
}

func (sc *StringConstraint) Regexp(rx *regexp.Regexp) *StringConstraint {
	sc.regexp = rx
	return sc
}

func String() *StringConstraint {
	return &StringConstraint{
		maxLength: -1,
	}
}
