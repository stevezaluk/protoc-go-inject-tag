package file

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stevezaluk/protoc-go-inject-tag/inject"
	"github.com/stevezaluk/protoc-go-inject-tag/verbose"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ParseFile(inputPath string, src interface{}, xxxSkip []string) (areas []inject.TextArea, err error) {
	verbose.Logf("parsing file %q for inject tag comments", inputPath)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, inputPath, src, parser.ParseComments)
	if err != nil {
		return
	}

	for _, decl := range f.Decls {
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
					verbose.Logf("warn: deprecated 'inject_tag' used")
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
	verbose.Logf("parsed file %q, number of fields to inject custom tags: %d", inputPath, len(areas))
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
		verbose.Logf("inject custom tag %q to expression %q", area.InjectTag, string(contents[area.Start-1:area.End-1]))
		contents = inject.InjectTag(contents, area, removeTagComment)
	}
	if err = os.WriteFile(inputPath, contents, 0o644); err != nil {
		return
	}

	if len(areas) > 0 {
		verbose.Logf("file %q is injected with custom tags", inputPath)
	}
	return
}

func IterFiles(inputPath string) {
	globs, err := filepath.Glob(inputPath)
	if err != nil {
		panic(err)
	}

	var fileCount int
	for _, path := range globs {
		finfo, err := os.Stat(path)
		if err != nil {
			panic(err)
		}

		if finfo.IsDir() {
			continue
		}

		if !strings.HasSuffix(strings.ToLower(finfo.Name()), ".go") {
			continue
		}

		fileCount++

		areas, err := ParseFile(path, nil, viper.GetStringSlice("tag.skip"))
		if err != nil {
			panic(err)
		}

		if err = WriteFile(path, areas, viper.GetBool("tag.remove-comments")); err != nil {
			panic(err)
		}
	}

	if fileCount == 0 {
		log.Fatalf("input %q matched no files; see -help", inputPath)
	}
}
