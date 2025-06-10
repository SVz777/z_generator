package strct

import (
	"bytes"
	"fmt"
	"go/format"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen/utils"
)

const GenTypeWrapper = "struct_wrapper"

func init() {
	gen.Register(GenTypeWrapper, NewWrapperGenerator())
}

type WrapperGenerator struct {
}

func NewWrapperGenerator() gen.IGen {
	return &WrapperGenerator{}
}

func (g *WrapperGenerator) Gen(v *parser.Struct) []byte {
	out := &bytes.Buffer{}

	for _, field := range v.Field {
		fieldName := field.Name
		var fieldTp string
		if field.Type[0] == '*' {
			fieldTp = field.Type[1:]
		} else {
			fieldTp = field.Type
		}
		out.WriteString(fmt.Sprintf("func(%s *%s) Empty%s() bool {\n", v.Receiver, v.Type, fieldName))
		out.WriteString(fmt.Sprintf("return %s.%s == nil\n", v.Receiver, fieldName))
		out.WriteString("}\n")

		out.WriteString(fmt.Sprintf("func(%s *%s) Get%s() (ret %s) {\n", v.Receiver, v.Type, fieldName, fieldTp))
		out.WriteString(fmt.Sprintf("if %s.%s == nil {\n ", v.Receiver, fieldName))
		out.WriteString("return\n")
		out.WriteString("}\n")

		if utils.IsValueType(fieldTp) {
			out.WriteString(fmt.Sprintf("return *%s.%s\n", v.Receiver, fieldName))
		} else {
			out.WriteString(fmt.Sprintf("return %s.%s\n", v.Receiver, fieldName))
		}
		out.WriteString("}\n")

		out.WriteString(fmt.Sprintf("func(%s *%s) Set%s(v %s){\n", v.Receiver, v.Type, fieldName, fieldTp))
		if utils.IsValueType(fieldTp) {
			out.WriteString(fmt.Sprintf("%s.%s = &v\n", v.Receiver, fieldName))
		} else {
			out.WriteString(fmt.Sprintf("%s.%s = v\n", v.Receiver, fieldName))
		}
		out.WriteString("}\n")
	}
	ret, _ := format.Source(out.Bytes())
	return ret
}
