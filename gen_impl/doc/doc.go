package doc

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/parser"
)

const GenType = "doc"

type Generator struct {
}

func NewGenerator() gen.IGen {
	return &Generator{}
}

func (g *Generator) Gen(v *parser.Struct) []byte {
	out := &bytes.Buffer{}
	jsonTag := "json"
	vTag := "validate"
	out.WriteString("| 参数名| 类型  |  是否必须    |备注|\n")
	out.WriteString("| ---- | ---- | ---- | ---- |\n")
	for _, field := range v.Field {
		fmt.Println(*field)
		if field.Type[0] == '*' {
			field.Type = field.Type[1:]
		}
		out.WriteString(fmt.Sprintf("|%s|%s|%s| |\n", strings.Split(field.Tag.Get(jsonTag), ",")[0], field.Type, field.Tag.Get(vTag)))
	}
	return out.Bytes()
}
