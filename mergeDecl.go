package gocoder

import (
	//"fmt"
	"go/ast"
	"reflect"
)

func (coder *GoCoder)mergeDecl(decl ast.Decl) {

	if reflect.TypeOf(decl).String() == "*ast.GenDecl" {
		declObj := decl.(*ast.GenDecl)
		if len(coder.genDecls[declObj.Tok]) > 0 {
			for _,spec := range declObj.Specs {
				had := false
				if reflect.TypeOf(spec).String() == "*ast.ImportSpec" {
					importSpec := spec.(*ast.ImportSpec)
					for _,decl2 := range coder.genDecls[declObj.Tok] {
						for _,spec2 := range decl2.Specs {
							importSpec2 := spec2.(*ast.ImportSpec)
							if importSpec.Path.Value == importSpec2.Path.Value {
								had = true
								break
							}
						}
						if had {
							break
						}
					}

				} else if reflect.TypeOf(spec).String() == "*ast.ValueSpec" {
					valueSpec := spec.(*ast.ValueSpec)
					for _,decl2 := range coder.genDecls[declObj.Tok] {
						for _,spec2 := range decl2.Specs {
							valueSpec2 := spec2.(*ast.ValueSpec)
							if valueSpec.Names[0].Name == valueSpec2.Names[0].Name {
								had = true
								break
							}
						}
						if had {
							break
						}
					}
					//coder.genDecls[declObj.Tok][0].Specs = append(coder.genDecls[declObj.Tok][0].Specs,valueSpec)
				} else if reflect.TypeOf(spec).String() == "*ast.TypeSpec" {
					typeSpec := spec.(*ast.TypeSpec)
					for _,decl2 := range coder.genDecls[declObj.Tok] {
						for _,spec2 := range decl2.Specs {
							typeSpec2 := spec2.(*ast.TypeSpec)
							if typeSpec.Name.Name == typeSpec2.Name.Name {
								had = true
								break
							}
						}
						if had {
							break
						}
					}
				}
				if !had {
					coder.genDecls[declObj.Tok][0].Specs = append(coder.genDecls[declObj.Tok][0].Specs,spec)
				}
			}
		} else {
			coder.genDecls[declObj.Tok] = append(coder.genDecls[declObj.Tok],declObj)
		}
	} else if "*ast.FuncDecl" == reflect.TypeOf(decl).String() {
		//fmt.Println(*decl.(*ast.FuncDecl))
		declObj := decl.(*ast.FuncDecl)
		for _,decl2 := range coder.funcDecls {
			if decl2.Recv == nil {
				if decl2.Name.Name == declObj.Name.Name {
					return
				}
			} else if decl2.Name.Name == declObj.Name.Name {
				if decl2.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name == declObj.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name {
					return
				}
			}
		}
		coder.funcDecls = append(coder.funcDecls,declObj)
	}
	return
}
