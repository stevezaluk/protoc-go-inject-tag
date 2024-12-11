package file

import (
	"errors"
	"fmt"
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

func ParseAST(astFile *ast.File, xxxSkip []string) (areas []inject.TextArea, err error) {
	for _, decl := range astFile.Decls {
		// check if is generic declaration
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		var typeSpec *ast.TypeSpec
		for _, spec := range genDecl.Specs {
			if ts, tsOK := spec.(*ast.TypeSpec); tsOK {
				typeSpec = ts
				break
			}
		}

		// skip if can't get type spec
		if typeSpec == nil {
			continue
		}

		// not a struct, skip
		structDecl, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		builder := strings.Builder{}
		if len(xxxSkip) > 0 {
			for i, skip := range xxxSkip {
				builder.WriteString(fmt.Sprintf("%s:\"-\"", skip))
				if i > 0 {
					builder.WriteString(",")
				}
			}
		}

		for _, field := range structDecl.Fields.List {
			// skip if field has no doc
			if len(field.Names) > 0 {
				name := field.Names[0].Name
				if len(xxxSkip) > 0 && strings.HasPrefix(name, "XXX") {
					currentTag := field.Tag.Value
					area := inject.TextArea{
						Start:      int(field.Pos()),
						End:        int(field.End()),
						CurrentTag: currentTag[1 : len(currentTag)-1],
						InjectTag:  builder.String(),
					}
					areas = append(areas, area)
				}
			}

			comments := []*ast.Comment{}

			if field.Doc != nil {
				comments = append(comments, field.Doc.List...)
			}

			// The "doc" field (above comment) is more commonly "free-form"
			// due to the ability to have a much larger comment without it
			// being unwieldy. As such, the "comment" field (trailing comment),
			// should take precedence if there happen to be multiple tags
			// specified, both in the field doc, and the field line. Whichever
			// comes last, will take precedence.
			if field.Comment != nil {
				comments = append(comments, field.Comment.List...)
			}

			for _, comment := range comments {
				tag := inject.TagFromComment(comment.Text)
				if tag == "" {
					continue
				}

				if strings.Contains(comment.Text, "inject_tag") {
					slog.Warn("warn: deprecated 'inject_tag' used")
				}

				currentTag := field.Tag.Value
				area := inject.TextArea{
					Start:        int(field.Pos()),
					End:          int(field.End()),
					CurrentTag:   currentTag[1 : len(currentTag)-1],
					InjectTag:    tag,
					CommentStart: int(comment.Pos()),
					CommentEnd:   int(comment.End()),
				}
				areas = append(areas, area)
			}
		}
	}
	return
}

func WriteFile(inputPath string, areas []inject.TextArea, removeTagComment bool) (err error) {
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
		contents = inject.InjectTag(contents, area, removeTagComment)
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
	areas, err := ParseAST(astFile, viper.GetStringSlice("tag.skip"))
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
