package jsval

import (
	"strings"
	"testing"

	"github.com/lestrrat/go-jsschema"
	"github.com/stretchr/testify/assert"
)

func TestObject(t *testing.T) {
	const src = `{
  "type": "object",
  "additionalItems": false,
  "properties": {
    "name": {
      "type": "string",
      "maxLength": 20,
      "pattern": "^[a-z ]+$"
    },
	  "age": {
		  "type": "integer",
	    "minimum": 0
	  },
	  "tags": {
      "type": "array",
	    "items": {
        "type": "string"
      }
    }
  }
}`

	s, err := schema.Read(strings.NewReader(src))
	if !assert.NoError(t, err, "reading schema should succeed") {
		return
	}

	c := Object()
	if !assert.NoError(t, c.FromSchema(s), "Object.FromSchema should succeed") {
		return
	}

	data := []interface{}{
		map[string]interface{}{"Name": "World"},
		map[string]interface{}{"name": "World"},
		map[string]interface{}{"name": "wooooooooooooooooooooooooooooooorld"},
		map[string]interface{}{
			"tags": []interface{}{ 1, "foo", false },
		},
	}
	for _, input := range data {
		t.Logf("Testing %#v (should FAIL)", input)
		if !assert.Error(t, c.Validate(input), "validation fails") {
			return
		}
	}

	data = []interface{}{
		map[string]interface{}{"name": "world"},
		map[string]interface{}{"tags": []interface{}{"foo", "bar", "baz"}},
	}
	for _, input := range data {
		t.Logf("Testing %#v (should PASS)", input)
		if !assert.NoError(t, c.Validate(input), "validation passes") {
			return
		}
	}
}

