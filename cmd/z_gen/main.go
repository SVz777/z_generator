package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/config"
	"github.com/SVz777/z_generator/gen/parser"
	_ "github.com/SVz777/z_generator/gen_impl/doc"
	_ "github.com/SVz777/z_generator/gen_impl/options"
	_ "github.com/SVz777/z_generator/gen_impl/php"
	_ "github.com/SVz777/z_generator/gen_impl/strct"
	_ "github.com/SVz777/z_generator/gen_impl/validate"
)

type Config struct {
	Typ        string
	FileName   string
	StructName string
	Output     string
	Inplace    bool
}

func main() {
	cfg := config.Config.Parse()
	if cfg.FileName == "" || cfg.StructName == "" {
		flag.Usage()
		os.Exit(-1)
	}
	p := parser.NewAST(cfg.FileName, cfg.StructName)
	g := gen.GetGenerator(cfg.Typ)
	s := "\n// generate by z_gen\n\n"
	s += string(g.Gen(p.Parser()))
	if cfg.Inplace {
		f, err := os.OpenFile(cfg.FileName, os.O_RDWR|os.O_APPEND, 0755)
		if err != nil {
			panic(err)
		}
		if _, err := fmt.Fprintln(f, s); err != nil {
			panic(err)
		}
	} else if cfg.Output != "" {
		f, err := os.OpenFile(cfg.Output, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
		if _, err := fmt.Fprintln(f, s); err != nil {
			panic(err)
		}
	} else {
		fmt.Println(s)
	}
}
