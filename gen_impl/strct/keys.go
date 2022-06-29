package strct

import (
	"bytes"
	"go/format"
	"strings"
	"text/template"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/parser"
)

const GenTypeKeys = "struct_keys"

type KeyData struct {
	Json string
}

type KeysGenerator struct {
	v *parser.Struct
}

func NewKeysGenerator() gen.IGen {
	return &KeysGenerator{}
}

func (g *KeysGenerator) Gen(v *parser.Struct) []byte {
	var (
		allKeys     []KeyData
		intKeys     []KeyData
		uint64Keys  []KeyData
		float64Keys []KeyData
	)
	out := &bytes.Buffer{}

	keysTemp := `
var KeysGenerator = set.Set{
{{range .allKeys}}"{{.Json}}":collections.Empty{},
{{end}}
}
var NeedTransformIntKey = set.Set{
{{range .intKeys}}"{{.Json}}":collections.Empty{},
{{end}}
}

var NeedTransformUint64Key = set.Set{
{{range .uint64Keys}}"{{.Json}}":collections.Empty{},
{{end}}
}

var NeedTransformFloat64Key = set.Set{
{{range .float64Keys}}"{{.Json}}":collections.Empty{},
{{end}}
}

`
	tmpl, _ := template.New("Tmpl").Parse(keysTemp)

	for _, field := range v.Field {
		json := field.Tag.Get("json")
		name := strings.Split(json, ",")[0]

		var tp string
		if field.Type[0] == '*' {
			tp = field.Type[1:]
		} else {
			tp = field.Type
		}

		switch tp {
		case "int":
			intKeys = append(intKeys, KeyData{name})
		case "uint64":
			uint64Keys = append(uint64Keys, KeyData{name})
		case "float64":
			float64Keys = append(float64Keys, KeyData{name})
		}
		allKeys = append(allKeys, KeyData{name})
	}

	if err := tmpl.Execute(out, map[string][]KeyData{
		"allKeys":     allKeys,
		"intKeys":     intKeys,
		"uint64Keys":  uint64Keys,
		"float64Keys": float64Keys,
	}); err != nil {
		panic(err)
	}

	ret, _ := format.Source(out.Bytes())
	return ret
}
