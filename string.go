package jsval

import (
	"errors"
	"net"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"time"

	"github.com/lestrrat/go-pdebug"
)

func (sc *StringConstraint) Default(v interface{}) *StringConstraint {
	sc.defaultValue.initialized = true
	sc.defaultValue.value = v
	return sc
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

	switch s.format{
	case "datetime":
		if _, err = time.Parse(time.RFC3339, str); err != nil {
			return errors.New("invalid datetime")
		}
	case "email":
		if _, err = mail.ParseAddress(str); err != nil {
			return errors.New("invalid email address: " + err.Error())
		}
	case "hostname":
		if !isDomainName(str) {
			return errors.New("invalid hostname")
		}
	case "ipv4":
		// Should only contain numbers and "."
		for _, r := range str {
			switch {
			case r == 0x2E || 0x30 <= r && r <= 0x39:
			default:
				return errors.New("invalid IPv4 address")
			}
		}
		if addr := net.ParseIP(str); addr == nil {
			return errors.New("invalid IPv4 address")
		}
	case "ipv6":
		// Should only contain numbers and ":"
		for _, r := range str {
			switch {
			case r == 0x3A || 0x30 <= r && r <= 0x39:
			default:
				return errors.New("invalid IPv6 address")
			}
		}
		if addr := net.ParseIP(str); addr == nil {
			return errors.New("invalid IPv6 address")
		}
	case "uri":
		if _, err = url.Parse(str); err != nil {
			return errors.New("invalid URI")
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

// stolen from src/net/dnsclient.go
func isDomainName(s string) bool {
	// See RFC 1035, RFC 3696.
	if len(s) == 0 {
		return false
	}
	if len(s) > 255 {
		return false
	}

	last := byte('.')
	ok := false // Ok once we've seen a letter.
	partlen := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		default:
			return false
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_':
			ok = true
			partlen++
		case '0' <= c && c <= '9':
			// fine
			partlen++
		case c == '-':
			// Byte before dash cannot be dot.
			if last == '.' {
				return false
			}
			partlen++
		case c == '.':
			// Byte before dot cannot be dot, dash.
			if last == '.' || last == '-' {
				return false
			}
			if partlen > 63 || partlen == 0 {
				return false
			}
			partlen = 0
		}
		last = c
	}
	if last == '-' || partlen > 63 {
		return false
	}

	return ok
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

func (sc *StringConstraint) Format(f string) *StringConstraint {
	sc.format = f
	return sc
}

func String() *StringConstraint {
	return &StringConstraint{
		maxLength: -1,
	}
}
