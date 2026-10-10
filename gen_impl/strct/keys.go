package strct

import (
	_ "embed"
	"strings"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen/render"
	"github.com/SVz777/z_generator/gen/utils"
)

const GenTypeKeys = "struct_keys"

func init() {
	gen.Register(GenTypeKeys, NewKeysGenerator())
}

//go:embed keys.tmpl
var keysTpl string

type KeyData struct {
	Json string
}

type KeysGenerator struct {
}

func NewKeysGenerator() gen.IGen {
	return &KeysGenerator{}
}

func (g *KeysGenerator) Gen(v *parser.Struct) ([]byte, error) {
	return render.Render(GenTypeKeys, keysTpl, g.getParams(v), render.WithFmt())
}

type keysParams struct {
	AllKeys     []KeyData
	IntKeys     []KeyData
	Uint64Keys  []KeyData
	Float64Keys []KeyData
}

func (g *KeysGenerator) getParams(v *parser.Struct) *keysParams {
	p := &keysParams{}
	for _, field := range v.Field {
		// 与改造前保持一致：json tag 为空时 key 就是空串，不回退到字段名
		key := KeyData{Json: strings.Split(field.Tag.Get("json"), ",")[0]}
		// 指针类型按值类型归类
		switch utils.TrimPtr(field.Type) {
		case "int":
			p.IntKeys = append(p.IntKeys, key)
		case "uint64":
			p.Uint64Keys = append(p.Uint64Keys, key)
		case "float64":
			p.Float64Keys = append(p.Float64Keys, key)
		}
		p.AllKeys = append(p.AllKeys, key)
	}
	return p
}
