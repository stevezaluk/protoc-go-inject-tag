package inject

import "regexp"

var (
	CommentRegex = regexp.MustCompile(`^//.*?@(?i:gotags?|inject_tags?):\s*(.*)$`)
	InjectRegex  = regexp.MustCompile("`.+`$")
	TagsRegex    = regexp.MustCompile(`[\w_]+:"[^"]+"`)
	AllRegex     = regexp.MustCompile(".*")
)

type TextArea struct {
	Start        int
	End          int
	CurrentTag   string
	InjectTag    string
	CommentStart int
	CommentEnd   int
}
