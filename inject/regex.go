package inject

import (
	"fmt"
	"github.com/spf13/viper"
	"regexp"
)

var (
	CommentRegex   = "tag.regex.comment"
	InjectionRegex = "tag.regex.inject"
	TagsRegex      = "tag.regex.tags"
	AllRegex       = "tag.regex.all"
)

/*
InitRegex Initialize regular expressions used for parsing comments pulled from the AST, and store them in viper
*/
func InitRegex() {
	viper.Set(CommentRegex, regexp.MustCompile(fmt.Sprintf(`^//.*?@(?i:%s?):\s*(.*)$`, viper.GetString("tag.comment-prefix"))))
	viper.Set(InjectionRegex, regexp.MustCompile("`.+`$"))
	viper.Set(TagsRegex, regexp.MustCompile(`[\w_]+:"[^"]+"`))
	viper.Set(AllRegex, regexp.MustCompile(".*"))
}

/*
GetRegex Fetch a value from viper and cast it as a Regexp pointer.
*/
func GetRegex(name string) *regexp.Regexp {
	regex := viper.Get(name)

	return regex.(*regexp.Regexp)
}
