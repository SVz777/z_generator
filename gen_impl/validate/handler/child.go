package handler

import (
	"fmt"
)

type Child struct {
}

func (g *Child) Check(typ string) bool {
	return true
}

func (g *Child) Handle(fieldName string, funcParams ...string) (condition string, err string) {
	return fmt.Sprintf("err := %s.Validate(); err != nil", fieldName),
		fmt.Sprintf(`fmt.Errorf("%s must validate err: %%w", err)`, fieldName)
}
