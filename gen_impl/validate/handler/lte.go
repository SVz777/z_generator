package handler

import (
	"fmt"

	"github.com/SVz777/z_generator/gen/parser"
)

type Lte struct {
}

func (g *Lte) Check(typ string) bool {
	return checkNumberType(typ)
}

func (g *Lte) Handle(field *parser.Field, funcParams ...string) error {
	expr := FieldExpr(field)
	if len(funcParams) != 1 {
		return fmt.Errorf("%w: lte 需要 1 个参数，收到 %d 个: %v", ErrHandle, len(funcParams), funcParams)
	}
	condition := fmt.Sprintf("%s > %s", expr, funcParams[0])
	return validationCode(OptionalCondition(field, condition),
		fmt.Sprintf(`fmt.Errorf("%s must be less than or equal to %s")`, expr, funcParams[0]))
}
