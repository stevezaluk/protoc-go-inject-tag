package inject

import (
	"fmt"
	"github.com/spf13/viper"
	"regexp"
)

/*
InitRegex Initialize regular expressions used for parsing comments pulled from the AST, and store them in viper
*/
func InitRegex() {
	viper.Set("tag.regex.comment", regexp.MustCompile(fmt.Sprintf(`^//.*?@(?i:%s?):\s*(.*)$`, viper.GetString("tag.comment-prefix"))))
	viper.Set("tag.regex.inject", regexp.MustCompile("`.+`$"))
	viper.Set("tag.regex.tags", regexp.MustCompile(`[\w_]+:"[^"]+"`))
	viper.Set("tag.regex.all", regexp.MustCompile(".*"))
}

/*
GetRegex Fetch a value from viper and cast it as a Regexp pointer.
*/
func GetRegex(name string) *regexp.Regexp {
	regex := viper.Get(name)

	return regex.(*regexp.Regexp)
}
