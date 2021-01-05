/**
 * @file    parser_test.go
 * @author
 *  ___  _  _  ____
 * / __)( \/ )(_   )
 * \__ \ \  /  / /_
 * (___/  \/  (____)
 * (903943711@qq.com)
 * @date    2020-01-06
 * @desc
 */
package parser_test

import (
	"testing"

	"github.com/SVz777/z_generator/gen_core/parser"
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
		t.Error("not equal")
	}
}
