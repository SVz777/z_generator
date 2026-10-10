package utils

import (
	"reflect"
	"testing"

	"github.com/SVz777/z_generator/gen/parser"
)

func TestFieldCanBeNil(t *testing.T) {
	cases := []struct {
		field   *parser.Field
		want    bool
		wantErr bool
	}{
		{&parser.Field{Type: "Reader", Kind: reflect.Interface}, true, false},
		{&parser.Field{Type: "Value", Kind: reflect.Struct}, false, false},
		{&parser.Field{Type: "interfaceName", Kind: reflect.Struct}, false, false},
		{&parser.Field{Type: "unsafe.Pointer", Kind: reflect.UnsafePointer}, true, false},
		{&parser.Field{Type: "[]int"}, true, false},
		{&parser.Field{Type: "int"}, false, false},
		{&parser.Field{Type: "Reader"}, false, true},
		{nil, false, true},
	}
	for _, tc := range cases {
		got, err := FieldCanBeNil(tc.field)
		if got != tc.want || (err != nil) != tc.wantErr {
			t.Errorf("%+v: got %v, %v", tc.field, got, err)
		}
	}
}

// TestCanBeNil 逐个类型断言"能否与 nil 比较"。
// 判错会让值类型字段生成 `x.F == nil` —— 编译不过的代码。
func TestCanBeNil(t *testing.T) {
	cases := map[string]bool{
		"int": false, "string": false, "bool": false, "float64": false,
		"[3]int": false, "time.Time": false, "[][]byte": true,
		"*int": true, "*[]int": true, "**string": true,
		"[]int": true, "map[string]int": true, "chan int": true, "<-chan int": true,
		"func() error": true, "interface{}": true, "interface{Foo() int}": true,
		"any": true, "error": true,
		"channel": false, "function": false, "interfaceName": false, "Reader": false,
	}
	for typ, want := range cases {
		if got := CanBeNil(typ); got != want {
			t.Errorf("CanBeNil(%q) = %v, want %v", typ, got, want)
		}
	}
}

// TestZeroValue 断言各类型的零值字面量，用于不能与 nil 比较的字段判空。
func TestZeroValue(t *testing.T) {
	cases := map[string]string{
		"int":        "0",
		"int64":      "0",
		"uint8":      "0",
		"byte":       "0",
		"float64":    "0",
		"complex128": "0",
		"string":     `""`,
		"bool":       "false",
		"[3]int":     "[3]int{}",
		"time.Time":  "time.Time{}",
	}
	for typ, want := range cases {
		if got := ZeroValue(typ); got != want {
			t.Errorf("ZeroValue(%q) = %q, want %q", typ, got, want)
		}
	}
}

// TestZeroValueUnknown 不认识的类型回退成复合字面量，而不是空串。
func TestZeroValueUnknown(t *testing.T) {
	if got := ZeroValue("pkg.Config"); got != "pkg.Config{}" {
		t.Errorf("ZeroValue(pkg.Config) = %q, want pkg.Config{}", got)
	}
}

func TestTrimPtr(t *testing.T) {
	cases := map[string]string{
		"*int":   "int",
		"**int":  "*int",
		"int":    "int",
		"[]*int": "[]*int",
		"":       "",
	}
	for in, want := range cases {
		if got := TrimPtr(in); got != want {
			t.Errorf("TrimPtr(%q) = %q, want %q", in, got, want)
		}
	}
}
