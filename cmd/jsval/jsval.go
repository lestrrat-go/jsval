package main

import (
	"io"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/lestrrat/go-jsschema"
	"github.com/lestrrat/go-jsval"
)

func main() {
	os.Exit(_main())
}

// jsval schema.json [ref]

type Options struct {
	Schema  string `short:"s" long:"schema" description:"the source JSON schema file"`
	OutFile string `short:"o" long:"outfile" description:"output file to generate"`
	Name    string `short:"n" long:"name" description:"name of the function"`
}

func _main() int {
	var opts Options
	if _, err := flags.Parse(&opts); err != nil {
		log.Printf("%s", err)
		return 1
	}

	s, err := schema.ReadFile(opts.Schema)
	if err != nil {
		log.Printf("%s", err)
		return 1
	}

	v := jsval.New()
	if err := v.Build(s); err != nil {
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