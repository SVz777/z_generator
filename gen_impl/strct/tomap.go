package strct

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen/render"
	"github.com/SVz777/z_generator/gen/utils"
)

const GenTypeToMap = "struct2map"

func init() {
	gen.Register(GenTypeToMap, NewToStringMapGenerator())
}

//go:embed tomap.tmpl
var toMapTpl string

type ToStringMapGenerator struct {
}

func NewToStringMapGenerator() gen.IGen {
	return &ToStringMapGenerator{}
}

func (g *ToStringMapGenerator) Gen(v *parser.Struct) ([]byte, error) {
	p, err := g.getParams(v)
	if err != nil {
		return nil, err
	}
	return render.Render(GenTypeToMap, toMapTpl, p, render.WithFmt())
}

type mapField struct {
	Name string
	Key  string
	Conv string
}

type toMapParams struct {
	Receiver string
	Type     string
	Fields   []mapField
}

func (g *ToStringMapGenerator) getParams(v *parser.Struct) (*toMapParams, error) {
	p := &toMapParams{
		Receiver: v.Receiver,
		Type:     v.Type,
	}
	for _, field := range v.Field {
		name := field.Name
		key := strings.Split(field.Tag.Get("json"), ",")[0]
		if len(key) == 0 {
			key = name
		}
		getter := fmt.Sprintf("%s.Get%s()", v.Receiver, name)
		if utils.IsPtr(field.Type) {
			getter = fmt.Sprintf("%s.GetValue%s()", v.Receiver, name)
		}
		var conv string
		switch utils.TrimPtr(field.Type) {
		case "int":
			conv = fmt.Sprintf("strconv.Itoa(%s)", getter)
		case "uint64":
			conv = fmt.Sprintf("strconv.FormatUint(%s,10)", getter)
		case "float64":
			conv = fmt.Sprintf("strconv.FormatFloat(%s,'f',-1,64)", getter)
		case "string":
			conv = getter
		default:
			return nil, fmt.Errorf("结构体 %s 的字段 %s 类型 %s 不支持转换为 map[string]string", v.Type, name, field.Type)
		}
		p.Fields = append(p.Fields, mapField{Name: name, Key: key, Conv: conv})
	}
	return p, nil
}
