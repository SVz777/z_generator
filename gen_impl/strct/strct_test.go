package strct_test

import (
	"go/ast"
	goparser "go/parser"
	"go/token"
	"go/types"
	"strings"
	"testing"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen_impl/strct"
)

type SS struct {
	A *int    `json:"a" validate:"required"`
	B float64 `json:"b" validate:"exists"`
}

type unsupportedMapField struct {
	Flag bool `json:"flag"`
}

func mustGen(t *testing.T, g gen.IGen, st *parser.Struct) string {
	t.Helper()
	code, err := g.Gen(st)
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}
	return string(code)
}

type wrapperInterface interface{ Read() }
type wrapperAlias = wrapperInterface
type wrapperValue struct{ N int }
type wrapperFixture struct {
	Reader wrapperInterface
	Alias  wrapperAlias
	Value  wrapperValue
}

func TestWrapperNamedInterfaceCompiles(t *testing.T) {
	const declarations = "package strct_test; type wrapperInterface interface{ Read() }; type wrapperAlias = wrapperInterface; type wrapperValue struct{ N int }; type wrapperFixture struct{ Reader wrapperInterface; Alias wrapperAlias; Value wrapperValue }"
	parsed, err := parser.StructFromFile("strct_test.go", "wrapperFixture")
	if err != nil {
		t.Fatal(err)
	}
	for _, st := range []*parser.Struct{parsed, parser.StructOf(wrapperFixture{})} {
		code := mustGen(t, strct.NewWrapperGenerator(), st)
		if !strings.Contains(code, ".Reader == nil") || !strings.Contains(code, ".Alias == nil") {
			t.Fatalf("接口应与 nil 比较:\n%s", code)
		}
		code = strings.ReplaceAll(code, "strct_test.wrapperInterface", "wrapperInterface")
		code = strings.ReplaceAll(code, "strct_test.wrapperValue", "wrapperValue")
		fset := token.NewFileSet()
		file, err := goparser.ParseFile(fset, "generated.go", declarations+"\n"+code, 0)
		if err != nil {
			t.Fatal(err)
		}
		conf := types.Config{}
		if _, err := conf.Check("fixture", fset, []*ast.File{file}, nil); err != nil {
			t.Fatalf("生成代码类型检查失败: %v\n%s", err, code)
		}
	}
}

func TestWrapperUnknownNamedType(t *testing.T) {
	st := &parser.Struct{Type: "S", Receiver: "s", Field: []*parser.Field{{Name: "Reader", Type: "UnknownReader"}}}
	_, err := strct.NewWrapperGenerator().Gen(st)
	if err == nil || !strings.Contains(err.Error(), "UnknownReader") {
		t.Fatalf("未知类型应明确报错: %v", err)
	}
}

func TestKeys(t *testing.T) {
	g := strct.NewKeysGenerator()
	st := parser.StructOf(SS{})
	t.Log("\n" + mustGen(t, g, st))

	st, err := parser.StructFromFile("./strct_test.go", "SS")
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	t.Log("\n" + mustGen(t, g, st))
}

func TestWrapper(t *testing.T) {
	g := strct.NewWrapperGenerator()
	st := parser.StructOf(SS{})
	t.Log("\n" + mustGen(t, g, st))

	st, err := parser.StructFromFile("./strct_test.go", "SS")
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	t.Log("\n" + mustGen(t, g, st))
}

// TestWrapperValueType 是回归测试：模板曾对所有字段一律生成 `== nil` 和 `= &v`，
// 遇到值类型字段（SS.B 是 float64）生成的代码编译不过。
func TestWrapperValueType(t *testing.T) {
	st := parser.StructOf(SS{})
	code := mustGen(t, strct.NewWrapperGenerator(), st)

	// 值类型：Empty 判零值、Set 直接赋值
	if !strings.Contains(code, ".B == 0") {
		t.Errorf("值类型字段 B 的 Empty 应判零值:\n%s", code)
	}
	if !strings.Contains(code, ".B = v") {
		t.Errorf("值类型字段 B 的 Set 应直接赋值:\n%s", code)
	}
	// 指针类型：Empty 判 nil、Set 取地址
	if !strings.Contains(code, ".A == nil") {
		t.Errorf("指针字段 A 的 Empty 应判 nil:\n%s", code)
	}
	if !strings.Contains(code, ".A = &v") {
		t.Errorf("指针字段 A 的 Set 应取地址:\n%s", code)
	}
}

func TestMapping(t *testing.T) {
	g := strct.NewToStringMapGenerator()
	st := parser.StructOf(SS{})
	t.Log("\n" + mustGen(t, g, st))

	st, err := parser.StructFromFile("./strct_test.go", "SS")
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	t.Log("\n" + mustGen(t, g, st))
}

func TestMappingPointerUsesValueGetter(t *testing.T) {
	code := mustGen(t, strct.NewToStringMapGenerator(), parser.StructOf(SS{}))
	if !strings.Contains(code, "strconv.Itoa(sS.GetValueA())") {
		t.Fatalf("指针数值字段应使用值 getter:\n%s", code)
	}
}

func TestMappingRejectsUnsupportedField(t *testing.T) {
	_, err := strct.NewToStringMapGenerator().Gen(parser.StructOf(unsupportedMapField{}))
	if err == nil {
		t.Fatal("不支持的字段类型应返回错误")
	}
	if !strings.Contains(err.Error(), "Flag") || !strings.Contains(err.Error(), "bool") {
		t.Fatalf("错误信息应包含字段名和类型: %v", err)
	}
}
