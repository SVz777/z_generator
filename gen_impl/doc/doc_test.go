package doc_test

import (
	"fmt"
	"testing"

	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen_impl/doc"
)

type SS struct {
	A *int    `json:"a" validate:"required"`
	B float64 `json:"b" validate:"exists"`
}

func TestDoc(t *testing.T) {
	var p parser.IParser
	p = parser.NewRuntime(SS{})
	fmt.Println(string(doc.NewGenerator().Gen(p.Parser())))

	filename := "./doc_test.go"
	stname := "SS"
	p = parser.NewAST(filename, stname)
	fmt.Println(string(doc.NewGenerator().Gen(p.Parser())))
}
