package validate_test

import (
	"strings"
	"testing"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen_impl/validate"
)

//go:generate z_gen -f=$GOFILE -s AA -tp=validate -o=AA.validate_test.go
type AA struct {
	A *int   `json:"a" validate:"range=3,7"`
	B *BB    `json:"b" validate:"child"`
	C string `json:"c" validate:"len=8"`
	D int64  `json:"d" validate:"range=300,700"`
}

type optionalPointerCase struct {
	Number *int    `validate:"range=3,7"`
	Text   *string `validate:"min_len=2|max_len=8|len=4"`
}

//go:generate z_gen -f=$GOFILE -s BB -tp=validate -o=BB.validate_test.go
type BB struct {
	C string `json:"c" validate:"len=8"`
}

// requiredCase 覆盖 required 在各类型下的零值判断。
type requiredCase struct {
	Ptr   *int           `validate:"required"`
	Str   string         `validate:"required"`
	Num   int64          `validate:"required"`
	Slice []string       `validate:"required"`
	Map   map[string]int `validate:"required"`
	Flag  bool           `validate:"required"`
	Iface interface{}    `validate:"required"`
}

// requiredUnsupported 的 Inner 是自定义 struct，required 判断不了它的零值，应当报错。
type requiredUnsupported struct {
	Inner requiredCase `validate:"required"`
}

// badRule 规则名不存在。改造前会降级成 `if false`，校验静默失效。
type badRule struct {
	A int `validate:"nosuchrule"`
}

// badParam gt 需要一个参数却没给。
type badParam struct {
	A int `validate:"gt"`
}

// badType 数值规则用在了 string 字段上。
type badType struct {
	A string `validate:"gt=3"`
}

type badChildType struct {
	A int `validate:"child"`
}

type badChildPointerType struct {
	A *[]int `validate:"child"`
}

func mustGen(t *testing.T, g gen.IGen, st *parser.Struct) string {
	t.Helper()
	code, err := g.Gen(st)
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}
	return string(code)
}

func genStruct(t *testing.T, name string) (string, error) {
	t.Helper()
	v, err := parser.StructFromFile("./validate_test.go", name)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	code, err := validate.NewGenerator().Gen(v)
	return string(code), err
}

func TestValidate(t *testing.T) {
	st, err := parser.StructFromFile("./validate_test.go", "AA")
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	code := mustGen(t, validate.NewGenerator(), st)
	t.Log("\n" + code)

	// 改造前 format.Source 失败时错误信息会被当成代码返回，
	// 这里断言生成结果里不含 go/parser 的典型错误文本。
	if strings.Contains(code, "expected") {
		t.Errorf("生成结果疑似包含语法错误信息:\n%s", code)
	}
	if !strings.Contains(code, "aA.B == nil") {
		t.Errorf("child 指针字段缺少 nil 防护:\n%s", code)
	}
}

// TestRequired 覆盖 required 在各类型下生成的零值判断。
// 指针必须不解引用（*p 在 p 为 nil 时会 panic），这是最容易写错的地方。
func TestOptionalPointerRules(t *testing.T) {
	code, err := genStruct(t, "optionalPointerCase")
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}

	want := []string{
		"optionalPointerCase.Number != nil && (*optionalPointerCase.Number < 3 || *optionalPointerCase.Number > 7)",
		"optionalPointerCase.Text != nil && (len(*optionalPointerCase.Text) < 2)",
		"optionalPointerCase.Text != nil && (len(*optionalPointerCase.Text) > 8)",
		"optionalPointerCase.Text != nil && (len(*optionalPointerCase.Text) != 4)",
	}
	for _, w := range want {
		if !strings.Contains(code, w) {
			t.Errorf("生成结果里缺少 %q:\n%s", w, code)
		}
	}
}

func TestRequired(t *testing.T) {
	code, err := genStruct(t, "requiredCase")
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}
	t.Log("\n" + code)

	want := []string{
		"requiredCase.Ptr == nil",
		`requiredCase.Str == ""`,
		"requiredCase.Num == 0",
		"len(requiredCase.Slice) == 0",
		"len(requiredCase.Map) == 0",
		"requiredCase.Flag == false",
		"requiredCase.Iface == nil",
	}
	for _, w := range want {
		if !strings.Contains(code, w) {
			t.Errorf("生成结果里缺少 %q:\n%s", w, code)
		}
	}
	// 指针绝不能被解引用来判断 nil
	if strings.Contains(code, "*requiredCase.Ptr") {
		t.Errorf("指针字段被解引用了，nil 时会 panic:\n%s", code)
	}
}

// TestGenError 覆盖几条出错路径：规则写错必须让生成失败，
// 而不是降级成恒不触发的校验代码。
func TestGenError(t *testing.T) {
	cases := []struct {
		name string
		st   string
	}{
		{"规则名不存在", "badRule"},
		{"缺少必要参数", "badParam"},
		{"类型不支持该规则", "badType"},
		{"child 不支持基础类型", "badChildType"},
		{"child 不支持复合类型指针", "badChildPointerType"},
		{"required 判断不了 struct 零值", "requiredUnsupported"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, err := genStruct(t, c.st)
			if err == nil {
				t.Fatalf("期望生成失败，实际成功:\n%s", code)
			}
			t.Logf("如期报错: %v", err)
		})
	}
}
