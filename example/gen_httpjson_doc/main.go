package main

import (
	"fmt"
	"os"

	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen_impl/doc"
)

type CancelOrder struct {
	A *uint64 `json:"a" validate:"required"`
	B *int    `json:"b" validate:"required"`

	C *uint64  `json:"c,omitempty" validate:"omitempty"`
	D *string  `json:"d,omitempty" validate:"omitempty"`
	E *float64 `json:"e,omitempty" validate:"omitempty"`
	F *float64 `json:"f,omitempty" validate:"omitempty"`
}

func main() {
	v := CancelOrder{}
	code, err := doc.NewGenerator().Gen(parser.StructOf(v))
	if err != nil {
		fmt.Fprintf(os.Stderr, "生成失败: %v\n", err)
		os.Exit(1)
	}
	_, _ = os.Stdout.Write(code)
}
