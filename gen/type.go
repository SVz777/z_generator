package gen

import (
	"github.com/SVz777/z_generator/gen/parser"
)

type IGen interface {
	Gen(*parser.Struct) []byte
}
