package options_test

import (
	"fmt"
	"testing"

	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen_impl/options"
)

type Options struct {
	A *int    `json:"a" validate:"required"`
	B float64 `json:"b" validate:"exists"`
}

type Option func(opts *Options)

func TestOptions(t *testing.T) {
	var p parser.IParser
	p = parser.NewRuntime(Options{})
	fmt.Println(string(options.NewGenerator().Gen(p.Parser())))

	filename := "./options_test.go"
	stname := "Options"
	p = parser.NewAST(filename, stname)
	fmt.Println(string(options.NewGenerator().Gen(p.Parser())))
}
