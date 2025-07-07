package options

import (
	"bytes"
	_ "embed"
	"go/format"
	"text/template"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/parser"
)

const GenType = "options"

func init() {
	gen.Register(GenType, NewGenerator())
}

//go:embed options.tmpl
var tpl string

type Generator struct {
}

func NewGenerator() gen.IGen {
	return &Generator{}
}

func (g *Generator) Gen(v *parser.Struct) []byte {
	p, err := g.getParams(v)
	out := bytes.NewBuffer(nil)
	ins, err := template.New("validate").Parse(tpl)
	if err != nil {
		return []byte(err.Error())
	}
	if err := ins.Execute(out, p); err != nil {
		return []byte(err.Error())
	}
	ret, err := format.Source(out.Bytes())
	if err != nil {
		return []byte(err.Error())
	}
	return ret
}

type params struct {
	Type   string
	Fields []*parser.Field
}

func (g *Generator) getParams(v *parser.Struct) (*params, error) {
	return &params{
		Type:   v.Type,
		Fields: v.Field,
	}, nil
}
