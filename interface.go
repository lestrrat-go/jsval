package jsval

import (
	"reflect"
	"regexp"
	"sync"

	"github.com/lestrrat/go-jsschema"
	"github.com/lestrrat/go-structinfo"
)

var zeroval = reflect.Value{}

type Constraint interface {
	DefaultValue() interface{}
	FromSchema(s *schema.Schema) error
	HasDefault() bool
	IsRequired() bool
	Required(bool)
	Validate(interface{}) error
}

type nilConstraint struct{}

var NilConstraint = nilConstraint{}

type required bool
type defaultValue struct {
	initialized bool
	value       interface{}
}

type BooleanConstraint struct {
	defaultValue
	required
}

type StringConstraint struct {
	defaultValue
	required
	maxLength int
	minLength int
	regexp    *regexp.Regexp
	enum      []interface{}
}

type NumberConstraint struct {
	defaultValue
	required
	applyMinimum bool
	applyMaximum bool
	minimum      float64
	maximum      float64
	enum         []interface{}
}

type IntegerConstraint struct {
	NumberConstraint
}

type ArrayConstraint struct {
	defaultValue
	required
	itemspec        Constraint
	positionalItems []Constraint
	additionalItems Constraint
}

var DefaultFieldNamesFromStruct = structinfo.JSONFieldsFromStruct

type ObjectConstraint struct {
	defaultValue
	required
	lock                 sync.Mutex
	properties           map[string]Constraint
	additionalItems      Constraint
	FieldNamesFromStruct func(reflect.Value) []string
}

type AnyConstraint struct {
	nilConstraint
	constraints []Constraint
}
