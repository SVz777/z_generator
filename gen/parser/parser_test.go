package parser_test

import (
	"testing"

	"github.com/SVz777/z_generator/gen/parser"
)

type SS struct {
	A *int    `json:"a" validate:"required"`
	B float64 `json:"b" validate:"required"`
}

// multi 的多名字段用来覆盖 `A, B int` 展开；放在文件顶层好让 AST 解析器也能取到。
type multi struct {
	A, B int
}

// typeAll 里的类型是手写类型还原时出过错的样本：
// 数组丢长度、chan/func/匿名结构体/泛型还原成空串、带方法的接口退化成
// interface{}、map 丢掉 value——都会生成出类型不对甚至编译不过的代码。
// 这个文件会被 AST 解析器读取，改字段名记得同步 TestTypeName 的期望值。
type typeAll struct {
	Ptr     **string
	Arr     [3]int
	Chan    chan int
	Recv    <-chan string
	Fn      func(a int) (error, bool)
	Anon    struct{ X int }
	Iface   interface{ Foo() int }
	MapVal  map[interface{}]struct{}
	SliceFn []func() error
	Gen     gen1[int]
}

type gen1[T any] struct{ V T }

// Get 用指针接收者，反射必须传指针才能拿到它（值类型的方法集不含指针接收者方法）。
func (s *SS) Get(a int, b string) (int, error) {
	return 0, nil
}

func (s *SS) Set(name string) {}

// Reset 的参数未命名，用来覆盖「只有类型没有名字」的分支。
func (s *SS) Reset(_ bool) {}

func TestParser(t *testing.T) {
	// 传指针：反射只有拿到指针类型才能看到指针接收者的方法，
	// 传值的话 AST 能解析到 3 个方法、反射只能看到 0 个。
	rt := parser.StructOf(&SS{})

	at, err := parser.StructFromFile("./parser_test.go", "SS")
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	if at.String() != rt.String() {
		t.Errorf("AST 与反射解析结果不一致:\n--- ast ---\n%s\n--- runtime ---\n%s", at, rt)
	}
}

// TestTypeName 逐个断言复杂类型的还原结果，防止类型还原退化成空串或丢掉一部分。
func TestTypeName(t *testing.T) {
	want := map[string]string{
		"Ptr":     "**string",
		"Arr":     "[3]int",
		"Chan":    "chan int",
		"Recv":    "<-chan string",
		"Fn":      "func(a int) (error, bool)",
		"Anon":    "struct{X int}",
		"Iface":   "interface{Foo() int}",
		"MapVal":  "map[interface{}]struct{}",
		"SliceFn": "[]func() error",
		"Gen":     "gen1[int]",
	}

	st, err := parser.StructFromFile("./parser_test.go", "typeAll")
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	if len(st.Field) != len(want) {
		t.Fatalf("字段数 = %d, want %d", len(st.Field), len(want))
	}
	for _, f := range st.Field {
		w, ok := want[f.Name]
		if !ok {
			t.Errorf("解析出预期外的字段 %q", f.Name)
			continue
		}
		if f.Type != w {
			t.Errorf("字段 %s 类型 = %q, want %q", f.Name, f.Type, w)
		}
	}
}

// TestParseFileError 覆盖源文件不存在时必须返回 error，而不是 nil 导致后续 panic。
func TestParseFileError(t *testing.T) {
	if _, err := parser.StructFromFile("./not_exist.go", "SS"); err == nil {
		t.Error("文件不存在时应返回 error")
	}
}

// TestStructFromMiss 覆盖结构体名写错时必须报错：
// 静默返回空 Struct 会让后续生成出一堆没有字段的空壳代码。
func TestStructFromMiss(t *testing.T) {
	file, err := parser.ParseFile("./parser_test.go")
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	if _, err := parser.StructFrom(file, "NoSuchStruct"); err == nil {
		t.Error("结构体不存在时应返回 error")
	}
}

// TestParseFileReuse 一次解析可以反复提取多个结构体，
// 这是把 *ast.File 直接交给调用方（而不是绑死单个结构体名）的直接收益。
func TestParseFileReuse(t *testing.T) {
	file, err := parser.ParseFile("./parser_test.go")
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	for _, name := range []string{"SS", "multi"} {
		st, err := parser.StructFrom(file, name)
		if err != nil {
			t.Fatalf("提取 %s 失败: %v", name, err)
		}
		if st.Type != name {
			t.Errorf("st.Type = %q, want %q", st.Type, name)
		}
		if len(st.Field) == 0 {
			t.Errorf("%s 没有解析到字段", name)
		}
	}
}

// TestMultiNameField 覆盖 `A, B int` 这种多名字段不会丢字段。
func TestMultiNameField(t *testing.T) {
	// 不写 tag，避免 vet 的 structtag 检查报重复 json tag
	if got := len(parser.StructOf(multi{}).Field); got != 2 {
		t.Errorf("字段数 = %d, want 2（多名字段应逐个展开）", got)
	}
}

// TestFuncFilled 覆盖 Struct.Func 会被填充：方法名、入参、出参都要解析出来。
func TestFuncFilled(t *testing.T) {
	st, err := parser.StructFromFile("./parser_test.go", "SS")
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	want := map[string]struct {
		in  []string
		out []string
	}{
		"Get":   {in: []string{"int", "string"}, out: []string{"int", "error"}},
		"Set":   {in: []string{"string"}, out: nil},
		"Reset": {in: []string{"bool"}, out: nil},
	}

	if len(st.Func) != len(want) {
		t.Fatalf("方法数 = %d, want %d: %s", len(st.Func), len(want), st)
	}

	for _, fn := range st.Func {
		w, ok := want[fn.Name]
		if !ok {
			t.Errorf("解析出预期外的方法 %q", fn.Name)
			continue
		}
		if got := types(fn.In); !equal(got, w.in) {
			t.Errorf("%s 入参 = %v, want %v", fn.Name, got, w.in)
		}
		if got := types(fn.Out); !equal(got, w.out) {
			t.Errorf("%s 出参 = %v, want %v", fn.Name, got, w.out)
		}
	}
}

// TestFuncNameFilled 覆盖 AST 解析器能拿到参数名（反射解析器拿不到，这是它独有的信息）。
func TestFuncNameFilled(t *testing.T) {
	st, err := parser.StructFromFile("./parser_test.go", "SS")
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	for _, fn := range st.Func {
		if fn.Name != "Get" {
			continue
		}
		if len(fn.In) != 2 {
			t.Fatalf("Get 入参数 = %d, want 2", len(fn.In))
		}
		if fn.In[0].Name != "a" || fn.In[1].Name != "b" {
			t.Errorf("Get 参数名 = %q, %q; want a, b", fn.In[0].Name, fn.In[1].Name)
		}
		return
	}
	t.Error("没有解析到 Get 方法")
}

func types(params []*parser.Param) []string {
	if len(params) == 0 {
		return nil
	}
	out := make([]string, 0, len(params))
	for _, p := range params {
		out = append(out, p.Type)
	}
	return out
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
