package parser

import (
	"reflect"
	"strings"
)

type IParser interface {
	Parser() (out *Struct)
}

type Field struct {
	Name string
	Type string
	Tag  reflect.StructTag
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

type Struct struct {
	Receiver string
	Type     string
	Field    []*Field
	Func     []*Func
}

func (s *Struct) String() string {
	ss := strings.Builder{}
	ss.WriteString("-----struct-----\n")
	ss.WriteString(s.Receiver + " " + s.Type + "\n")
	ss.WriteString("-----field-----\n")
	for _, field := range s.Field {
		ss.WriteString(field.Name + "\t" + field.Type + "\t" + string(field.Tag) + "\n")
	}
	ss.WriteString("-----func-----\n")
	for _, fun := range s.Func {
		ss.WriteString(fun.Name + ":\n")
		ss.WriteString("-----in-----\n")
		for _, param := range fun.In {
			ss.WriteString(param.Name + "\t" + param.Type + "\n")
		}
		ss.WriteString("-----out-----\n")
		for _, param := range fun.Out {
			ss.WriteString(param.Name + "\t" + param.Type + "\n")
		}
	}
	return ss.String()
}
