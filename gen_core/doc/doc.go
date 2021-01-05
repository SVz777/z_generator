package doc

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/SVz777/z_generator/gen_core/parser"
)

type Doc struct {
	v *parser.Struct
}

func NewDoc(v *parser.Struct) *Doc {
	return &Doc{
		v: v,
	}
}

func (doc *Doc) Gen() (out *bytes.Buffer) {
	out = &bytes.Buffer{}
	jsonTag := "json"
	vTag := "validate"
	fmt.Fprintf(out, "| 参数名| 类型  |  是否必须    |备注|\n")
	fmt.Fprintf(out, "| ---- | ---- | ---- | ---- |\n")
	for _, field := range doc.v.Field {
		//f :=  vof.Field(i)
		if field.Type[0] == '*' {
			field.Type = field.Type[1:]
		}
		fmt.Fprintf(out, "|%s|%s|%s| |\n", strings.Split(field.Tag.Get(jsonTag), ",")[0], field.Type, field.Tag.Get(vTag))
	}
	return
}
