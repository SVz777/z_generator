package options

import (
	"bytes"
	"fmt"
	"go/format"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/parser"
)

const GenType = "options"

type Generator struct {
}

func NewGenerator() gen.IGen {
	return &Generator{}
}

func (g *Generator) Gen(v *parser.Struct) []byte {
	out := &bytes.Buffer{}

	out.WriteString(fmt.Sprintf("func (opts *%s) Update(opt ...Option) {\n", v.Type))
	out.WriteString("for _, o := range opt {\n")
	out.WriteString("o(opts)")
	out.WriteString("}\n")
	out.WriteString("}\n")
	for _, field := range v.Field {
		out.WriteString(fmt.Sprintf("func With%s(v %s) Option {\n", field.Name, field.Type))
		out.WriteString(fmt.Sprintf("return func(opts *%s) {\n", v.Type))
		out.WriteString(fmt.Sprintf("opts.%s = v\n", field.Name))
		out.WriteString("}\n")
		out.WriteString("}\n")
	}
	ret, _ := format.Source(out.Bytes())
	return ret
}
