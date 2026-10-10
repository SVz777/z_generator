package handler

import (
	"fmt"
	"strings"

	"github.com/SVz777/z_generator/gen/parser"
)

// Required 校验字段是否填了值（不为零值）。
//
// 零值判断依赖字段类型：指针/接口 == nil、string == ""、数字 == 0、
// bool == false、slice/map len() == 0。判断不了的类型（比如自定义 struct）
// Check 直接返回 false，由调用方报错——总比生成一段恒不触发的校验代码安全。
type Required struct{}

func (r *Required) Check(typ string) bool {
	_, ok := zeroCondition("", typ)
	return ok
}

func (r *Required) Handle(field *parser.Field, funcParams ...string) error {
	if len(funcParams) != 0 {
		return fmt.Errorf("%w: required 不接受参数，收到 %v", ErrHandle, funcParams)
	}
	// 关键：用未解引用的表达式。指针为 nil 时 *req.Age 会直接 panic，
	// 判断"有没有值"必须看指针本身。
	expr := RawExpr(field)
	cond, ok := zeroCondition(expr, field.Type)
	if !ok {
		return fmt.Errorf("%w: required 无法判断类型 %s 的零值", ErrHandle, field.Type)
	}
	return validationCode(cond, fmt.Sprintf(`fmt.Errorf("%s is required")`, expr))
}

// zeroCondition 生成"字段是否为零值"的判断表达式。
// ok 为 false 表示无法为该类型生成可靠的零值判断。
func zeroCondition(expr, typ string) (string, bool) {
	// 反射解析器给出的类型带空格（interface {}），AST 解析器不带（interface{}），
	// 去掉空格统一比对。
	norm := strings.ReplaceAll(strings.TrimPrefix(typ, "*"), " ", "")
	switch {
	case strings.HasPrefix(typ, "*"): // 指针，零值是 nil
		return expr + " == nil", true
	case norm == "string":
		return expr + ` == ""`, true
	case norm == "bool":
		return expr + " == false", true
	case numberTypes[norm]:
		return expr + " == 0", true
	case norm == "interface{}" || norm == "any":
		return expr + " == nil", true
	case strings.HasPrefix(norm, "[]") || strings.HasPrefix(norm, "map["):
		return "len(" + expr + ") == 0", true
	}
	return "", false
}
