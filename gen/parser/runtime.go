package parser

import (
	"reflect"
	"strings"

	"github.com/SVz777/stringchange"
)

type Runtime struct {
	v interface{}
}

func NewRuntime(v interface{}) *Runtime {
	return &Runtime{v: v}
}

func (rt *Runtime) Parser() (out *Struct) {
	out = &Struct{}

	t := reflect.TypeOf(rt.v)

	out.Type = t.Name()
	pkgPath := strings.Split(t.PkgPath(), "/")
	out.Package = pkgPath[len(pkgPath)-1] //t.PkgPath()
	out.Receiver = stringchange.ToCamelCase(out.Type)

	n := t.NumField()

	for i := 0; i < n; i++ {
		field := t.Field(i)
		fieldName := field.Name
		fieldTp := field.Type.String()
		ff := &Field{}
		ff.Type = fieldTp
		ff.Receiver = out.Receiver
		ff.Name = fieldName
		ff.Tag = field.Tag
		out.Field = append(out.Field, ff)
	}
	return
}
