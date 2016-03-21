package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"sort"
)

const preamble = `package jsval

import (
	"bytes"
	"encoding/json"
	"reflect"
	"time"
)

type ErrInvalidMaybeValue struct {
	Value interface{}
}

func (e ErrInvalidMaybeValue) Error() string {
	buf := bytes.Buffer{}
	buf.WriteString("invalid Maybe value: ")
	t := reflect.TypeOf(e.Value)
	switch t {
	case nil:
		buf.WriteString("(nil)")
	default:
		buf.WriteByte('(')
		buf.WriteString(t.String())
		buf.WriteByte(')')
	}

	return buf.String()
}

// Maybe is an interface that can be used for struct fields which
// want to differentiate between initialized and uninitialized state.
// For example, a string field, if uninitialized, will contain the zero
// value of "", but that empty string *could* be a valid value for
// our validation purposes.
//
// To differentiate between an uninitialized string and an empty string,
// you should wrap it with a wrapper that implements the Maybe interface
// and JSVal will do its best to figure this out
type Maybe interface {
	// Valid should return true if this value has been properly initialized.
	// If this returns false, JSVal will treat as if the field is has not been
	// provided at all.
	Valid() bool

	// Value should return whatever the underlying value is.
	Value() interface{}

	// Set sets a value to this Maybe value, and turns on the Valid flag.
	// An error may be returned if the value could not be set (e.g.
	// you provided a value with the wrong type)
	Set(interface{}) error

	// Reset clears the Maybe value, and sets the Valid flag to false.
	Reset()
}

type ValidFlag bool

func (v *ValidFlag) Reset() {
	*v = false
}

func (v ValidFlag) Valid() bool {
	return bool(v)
}

`

func main() {
	types := map[string]string{
		"String": "string",
		"Int":    "int64",
		"Uint":   "uint64",
		"Float":  "float64",
		"Bool":   "bool",
		"Time":   "time.Time",
	}

	typenames := make([]string, 0, len(types))
	for t := range types {
		typenames = append(typenames, t)
	}
	sort.Strings(typenames)

	buf := bytes.Buffer{}
	buf.WriteString(preamble)
	for _, t := range typenames {
		bt := types[t]

		fmt.Fprintf(&buf, "\n\ntype Maybe%s struct{", t)
		buf.WriteString("\nValidFlag")
		fmt.Fprintf(&buf, "\n%s %s", t, bt)
		buf.WriteString("\n}")
		fmt.Fprintf(&buf, "\n\nfunc (v *Maybe%s) Set(x interface{}) error {", t)
		// Numeric types are special, because they can be converted.
		// float64 is included in the int/uint because JSON uses float64
		// to express numeric values, and we work with a lot of JSON
		switch t {
		case "Int":
			buf.WriteString("\nswitch x.(type) {")
			for _, ct := range []string{"int", "int8", "int16", "int32", "float64"} {
				fmt.Fprintf(&buf, "\ncase %s:", ct)
				fmt.Fprintf(&buf, "\nv.%s = int64(x.(%s))", t, ct)
			}
			buf.WriteString("\ncase int64:")
			fmt.Fprintf(&buf, "\nv.%s = x.(int64)", t)
			buf.WriteString("\ndefault:")
			buf.WriteString("\nreturn ErrInvalidMaybeValue{Value: x}")
			buf.WriteString("\n}")
			buf.WriteString("\nv.ValidFlag = true")
			buf.WriteString("\nreturn nil")
			buf.WriteString("\n}")
		case "Uint":
			buf.WriteString("\nswitch x.(type) {")
			for _, ct := range []string{"uint", "uint8", "uint16", "uint32", "float64"} {
				fmt.Fprintf(&buf, "\ncase %s:", ct)
				fmt.Fprintf(&buf, "\nv.%s = uint64(x.(%s))", t, ct)
			}
			buf.WriteString("\ncase uint64:")
			fmt.Fprintf(&buf, "\nv.%s = x.(uint64)", t)
			buf.WriteString("\ndefault:")
			buf.WriteString("\nreturn ErrInvalidMaybeValue{Value: x}")
			buf.WriteString("\n}")
			buf.WriteString("\nv.ValidFlag = true")
			buf.WriteString("\nreturn nil")
			buf.WriteString("\n}")
		case "Float":
			buf.WriteString("\nswitch x.(type) {")
			buf.WriteString("\ncase float32:")
			fmt.Fprintf(&buf, "\nv.%s = float64(x.(float32))", t)
			buf.WriteString("\ncase float64:")
			fmt.Fprintf(&buf, "\nv.%s = x.(float64)", t)
			buf.WriteString("\ndefault:")
			buf.WriteString("\nreturn ErrInvalidMaybeValue{Value: x}")
			buf.WriteString("\n}")
			buf.WriteString("\nv.ValidFlag = true")
			buf.WriteString("\nreturn nil")
			buf.WriteString("\n}")
		default:
			fmt.Fprintf(&buf, "\ns, ok := x.(%s)", bt)
			buf.WriteString("\nif !ok {")
			buf.WriteString("\nreturn ErrInvalidMaybeValue{Value: x}")
			buf.WriteString("\n}")
			buf.WriteString("\nv.ValidFlag = true")
			fmt.Fprintf(&buf, "\nv.%s = s", t)
			buf.WriteString("\nreturn nil")
			buf.WriteString("\n}")
		}
		fmt.Fprintf(&buf, "\n\nfunc (v Maybe%s) Value() interface{} {", t)
		fmt.Fprintf(&buf, "\nreturn v.%s", t)
		buf.WriteString("\n}")
		fmt.Fprintf(&buf, "\n\nfunc (v Maybe%s) MarshalJSON() ([]byte, error) {", t)
		fmt.Fprintf(&buf, "\nreturn json.Marshal(v.%s)", t)
		buf.WriteString("\n}")
		fmt.Fprintf(&buf, "\n\nfunc (v *Maybe%s) UnmarshalJSON(data []byte) error {", t)
		fmt.Fprintf(&buf, "\nvar in %s", bt)
		buf.WriteString("\nif err := json.Unmarshal(data, &in); err != nil {")
		buf.WriteString("\nreturn err")
		buf.WriteString("\n}")
		buf.WriteString("\nreturn v.Set(in)")
		buf.WriteString("\n}")
	}

	fsrc, err := format.Source(buf.Bytes())
	if err != nil {
		log.Printf("Error formatting: %s", err)
	}

	fmt.Printf("%s", fsrc)
}