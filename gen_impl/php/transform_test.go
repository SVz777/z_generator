package php_test

import (
	"testing"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen_impl/php"
)

type SS struct {
	A *int    `json:"a" validate:"required"`
	B float64 `json:"b" validate:"exists"`
}

func mustGen(t *testing.T, g gen.IGen, st *parser.Struct) string {
	t.Helper()
	code, err := g.Gen(st)
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}
	return string(code)
}

func TestTransform(t *testing.T) {
	g := php.NewGenerator()
	st := parser.StructOf(SS{})
	t.Log("\n" + mustGen(t, g, st))

	st, err := parser.StructFromFile("./transform_test.go", "SS")
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	t.Log("\n" + mustGen(t, g, st))
}
