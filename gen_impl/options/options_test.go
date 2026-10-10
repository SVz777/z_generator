package options_test

import (
	"testing"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen_impl/options"
)

type Options struct {
	A *int    `json:"a" validate:"required"`
	B float64 `json:"b" validate:"exists"`
}

type Option func(opts *Options)

func mustGen(t *testing.T, g gen.IGen, st *parser.Struct) string {
	t.Helper()
	code, err := g.Gen(st)
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}
	return string(code)
}

func TestOptions(t *testing.T) {
	g := options.NewGenerator()
	st := parser.StructOf(Options{})
	t.Log("\n" + mustGen(t, g, st))

	st, err := parser.StructFromFile("./options_test.go", "Options")
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	t.Log("\n" + mustGen(t, g, st))
}
