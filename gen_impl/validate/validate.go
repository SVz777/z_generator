package validate

import (
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/args"
	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen/render"
	"github.com/SVz777/z_generator/gen_impl/validate/handler"
)

const GenType = "validate"

func init() {
	gen.Register(GenType, NewGenerator())
	// 声明本生成器支持的扩展参数，只影响 -h 展示，不影响传参。
	// 带上 GenType，这样只有 -tp validate -h 时才展示它们。
	args.Declare(GenType, "validate_tag", "校验用的 struct tag 名，默认 validate")
	args.Declare(GenType, "validate_tag_sep", "tag 内多条规则的分割符，默认 |")
}

//go:embed validate.tmpl
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
	SingleFile bool
	Package    string

	ValidateTag    string
	ValidateTagSep string
	Receiver       string
	Type           string
	Fields         []*parser.Field
}

// GetFunc 把一条 tag 规则翻译成完整校验代码，供模板直接嵌入。
//
// 规则名写错、类型不支持、参数不合法一律返回 error，不做降级。
// 早先的做法是降级成 `if false { ... }`，结果 tag 写错时生成的校验代码永远不触发——
// 代码看着正常、能编译、运行时什么都不校验，比直接报错危险得多。
func (p *params) GetFunc(field *parser.Field, tag string) (ret string, err error) {
	funcName, value, hasValue := strings.Cut(tag, "=")
	if funcName == "" {
		return ret, fmt.Errorf("字段 %s 的 tag %q 不是 <规则>[=<参数>] 格式", field.Name, tag)
	}

	h := handler.GetHandler(funcName)
	if h == nil {
		return ret, fmt.Errorf("字段 %s 的 tag %q 没有对应的校验处理器", field.Name, tag)
	}
	if !h.Check(field.Type) {
		return ret, fmt.Errorf("字段 %s 的类型 %s 不支持校验规则 %s", field.Name, field.Type, funcName)
	}

	var funcParams []string
	if hasValue {
		funcParams = []string{value}
	}
	// 是否解引用由各 handler 自己决定（required 判断 nil 时不能解引用）
	handleErr := h.Handle(field, funcParams...)
	if errors.Is(handleErr, handler.ErrHandle) {
		return ret, fmt.Errorf("字段 %s 的 tag %q: %w", field.Name, tag, handleErr)
	}
	if handleErr == nil {
		return ret, fmt.Errorf("字段 %s 的 tag %q: handler 未返回生成结果", field.Name, tag)
	}
	return handleErr.Error(), nil
}

func (g *Generator) getParams(v *parser.Struct) (*params, error) {
	return &params{
		SingleFile: !args.Args.Inplace,
		Package:    v.Package,
		// 从多余参数里取，未传时回落到默认值
		ValidateTag:    args.Args.GetOr("validate_tag", "validate"),
		ValidateTagSep: args.Args.GetOr("validate_tag_sep", "|"),
		Receiver:       args.Args.GetOr("receiver", v.Receiver),
		Type:           v.Type,
		Fields:         v.Field,
	}, nil
}
