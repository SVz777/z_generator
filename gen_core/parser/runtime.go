/**
 * @file    runtime.go
 * @author
 *  ___  _  _  ____
 * / __)( \/ )(_   )
 * \__ \ \  /  / /_
 * (___/  \/  (____)
 * (903943711@qq.com)
 * @date    2020-01-06
 * @desc
 */
package parser

import (
	"reflect"
	"strings"
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
	out.Receiver = strings.ToLower(out.Type)

	n := t.NumField()

	for i := 0; i < n; i++ {
		field := t.Field(i)
		fieldName := field.Name
		fieldTp := field.Type.String()
		ff := &Field{}
		ff.Type = fieldTp
		ff.Name = fieldName
		ff.Tag = field.Tag
		out.Field = append(out.Field, ff)
	}
	return
}
