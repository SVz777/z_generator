/**
 * @file    main.go
 * @author
 *  ___  _  _  ____
 * / __)( \/ )(_   )
 * \__ \ \  /  / /_
 * (___/  \/  (____)
 * (903943711@qq.com)
 * @date    2020-01-06
 * @desc
 */
package main

import (
	"flag"
	"fmt"
	"go/format"
	"os"

	"github.com/SVz777/z_generator/gen_core/parser"
	"github.com/SVz777/z_generator/gen_core/strct"
)

type Config struct {
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
	flag.Parse()
	if config.FileName == "" || config.StructName == "" {
		flag.Usage()
		os.Exit(-1)
	}
	p := parser.NewAST(config.FileName, config.StructName)
	g := strct.NewStructWrapper(p.Parser())
	bSource, _ := format.Source(g.Gen().Bytes())
	s := "// generate by svz_generator\n"
	s += string(bSource)
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
