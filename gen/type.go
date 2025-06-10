package gen

import (
	"github.com/SVz777/z_generator/gen/parser"
)

type IGen interface {
	Gen(*parser.Struct) []byte
}

var generatorMap = map[string]IGen{}

func Register(name string, gen IGen) {
	generatorMap[name] = gen
}

func GetGenerator(name string) IGen {
	return generatorMap[name]
}
