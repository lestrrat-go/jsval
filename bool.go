package jsval

import (
	"errors"
	"reflect"

	"github.com/lestrrat/go-jsschema"
)

func (c *BooleanConstraint) buildFromSchema(_ *buildctx, _ *schema.Schema) error {
	return nil
}

func Boolean() *BooleanConstraint {
	return &BooleanConstraint{}
}

func (b *BooleanConstraint) Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Bool:
	default:
		return errors.New("value is not a boolean")
	}
	return nil
}
