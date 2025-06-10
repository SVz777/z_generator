package validate_test

import (
	"fmt"
	"testing"

	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen_impl/validate"
)

//go:generate z_gen -f=$GOFILE -s AA -tp=validate -o=AA.validate_test.go
type AA struct {
	A *int   `json:"a" validate:"range=3,7"`
	B *BB    `json:"b" validate:"child"`
	C string `json:"c" validate:"len=8"`
	D int64  `json:"d" validate:"range=300,700"`
}

//go:generate z_gen -f=$GOFILE -s BB -tp=validate -o=BB.validate_test.go
type BB struct {
	C string `json:"c" validate:"lt=3|gt=6"`
	D string `json:"d" validate:"len=8"`
}

func TestValidate(t *testing.T) {
	var p parser.IParser
	//p = parser.NewRuntime(AA{})
	//fmt.Println(string(validate.NewGenerator().Gen(p.Parser())))

	filename := "./validate_test.go"
	stname := "AA"
	p = parser.NewAST(filename, stname)
	ret := validate.NewGenerator().Gen(p.Parser())
	fmt.Printf("%s\n", ret)
}
