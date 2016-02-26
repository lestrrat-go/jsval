# go-jsval

[![Build Status](https://travis-ci.org/lestrrat/go-jsval.svg?branch=master)](https://travis-ci.org/lestrrat/go-jsval)

[![GoDoc](https://godoc.org/github.com/lestrrat/go-jsval?status.svg)](https://godoc.org/github.com/lestrrat/go-jsval)

Validator toolset, aimed to be used with JSON Schema

# Description

The `go-jsval` package is a data validation toolset, with
a tool to generate validators in Go from JSON schemas.

# Synopsis

Read a schema file and create a validator:

```go
package jsval_test

import (
  "log"

  "github.com/lestrrat/go-jsschema"
  "github.com/lestrrat/go-jsval/builder"
)

func ExampleBuild() {
  s, err := schema.ReadFile(`/path/to/schema.json`)
  if err != nil {
    log.Printf("failed to open schema: %s", err)
    return
  }

  b := builder.New()
  v, err := b.Build(s)
  if err != nil {
    log.Printf("failed to build validator: %s", err)
    return
  }

  var input interface{}
  if err := v.Validate(input); err != nil {
    log.Printf("validation failed: %s", err)
    return
  }
}
```

Build a validator by hand:

```go
func ExampleManual() {
  v := jsval.Object().
    AddProp(`zip`, jsval.String().RegexpString(`^\d{5}$`)).
    AddProp(`address`, jsval.String()).
    AddProp(`name`, jsval.String()).
    AddProp(`phone_number`, jsval.String().RegexpString(`^[\d-]+$`)).
    Required(`zip`, `address`, `name`)

  var input interface{}
  if err := v.Validate(input); err != nil {
    log.Printf("validation failed: %s", err)
    return
  }
}
```

# Install

```
go get -u github.com/lestrrat/go-jsval
```

If you want to install the `jsval` tool, do

```
go get -u github.com/lestrrat/go-jsval/cmd/jsval
```

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

If your document isn't a real JSON schema, but contains a
JSON schema (like JSON Hyper Schema), you can use the `-ptr`
argument to access a specific portion of a JSON document:

```
jsval -schema hyper.json -name Foo -ptr "#/links/0"
```

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

# References

| Name                                                     | Notes                            |
|:--------------------------------------------------------:|:---------------------------------|
| [go-jsschema](https://github.com/lestrrat/go-jsschema)   | JSON Schema implementation       |
| [go-jshschema](https://github.com/lestrrat/go-jshschema) | JSON Hyper Schema implementation |
| [go-jsref](https://github.com/lestrrat/go-jsref)         | JSON Reference implementation    |
| [go-jspointer](https://github.com/lestrrat/go-jspointer) | JSON Pointer implementations     |

