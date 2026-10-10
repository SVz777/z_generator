package doc

import (
	_ "embed"
	"strings"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen/render"
	"github.com/SVz777/z_generator/gen/utils"
)

const GenType = "doc"

func init() {
	gen.Register(GenType, NewGenerator())
}

//go:embed doc.tmpl
var docTpl string

type Generator struct {
}

func NewGenerator() gen.IGen {
	return &Generator{}
}

func (g *Generator) Gen(v *parser.Struct) ([]byte, error) {
	// 生成的不是 Go 代码，不能用 gofmt，只做渲染
	return render.Render(GenType, docTpl, g.getParams(v))
}

type docField struct {
	Name     string
	Type     string
	Validate string
}

type docParams struct {
	Fields []docField
}

func (g *Generator) getParams(v *parser.Struct) *docParams {
	const (
		jsonTag = "json"
		vTag    = "validate"
	)
	p := &docParams{}
	for _, field := range v.Field {
		p.Fields = append(p.Fields, docField{
			Name: strings.Split(field.Tag.Get(jsonTag), ",")[0],
			// 注意：这里必须去掉指针符号的副本，不能直接改 field.Type，
			// 否则会污染入参——同一个 *parser.Struct 可能被多个生成器共用。
			Type:     utils.TrimPtr(field.Type),
			Validate: field.Tag.Get(vTag),
		})
	}
	return p
}
