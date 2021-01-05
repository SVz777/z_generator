/**
 * @file    main.go
 * @author
 *  ___  _  _  ____
 * / __)( \/ )(_   )
 * \__ \ \  /  / /_
 * (___/  \/  (____)
 * (903943711@qq.com)
 * @date    2020-01-06
 * @desc
 */
package main

import (
	"fmt"
	"go/format"

	"github.com/SVz777/z_generator/gen_core/parser"
	"github.com/SVz777/z_generator/gen_core/strct"
)

type Create struct {
	A int    `json:"a"`
	B string `json:"B"`
}

func main() {
	p := parser.NewRuntime(Create{})
	g := strct.NewStructWrapper(p.Parser())
	s, _ := format.Source(g.Gen().Bytes())
	fmt.Println(string(s))
}
