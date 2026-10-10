package utils

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/SVz777/z_generator/gen/parser"
)

var (
	// Funcs 模板函数表
	Funcs = template.FuncMap{
		"canBeNil": FieldCanBeNil,
		"zeroOf":   ZeroValue,
		"isPtr":    IsPtr,
		"trimPtr":  TrimPtr,
		"split":    strings.Split,
	}
)

// numberTypes 是所有数值类型，零值统一是 0。
var numberTypes = map[string]bool{
	"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
	"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
	"uintptr": true, "byte": true, "rune": true,
	"float32": true, "float64": true, "complex64": true, "complex128": true,
}

// CanBeNil 仅识别内置类型和类型字面量，不解析自定义命名类型。
func CanBeNil(typeName string) bool {
	return kindCanBeNil(parser.KindOfTypeName(typeName))
}

func FieldCanBeNil(field *parser.Field) (bool, error) {
	if field == nil {
		return false, fmt.Errorf("缺少字段类型信息")
	}
	kind := field.Kind
	if kind == reflect.Invalid {
		kind = parser.KindOfTypeName(field.Type)
	}
	if kind == reflect.Invalid {
		return false, fmt.Errorf("无法确认字段 %s 的类型 %s 是否可与 nil 比较", field.Name, field.Type)
	}
	return kindCanBeNil(kind), nil
}

func kindCanBeNil(kind reflect.Kind) bool {
	switch kind {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan,
		reflect.Func, reflect.Interface, reflect.UnsafePointer:
		return true
	default:
		return false
	}
}

// ZeroValue 给出类型的零值字面量，用于 CanBeNil 为 false 时判空。
// 认识的类型给出 0 / "" / false，其余（数组、结构体）用复合字面量 T{}。
func ZeroValue(typeName string) string {
	switch typeName {
	case "string":
		return `""`
	case "bool":
		return "false"
	}
	if numberTypes[typeName] {
		return "0"
	}
	return typeName + "{}"
}

// IsPtr 判断类型名是否为指针类型。
func IsPtr(typeName string) bool {
	return strings.HasPrefix(typeName, "*")
}

// TrimPtr 去掉指针前置*
func TrimPtr(s string) string {
	if len(s) > 0 && s[0] == '*' {
		return s[1:]
	}
	return s
}
