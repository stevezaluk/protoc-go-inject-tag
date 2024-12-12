package inject

import (
	"go/ast"
	"log/slog"
	"strings"
)

type TextArea struct {
	Start        int
	End          int
	CurrentTag   string
	InjectTag    string
	CommentStart int
	CommentEnd   int
}

/*
NewTextAreas Create a slice of text area pointers for a struct field
*/
func NewTextAreas(field *ast.Field) (areas []*TextArea) {
	var comments []*ast.Comment

	if field.Doc != nil {
		comments = append(comments, field.Doc.List...)
	}

	if field.Comment != nil {
		comments = append(comments, field.Comment.List...)
	}

	for _, comment := range comments {
		commentText := comment.Text
		tag := TagFromComment(commentText)
		if tag == "" {
			continue
		}

		// to be removed
		if strings.Contains(commentText, "inject_tag") {
			slog.Warn("warn: deprecated 'inject_tag' used")
		}

		area := TextArea{
			Start:        int(field.Pos()),
			End:          int(field.End()),
			CurrentTag:   field.Tag.Value[1 : len(field.Tag.Value)-1],
			InjectTag:    tag,
			CommentStart: int(comment.Pos()),
			CommentEnd:   int(comment.End()),
		}

		areas = append(areas, &area)
	}

	return areas
}
