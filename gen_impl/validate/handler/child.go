package handler

import (
	"fmt"
	"strings"

	"github.com/SVz777/z_generator/gen/parser"
)

type Child struct {
}

func (g *Child) Check(typ string) bool {
	base := strings.TrimPrefix(typ, "*")
	if base == "" || strings.ContainsAny(base, "[]{}() ") {
		return false
	}
	return !baseTypes[base]
}

func (g *Child) Handle(field *parser.Field, funcParams ...string) error {
	expr := RawExpr(field)
	if len(funcParams) != 0 {
		return fmt.Errorf("%w: child 不接受参数，收到 %v", ErrHandle, funcParams)
	}
	condition := fmt.Sprintf("err := %s.Validate(); err != nil", expr)
	if strings.HasPrefix(field.Type, "*") {
		condition = fmt.Sprintf("err := func() error { if %s == nil { return nil }; return %s.Validate() }(); err != nil", expr, expr)
	}
	return validationCode(condition, fmt.Sprintf(`fmt.Errorf("%s must validate err: %%w", err)`, expr))
}
