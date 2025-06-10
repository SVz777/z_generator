package handler

import (
	"fmt"
)

type Lt struct {
}

func (g *Lt) Check(typ string) bool {
	return checkNumberType(typ)
}

func (g *Lt) Handle(fieldName string, funcParams ...string) (condition string, err string) {
	if len(funcParams) != 1 {
		return "false", fmt.Sprintf(`fmt.Errorf("%s must len 1")`, funcParams)
	}
	return fmt.Sprintf("%s >= %s", fieldName, funcParams[0]), fmt.Sprintf(`fmt.Errorf("%s must be less than %s")`, fieldName, funcParams[0])
}
