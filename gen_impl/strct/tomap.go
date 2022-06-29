package strct

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/parser"
)

const GenTypeToMap = "struct_tomap"

type ToStringMapGenerator struct {
}

func NewToStringMapGenerator() gen.IGen {
	return &ToStringMapGenerator{}
}

func (g *ToStringMapGenerator) Gen(v *parser.Struct) []byte {
	out := &bytes.Buffer{}

	out.WriteString(fmt.Sprintf("func (%s *%s) Output() map[string]string {\n", v.Receiver, v.Type))
	out.WriteString(fmt.Sprintf("ret := make(map[string]string,%d)\n", len(v.Field)))
	for _, field := range v.Field {
		fieldName := field.Name
		json := field.Tag.Get("json")
		jsonName := strings.Split(json, ",")[0]
		if len(jsonName) == 0 {
			jsonName = fieldName
		}

		var fieldTp string
		if field.Type[0] == '*' {
			fieldTp = field.Type[1:]
		} else {
			fieldTp = field.Type
		}
		switch fieldTp {
		case "int":
			out.WriteString(fmt.Sprintf("if !%s.Empty%s() {\n", v.Receiver, fieldName))
			out.WriteString(fmt.Sprintf("ret[\"%s\"] = strconv.Itoa(%s.Get%s())\n", jsonName, v.Receiver, fieldName))
			out.WriteString("}\n")
		case "uint64":
			out.WriteString(fmt.Sprintf("if !%s.Empty%s() {\n", v.Receiver, fieldName))
			out.WriteString(fmt.Sprintf("ret[\"%s\"] = strconv.FormatUint(%s.Get%s(),10)\n", jsonName, v.Receiver, fieldName))
			out.WriteString("}\n")
		case "float64":
			out.WriteString(fmt.Sprintf("if !%s.Empty%s() {\n", v.Receiver, fieldName))
			out.WriteString(fmt.Sprintf("ret[\"%s\"] = strconv.FormatFloat(%s.Get%s(),'f',-1,64)\n", jsonName, v.Receiver, fieldName))
			out.WriteString("}\n")
		case "string":
			out.WriteString(fmt.Sprintf("if !%s.Empty%s() {\n", v.Receiver, fieldName))
			out.WriteString(fmt.Sprintf("ret[\"%s\"] = %s.Get%s()\n", jsonName, v.Receiver, fieldName))
			out.WriteString("}\n")
		default:
			fmt.Println("unsupport tp:" + fieldTp)
		}

	}

	out.WriteString("return ret")
	out.WriteString("}\n")

	out.WriteString(fmt.Sprintf("func (%s %s) MarshalJSON() ([]byte, error) {\n", v.Receiver, v.Type))
	out.WriteString(fmt.Sprintf("return json.Marshal(%s.Output())\n", v.Receiver))
	out.WriteString("}\n")
	ret, _ := format.Source(out.Bytes())
	return ret
}
