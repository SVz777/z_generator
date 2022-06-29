package main

import (
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
	var p parser.IParser = parser.NewRuntime(v)
	_, _ = os.Stdout.Write(php.NewGenerator().Gen(p.Parser()))
}
