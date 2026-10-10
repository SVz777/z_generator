// Package parser 把 Go 源码（go/ast）或运行时类型（reflect）统一解析成 Struct 描述。
//
// 本包不另造一套 AST：解析源文件的产物就是标准库的 *ast.File，
// 调用方可以直接用 go/ast、go/types 的全部能力继续加工。
// 这里只提供两类东西：
//
//   - 把 *ast.File 转成 Struct 的入口：StructFromFile / StructFrom
//   - 可单独复用的表达式辅助函数：TypeName / ParamList / FieldsOf / MethodOf
package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/SVz777/tk/convert"
)

// ParseFile 读取并解析源文件，返回标准库的 *ast.File。
// 一次解析可以反复喂给 StructFrom 提取多个结构体。
func ParseFile(filename string) (*ast.File, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取源文件 %s: %w", filename, err)
	}
	file, err := parser.ParseFile(&token.FileSet{}, filename, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("解析源文件 %s: %w", filename, err)
	}
	return file, nil
}

// StructFromFile 加载文件所属包的 AST 和类型信息，再提取指定结构体。
func StructFromFile(filename, name string) (*Struct, error) {
	file, info, err := loadFileTypes(filename)
	if err != nil {
		return nil, err
	}
	return StructFrom(file, name, info)
}

// StructFrom 从已解析的文件中提取指定结构体及其方法。
// 名字不存在时返回 error——否则会拿到一个空 Struct 继续生成出无意义的代码。
func StructFrom(file *ast.File, name string, infos ...*types.Info) (*Struct, error) {
	var info *types.Info
	if len(infos) > 0 {
		info = infos[0]
	}
	out := &Struct{
		Package:  file.Name.String(),
		Type:     name,
		Receiver: convert.ToCamelCase(name),
	}

	var resolver *kindResolver
	if info == nil {
		resolver = newKindResolver(file)
	}
	found := false
	// 结构体和方法都只在文件顶层声明，直接遍历 Decls 就够了；
	// 用 ast.Inspect 会把每个函数体也递归一遍，纯属白跑。
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			// 单独声明和 `type ( ... )` 分组声明都收在 GenDecl.Specs 里
			for _, spec := range d.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok || ts.Name.String() != name {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					return nil, fmt.Errorf("%s 不是结构体", name)
				}
				fields, err := fieldsOf(st, out.Receiver, resolver, info)
				if err != nil {
					return nil, fmt.Errorf("结构体 %s: %w", name, err)
				}
				out.Field = fields
				found = true
			}
		case *ast.FuncDecl:
			fn, err := MethodOf(d, name)
			if err != nil {
				return nil, err
			}
			if fn != nil {
				out.Func = append(out.Func, fn)
			}
		}
	}
	if !found {
		return nil, fmt.Errorf("在包 %s 中找不到结构体 %s", out.Package, name)
	}

	// 源码顺序和反射顺序（字典序）不同，排一下保证两种解析结果可直接比对。
	sort.Slice(out.Func, func(i, j int) bool { return out.Func[i].Name < out.Func[j].Name })
	return out, nil
}

// FieldsOf 展开结构体字段列表。
// 匿名字段（嵌入字段）跳过：反射侧也拿不到它的完整类型信息，两侧保持一致。
// 字段类型还原不出来时返回 error——空类型名会生成出编译不过的代码，不能静默放过。
func FieldsOf(st *ast.StructType, receiver string, infos ...*types.Info) ([]*Field, error) {
	var info *types.Info
	if len(infos) > 0 {
		info = infos[0]
	}
	return fieldsOf(st, receiver, nil, info)
}

func fieldsOf(st *ast.StructType, receiver string, resolver *kindResolver, info *types.Info) ([]*Field, error) {
	if st.Fields == nil {
		return nil, nil
	}
	var out []*Field
	for _, field := range st.Fields.List {
		if field.Names == nil {
			continue
		}
		typ := TypeName(field.Type)
		if typ == "" {
			return nil, fmt.Errorf("无法还原字段 %s 的类型", field.Names[0].Name)
		}
		var kind reflect.Kind
		if info != nil {
			kind = kindOfType(info.TypeOf(field.Type))
		} else {
			if resolver == nil {
				resolver = newKindResolver(nil)
			}
			kind = resolver.resolve(field.Type)
		}
		// 一个字段声明可能有多个名字：`A, B int`，要逐个展开，
		// 只取 Names[0] 会静默丢字段。
		for _, name := range field.Names {
			out = append(out, &Field{
				Receiver: receiver,
				Name:     name.String(),
				Type:     typ,
				Tag:      TagOf(field.Tag),
				Kind:     kind,
			})
		}
	}
	return out, nil
}

// TagOf 解码结构体字段的 tag 字面量，没有 tag 或字面量无效时返回空。
func TagOf(lit *ast.BasicLit) reflect.StructTag {
	if lit == nil || lit.Kind != token.STRING {
		return ""
	}
	value, err := strconv.Unquote(lit.Value)
	if err != nil {
		return ""
	}
	return reflect.StructTag(value)
}

// MethodOf 提取以 recv 为接收者的方法。
// 普通函数或接收者不匹配时返回 (nil, nil)；类型还原失败时返回 error。
func MethodOf(decl *ast.FuncDecl, recv string) (*Func, error) {
	if decl.Recv == nil || len(decl.Recv.List) == 0 {
		return nil, nil
	}
	// 接收者可能是值类型也可能是指针，统一去掉 * 再比对
	if strings.TrimPrefix(TypeName(decl.Recv.List[0].Type), "*") != recv {
		return nil, nil
	}
	fn := &Func{Name: decl.Name.String()}
	if decl.Type.Params != nil {
		in, err := ParamList(decl.Type.Params.List)
		if err != nil {
			return nil, fmt.Errorf("方法 %s 的入参: %w", fn.Name, err)
		}
		fn.In = in
	}
	if decl.Type.Results != nil {
		out, err := ParamList(decl.Type.Results.List)
		if err != nil {
			return nil, fmt.Errorf("方法 %s 的返回值: %w", fn.Name, err)
		}
		fn.Out = out
	}
	return fn, nil
}

// ParamList 展开参数或返回值列表。
// 未命名参数（如 func(int)）只有类型没有名字；`a, b int` 要逐个展开。
func ParamList(list []*ast.Field) ([]*Param, error) {
	var out []*Param
	for _, f := range list {
		typ := TypeName(f.Type)
		if typ == "" {
			return nil, fmt.Errorf("无法还原类型（参数 %d）", len(out))
		}
		if len(f.Names) == 0 {
			out = append(out, &Param{Type: typ})
			continue
		}
		for _, name := range f.Names {
			out = append(out, &Param{Name: name.String(), Type: typ})
		}
	}
	return out, nil
}

// TypeName 把类型表达式还原成源码里的字符串形式；参数为 nil 时返回空串。
//
// 直接委托 go/types.ExprString，不手写 switch：类型表达式的形态很多
// （数组长度、chan 方向、函数签名、泛型实参、匿名结构体和接口……），
// 手写必然漏，而且漏掉时不是报错，是还原成空串或丢掉一部分——
// 那样生成出的代码要么编译不过，要么类型悄悄变了，比直接失败危险得多。
//
// 需要支持新语法时标准库会跟着 go/parser 一起升级，这里不用改。
func TypeName(expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	return types.ExprString(expr)
}

// TypeKind 获取类型种类
func TypeKind(expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	kind := newKindResolver(nil).resolve(expr)
	if kind == reflect.Invalid {
		return ""
	}
	return kind.String()
}
