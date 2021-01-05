/**
 * @file    wapper.go
 * @author
 *  ___  _  _  ____
 * / __)( \/ )(_   )
 * \__ \ \  /  / /_
 * (___/  \/  (____)
 * (903943711@qq.com)
 * @date    2019-10-31
 * @desc
 */
package strct

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/SVz777/z_generator/gen_core/parser"
	"github.com/SVz777/z_generator/gen_core/utils"
)

type StructWrapper struct {
	v *parser.Struct
}

func NewStructWrapper(v *parser.Struct) *StructWrapper {
	return &StructWrapper{v: v}
}

func (sw *StructWrapper) Gen() (out *bytes.Buffer) {
	out = &bytes.Buffer{}

	for _, field := range sw.v.Field {
		fieldName := field.Name
		var fieldTp string
		if field.Type[0] == '*' {
			fieldTp = field.Type[1:]
		} else {
			fieldTp = field.Type
		}

		fmt.Fprintf(out, "func(%s *%s) Empty%s() bool {\n", sw.v.Receiver, sw.v.Type, fieldName)
		fmt.Fprintf(out, "return %s.%s == nil\n", sw.v.Receiver, fieldName)
		fmt.Fprintln(out, "}")

		fmt.Fprintf(out, "func(%s *%s) Get%s() (ret %s) {\n", sw.v.Receiver, sw.v.Type, fieldName, fieldTp)
		fmt.Fprintf(out, "if %s.%s == nil {\n ", sw.v.Receiver, fieldName)
		fmt.Fprintln(out, "return")
		fmt.Fprintln(out, "}")

		if utils.IsValueType(fieldTp) {
			fmt.Fprintf(out, "return *%s.%s\n", sw.v.Receiver, fieldName)
		} else {
			fmt.Fprintf(out, "return %s.%s\n", sw.v.Receiver, fieldName)
		}
		fmt.Fprintln(out, "}")

		fmt.Fprintf(out, "func(%s *%s) Set%s(v %s){\n", sw.v.Receiver, sw.v.Type, fieldName, fieldTp)
		if utils.IsValueType(fieldTp) {
			fmt.Fprintf(out, "%s.%s = &v\n", sw.v.Receiver, fieldName)
		} else {
			fmt.Fprintf(out, "%s.%s = v\n", sw.v.Receiver, fieldName)
		}
		fmt.Fprintln(out, "}")
	}
	return
}

type StructToStringMap struct {
	v *parser.Struct
}

func NewStructToStringMap(v *parser.Struct) *StructToStringMap {
	return &StructToStringMap{v: v}
}

func (ssm *StructToStringMap) Gen() (out *bytes.Buffer) {
	out = &bytes.Buffer{}

	fmt.Fprintf(out, "func (%s *%s) Output() map[string]string {\n", ssm.v.Receiver, ssm.v.Type)
	fmt.Fprintf(out, "ret := make(map[string]string,%d)\n", len(ssm.v.Field))
	for _, field := range ssm.v.Field {
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
			fmt.Fprintf(out, "if !%s.Empty%s() {\n", ssm.v.Receiver, fieldName)
			fmt.Fprintf(out, "ret[\"%s\"] = strconv.Itoa(%s.Get%s())\n", jsonName, ssm.v.Receiver, fieldName)
			fmt.Fprintln(out, "}")
		case "uint64":
			fmt.Fprintf(out, "if !%s.Empty%s() {\n", ssm.v.Receiver, fieldName)
			fmt.Fprintf(out, "ret[\"%s\"] = strconv.FormatUint(%s.Get%s(),10)\n", jsonName, ssm.v.Receiver, fieldName)
			fmt.Fprintln(out, "}")
		case "float64":
			fmt.Fprintf(out, "if !%s.Empty%s() {\n", ssm.v.Receiver, fieldName)
			fmt.Fprintf(out, "ret[\"%s\"] = strconv.FormatFloat(%s.Get%s(),'f',-1,64)\n", jsonName, ssm.v.Receiver, fieldName)
			fmt.Fprintln(out, "}")
		case "string":
			fmt.Fprintf(out, "if !%s.Empty%s() {\n", ssm.v.Receiver, fieldName)
			fmt.Fprintf(out, "ret[\"%s\"] = %s.Get%s()\n", jsonName, ssm.v.Receiver, fieldName)
			fmt.Fprintln(out, "}")
		default:
			fmt.Println("unsupport tp:" + fieldTp)
		}

	}
	fmt.Fprintln(out, "return ret")
	fmt.Fprintln(out, "}")

	fmt.Fprintf(out, "func (%s %s) MarshalJSON() ([]byte, error) {\n", ssm.v.Receiver, ssm.v.Type)
	fmt.Fprintf(out, "return json.Marshal(%s.Output())\n", ssm.v.Receiver)
	fmt.Fprintln(out, "}")
	return
}
