package jsval

import (
	"strings"
	"testing"

	"github.com/lestrrat/go-jsschema"
	"github.com/stretchr/testify/assert"
)

func TestNumberFromSchema(t *testing.T) {
	const src = `{
  "type": "number",
  "minimum": 5,
  "maximum": 15,
  "default": 10
}`

	s, err := schema.Read(strings.NewReader(src))
	if !assert.NoError(t, err, "schema.Read should succeed") {
		return
	}

	v := New()
	if !assert.NoError(t, v.Build(s), "Validator.Build should succeed") {
		return
	}

	c2 := Number()
	c2.Default(float64(10)).Maximum(15).Minimum(5)
	if !assert.Equal(t, c2, v.root, "constraints are equal") {
		return
	}
}

func TestNumber(t *testing.T) {
	c := Number()
	c.Default(float64(10)).Maximum(15)

	if !assert.True(t, c.HasDefault(), "HasDefault is true") {
		return
	}

	if !assert.Equal(t, c.DefaultValue(), float64(10), "DefaultValue returns expected value") {
		return
	}

	var s float64
	if !assert.NoError(t, c.Validate(s), "validate should succeed") {
		return
	}

	c.Minimum(5)
	if !assert.Error(t, c.Validate(s), "validate should fail") {
		return
	}

	s = 10
	if !assert.NoError(t, c.Validate(s), "validate should succeed") {
		return
	}
}