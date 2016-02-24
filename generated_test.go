package jsval_test

import "testing"

func TestGenerated(t *testing.T) {
	V := JSValFoo()
	err := V.Validate(map[string]interface{}{
		"minItems": -1,
	})
	if err == nil {
		t.Errorf("Validation failed: %s", err)
	}
}
