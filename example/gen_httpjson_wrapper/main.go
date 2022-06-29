package main

import (
	"fmt"

	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen_impl/strct"
)

type Create struct {
	A int    `json:"a"`
	B string `json:"B"`
}

func main() {
	p := parser.NewRuntime(Create{})
	fmt.Println(string(strct.NewWrapperGenerator().Gen(p.Parser())))
}
