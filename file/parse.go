package file

import (
	"github.com/stevezaluk/protoc-go-inject-tag/inject"
	"go/ast"
	"go/parser"
	"go/token"
)

/*
GenerateAST Take the file path and generate an AST representation of the file
*/
func GenerateAST(path string) (*ast.File, error) {
	fileSet := token.NewFileSet()

	fileAst, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return fileAst, err
	}

	return fileAst, nil
}

/*
GetStructFields Fetch fields for all struct declarations in a file
*/
func GetStructFields(astFile *ast.File) (ret []*ast.Field) {
	for _, decl := range astFile.Decls {
		generic, ok := decl.(*ast.GenDecl)
		if !ok {
			continue // skip if we fail to cast here
		}

		// iterate through each spec until we can successfully cast it
		// to a TypeSpec. This should be re-worked eventually to see if
		// there is an easier way to do this without relying on an O(n)
		// solution. Not sure if that is possible, so it stays for now
		var typeSpec *ast.TypeSpec
		for _, spec := range generic.Specs {
			if potentialSpec, ok := spec.(*ast.TypeSpec); ok {
				typeSpec = potentialSpec
				break
			}
		}

		if typeSpec == nil {
			continue
		}

		structDecl, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		ret = append(ret, structDecl.Fields.List...)
	}

	return ret
}

/*
ParseTextAreas Fetch a slice of struct fields from the AST file and generate text areas to be used for injection
*/
func ParseTextAreas(astFile *ast.File) (areas []*inject.TextArea, err error) {
	for _, field := range GetStructFields(astFile) {
		areas = append(areas, inject.NewTextAreas(field)...)
	}
	return
}
