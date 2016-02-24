package jsval

import (
	"errors"
	"reflect"
	"regexp"

	"github.com/lestrrat/go-jsschema"
	"github.com/lestrrat/go-pdebug"
)

func (sc *StringConstraint) Default(v interface{}) *StringConstraint {
	sc.defaultValue.initialized = true
	sc.defaultValue.value = v
	return sc
}

func (c *StringConstraint) buildFromSchema(ctx *buildctx, s *schema.Schema) error {
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
func (s *StringConstraint) Validate(v interface{}) (err error) {
	if pdebug.Enabled {
		g := pdebug.IPrintf("START StringConstraint.Validate")
		defer func() {
			if err == nil {
				g.IRelease("END StringConstraint.Validate (PASS)")
			} else {
				g.IRelease("END StringConstraint.Validate (FAIL): %s", err)
			}
		}()
	}
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
	if s.maxLength > 0 {
		if pdebug.Enabled {
			pdebug.Printf("Checking MaxLength (%d)", s.maxLength)
		}
		if ls > s.maxLength {
			return errors.New("string longer than maxLength")
		}
	}

	if s.minLength > -1 {
		if pdebug.Enabled {
			pdebug.Printf("Checking MinLength (%d)", s.minLength)
		}
		if ls < s.minLength {
			return errors.New("string shorter than minLength")
		}
	}

	if rx := s.regexp; rx != nil {
		if pdebug.Enabled {
			pdebug.Printf("Checking Regexp")
		}
		if !rx.MatchString(str) {
			return errors.New("string does not match regular expression")
		}
	}

	if enum := s.enums; enum != nil {
		if err := enum.Validate(str); err != nil {
			return err
		}
	}

	return nil
}

func (sc *StringConstraint) Enum(l []interface{}) *StringConstraint {
	if sc.enums == nil {
		sc.enums = Enum()
	}
	sc.enums.Enum(l)
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
