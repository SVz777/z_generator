package options

import (
	_ "embed"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen/render"
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

func (g *Generator) Gen(v *parser.Struct) ([]byte, error) {
	p, err := g.getParams(v)
	if err != nil {
		return nil, err
	}
	return render.Render(GenType, tpl, p, render.WithFmt())
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
