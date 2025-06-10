package config

import (
	"flag"
)

type config struct {
	Typ        string
	FileName   string
	StructName string
	Output     string
	Inplace    bool
}

func (c *config) Parse() *config {
	if flag.Parsed() {
		return c
	}
	flag.StringVar(&c.FileName, "filename", "", "filename")
	flag.StringVar(&c.FileName, "f", "", "filename")
	flag.StringVar(&c.StructName, "structname", "", "structname")
	flag.StringVar(&c.StructName, "s", "", "structname")
	flag.StringVar(&c.Output, "output", "", "output")
	flag.StringVar(&c.Output, "o", "", "output")
	flag.BoolVar(&c.Inplace, "inplace", false, "inplace")

	flag.StringVar(&c.Typ, "tp", "", "gen type")

	flag.Parse()
	return c
}

var Config config
