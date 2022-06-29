package parser_test

import (
	"fmt"
	"testing"

	"github.com/SVz777/z_generator/gen/parser"
)

type SS struct {
	A *int    `json:"a" validate:"required"`
	B float64 `json:"b" validate:"exists"`
}

func TestParser(t *testing.T) {
	var p parser.IParser
	p = parser.NewRuntime(SS{})
	rt := p.Parser()
	filename := "./parser_test.go"
	stname := "SS"
	p = parser.NewAST(filename, stname)
	at := p.Parser()

	if at.String() != rt.String() {
		fmt.Println("----svz----")
		t.Error("not equal")
	}
}
