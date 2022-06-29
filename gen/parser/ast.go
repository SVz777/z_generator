package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"reflect"

	"github.com/SVz777/stringchange"
)

type AST struct {
	fileNode   *ast.File
	structName string
}

func NewAST(filename string, structNames string) *AST {
	filenode := parseFile(filename)
	return &AST{
		fileNode:   filenode,
		structName: structNames,
	}
}

func (t *AST) Parser() (out *Struct) {
	out = &Struct{}

	ast.Inspect(t.fileNode, func(node ast.Node) bool {
		switch tp := node.(type) {
		//case ast.Expr:
		case *ast.TypeSpec:
			if structTp, ok := tp.Type.(*ast.StructType); ok {
				if t.structName == tp.Name.String() {
					out.Type = t.structName
					out.Receiver = stringchange.ToCamelCase(t.structName)

					for _, field := range structTp.Fields.List {
						if field.Names == nil {
							//匿名结构体不处理
							continue
						}

						ff := &Field{}
						ff.Type = t.getTypeName(field.Type)
						ff.Name = field.Names[0].String()
						if field.Tag == nil {
							ff.Tag = reflect.StructTag("")
						} else {
							ff.Tag = reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
						}
						out.Field = append(out.Field, ff)
					}
				}
			}
			//case ast.Stmt:
			//case ast.Decl:
		}
		return true
	})
	return
}

func (t *AST) getTypeName(expr ast.Expr) string {
	switch node := expr.(type) {
	case *ast.Ident:
		//非指针
		return node.Name
	case *ast.StarExpr:
		return "*" + t.getTypeName(node.X)
	case *ast.SelectorExpr:
		return t.getTypeName(node.X) + "." + t.getTypeName(node.Sel)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.ArrayType:
		//数组
		return "[]" + t.getTypeName(node.Elt)
	case *ast.MapType:
		//map
		return "map[" + t.getTypeName(node.Key) + "]" + t.getTypeName(node.Value)
	}
	return ""
}

func parseFile(filename string) *ast.File {
	fs := &token.FileSet{}
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	fileNode, err := parser.ParseFile(fs, filename, src, parser.ParseComments)
	if err != nil {
		return nil
	}
	return fileNode
}
