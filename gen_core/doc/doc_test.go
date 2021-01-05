/**
 * @file    doc_test.go
 * @author
 *  ___  _  _  ____
 * / __)( \/ )(_   )
 * \__ \ \  /  / /_
 * (___/  \/  (____)
 * (903943711@qq.com)
 * @date    2020-01-06
 * @desc
 */
package doc_test

import (
	"fmt"
	"testing"

	"github.com/SVz777/z_generator/gen_core/doc"
	"github.com/SVz777/z_generator/gen_core/gen"
	"github.com/SVz777/z_generator/gen_core/parser"
)

type SS struct {
	A *int    `json:"a" validate:"required"`
	B float64 `json:"b" validate:"exists"`
}

func TestDoc(t *testing.T) {
	var g gen.IGen
	var p parser.IParser
	p = parser.NewRuntime(SS{})
	g = doc.NewDoc(p.Parser())
	fmt.Println(g.Gen().String())

	filename := "./doc_test.go"
	stname := "SS"
	p = parser.NewAST(filename, stname)
	g = doc.NewDoc(p.Parser())
	fmt.Println(g.Gen().String())
}
