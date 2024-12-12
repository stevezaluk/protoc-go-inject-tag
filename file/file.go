package file

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/stevezaluk/protoc-go-inject-tag/inject"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
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

func WriteFile(inputPath string, areas []*inject.TextArea, removeTagComment bool) (err error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return
	}

	contents, err := io.ReadAll(f)
	if err != nil {
		return
	}

	if err = f.Close(); err != nil {
		return
	}

	// inject custom tags from tail of file first to preserve order
	for i := range areas {
		area := areas[len(areas)-i-1]
		slog.Debug("injected custom tag to expression", "tag", area.InjectTag, "expr", string(contents[area.Start-1:area.End-1]))
		contents = inject.InjectTag(contents, *area, removeTagComment)
	}
	if err = os.WriteFile(inputPath, contents, 0o644); err != nil {
		return
	}

	if len(areas) > 0 {
		slog.Debug("file is injected with custom tags", "file", inputPath)
	}
	return
}

/*
IsFileProtobuf Validates that the path is both a file and has the extension passed in tag.file-ext.
tag.file-ext defaults to ".pb.go" if one is not passed
*/
func IsFileProtobuf(path string) bool {
	if strings.HasSuffix(strings.ToLower(path), viper.GetString("tag.file-ext")) {
		return true
	}

	return false
}

/*
ProcessFile Converts the file passed in path to an AST and returns text areas to be injected
*/
func ProcessFile(path string) {
	slog.Debug("Generating AST for file", "file", path)
	astFile, err := GenerateAST(path)
	if err != nil {
		slog.Error("Failed to generate AST for file", "file", path, "err", err)
		return
	}

	slog.Debug("Parsing AST for file", "file", path)
	areas, err := ParseTextAreas(astFile)
	if err != nil {
		slog.Error("Error while parsing AST file", "file", path, "err", err)
		return
	}

	if err = WriteFile(path, areas, viper.GetBool("tag.remove-comments")); err != nil {
		slog.Error("Error while writing file to disk or injecting tags", "err", err)
		return
	}
}

/*
WalkFunc Handler function that gets called for each discovered file in WalkDir. If the path is a file, then it processes
it.
*/
func walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.IsDir() {
		if IsFileProtobuf(path) {
			slog.Debug("Processing file at path", "path", path)
			ProcessFile(path)
		}
	}

	return nil
}

/*
WalkDir Primary entrypoint for our application. Converts the UNIX path provided to an absolute path, processes the file
if the path is a single file, and recursively walks the path if it is a directory
*/
func WalkDir(path string) {
	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", os.Getenv("HOME"), -1)
		return
	}

	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		slog.Error("input path does not exist", "path", path)
		return
	}

	if !info.IsDir() { // path is a single file, proceed directly to processing
		if !IsFileProtobuf(path) {
			slog.Error("input does not match desired extension", "ext", viper.GetString("tag.file-ext"))
			return
		}

		ProcessFile(path)
		return
	}

	err = filepath.Walk(path, walkFunc) // walk the directory and look for files
	if err != nil {
		slog.Error("error while walking directory", "path", path)
	}
}
