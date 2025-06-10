package validate

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/config"
	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen_impl/validate/handler"
)

const GenType = "validate"

func init() {
	gen.Register(GenType, NewGenerator())
}

//go:embed validate.tmpl
var tpl string

type Generator struct {
}

func NewGenerator() gen.IGen {
	return &Generator{}
}

func (g *Generator) Gen(v *parser.Struct) []byte {
	p, err := g.getParams(v)
	if err != nil {
		return []byte(err.Error())
	}
	out := bytes.NewBuffer(nil)
	ins, err := template.New("validate").Funcs(template.FuncMap{
		"split": strings.Split,
	}).Parse(tpl)
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
	SingleFile bool
	Package    string

	ValidateTag    string
	ValidateTagSep string
	Receiver       string
	Type           string
	Fields         []*parser.Field
}

func (p *params) GetFunc(field *parser.Field, tag string, skipError bool) (ret [2]string, err error) {
	defer func() {
		if err != nil && skipError {
			ret[0] = "false"
			ret[1] = fmt.Sprintf(`fmt.Errorf("%s")`, err.Error())
			err = nil
		}
	}()

	fieldName := getFieldName(field)
	kvs := strings.SplitN(tag, "=", 2)
	if len(kvs) < 1 {
		err = fmt.Errorf("%s'tag %s must be in format", fieldName, tag)
		return
	}
	funcName := kvs[0]
	if funcName == "" {
		err = fmt.Errorf("%s'tag %s must be in format", fieldName, tag)
		return
	}

	h := handler.GetHandler(funcName)
	if h == nil {
		err = fmt.Errorf("%s'tag %s not found handler", fieldName, tag)
		return
	}

	if !h.Check(field.Type) {
		err = fmt.Errorf("// %s: %s check error", field.Name, field.Type)
		return
	}
	if len(kvs) > 1 {
		// 有参
		ret[0], ret[1] = h.Handle(fieldName, kvs[1:]...)
		return
	}
	// 无参数
	ret[0], ret[1] = h.Handle(fieldName)
	return
}

var baseType = map[string]bool{
	"uint":    true,
	"uint8":   true,
	"uint16":  true,
	"uint32":  true,
	"uint64":  true,
	"int":     true,
	"int8":    true,
	"int16":   true,
	"int32":   true,
	"int64":   true,
	"float32": true,
	"float64": true,
	"string":  true,
}

func getFieldName(field *parser.Field) string {
	name := field.Receiver + "." + field.Name
	if field.Type[0] == '*' && baseType[field.Type[1:]] {
		return "*" + name
	}
	return name
}

func (g *Generator) getParams(v *parser.Struct) (*params, error) {
	return &params{
		SingleFile:     !config.Config.Inplace,
		Package:        v.Package,
		ValidateTag:    "validate",
		ValidateTagSep: "|",
		Receiver:       v.Receiver,
		Type:           v.Type,
		Fields:         v.Field,
	}, nil
}
