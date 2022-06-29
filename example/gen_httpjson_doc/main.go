package main

import (
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
	var p parser.IParser = parser.NewRuntime(v)
	_, _ = os.Stdout.Write(doc.NewGenerator().Gen(p.Parser()))
}
