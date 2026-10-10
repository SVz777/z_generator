package handler

import (
	"fmt"

	"github.com/SVz777/z_generator/gen/parser"
)

type MaxLen struct {
}

func (g *MaxLen) Check(typ string) bool {
	return checkLenType(typ)
}

func (g *MaxLen) Handle(field *parser.Field, funcParams ...string) error {
	expr := FieldExpr(field)
	if len(funcParams) != 1 {
		return fmt.Errorf("%w: max_len 需要 1 个参数，收到 %d 个: %v", ErrHandle, len(funcParams), funcParams)
	}
	condition := fmt.Sprintf("len(%s) > %s", expr, funcParams[0])
	return validationCode(OptionalCondition(field, condition),
		fmt.Sprintf(`fmt.Errorf("%s must be at most %s characters long")`, expr, funcParams[0]))
}

type MinLen struct {
}

func (g *MinLen) Check(typ string) bool {
	return checkLenType(typ)
}

func (g *MinLen) Handle(field *parser.Field, funcParams ...string) error {
	expr := FieldExpr(field)
	if len(funcParams) != 1 {
		return fmt.Errorf("%w: min_len 需要 1 个参数，收到 %d 个: %v", ErrHandle, len(funcParams), funcParams)
	}
	condition := fmt.Sprintf("len(%s) < %s", expr, funcParams[0])
	return validationCode(OptionalCondition(field, condition),
		fmt.Sprintf(`fmt.Errorf("%s must be at least %s characters long")`, expr, funcParams[0]))
}

type Len struct {
}

func (g *Len) Check(typ string) bool {
	return checkLenType(typ)
}

func (g *Len) Handle(field *parser.Field, funcParams ...string) error {
	expr := FieldExpr(field)
	if len(funcParams) != 1 {
		return fmt.Errorf("%w: len 需要 1 个参数，收到 %d 个: %v", ErrHandle, len(funcParams), funcParams)
	}
	condition := fmt.Sprintf("len(%s) != %s", expr, funcParams[0])
	return validationCode(OptionalCondition(field, condition),
		fmt.Sprintf(`fmt.Errorf("%s must be %s characters long")`, expr, funcParams[0]))
}
