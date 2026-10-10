package parser

import (
	"reflect"
	"sort"
	"strings"

	"github.com/SVz777/tk/convert"
)

// StructOf 从运行时值反射出结构体描述，与 StructFrom（源码）产出同一套结构。
//
// 通常要传指针：值类型的方法集不含指针接收者的方法，
// 传值会漏掉所有 `func (s *T) Xxx()`。
func StructOf(v any) *Struct {
	out := &Struct{}

	mt := reflect.TypeOf(v)
	// 传指针时类型名和包路径都在 Elem() 上，但方法集要用原类型：
	// 指针的方法集才包含指针接收者的方法。
	t := mt
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	out.Type = t.Name()
	pkgPath := strings.Split(t.PkgPath(), "/")
	out.Package = pkgPath[len(pkgPath)-1]
	out.Receiver = convert.ToCamelCase(out.Type)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		out.Field = append(out.Field, &Field{
			Receiver: out.Receiver,
			Name:     field.Name,
			Type:     field.Type.String(),
			Tag:      field.Tag,
			Kind:     field.Type.Kind(),
		})
	}

	// 反射只能看到导出方法，且拿不到参数名——这是反射的固有边界，
	// 与 AST 解析器对比时只能在导出方法上保持一致。
	for i := 0; i < mt.NumMethod(); i++ {
		m := mt.Method(i)
		fn := &Func{Name: m.Name}
		// 非接口类型的 Method.Type 第一个入参是接收者，跳过
		for j := 1; j < m.Type.NumIn(); j++ {
			fn.In = append(fn.In, &Param{Type: m.Type.In(j).String()})
		}
		for j := 0; j < m.Type.NumOut(); j++ {
			fn.Out = append(fn.Out, &Param{Type: m.Type.Out(j).String()})
		}
		out.Func = append(out.Func, fn)
	}
	// reflect 的 NumMethod 已按名字排序，这里显式排一次是为了不依赖这个实现细节。
	sort.Slice(out.Func, func(i, j int) bool { return out.Func[i].Name < out.Func[j].Name })

	return out
}
