package handler

import (
	"fmt"
	"strings"
)

type Range struct {
}

func (g *Range) Check(typ string) bool {
	return checkNumberType(typ)
}

func (g *Range) Handle(fieldName string, funcParams ...string) (condition string, err string) {
	if len(funcParams) != 1 {
		return "false", fmt.Sprintf(`fmt.Errorf("%s must len 1")`, funcParams)
	}
	funcParams = strings.Split(funcParams[0], ",")
	if len(funcParams) != 2 {
		return "false", fmt.Sprintf(`fmt.Errorf("%s must len 2")`, funcParams)
	}
	return fmt.Sprintf("%s <= %s || %s >= %s", fieldName, funcParams[0], fieldName, funcParams[1]), fmt.Sprintf(`fmt.Errorf("%s must be between %s and %s")`, fieldName, funcParams[0], funcParams[1])
}
