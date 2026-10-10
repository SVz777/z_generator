package parser

import (
	"go/ast"
	goparser "go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

type namedInterface interface{ Read() }
type interfaceAlias = namedInterface
type namedSlice []int
type namedStruct struct{ Value int }
type kindFixture struct {
	Interface namedInterface
	Alias     interfaceAlias
	Slice     namedSlice
	Struct    namedStruct
	Pointer   *namedStruct
	Number    int
}

func TestFieldKindsRuntimeAndLocalAST(t *testing.T) {
	file, err := ParseFile("kind_test.go")
	if err != nil {
		t.Fatal(err)
	}
	source, err := StructFrom(file, "kindFixture")
	if err != nil {
		t.Fatal(err)
	}
	runtime := StructOf(kindFixture{})
	for i, field := range source.Field {
		if field.Kind != runtime.Field[i].Kind {
			t.Errorf("%s: source %s, runtime %s", field.Name, field.Kind, runtime.Field[i].Kind)
		}
	}
}

func TestFieldsOfTypesInfo(t *testing.T) {
	fset := token.NewFileSet()
	file, err := goparser.ParseFile(fset, "fixture.go", "package fixture; type I interface{ Read() }; type S struct { A, B I }", 0)
	if err != nil {
		t.Fatal(err)
	}
	info := &types.Info{Types: map[ast.Expr]types.TypeAndValue{}}
	conf := types.Config{}
	if _, err := conf.Check("fixture", fset, []*ast.File{file}, info); err != nil {
		t.Fatal(err)
	}
	st := file.Decls[1].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType)
	fields, err := FieldsOf(st, "s", info)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 2 {
		t.Fatalf("got %d fields", len(fields))
	}
	for _, field := range fields {
		if field.Kind != reflect.Interface {
			t.Errorf("%s: %s", field.Name, field.Kind)
		}
	}
}

func TestFileFieldKinds(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"go.mod":               "module fixture\n\ngo 1.18\n",
		"types.go":             "package fixture; type Reader interface{ Read() }; type Alias = Reader; type Value struct{ N int }",
		"fixture.go":           "package fixture; import stream \"io\"; import \"fixture/dependency\"; type S struct { Local Reader; Alias Alias; Imported stream.Reader; Value Value; Dependency dependency.Reader }",
		"excluded.go":          "//go:build never_selected_kind_test\n\npackage fixture; type Reader struct{}",
		"dependency/reader.go": "package dependency; type Reader interface{ Read() }", "fixture_test.go": "package fixture; import \"io\"; type TestS struct { Local Reader; Imported io.Reader }",
		"external_test.go": "package fixture_test; import \"io\"; type ExternalS struct { Imported io.Reader }",
	}
	for name, content := range files {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	for _, tc := range []struct {
		file, name string
		kinds      []reflect.Kind
	}{
		{"fixture.go", "S", []reflect.Kind{reflect.Interface, reflect.Interface, reflect.Interface, reflect.Struct, reflect.Interface}},
		{"fixture_test.go", "TestS", []reflect.Kind{reflect.Interface, reflect.Interface}},
		{"external_test.go", "ExternalS", []reflect.Kind{reflect.Interface}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			st, err := StructFromFile(filepath.Join(dir, tc.file), tc.name)
			if err != nil {
				t.Fatal(err)
			}
			for i, field := range st.Field {
				if field.Kind != tc.kinds[i] {
					t.Errorf("%s: got %s, want %s", field.Name, field.Kind, tc.kinds[i])
				}
			}
		})
	}
}

func TestPackageASTLoading(t *testing.T) {
	cases := []struct {
		name    string
		source  string
		wantErr bool
		kinds   []reflect.Kind
	}{
		{"semantic", "package fixture; type Box[T any] struct{ Value T }; type S struct { Box Box[int] `json:\"box\"`; Number int }", false, []reflect.Kind{reflect.Struct, reflect.Int}},
		{"partial", "package fixture; type S struct { Missing Unknown; Number int }", false, []reflect.Kind{reflect.Invalid, reflect.Int}},
		{"syntax", "package fixture; type S struct {", true, nil},
		{"excluded", "//go:build never_selected_kind_test\n\npackage fixture; type S struct{}", true, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			for name, content := range map[string]string{"go.mod": "module fixture\n\ngo 1.23.0\n", "fixture.go": tc.source} {
				if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0600); err != nil {
					t.Fatal(err)
				}
			}
			filename := filepath.Join(dir, "fixture.go")
			st, err := StructFromFile(filename, "S")
			if (err != nil) != tc.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantErr {
				return
			}
			if len(st.Field) != len(tc.kinds) {
				t.Fatalf("got %d fields, want %d", len(st.Field), len(tc.kinds))
			}
			for i, field := range st.Field {
				if field.Kind != tc.kinds[i] {
					t.Errorf("%s: got %s, want %s", field.Name, field.Kind, tc.kinds[i])
				}
			}
			if tc.name == "semantic" {
				if st.Field[0].Tag.Get("json") != "box" {
					t.Fatal("field tag was lost")
				}
				file, info, err := loadFileTypes(filename)
				if err != nil {
					t.Fatal(err)
				}
				loaded, err := StructFrom(file, "S", info)
				if err != nil || loaded.Field[0].Kind != reflect.Struct {
					t.Fatalf("AST and type information do not match: %v", err)
				}
				if _, err := StructFromFile(filename, "Missing"); err == nil {
					t.Fatal("missing struct should return an error")
				}
			}
		})
	}
}

func TestStandaloneFileTypes(t *testing.T) {
	t.Setenv("GOWORK", "off")
	t.Setenv("GO111MODULE", "on")
	for _, tc := range []struct {
		name    string
		source  string
		wantErr bool
	}{
		{"valid", "package main; import stream \"io\"; type Reader = stream.Reader; type S struct { Local Reader; Imported stream.Reader; Number int `json:\"number\"` }; func (*S) Read() {}", false},
		{"syntax", "package main; type S struct {", true},
		{"excluded", "//go:build never_selected_kind_test\n\npackage main; type S struct{}", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			filename := filepath.Join(dir, "t.go")
			if err := os.WriteFile(filename, []byte(tc.source), 0600); err != nil {
				t.Fatal(err)
			}
			cwd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Chdir(dir); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := os.Chdir(cwd); err != nil {
					t.Error(err)
				}
			})
			for _, path := range []string{filename, "t.go"} {
				st, err := StructFromFile(path, "S")
				if (err != nil) != tc.wantErr {
					t.Fatalf("%s: unexpected error: %v", path, err)
				}
				if tc.wantErr {
					continue
				}
				if len(st.Field) != 3 {
					t.Fatalf("got %d fields, want 3", len(st.Field))
				}
				for i, kind := range []reflect.Kind{reflect.Interface, reflect.Interface, reflect.Int} {
					if st.Field[i].Kind != kind {
						t.Errorf("%s: got %s, want %s", st.Field[i].Name, st.Field[i].Kind, kind)
					}
				}
				if st.Field[2].Tag.Get("json") != "number" || len(st.Func) != 1 || st.Func[0].Name != "Read" {
					t.Fatal("field tag or method was lost")
				}
			}
		})
	}
}

func TestSourceFileIdentity(t *testing.T) {
	dir := t.TempDir()
	for name, content := range map[string]string{
		"go.mod":    "module identity\n\ngo 1.23.0\n",
		"types.go":  "package identity; type Reader interface{ Read() }",
		"source.go": "package identity; type S struct { Local Reader }",
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	filename := filepath.Join(dir, "source.go")
	t.Run("different", func(t *testing.T) {
		if sameSourceFile(filename, filepath.Join(dir, "types.go")) || sameSourceFile(filename, filepath.Join(dir, "missing.go")) {
			t.Fatal("different files must not match")
		}
	})
	t.Run("windows case", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("Windows path casing")
		}
		path := strings.ToUpper(filename[:1]) + filename[1:]
		if path == filename {
			path = strings.ToLower(filename[:1]) + filename[1:]
		}
		if !sameSourceFile(path, filename) {
			t.Fatal("same file with different drive casing must match")
		}
		st, err := StructFromFile(path, "S")
		if err != nil || st.Field[0].Kind != reflect.Interface {
			t.Fatalf("cross-file type resolution was lost: %v", err)
		}
	})
	t.Run("symlink", func(t *testing.T) {
		path := filepath.Join(dir, "linked.go")
		if err := os.Symlink(filename, path); err != nil {
			t.Skipf("symlink unavailable: %v", err)
		}
		if !sameSourceFile(path, filename) {
			t.Fatal("symlink must match its target")
		}
	})
}

func TestStandaloneBuildConfiguration(t *testing.T) {
	t.Setenv("GOWORK", "off")
	t.Setenv("GO111MODULE", "on")
	for _, tc := range []struct {
		name, filename, source, flags, cgo string
		wantErr                            bool
	}{
		{"internal test", "t_test.go", "package fixture; import \"io\"; type S struct { Reader io.Reader }", "", "0", false},
		{"external test", "t_test.go", "package fixture_test; import \"io\"; type S struct { Reader io.Reader }", "", "0", false},
		{"tag enabled", "t.go", "//go:build custom_fixture_tag\n\npackage fixture; import \"io\"; type S struct { Reader io.Reader }", "-tags=custom_fixture_tag", "0", false},
		{"tag disabled", "t.go", "//go:build custom_fixture_tag\n\npackage fixture; type S struct{}", "", "0", true},
		{"cgo disabled", "t.go", "package fixture; import \"C\"; type S struct{}", "", "0", true},
		{"other platform", "t_plan9.go", "package fixture; type S struct{}", "", "0", true},
		{"syntax diagnostics", "t.go", "package fixture; type S struct {", "", "0", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("GOFLAGS", tc.flags)
			t.Setenv("CGO_ENABLED", tc.cgo)
			filename := filepath.Join(t.TempDir(), tc.filename)
			if err := os.WriteFile(filename, []byte(tc.source), 0600); err != nil {
				t.Fatal(err)
			}
			st, err := StructFromFile(filename, "S")
			if (err != nil) != tc.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.name == "syntax diagnostics" {
				joined, ok := err.(interface{ Unwrap() []error })
				if !ok || len(joined.Unwrap()) != 2 || !strings.Contains(joined.Unwrap()[0].Error(), filename) || !strings.Contains(joined.Unwrap()[1].Error(), "解析源文件") {
					t.Fatalf("package or standalone diagnostics were lost: %v", err)
				}
			}
			if !tc.wantErr && (len(st.Field) != 1 || st.Field[0].Kind != reflect.Interface) {
				t.Fatalf("imported type resolution was lost: %+v", st.Field)
			}
		})
	}
}

func TestTagLiterals(t *testing.T) {
	for _, value := range []string{"`json:\"name,omitempty\"`", strconv.Quote(`json:"name,omitempty"`)} {
		file, err := goparser.ParseFile(token.NewFileSet(), "tags.go", "package fixture; type S struct { Name string "+value+" }", 0)
		if err != nil {
			t.Fatal(err)
		}
		st, err := StructFrom(file, "S")
		if err != nil || st.Field[0].Tag.Get("json") != "name,omitempty" {
			t.Fatalf("tag decoding failed for %s: %v", value, err)
		}
	}
	for _, lit := range []*ast.BasicLit{nil, {Kind: token.STRING, Value: "bad"}, {Kind: token.INT, Value: "1"}} {
		if TagOf(lit) != "" {
			t.Fatal("invalid tag must be empty")
		}
	}
}

func TestUnknownAndCyclicKinds(t *testing.T) {
	fset := token.NewFileSet()
	file, err := goparser.ParseFile(fset, "fixture.go", "package fixture; type A = B; type B = A; type S struct{ Cycle A; Missing Other; Channel chan int }", 0)
	if err != nil {
		t.Fatal(err)
	}
	st, err := StructFrom(file, "S")
	if err != nil {
		t.Fatal(err)
	}
	if st.Field[0].Kind != reflect.Invalid || st.Field[1].Kind != reflect.Invalid || st.Field[2].Kind != reflect.Chan {
		t.Fatalf("unexpected field kinds: %+v", st.Field)
	}
	for _, name := range []string{"channel", "function", "interfaceName", "Reader"} {
		if got := KindOfTypeName(name); got != reflect.Invalid {
			t.Errorf("%s: %s", name, got)
		}
	}
}
