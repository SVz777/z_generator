/**
 * @file    transform.go
 * @author
 *  ___  _  _  ____
 * / __)( \/ )(_   )
 * \__ \ \  /  / /_
 * (___/  \/  (____)
 * (903943711@qq.com)
 * @date    2019-12-11
 * @desc
 */
package php

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/SVz777/z_generator/gen_core/parser"
)

type Type struct {
	v *parser.Struct
}

func NewType(v *parser.Struct) *Type {
	return &Type{v: v}
}

func (tp *Type) Gen() (out *bytes.Buffer) {
	out = &bytes.Buffer{}

	mm := make(map[string]string, len(tp.v.Field))
	fmt.Fprintf(out, "\"%s\"=>[\n", tp.v.Type)
	for _, field := range tp.v.Field {
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
		fmt.Fprintf(out, "'%s'=>'%s',\n", jsonName, tp)

	}
	fmt.Fprintln(out, `],`)
	return
}
