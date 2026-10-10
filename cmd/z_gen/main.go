package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/args"
	"github.com/SVz777/z_generator/gen/parser"
	_ "github.com/SVz777/z_generator/gen_impl/doc"
	_ "github.com/SVz777/z_generator/gen_impl/options"
	_ "github.com/SVz777/z_generator/gen_impl/php"
	_ "github.com/SVz777/z_generator/gen_impl/strct"
	_ "github.com/SVz777/z_generator/gen_impl/validate"
)

func main() {
	// 让 -tp 的帮助信息列出所有已注册的生成器类型
	args.SetTypeLister(gen.ListTypes)

	cfg := args.Args.Parse()
	if cfg.FileName == "" || cfg.StructName == "" {
		flag.Usage()
		os.Exit(1)
	}
	g := gen.GetGenerator(cfg.Typ)
	// 类型不存在时给出可用列表，否则下面会直接 nil panic
	if g == nil {
		fmt.Fprintf(os.Stderr, "未知的生成器类型: %q\n可用类型: %s\n",
			cfg.Typ, strings.Join(gen.ListTypes(), ", "))
		os.Exit(1)
	}
	st, err := parser.StructFromFile(cfg.FileName, cfg.StructName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	code, err := g.Gen(st)
	if err != nil {
		fmt.Fprintf(os.Stderr, "生成失败: %v\n", err)
		os.Exit(1)
	}
	s := "\n// generate by z_gen\n\n"
	s += string(code)
	switch {
	case cfg.Inplace:
		err = writeFile(cfg.FileName, os.O_RDWR|os.O_APPEND, 0755, s)
	case cfg.Output != "":
		err = writeFile(cfg.Output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666, s)
	default:
		fmt.Println(s)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// writeFile 把生成结果写入文件。
func writeFile(path string, flag int, perm os.FileMode, s string) error {
	f, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return fmt.Errorf("打开 %s 失败: %w", path, err)
	}
	defer f.Close()
	if _, err := fmt.Fprintln(f, s); err != nil {
		return fmt.Errorf("写入 %s 失败: %w", path, err)
	}
	return nil
}
