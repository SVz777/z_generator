/**
 * @file    strct_test.go
 * @author
 *  ___  _  _  ____
 * / __)( \/ )(_   )
 * \__ \ \  /  / /_
 * (___/  \/  (____)
 * (903943711@qq.com)
 * @date    2020-01-06
 * @desc
 */
package strct_test

import (
	"fmt"
	"testing"

	"github.com/SVz777/z_generator/gen_core/gen"
	"github.com/SVz777/z_generator/gen_core/parser"
	"github.com/SVz777/z_generator/gen_core/strct"
)

type SS struct {
	A *int    `json:"a" validate:"required"`
	B float64 `json:"b" validate:"exists"`
}

func TestKeys(t *testing.T) {
	var g gen.IGen
	var p parser.IParser
	p = parser.NewRuntime(SS{})
	g = strct.NewKeys(p.Parser())
	fmt.Println(g.Gen().String())

	filename := "./strct_test.go"
	stname := "SS"
	p = parser.NewAST(filename, stname)
	g = strct.NewKeys(p.Parser())
	fmt.Println(g.Gen().String())
}

func TestWrapper(t *testing.T) {
	var g gen.IGen
	var p parser.IParser
	p = parser.NewRuntime(SS{})
	g = strct.NewStructWrapper(p.Parser())
	fmt.Println(g.Gen().String())

	filename := "./strct_test.go"
	stname := "SS"
	p = parser.NewAST(filename, stname)
	g = strct.NewStructWrapper(p.Parser())
	fmt.Println(g.Gen().String())
}

func TestMapping(t *testing.T) {
	var g gen.IGen
	var p parser.IParser
	p = parser.NewRuntime(SS{})
	g = strct.NewStructToStringMap(p.Parser())
	fmt.Println(g.Gen().String())

	filename := "./strct_test.go"
	stname := "SS"
	p = parser.NewAST(filename, stname)
	g = strct.NewStructToStringMap(p.Parser())
	fmt.Println(g.Gen().String())
}
