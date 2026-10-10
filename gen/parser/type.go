package parser

import (
	"reflect"
	"strings"
)

type Field struct {
	Receiver string
	Name     string
	Type     string
	Tag      reflect.StructTag
	Kind     reflect.Kind // reflect.Invalid 表示尚未确认底层类型。
}

type Param struct {
	Name string
	Type string
}

type Func struct {
	Name string
	In   []*Param
	Out  []*Param
}

// Struct 是解析的统一产物：源码走 StructFrom / StructFromFile，运行时走 StructOf。
type Struct struct {
	Package  string
	Receiver string
	Type     string
	Field    []*Field
	Func     []*Func
}

func (s *Struct) String() string {
	ss := strings.Builder{}
	ss.WriteString("-----package-----\n")
	ss.WriteString(s.Package + "\n")
	ss.WriteString("-----struct-----\n")
	ss.WriteString(s.Receiver + " " + s.Type + "\n")
	ss.WriteString("-----field-----\n")
	for _, field := range s.Field {
		ss.WriteString(field.Name + "\t" + field.Type + "\t" + string(field.Tag) + "\n")
	}
	ss.WriteString("-----func-----\n")
	for _, fun := range s.Func {
		ss.WriteString(fun.Name + ":\n")
		// 只打印类型不打印参数名：反射解析器拿不到参数名，
		// 打出来的话 AST 与 Runtime 两种结果就无法比对了。
		ss.WriteString("-----in-----\n")
		for _, param := range fun.In {
			ss.WriteString(param.Type + "\n")
		}
		ss.WriteString("-----out-----\n")
		for _, param := range fun.Out {
			ss.WriteString(param.Type + "\n")
		}
	}
	return ss.String()
}
