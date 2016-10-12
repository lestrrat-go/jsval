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

	var s TestMaybeStruct
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

	var s TestMaybeStruct
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

	var s TestMaybeStruct
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

func TestMaybeInt(t *testing.T) {
	var i jsval.MaybeInt

	if !assert.NoError(t, i.Set(10), "const 10 can be set to MaybeInt (coersion takes place)") {
		return
	}

	if !assert.NoError(t, i.Set(10.0), "const 10.0 can be set to MaybeInt (coersion takes place)") {
		return
	}
}
