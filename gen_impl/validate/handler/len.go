package handler

import (
	"fmt"
)

type MaxLen struct {
}

func (g *MaxLen) Check(typ string) bool {
	return checkLenType(typ)
}

func (g *MaxLen) Handle(fieldName string, funcParams ...string) (condition string, err string) {
	if len(funcParams) != 1 {
		return "false", fmt.Sprintf(`fmt.Errorf("%s must len 1")`, funcParams)
	}
	return fmt.Sprintf("len(%s) > %s", fieldName, funcParams[0]), fmt.Sprintf(`fmt.Errorf("%s must be at most %s characters long")`, fieldName, funcParams[0])
}

type MinLen struct {
}

func (g *MinLen) Check(typ string) bool {
	return checkLenType(typ)
}

func (g *MinLen) Handle(fieldName string, funcParams ...string) (condition string, err string) {
	if len(funcParams) != 1 {
		return "false", fmt.Sprintf(`fmt.Errorf("%s must len 1")`, funcParams)
	}
	return fmt.Sprintf("len(%s) < %s", fieldName, funcParams[0]), fmt.Sprintf(`fmt.Errorf("%s must be at least %s characters long")`, fieldName, funcParams[0])
}

type Len struct {
}

func (g *Len) Check(typ string) bool {
	return checkLenType(typ)
}

func (g *Len) Handle(fieldName string, funcParams ...string) (condition string, err string) {
	if len(funcParams) != 1 {
		return "false", fmt.Sprintf(`fmt.Errorf("%s must len 1")`, funcParams)
	}
	return fmt.Sprintf("len(%s) != %s", fieldName, funcParams[0]), fmt.Sprintf(`fmt.Errorf("%s must be %s characters long")`, fieldName, funcParams[0])
}
