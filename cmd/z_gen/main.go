package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/SVz777/z_generator/gen"
	"github.com/SVz777/z_generator/gen/parser"
	"github.com/SVz777/z_generator/gen_impl/doc"
	"github.com/SVz777/z_generator/gen_impl/options"
	"github.com/SVz777/z_generator/gen_impl/php"
	"github.com/SVz777/z_generator/gen_impl/strct"
)

type Config struct {
	Typ        string
	FileName   string
	StructName string
	Inplace    bool
}

func main() {
	config := &Config{}
	flag.StringVar(&config.FileName, "filename", "", "filename")
	flag.StringVar(&config.FileName, "f", "", "filename")
	flag.StringVar(&config.StructName, "structname", "", "structname")
	flag.StringVar(&config.StructName, "s", "", "structname")
	flag.BoolVar(&config.Inplace, "inplace", false, "inplace")

	flag.StringVar(&config.Typ, "tp", "", "gen type")

	flag.Parse()
	fmt.Println(config)
	if config.FileName == "" || config.StructName == "" {
		flag.Usage()
		os.Exit(-1)
	}
	p := parser.NewAST(config.FileName, config.StructName)
	g := NewGenerator(config.Typ)
	s := "\n// generate by svz_generator\n\n"
	s += string(g.Gen(p.Parser()))
	if config.Inplace {
		f, err := os.OpenFile(config.FileName, os.O_RDWR|os.O_APPEND, 0755)
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

func NewGenerator(genType string) gen.IGen {
	switch genType {
	case doc.GenType:
		return doc.NewGenerator()
	case options.GenType:
		return options.NewGenerator()
	case php.GenType:
		return php.NewGenerator()
	case strct.GenTypeToMap:
		return strct.NewToStringMapGenerator()
	case strct.GenTypeKeys:
		return strct.NewKeysGenerator()
	default:
		return strct.NewWrapperGenerator()
	}
}
