package php

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/parser"
)

const GenType = "php"

type Generator struct {
}

func NewGenerator() gen.IGen {
	return &Generator{}
}

func (g *Generator) Gen(v *parser.Struct) []byte {
	out := &bytes.Buffer{}
	mm := make(map[string]string, len(v.Field))

	out.WriteString(fmt.Sprintf("\"%s\"=>[\n", v.Type))

	for _, field := range v.Field {
		jsonTag := field.Tag.Get("json")
		jsonName := strings.Split(jsonTag, ",")[0]
		tp := field.Type
		switch tp {
		case "*float64":
			tp = "double"
		case "[]string":
		case "map[string]interface {}":
			tp = "map"
		case "[]int":
			fallthrough
		case "[]uint64":
			tp = "[]int"
		case "*uint64":
			fallthrough
		case "*int":
			tp = "int"
		case "*string":
			tp = "string"
		}
		mm[jsonName] = tp
		out.WriteString(fmt.Sprintf("'%s'=>'%s',\n", jsonName, tp))

	}
	out.WriteString("],")
	return out.Bytes()
}
