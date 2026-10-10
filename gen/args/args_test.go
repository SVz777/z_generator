package args

import (
	"reflect"
	"strings"
	"testing"
)

func TestSplitArgs(t *testing.T) {
	known := map[string]bool{"f": true, "s": true, "tp": true, "inplace": true}

	cases := []struct {
		name  string
		args  []string
		rest  []string
		extra map[string]string
	}{
		{
			name:  "混合写法",
			args:  []string{"-f", "a.go", "-s", "Foo", "-tp", "validate", "-vt", "custom", "-num", "3", "--flagx"},
			rest:  []string{"-f", "a.go", "-s", "Foo", "-tp", "validate"},
			extra: map[string]string{"vt": "custom", "num": "3", "flagx": "true"},
		},
		{
			// 未声明的参数没法像 flag 包那样靠类型判断有没有值，
			// 所以未知 flag 后面紧跟的非 flag 参数一律当作它的值。
			// 位置参数要么放在未知 flag 之前，要么用 -- 显式分隔。
			name:  "未知 flag 后的位置参数会被当成它的值",
			args:  []string{"-f", "a.go", "--flagx", "tailpos"},
			rest:  []string{"-f", "a.go"},
			extra: map[string]string{"flagx": "tailpos"},
		},
		{
			// 与标准 flag 包保持一致：碰到第一个位置参数就停止解析，
			// 后面的 -vt custom 会被当作位置参数原样保留。
			name:  "位置参数之后不再解析后续 flag",
			args:  []string{"-f", "a.go", "tailpos", "-vt", "custom"},
			rest:  []string{"-f", "a.go", "tailpos", "-vt", "custom"},
			extra: map[string]string{},
		},
		{
			name:  "用双横线分隔出位置参数",
			args:  []string{"-f", "a.go", "-vt", "custom", "--", "tailpos", "-x"},
			rest:  []string{"-f", "a.go", "--", "tailpos", "-x"},
			extra: map[string]string{"vt": "custom"},
		},
		{
			name:  "等号写法与单双横线等价",
			args:  []string{"-k=v1", "--k2=v2"},
			rest:  []string{},
			extra: map[string]string{"k": "v1", "k2": "v2"},
		},
		{
			name:  "bool 开关不带值",
			args:  []string{"-strict", "-f", "a.go"},
			rest:  []string{"-f", "a.go"},
			extra: map[string]string{"strict": "true"},
		},
		{
			name:  "未知开关后紧跟已知参数",
			args:  []string{"-f", "a.go", "-strict", "-s", "Foo"},
			rest:  []string{"-f", "a.go", "-s", "Foo"},
			extra: map[string]string{"strict": "true"},
		},
		{
			name:  "内置 bool 参数后的位置参数不被误吞",
			args:  []string{"-f", "a.go", "-s", "Foo", "-inplace", "mydir"},
			rest:  []string{"-f", "a.go", "-s", "Foo", "-inplace", "mydir"},
			extra: map[string]string{},
		},
		{
			name:  "负数值不被当作 flag",
			args:  []string{"-f", "a.go", "-min", "-3"},
			rest:  []string{"-f", "a.go"},
			extra: map[string]string{"min": "-3"},
		},
		{
			name:  "双横线终止符后的内容原样保留",
			args:  []string{"-f", "a.go", "--", "-notaflag", "x"},
			rest:  []string{"-f", "a.go", "--", "-notaflag", "x"},
			extra: map[string]string{},
		},
		{
			name:  "没有多余参数",
			args:  []string{"-f", "a.go", "-s", "Foo"},
			rest:  []string{"-f", "a.go", "-s", "Foo"},
			extra: map[string]string{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			extra := map[string]string{}
			rest := splitArgs(c.args, known, &extra)
			if !reflect.DeepEqual(rest, c.rest) {
				t.Errorf("rest 不符\n got=%v\nwant=%v", rest, c.rest)
			}
			if !reflect.DeepEqual(extra, c.extra) {
				t.Errorf("extra 不符\n got=%v\nwant=%v", extra, c.extra)
			}
		})
	}
}

func TestGet(t *testing.T) {
	c := &argSet{Extra: map[string]string{
		"validate_tag":     "check",
		"validate_tag_sep": ",",
		"strict":           "true",
		"off":              "false",
		"num":              "12",
		"bad":              "abc",
		"empty":            "",
	}}

	if got := c.GetOr("validate_tag", "validate"); got != "check" {
		t.Errorf("GetOr 已存在 = %q", got)
	}
	if got := c.GetOr("validate_tag_sep", "|"); got != "," {
		t.Errorf("GetOr 分隔符 = %q", got)
	}
	if got := c.GetOr("missing", "dft"); got != "dft" {
		t.Errorf("GetOr 缺失 = %q, want dft", got)
	}
	if got := c.GetOr("empty", "dft"); got != "dft" {
		t.Errorf("GetOr 空值 = %q, want dft", got)
	}
	if got := c.Get("missing"); got != "" {
		t.Errorf("Get 缺失 = %q, want 空", got)
	}
	if !c.GetBool("strict") {
		t.Error("GetBool strict 应为 true")
	}
	if c.GetBool("off") {
		t.Error("GetBool off 应为 false")
	}
	if c.GetBool("bad") {
		t.Error("GetBool 解析失败应返回 false")
	}
	if c.GetBool("missing") {
		t.Error("GetBool 缺失应返回 false")
	}
	if got := c.GetInt("num", 0); got != 12 {
		t.Errorf("GetInt num = %d, want 12", got)
	}
	if got := c.GetInt("bad", 7); got != 7 {
		t.Errorf("GetInt 解析失败 = %d, want 7", got)
	}
	if got := c.GetInt("missing", 7); got != 7 {
		t.Errorf("GetInt 缺失 = %d, want 7", got)
	}
}

func TestSub(t *testing.T) {
	c := &argSet{Extra: map[string]string{
		"php.pkg":    "model",
		"php.strict": "true",
		"validate":   "x",
	}}
	got := c.Sub("php.")
	want := map[string]string{"pkg": "model", "strict": "true"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Sub = %v, want %v", got, want)
	}
	if len(c.Sub("none.")) != 0 {
		t.Error("Sub 无匹配应返回空 map")
	}
}

func TestPeekFlagValue(t *testing.T) {
	cases := []struct {
		args []string
		want string
	}{
		{[]string{"-tp", "validate", "-f", "a.go"}, "validate"},
		{[]string{"-f", "a.go", "-tp", "validate"}, "validate"},
		{[]string{"-tp=validate"}, "validate"},
		{[]string{"--tp", "validate"}, "validate"},
		{[]string{"-f", "a.go"}, ""},
		{[]string{}, ""},
		{[]string{"-tp"}, ""},
	}
	for _, c := range cases {
		if got := peekFlagValue(c.args, "tp"); got != c.want {
			t.Errorf("peekFlagValue(%v) = %q, want %q", c.args, got, c.want)
		}
	}
	// -h 写在 -tp 之前也能取到，这是提前预扫描的意义所在
	if got := peekFlagValue([]string{"-h", "-tp", "validate"}, "tp"); got != "validate" {
		t.Errorf("-h 在前时应仍能取到 tp，got %q", got)
	}
}

func TestFilterDeclared(t *testing.T) {
	list := []declaredFlag{
		{GenTyp: "php", Key: "php.pkg", Usage: "u1"},
		{GenTyp: "validate", Key: "validate_tag", Usage: "u2"},
		{GenTyp: "validate", Key: "validate_tag_sep", Usage: "u3"},
	}

	if got := filterDeclared(list, "validate"); len(got) != 2 {
		t.Fatalf("filterDeclared(validate) 应返回 2 条, got %d", len(got))
	} else {
		if got[0].Key != "validate_tag" || got[1].Key != "validate_tag_sep" {
			t.Errorf("filterDeclared(validate) 内容不符: %v", got)
		}
	}
	if got := filterDeclared(list, "php"); len(got) != 1 || got[0].Key != "php.pkg" {
		t.Errorf("filterDeclared(php) = %v", got)
	}
	if got := filterDeclared(list, "doc"); len(got) != 0 {
		t.Errorf("未声明的生成器应返回空, got %v", got)
	}
	if got := filterDeclared(list, ""); len(got) != 3 {
		t.Errorf("genTyp 为空应返回全部, got %d", len(got))
	}
}

func TestDeclaredGenTypes(t *testing.T) {
	list := []declaredFlag{
		{GenTyp: "validate", Key: "validate_tag"},
		{GenTyp: "php", Key: "php.pkg"},
		{GenTyp: "validate", Key: "validate_tag_sep"},
		{GenTyp: "php", Key: "php.strict"},
	}
	got := declaredGenTypes(list)
	want := []string{"php", "validate"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("declaredGenTypes = %v, want %v (应去重且按字典序)", got, want)
	}
	// GenTyp 为空的声明不参与索引
	if got := declaredGenTypes([]declaredFlag{{Key: "orphan"}}); len(got) != 0 {
		t.Errorf("GenTyp 为空应被忽略, got %v", got)
	}
	if got := declaredGenTypes(nil); len(got) != 0 {
		t.Errorf("空输入应返回空, got %v", got)
	}
}

func TestWrapList(t *testing.T) {
	cases := []struct {
		name  string
		items []string
		width int
		want  string
	}{
		{
			name:  "不折行",
			items: []string{"doc", "php"},
			width: 60,
			want:  "doc, php",
		},
		{
			name:  "超出宽度折行",
			items: []string{"struct2map", "struct_keys", "struct_wrapper", "validate"},
			width: 30,
			want:  "struct2map, struct_keys,\nstruct_wrapper, validate",
		},
		{
			name:  "单个超长元素不拆开",
			items: []string{"a_very_long_type_name", "b"},
			width: 10,
			want:  "a_very_long_type_name,\nb",
		},
		{
			name:  "空列表",
			items: nil,
			width: 60,
			want:  "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := wrapList(tc.items, tc.width); got != tc.want {
				t.Errorf("wrapList = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestTpUsage(t *testing.T) {
	// 未注入时退化成一句话，不能是空的
	old := typeLister
	typeLister = nil
	defer func() { typeLister = old }()
	if got := tpUsage(); got != "生成类型" {
		t.Errorf("未注入时 tpUsage = %q, want %q", got, "生成类型")
	}

	typeLister = func() []string { return []string{"doc", "php"} }
	got := tpUsage()
	if !strings.HasPrefix(got, "生成类型，可选:\n") {
		t.Errorf("tpUsage 应列出可选类型, got %q", got)
	}
	if !strings.Contains(got, "doc, php") {
		t.Errorf("tpUsage 应包含注入的类型, got %q", got)
	}

	// 注入了但列表为空，同样退化
	typeLister = func() []string { return nil }
	if got := tpUsage(); got != "生成类型" {
		t.Errorf("空列表时 tpUsage = %q, want %q", got, "生成类型")
	}
}

func TestIsFlagLike(t *testing.T) {
	cases := map[string]bool{
		"-f":     true,
		"--long": true,
		"-3":     false,
		"a.go":   false,
		"":       false,
		"-":      false,
	}
	for in, want := range cases {
		if got := isFlagLike(in); got != want {
			t.Errorf("isFlagLike(%q) = %v, want %v", in, got, want)
		}
	}
}
