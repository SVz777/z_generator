package handler

const (
	TagGt     = "gt"
	TagGte    = "gte"
	TagLt     = "lt"
	TagLte    = "lte"
	TagRange  = "range"
	TagMinLen = "min_len"
	TagMaxLen = "max_len"
	TagLen    = "len"
	TagChild  = "child"
)

func init() {
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

type Handler interface {
	Check(typ string) bool
	Handle(fieldName string, funcParams ...string) (condition string, err string)
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
