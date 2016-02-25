package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/lestrrat/go-jsref"
	"github.com/lestrrat/go-jsschema"
	"github.com/lestrrat/go-jsval"
	"github.com/lestrrat/go-jsval/builder"
)

func main() {
	os.Exit(_main())
}

// jsval schema.json [ref]

type Options struct {
	Schema  string `short:"s" long:"schema" description:"the source JSON schema file"`
	OutFile string `short:"o" long:"outfile" description:"output file to generate"`
	Name    string `short:"n" long:"name" description:"name of the function"`
	Pointer string `short:"p" long:"ptr" description:"JSON pointer within the document"`
}

func _main() int {
	var opts Options
	if _, err := flags.Parse(&opts); err != nil {
		log.Printf("%s", err)
		return 1
	}

	f, err := os.Open(opts.Schema)
	if err != nil {
		log.Printf("%s", err)
		return 1
	}
	defer f.Close()

	var m map[string]interface{}
	if err := json.NewDecoder(f).Decode(&m); err != nil {
		log.Printf("%s", err)
		return 1
	}

	var s *schema.Schema
	if ptr := opts.Pointer; ptr != "" {
		resolver := jsref.New()
		x, err := resolver.Resolve(m, ptr)
		if err != nil {
			log.Printf("%s", err)
			return 1
		}

		m2, ok := x.(map[string]interface{})
		if !ok {
			log.Printf("Expected map")
			return 1
		}

		s = schema.New()
		if err := s.Extract(m2); err != nil {
			log.Printf("%s", err)
			return 1
		}
	} else {
		s, err = schema.ReadFile(opts.Schema)
		if err != nil {
			log.Printf("%s", err)
			return 1
		}
	}

	b := builder.New()
	v, err := b.Build(s)
	if err != nil {
		log.Printf("%s", err)
		return 1
	}

	var out io.Writer

	out = os.Stdout
	if fn := opts.OutFile; fn != "" {
		f, err := os.Create(fn)
		if err != nil {
			log.Printf("%s", err)
			return 1
		}
		defer f.Close()

		out = f
	}

	g := jsval.NewGenerator()
	if err := g.Process(out, v, opts.Name); err != nil {
		log.Printf("%s", err)
		return 1
	}

	return 0
}