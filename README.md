# go-jsval

[![Build Status](https://travis-ci.org/lestrrat/go-jsval.svg?branch=master)](https://travis-ci.org/lestrrat/go-jsval)

[![GoDoc](https://godoc.org/github.com/lestrrat/go-jsval?status.svg)](https://godoc.org/github.com/lestrrat/go-jsval)

Validator toolset, aimed to be used with JSON Schema

# Description

The `go-jsval` package is a 

# Features

## Can generate validators from JSON Schema definition

The following command creates a file named `foo_jsval.go` 
which contains a function named `JSvalFoo()`, which
returns a validator created from the the schema:

```
jsval -schema schema.json -name Foo -o foo_jsval.go
```

See the file `generated_validator_test.go` for a sample
generated from JSON Schema schema.

## Can handle JSON References in JSON Schema definitions

Note: Not very well tested. Test cases welcome

This packages tries to handle JSON References properly.
For example, in the schema below, "age" input is validated
against the `positiveInteger` schema:

```json
{
  "definitions": {
    "positiveInteger": {
      "type": "integer",
      "minimum": 0,
    }
  },
  "properties": {
    "age": { "$ref": "#/definitions/positiveInteger" }
  }
}
```

# TODO

* More complete coverage of JSON Schema. Many validation statements are still not implmented (Please file issues if you find any!)
* Work with schemas that does not present a "type" field.

