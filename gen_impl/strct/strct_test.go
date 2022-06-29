package strct_test

import (
	"fmt"
	"testing"

	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen_impl/strct"
)

type SS struct {
	A *int    `json:"a" validate:"required"`
	B float64 `json:"b" validate:"exists"`
}

func TestKeys(t *testing.T) {
	var p parser.IParser
	p = parser.NewRuntime(SS{})
	fmt.Println(string(strct.NewKeysGenerator().Gen(p.Parser())))

	filename := "./strct_test.go"
	stname := "SS"
	p = parser.NewAST(filename, stname)
	fmt.Println(string(strct.NewKeysGenerator().Gen(p.Parser())))
}

func TestWrapper(t *testing.T) {
	var p parser.IParser
	p = parser.NewRuntime(SS{})
	fmt.Println(string(strct.NewWrapperGenerator().Gen(p.Parser())))

	filename := "./strct_test.go"
	stname := "SS"
	p = parser.NewAST(filename, stname)
	fmt.Println(string(strct.NewWrapperGenerator().Gen(p.Parser())))
}

func TestMapping(t *testing.T) {
	var p parser.IParser
	p = parser.NewRuntime(SS{})
	fmt.Println(string(strct.NewToStringMapGenerator().Gen(p.Parser())))

	filename := "./strct_test.go"
	stname := "SS"
	p = parser.NewAST(filename, stname)
	fmt.Println(string(strct.NewToStringMapGenerator().Gen(p.Parser())))
}
