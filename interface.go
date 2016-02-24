package jsval

import (
	"reflect"
	"regexp"
	"sync"

	"github.com/lestrrat/go-jsref"
	"github.com/lestrrat/go-jsschema"
	"github.com/lestrrat/go-structinfo"
)

var zeroval = reflect.Value{}

type Validator interface {
	Validate(interface{}) error
}

type JSVal struct {
	root    Constraint
	reflock sync.Mutex
	refs    map[string]Constraint
	resolver *jsref.Resolver
}

type Constraint interface {
	buildFromSchema(*buildctx, *schema.Schema) error
	DefaultValue() interface{}
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
	enums     *EnumConstraint
	maxLength int
	minLength int
	regexp    *regexp.Regexp
}

type NumberConstraint struct {
	defaultValue
	required
	applyMinimum     bool
	applyMaximum     bool
	minimum          float64
	maximum          float64
	exclusiveMinimum bool
	exclusiveMaximum bool
	enums            *EnumConstraint
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
	minItems        int
	maxItems        int
	uniqueItems     bool
}

var DefaultFieldNamesFromStruct = structinfo.JSONFieldsFromStruct

type ObjectConstraint struct {
	defaultValue
	required
	lock                 sync.Mutex
	properties           map[string]Constraint
	additionalProperties Constraint
	FieldNamesFromStruct func(reflect.Value) []string
}

type EnumConstraint struct {
	nilConstraint
	enums []interface{}
}

type CombinedConstraint interface {
	Constraint
	Constraints() []Constraint
}

type comboconstraint struct {
	nilConstraint
	constraints []Constraint
}

type AnyConstraint struct {
	comboconstraint
}

type AllConstraint struct {
	comboconstraint
}
