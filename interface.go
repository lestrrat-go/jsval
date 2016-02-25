package jsval

import (
	"reflect"
	"regexp"
	"sync"

	"github.com/lestrrat/go-jsref"
	"github.com/lestrrat/go-structinfo"
)

var zeroval = reflect.Value{}

type Validator interface {
	Validate(interface{}) error
}

type JSVal struct {
	root     Constraint
	reflock  sync.Mutex
	refs     map[string]Constraint
	resolver *jsref.Resolver
}

type Constraint interface {
	DefaultValue() interface{}
	HasDefault() bool
	Validate(interface{}) error
}

type emptyConstraint struct{}

// EmptyConstraint is a constraint that returns true for any value
var EmptyConstraint = emptyConstraint{}

type nullConstraint struct{}

// NullConstraint is a constraint that only matches the JSON
// "null" value, or "nil" in golang
var NullConstraint = nullConstraint{}

type defaultValue struct {
	initialized bool
	value       interface{}
}

type BooleanConstraint struct {
	defaultValue
}

type StringConstraint struct {
	defaultValue
	enums     *EnumConstraint
	maxLength int
	minLength int
	regexp    *regexp.Regexp
	format    string
}

type NumberConstraint struct {
	defaultValue
	applyMinimum     bool
	applyMaximum     bool
	applyMultipleOf  bool
	minimum          float64
	maximum          float64
	multipleOf       float64
	exclusiveMinimum bool
	exclusiveMaximum bool
	enums            *EnumConstraint
}

type IntegerConstraint struct {
	NumberConstraint
}

type ArrayConstraint struct {
	defaultValue
	items           Constraint
	positionalItems []Constraint
	additionalItems Constraint
	minItems        int
	maxItems        int
	uniqueItems     bool
}

var DefaultFieldNamesFromStruct = structinfo.JSONFieldsFromStruct

type ObjectConstraint struct {
	defaultValue
	additionalProperties Constraint
	deplock              sync.Mutex
	patternProperties    map[*regexp.Regexp]Constraint
	proplock             sync.Mutex
	properties           map[string]Constraint
	propdeps             map[string][]string
	reqlock              sync.Mutex
	required             map[string]struct{}
	maxProperties        int
	minProperties        int
	schemadeps           map[string]Constraint
	FieldNamesFromStruct func(reflect.Value) []string
}

type EnumConstraint struct {
	emptyConstraint
	enums []interface{}
}

type CombinedConstraint interface {
	Constraint
	Constraints() []Constraint
}

type comboconstraint struct {
	emptyConstraint
	constraints []Constraint
}

type AnyConstraint struct {
	comboconstraint
}

type AllConstraint struct {
	comboconstraint
}

type OneOfConstraint struct {
	comboconstraint
}

type NotConstraint struct {
	child Constraint
}
