package handler

import (
	"fmt"

	"github.com/SVz777/z_generator/gen/parser"
)

type Gt struct{}

func (g *Gt) Check(typ string) bool {
	return checkNumberType(typ)
}

func (g *Gt) Handle(field *parser.Field, funcParams ...string) error {
	expr := FieldExpr(field)
	if len(funcParams) != 1 {
		return fmt.Errorf("%w: gt 需要 1 个参数，收到 %d 个: %v", ErrHandle, len(funcParams), funcParams)
	}
	// 失败条件用 <=：等于阈值也不算"大于"
	condition := fmt.Sprintf("%s <= %s", expr, funcParams[0])
	return validationCode(OptionalCondition(field, condition),
		fmt.Sprintf(`fmt.Errorf("%s must be greater than %s")`, expr, funcParams[0]))
}
