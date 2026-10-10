package handler

import (
	"errors"
	"strings"

	"github.com/SVz777/z_generator/gen/parser"
)

const (
	TagRequired = "required"
	TagGt       = "gt"
	TagGte      = "gte"
	TagLt       = "lt"
	TagLte      = "lte"
	TagRange    = "range"
	TagMinLen   = "min_len"
	TagMaxLen   = "max_len"
	TagLen      = "len"
	TagChild    = "child"
)

func init() {
	RegisterHandler(TagRequired, &Required{})
	RegisterHandler(TagGt, &Gt{})
	RegisterHandler(TagGte, &Gte{})
	RegisterHandler(TagLt, &Lt{})
	RegisterHandler(TagLte, &Lte{})
	RegisterHandler(TagRange, &Range{})
	RegisterHandler(TagMinLen, &MinLen{})
	RegisterHandler(TagMaxLen, &MaxLen{})
	RegisterHandler(TagLen, &Len{})
	RegisterHandler(TagChild, &Child{})
}

var ErrHandle = errors.New("validation handler failed")

// Handler 把一条校验规则翻译成可直接嵌入生成代码的 Go 表达式。
type Handler interface {
	// Check 判断该处理器是否支持 typ 类型的字段。
	Check(typ string) bool

	// Handle 返回承载完整校验代码文本的 error。
	// 规则参数不合法时用 fmt.Errorf 包装 ErrHandle，调用方通过 errors.Is 判断。
	//
	// 普通 error 的 Error() 文本会被写进生成文件，并不代表处理失败。
	// 传整个 field 而不是预处理好的表达式，是因为要不要解引用取决于规则本身：
	// 比较大小需要 *req.Age，但判断是否为 nil 必须是 req.Age——解引用 nil 指针会 panic。
	Handle(field *parser.Field, funcParams ...string) error
}

func validationCode(condition, errExpr string) error {
	return errors.New("if " + condition + " {\nreturn " + errExpr + "\n}")
}

// baseTypes 基本类型，指针解引用只对它们安全。
var baseTypes = map[string]bool{
	"uint":    true,
	"uint8":   true,
	"uint16":  true,
	"uint32":  true,
	"uint64":  true,
	"int":     true,
	"int8":    true,
	"int16":   true,
	"int32":   true,
	"int64":   true,
	"float32": true,
	"float64": true,
	"string":  true,
	"bool":    true,
}

// FieldExpr 返回字段用于比较运算的访问表达式。
// 基本类型的指针要解引用才能参与比较：*int 不能直接和整数比大小。
func FieldExpr(field *parser.Field) string {
	expr := field.Receiver + "." + field.Name
	if strings.HasPrefix(field.Type, "*") && baseTypes[strings.TrimPrefix(field.Type, "*")] {
		return "*" + expr
	}
	return expr
}

// RawExpr 返回字段未经解引用的访问表达式，判断指针是否为 nil 时用它。
func RawExpr(field *parser.Field) string {
	return field.Receiver + "." + field.Name
}

// OptionalCondition 让可选指针字段在 nil 时跳过非 required 校验，避免解引用 nil 指针。
func OptionalCondition(field *parser.Field, condition string) string {
	if strings.HasPrefix(field.Type, "*") {
		return RawExpr(field) + " != nil && (" + condition + ")"
	}
	return condition
}

var handlers = map[string]Handler{}

func GetHandler(funcName string) Handler {
	return handlers[funcName]
}

func RegisterHandler(funcName string, handler Handler) {
	handlers[funcName] = handler
}

var numberTypes = map[string]bool{
	"uint":    true,
	"uint8":   true,
	"uint16":  true,
	"uint32":  true,
	"uint64":  true,
	"int":     true,
	"int8":    true,
	"int16":   true,
	"int32":   true,
	"int64":   true,
	"float32": true,
	"float64": true,
}

func checkNumberType(typ string) bool {
	if len(typ) == 0 {
		return false
	}
	if typ[0] == '*' {
		typ = typ[1:]
	}
	return numberTypes[typ]
}

var lenTypes = map[string]bool{
	"string": true,
	"[]byte": true,
}

func checkLenType(typ string) bool {
	if len(typ) == 0 {
		return false
	}
	if typ[0] == '*' {
		typ = typ[1:]
	}
	return lenTypes[typ]
}
