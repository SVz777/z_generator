package handler

import (
	"fmt"
	"strings"

	"github.com/SVz777/z_generator/gen/parser"
)

type Range struct {
}

func (g *Range) Check(typ string) bool {
	return checkNumberType(typ)
}

func (g *Range) Handle(field *parser.Field, funcParams ...string) error {
	expr := FieldExpr(field)
	if len(funcParams) != 1 {
		return fmt.Errorf("%w: range 需要 1 个参数（形如 1,10），收到 %d 个: %v", ErrHandle, len(funcParams), funcParams)
	}
	bounds := strings.Split(funcParams[0], ",")
	if len(bounds) != 2 {
		return fmt.Errorf("%w: range 需要形如 1,10 的两个边界值，收到 %q", ErrHandle, funcParams[0])
	}
	// 闭区间 [min, max]：越界（含边界外）才算失败
	condition := fmt.Sprintf("%s < %s || %s > %s", expr, bounds[0], expr, bounds[1])
	return validationCode(OptionalCondition(field, condition),
		fmt.Sprintf(`fmt.Errorf("%s must be between %s and %s")`, expr, bounds[0], bounds[1]))
}
