package doc_test

import (
	"testing"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen_impl/doc"
)

type SS struct {
	A *int    `json:"a" validate:"required"`
	B float64 `json:"b" validate:"exists"`
}

// mustGen 生成失败时直接终止测试，避免把错误信息当代码继续断言。
func mustGen(t *testing.T, g gen.IGen, st *parser.Struct) string {
	t.Helper()
	code, err := g.Gen(st)
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}
	return string(code)
}

func TestDoc(t *testing.T) {
	g := doc.NewGenerator()
	st := parser.StructOf(SS{})
	t.Log("\n" + mustGen(t, g, st))

	st, err := parser.StructFromFile("./doc_test.go", "SS")
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	t.Log("\n" + mustGen(t, g, st))
}

// TestDocNotModifyInput 文档生成器曾经直接改写 field.Type 去掉指针符号，
// 污染共享的 *parser.Struct。这里确保它不再修改入参。
func TestDocNotModifyInput(t *testing.T) {
	v := parser.StructOf(SS{})
	if _, err := doc.NewGenerator().Gen(v); err != nil {
		t.Fatalf("生成失败: %v", err)
	}
	if v.Field[0].Type != "*int" {
		t.Errorf("入参被修改: field[0].Type = %q, want %q", v.Field[0].Type, "*int")
	}
}
