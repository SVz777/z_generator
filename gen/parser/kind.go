package parser

import (
	"go/ast"
	"go/parser"
	"go/types"
	"reflect"
)

func kindOfType(typ types.Type) reflect.Kind {
	if typ == nil {
		return reflect.Invalid
	}
	if _, ok := typ.(*types.TypeParam); ok {
		return reflect.Invalid
	}
	switch t := typ.Underlying().(type) {
	case *types.Interface:
		return reflect.Interface
	case *types.Struct:
		return reflect.Struct
	case *types.Pointer:
		return reflect.Ptr
	case *types.Slice:
		return reflect.Slice
	case *types.Map:
		return reflect.Map
	case *types.Chan:
		return reflect.Chan
	case *types.Signature:
		return reflect.Func
	case *types.Array:
		return reflect.Array
	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			return reflect.Bool
		case types.Int:
			return reflect.Int
		case types.Int8:
			return reflect.Int8
		case types.Int16:
			return reflect.Int16
		case types.Int32:
			return reflect.Int32
		case types.Int64:
			return reflect.Int64
		case types.Uint:
			return reflect.Uint
		case types.Uint8:
			return reflect.Uint8
		case types.Uint16:
			return reflect.Uint16
		case types.Uint32:
			return reflect.Uint32
		case types.Uint64:
			return reflect.Uint64
		case types.Uintptr:
			return reflect.Uintptr
		case types.Float32:
			return reflect.Float32
		case types.Float64:
			return reflect.Float64
		case types.Complex64:
			return reflect.Complex64
		case types.Complex128:
			return reflect.Complex128
		case types.String:
			return reflect.String
		case types.UnsafePointer:
			return reflect.UnsafePointer
		}
	}
	return reflect.Invalid
}

// KindOfTypeName only resolves built-in types and type literals, not named types.
func KindOfTypeName(name string) reflect.Kind {
	expr, err := parser.ParseExpr(name)
	if err != nil {
		return reflect.Invalid
	}
	return newKindResolver(nil).resolve(expr)
}

type kindResolver struct {
	decls  map[string]ast.Expr
	cache  map[string]reflect.Kind
	active map[string]bool
}

func newKindResolver(file *ast.File) *kindResolver {
	r := &kindResolver{decls: map[string]ast.Expr{}, cache: map[string]reflect.Kind{}, active: map[string]bool{}}
	if file != nil {
		for _, decl := range file.Decls {
			if gen, ok := decl.(*ast.GenDecl); ok {
				for _, spec := range gen.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						r.decls[ts.Name.Name] = ts.Type
					}
				}
			}
		}
	}
	return r
}

func (r *kindResolver) resolve(expr ast.Expr) reflect.Kind {
	switch t := expr.(type) {
	case *ast.ParenExpr:
		return r.resolve(t.X)
	case *ast.StarExpr:
		return reflect.Ptr
	case *ast.ArrayType:
		if t.Len == nil {
			return reflect.Slice
		}
		return reflect.Array
	case *ast.MapType:
		return reflect.Map
	case *ast.ChanType:
		return reflect.Chan
	case *ast.FuncType:
		return reflect.Func
	case *ast.InterfaceType:
		return reflect.Interface
	case *ast.StructType:
		return reflect.Struct
	case *ast.Ident:
		if kind, ok := r.cache[t.Name]; ok {
			return kind
		}
		if decl, ok := r.decls[t.Name]; ok {
			if r.active[t.Name] {
				return reflect.Invalid
			}
			r.active[t.Name] = true
			kind := r.resolve(decl)
			delete(r.active, t.Name)
			r.cache[t.Name] = kind
			return kind
		}
		if obj := types.Universe.Lookup(t.Name); obj != nil {
			if _, ok := obj.(*types.TypeName); ok {
				return kindOfType(obj.Type())
			}
		}
	}
	return reflect.Invalid
}
