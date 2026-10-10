package strct

import (
	_ "embed"
	"strings"
	"text/template"

	"github.com/SVz777/tk/convert"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/args"
	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen/render"
)

const GenTypeWrapper = "struct_wrapper"

func init() {
	gen.Register(GenTypeWrapper, NewWrapperGenerator())
	args.Declare(GenTypeWrapper, "empty", "是否生成空值校验")
	args.Declare(GenTypeWrapper, "get", "是否生成 get 方法")
	args.Declare(GenTypeWrapper, "get_value", "是否生成 getValue 方法")
	args.Declare(GenTypeWrapper, "set", "是否生成 set 方法")
}

//go:embed wrapper.tmpl
var tpl string

var funcs = template.FuncMap{
	"ToCamelCase": func(s string) string {
		r := convert.ToCamelCase(s)
		return strings.ToTitle(r)
	},
}

type WrapperGenerator struct {
}

func NewWrapperGenerator() gen.IGen {
	return &WrapperGenerator{}
}

func (g *WrapperGenerator) Gen(v *parser.Struct) ([]byte, error) {
	p := g.getParams(v)
	return render.Render(GenTypeWrapper, tpl, p, render.WithFuncs(funcs), render.WithFmt())
}

type params struct {
	Receiver string
	Type     string
	Fields   []*parser.Field
	Gen      genFuncs
}

type genFuncs struct {
	Empty    bool
	Get      bool
	GetValue bool
	Set      bool
}

func (g *WrapperGenerator) getParams(v *parser.Struct) *params {
	return &params{
		Receiver: v.Receiver,
		Type:     v.Type,
		Fields:   v.Field,
		// 四项默认全生成，可用 -empty=false -set=false 这类写法按需关闭
		Gen: genFuncs{
			Empty:    args.Args.GetBoolOr("empty", true),
			Get:      args.Args.GetBoolOr("get", true),
			GetValue: args.Args.GetBoolOr("get_value", true),
			Set:      args.Args.GetBoolOr("set", true),
		},
	}
}
