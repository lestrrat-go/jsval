package jsval

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Process(out io.Writer, v *JSVal, name string) error {
	ctx := genctx{
		pkgname: "jsval",
		vname:   "V",
	}

	fmt.Fprintf(out, "func JSVal%s() *%s.JSVal {\n", name, ctx.pkgname)
	g1 := ctx.Indent()
	defer g1()
	generateCode(&ctx, out, v)
	fmt.Fprintf(out, "\n\treturn %s\n}\n", ctx.vname)

	return nil
}

type genctx struct {
	prefix  []byte
	pkgname string
	vname   string
}

func (ctx *genctx) Prefix() []byte {
	return ctx.prefix
}

func (ctx *genctx) Indent() func() {
	ctx.prefix = append(ctx.prefix, '\t')
	return func() {
		l := len(ctx.prefix)
		if l == 0 {
			return
		}
		ctx.prefix = ctx.prefix[:l-1]
	}
}

func generateNilCode(ctx *genctx, out io.Writer, c nilConstraint) error {
	fmt.Fprintf(out, "%s%s.NilConstraint", ctx.Prefix(), ctx.pkgname)
	return nil
}
func generateValidatorCode(ctx *genctx, out io.Writer, v *JSVal) error {
	vname := ctx.vname
	p := ctx.Prefix()

	fmt.Fprintf(out, "%s%s := %s.New()", p, vname, ctx.pkgname)
	fmt.Fprintf(out, "\n%s%s.SetRoot(\n", p, vname)
	g := ctx.Indent()
	if err := generateCode(ctx, out, v.root); err != nil {
		g()
		return err
	}
	g()
	fmt.Fprintf(out, ",\n%s)\n", p)

	refs := make([]string, 0, len(v.refs))
	for ref := range v.refs {
		refs = append(refs, ref)
	}
	sort.Strings(refs)

	for _, ref := range refs {
		c := v.refs[ref]
		fmt.Fprintf(out, "\n%s%s.SetReference(\n%s\t`%s`,\n", p, vname, p, ref)
		g1 := ctx.Indent()
		if ref == "#" {
			fmt.Fprintf(out, "%s%s.Root()", ctx.Prefix(), ctx.vname)
		} else {
			if err := generateCode(ctx, out, c); err != nil {
				g1()
				return err
			}
		}
		fmt.Fprintf(out, ",\n%s)", p)
		g1()
	}

	return nil
}

func generateCode(ctx *genctx, out io.Writer, c Validator) error {
	buf := &bytes.Buffer{}

	switch c.(type) {
	case nilConstraint:
		generateNilCode(ctx, buf, c.(nilConstraint))
	case *JSVal:
		generateValidatorCode(ctx, buf, c.(*JSVal))
	case *AnyConstraint:
		generateAnyCode(ctx, buf, c.(*AnyConstraint))
	case *AllConstraint:
		generateAllCode(ctx, buf, c.(*AllConstraint))
	case *BooleanConstraint:
		generateBooleanCode(ctx, buf, c.(*BooleanConstraint))
	case *StringConstraint:
		generateStringCode(ctx, buf, c.(*StringConstraint))
	case *IntegerConstraint:
		generateIntegerCode(ctx, buf, c.(*IntegerConstraint))
	case *NumberConstraint:
		generateNumberCode(ctx, buf, c.(*NumberConstraint))
	case *ReferenceConstraint:
		generateReferenceCode(ctx, buf, c.(*ReferenceConstraint))
	case *ArrayConstraint:
		generateArrayCode(ctx, buf, c.(*ArrayConstraint))
	case *ObjectConstraint:
		if err := generateObjectCode(ctx, buf, c.(*ObjectConstraint)); err != nil {
			return err
		}
	}

	s := buf.String()
	s = strings.TrimSuffix(s, ".\n")
	fmt.Fprintf(out, s)

	return nil
}

func generateReferenceCode(ctx *genctx, out io.Writer, c *ReferenceConstraint) error {
	fmt.Fprintf(out, "%s%s.Reference(%s).RefersTo(`%s`)", ctx.Prefix(), ctx.pkgname, ctx.vname, c.reference)

	return nil
}

func generateAnyCode(ctx *genctx, out io.Writer, c *AnyConstraint) error {
	p := ctx.Prefix()
	fmt.Fprintf(out, "%s%s.Any()", p, ctx.pkgname)
	for _, c1 := range c.constraints {
		g1 := ctx.Indent()
		fmt.Fprintf(out, ".\n%sAdd(\n", ctx.Prefix())
		g2 := ctx.Indent()
		if err := generateCode(ctx, out, c1); err != nil {
			g2()
			g1()
			return err
		}
		g2()
		fmt.Fprintf(out, ",\n%s)", ctx.Prefix())
		g1()
	}
	return nil
}

func generateAllCode(ctx *genctx, out io.Writer, c *AllConstraint) error {
	if len(c.constraints) == 0 {
		return generateNilCode(ctx, out, NilConstraint)
	}

	p := ctx.Prefix()
	fmt.Fprintf(out, "%s%s.All()", p, ctx.pkgname)
	for _, c1 := range c.constraints {
		g1 := ctx.Indent()
		fmt.Fprintf(out, ".\n%sAdd(\n", ctx.Prefix())
		g2 := ctx.Indent()
		if err := generateCode(ctx, out, c1); err != nil {
			g2()
			g1()
			return err
		}
		g2()
		fmt.Fprintf(out, ",\n%s)", ctx.Prefix())
		g1()
	}
	return nil
}

func generateIntegerCode(ctx *genctx, out io.Writer, c *IntegerConstraint) error {
	fmt.Fprintf(out, "%s%s.Integer()", ctx.Prefix(), ctx.pkgname)

	if c.applyMinimum {
		fmt.Fprintf(out, ".Minimum(%d)", int(c.minimum))
	}

	if c.applyMaximum {
		fmt.Fprintf(out, ".Maximum(%d)", int(c.maximum))
	}

	return nil
}

func generateNumberCode(ctx *genctx, out io.Writer, c *NumberConstraint) error {
	fmt.Fprintf(out, "%s%s.Number()", ctx.Prefix(), ctx.pkgname)

	if c.applyMinimum {
		fmt.Fprintf(out, ".Minimum(%f)", c.minimum)
	}

	if c.exclusiveMinimum {
		fmt.Fprintf(out, ".ExclusiveMinimum(true)")
	}

	if c.applyMaximum {
		fmt.Fprintf(out, ".Maximum(%f)", c.maximum)
	}

	if c.exclusiveMaximum {
		fmt.Fprintf(out, ".ExclusiveMaximum(true)")
	}

	if c.HasDefault() {
		fmt.Fprintf(out, ".Default(%f)", c.DefaultValue())
	}

	return nil
}

func generateEnumCode(ctx *genctx, out io.Writer, c *EnumConstraint) error {
	fmt.Fprintf(out, "[]interface{}{")
	l := len(c.enums)
	for i, v := range c.enums {
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.String:
			fmt.Fprintf(out, "%s", strconv.Quote(rv.String()))
		}
		if i < l-1 {
			fmt.Fprintf(out, ", ")
		}
	}
	fmt.Fprintf(out, "}")

	return nil
}

func generateStringCode(ctx *genctx, out io.Writer, c *StringConstraint) error {
	fmt.Fprintf(out, "%s%s.String()", ctx.Prefix(), ctx.pkgname)

	if c.maxLength > -1 {
		fmt.Fprintf(out, ".MaxLength(%d)", c.maxLength)
	}

	if c.minLength > 0 {
		fmt.Fprintf(out, ".MinLength(%d)", c.minLength)
	}

	if f := c.format; f != "" {
		fmt.Fprintf(out, ".Format(%s)", strconv.Quote(string(f)))
	}

	if rx := c.regexp; rx != nil {
		fmt.Fprintf(out, ".Regexp(`%s`)", rx.String())
	}

	if enum := c.enums; enum != nil {
		fmt.Fprintf(out, ".Enum(")
		if err := generateEnumCode(ctx, out, enum); err != nil {
			return err
		}
		fmt.Fprintf(out, ",)")
	}

	return nil
}

func generateObjectCode(ctx *genctx, out io.Writer, c *ObjectConstraint) error {
	fmt.Fprintf(out, "%s%s.Object()", ctx.Prefix(), ctx.pkgname)

	// object code usually becomes quite nested, so we indent one level
	// to begin with
	g1 := ctx.Indent()
	defer g1()
	p := ctx.Prefix()

	if c.HasDefault() {
		fmt.Fprintf(out, ".\n%sDefault(%s)", p, c.DefaultValue())
	}

	if len(c.required) > 0 {
		fmt.Fprintf(out, ".\n%sRequired([]string{", p)
		l := len(c.required)
		pnames := make([]string, 0, l)
		for pname := range c.required {
			pnames = append(pnames, pname)
		}
		sort.Strings(pnames)
		for i, pname := range pnames {
			fmt.Fprint(out, strconv.Quote(pname))
			if i < l - 1 {
				fmt.Fprint(out, ", ")
			}
		}
		fmt.Fprint(out, "})")
	}


	if aprop := c.additionalProperties; aprop != nil {
		fmt.Fprintf(out, ".\n%sAdditionalProperties(\n", p)
		g := ctx.Indent()
		if err := generateCode(ctx, out, aprop); err != nil {
			g()
			return err
		}
		fmt.Fprintf(out, ",\n%s)", p)
		g()
	}

	pnames := make([]string, 0, len(c.properties))
	for pname := range c.properties {
		pnames = append(pnames, pname)
	}
	sort.Strings(pnames)

	for _, pname := range pnames {
		pdef := c.properties[pname]

		g := ctx.Indent()
		fmt.Fprintf(out, ".\n%sAddProp(\n%s\t`%s`,\n", p, p, pname)
		if err := generateCode(ctx, out, pdef); err != nil {
			g()
			return err
		}
		fmt.Fprintf(out, ",\n%s)", p)
		g()
	}

	if m := c.propdeps; len(m) > 0 {
		for from, deplist := range m {
			for _, to := range deplist {
				fmt.Fprintf(out, ".\n%sPropDependency(%s, %s)", ctx.Prefix(), strconv.Quote(from), strconv.Quote(to))
			}
		}
	}

	return nil
}

func generateArrayCode(ctx *genctx, out io.Writer, c *ArrayConstraint) error {
	fmt.Fprintf(out, "%s%s.Array()", ctx.Prefix(), ctx.pkgname)
	if c.minItems > -1 {
		fmt.Fprintf(out, ".MinItems(%d)", c.minItems)
	}
	if c.maxItems > -1 {
		fmt.Fprintf(out, ".MaxItems(%d)", c.maxItems)
	}
	if c.uniqueItems {
		fmt.Fprintf(out, ".UniqueItems()")
	}
	return nil
}

func generateBooleanCode(ctx *genctx, out io.Writer, c *BooleanConstraint) error {
	fmt.Fprintf(out, "%s%s.Boolean()", ctx.Prefix(), ctx.pkgname)
	if c.HasDefault() {
		fmt.Fprintf(out, ".Default(%t)", c.DefaultValue())
	}
	return nil
}
