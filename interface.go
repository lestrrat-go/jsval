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
	root     Constraint
	reflock  sync.Mutex
	refs     map[string]Constraint
	resolver *jsref.Resolver
}

type Constraint interface {
	buildFromSchema(*buildctx, *schema.Schema) error
	DefaultValue() interface{}
	HasDefault() bool
	Validate(interface{}) error
}

type nilConstraint struct{}

var NilConstraint = nilConstraint{}

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
	format    schema.Format
}

type NumberConstraint struct {
	defaultValue
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
	additionalProperties Constraint
	deplock              sync.Mutex
	proplock             sync.Mutex
	properties           map[string]Constraint
	propdeps             map[string][]string
	reqlock              sync.Mutex
	required             map[string]struct{}
	schemadeps           map[string]Constraint
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
