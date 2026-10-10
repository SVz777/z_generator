package main

import (
	"fmt"
	"os"

	"github.com/SVz777/z_generator/gen/parser"

	"github.com/SVz777/z_generator/gen_impl/php"
)

type CancelOrder struct {
	OrderId    *uint64 `json:"order_id" validate:"required"`
	CancelType *int    `json:"cancel_type" validate:"required"`

	PassengerId *uint64  `json:"passenger_id,omitempty" validate:"omitempty"`
	OrderExpand *string  `json:"order_expand,omitempty" validate:"omitempty"`
	Lng         *float64 `json:"lng,omitempty" validate:"omitempty"`
	Lat         *float64 `json:"lat,omitempty" validate:"omitempty"`
}

func main() {
	v := CancelOrder{}
	code, err := php.NewGenerator().Gen(parser.StructOf(v))
	if err != nil {
		fmt.Fprintf(os.Stderr, "生成失败: %v\n", err)
		os.Exit(1)
	}
	_, _ = os.Stdout.Write(code)
}
