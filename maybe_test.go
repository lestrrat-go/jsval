package jsval_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/lestrrat/go-jsval"
	"github.com/stretchr/testify/assert"
)

type TestMaybeStruct struct {
	Name jsval.MaybeString `json:"name"`
	Age int `json:"age"`
}

func TestMaybeString_Empty(t *testing.T) {
	const src = `{"age": 10}`

	s := TestMaybeStruct{}
	if !assert.NoError(t, json.NewDecoder(strings.NewReader(src)).Decode(&s), "Decode works") {
		return
	}

	v := jsval.New().SetRoot(
		jsval.Object().
			AddProp("name", jsval.String()).
			AddProp("age", jsval.Integer()),
	)
	if !assert.NoError(t, v.Validate(&s), "Validate succeeds") {
		return
	}
}

func TestMaybeString_Populated(t *testing.T) {
	const src = `{"age": 10, "name": "John Doe"}`

	s := TestMaybeStruct{}
	if !assert.NoError(t, json.NewDecoder(strings.NewReader(src)).Decode(&s), "Decode works") {
		return
	}

	v := jsval.New().SetRoot(
		jsval.Object().
			AddProp("name", jsval.String()).
			AddProp("age", jsval.Integer()),
	)
	if !assert.NoError(t, v.Validate(&s), "Validate succeeds") {
		return
	}
}

func TestMaybeString_EmptyDefault(t *testing.T) {
	const src = `{"age": 10}`

	s := TestMaybeStruct{}
	if !assert.NoError(t, json.NewDecoder(strings.NewReader(src)).Decode(&s), "Decode works") {
		return
	}

	v := jsval.New().SetRoot(
		jsval.Object().
			AddProp("name", jsval.String().Default("John Doe")).
			AddProp("age", jsval.Integer()),
	)
	if !assert.NoError(t, v.Validate(&s), "Validate succeeds") {
		return
	}

	if !assert.Equal(t, s.Name.Value().(string), "John Doe", "Should have default value") {
		return
	}
}
