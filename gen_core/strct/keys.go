/**
 * @file    keys.go
 * @author
 *  ___  _  _  ____
 * / __)( \/ )(_   )
 * \__ \ \  /  / /_
 * (___/  \/  (____)
 * (903943711@qq.com)
 * @date    2019-12-09
 * @desc
 */
package strct

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/SVz777/z_generator/gen_core/parser"
)

type KeyData struct {
	Json string
}

type Keys struct {
	v *parser.Struct

	AllKeys     []KeyData
	IntKeys     []KeyData
	Uint64Keys  []KeyData
	Float64Keys []KeyData
}

func NewKeys(v *parser.Struct) *Keys {
	return &Keys{
		v:           v,
		AllKeys:     make([]KeyData, 0),
		IntKeys:     make([]KeyData, 0),
		Uint64Keys:  make([]KeyData, 0),
		Float64Keys: make([]KeyData, 0),
	}
}

func (keys *Keys) Gen() (out *bytes.Buffer) {
	out = &bytes.Buffer{}

	//s := strings.Builder{}
	keysTemp := `
var Keys = set.Set{
{{range .AllKeys}}"{{.Json}}":collections.Empty{},
{{end}}
}
var NeedTransformIntKey = set.Set{
{{range .IntKeys}}"{{.Json}}":collections.Empty{},
{{end}}
}

var NeedTransformUint64Key = set.Set{
{{range .Uint64Keys}}"{{.Json}}":collections.Empty{},
{{end}}
}

var NeedTransformFloat64Key = set.Set{
{{range .Float64Keys}}"{{.Json}}":collections.Empty{},
{{end}}
}

`
	tmpl, _ := template.New("Tmpl").Parse(keysTemp)

	for _, field := range keys.v.Field {
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
			keys.IntKeys = append(keys.IntKeys, KeyData{name})
		case "uint64":
			keys.Uint64Keys = append(keys.Uint64Keys, KeyData{name})
		case "float64":
			keys.Float64Keys = append(keys.Float64Keys, KeyData{name})
		}
		keys.AllKeys = append(keys.AllKeys, KeyData{name})
	}

	if err := tmpl.Execute(out, keys); err != nil {
		panic(err)
	}

	return
}
