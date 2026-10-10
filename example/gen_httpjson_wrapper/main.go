package main

import (
	"fmt"
	"os"

	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen_impl/strct"
)

type Create struct {
	A int    `json:"a"`
	B string `json:"B"`
}

func main() {
	code, err := strct.NewWrapperGenerator().Gen(parser.StructOf(Create{}))
	if err != nil {
		fmt.Fprintf(os.Stderr, "生成失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(code))
}
