package php_test

import (
	"fmt"
	"testing"

	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen_impl/php"
)

type SS struct {
	A *int    `json:"a" validate:"required"`
	B float64 `json:"b" validate:"exists"`
}

func TestTransform(t *testing.T) {
	var p parser.IParser
	p = parser.NewRuntime(SS{})
	fmt.Println(string(php.NewGenerator().Gen(p.Parser())))

	filename := "./transform_test.go"
	stname := "SS"
	p = parser.NewAST(filename, stname)
	fmt.Println(string(php.NewGenerator().Gen(p.Parser())))
}
