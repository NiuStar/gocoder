package gocoder

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"reflect"
)

type GoCoder struct {
	file			string
	f 				*ast.File
	fset			*token.FileSet
	genDecls 			map[token.Token][]*ast.GenDecl
	funcDecls 			[]*ast.FuncDecl
	//importDecl		*ast.GenDecl
}

func NewCoder(content string) (coder *GoCoder,err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return nil,err
	}
	//,funcDecls:make(map[string][]*ast.FuncDecl)
	coder = &GoCoder{f:f,genDecls:make(map[token.Token][]*ast.GenDecl),fset:fset}
	for _,decl := range coder.f.Decls {
		if reflect.TypeOf(decl).String() == "*ast.GenDecl" {
			declObj := decl.(*ast.GenDecl)
			coder.genDecls[declObj.Tok] = append(coder.genDecls[declObj.Tok],declObj)
		} else if "*ast.FuncDecl" == reflect.TypeOf(decl).String() {
			//fmt.Println(*decl.(*ast.FuncDecl))
			declObj := decl.(*ast.FuncDecl)
			//coder.funcDecls[declObj.Name.String()] = append(coder.funcDecls[declObj.Name.String()],declObj)
			coder.funcDecls = append(coder.funcDecls,declObj)
		}
		//coder.decls[reflect.TypeOf(decl).String()] = decl
	}
	return coder,nil
}
func NewCoderWtihFile(file string) (coder *GoCoder,err error) {
	body,_ := ioutil.ReadFile(file)
	coder,err = NewCoder(string(body))
	if err != nil {
		return nil,err
	}
	coder.file = file
	return coder,nil
}


func (coder *GoCoder)Merge(coder2 *GoCoder) {
	for _,decl := range coder2.f.Decls {
		coder.mergeDecl(decl)
	}
	//coder2.fset.
	/*
	for _,comment := range coder2.f.Comments {
		had := false
		for _,comment2 := range coder.f.Comments {
			if comment.Text() == comment2.Text() {
				had = true
				break
			}
		}
		if !had {
			coder.f.Comments = append(coder.f.Comments,comment)
		}
	}*/
}

func (coder *GoCoder)Save(path string) error {
	//coder.f.Comments = nil
	coder.f.Decls = []ast.Decl{}
	{
		for _,decl := range coder.genDecls[token.IMPORT] {
			coder.f.Decls = append(coder.f.Decls,decl)
			if decl.Doc != nil {
				//decl.
			}
		}
	}
	{
		for _,decl := range coder.genDecls[token.CONST] {
			coder.f.Decls = append(coder.f.Decls,decl)
		}
	}
	{
		for _,decl := range coder.genDecls[token.TYPE] {
			coder.f.Decls = append(coder.f.Decls,decl)
		}
	}
	for tok,decls := range coder.genDecls {
		if tok == token.IMPORT || tok == token.CONST || tok == token.TYPE {
			continue
		}
		for _,decl := range decls {
			coder.f.Decls = append(coder.f.Decls,decl)
		}
	}
	for _,decl := range coder.funcDecls {
		coder.f.Decls = append(coder.f.Decls,decl)
	}


	var buf bytes.Buffer
	//fset := token.NewFileSet() // use the wrong file set
	if err := printer.Fprint(&buf, coder.fset, coder.f);err != nil {
		return err
	}

	return ioutil.WriteFile(path,buf.Bytes(),0777)
}

func (coder *GoCoder)fieldFilter(name string, value reflect.Value) bool {
	fmt.Println(name,value)
	return false
}

func (coder *GoCoder)addComments() {

}

func (coder *GoCoder)Export() (error) {
	if len(coder.file) == 0 {
		return errors.New("导出失败，文件路径为空，请使用Save方法")
	}
	return coder.Save(coder.file)
}
