package inject

import (
	"fmt"
	"github.com/spf13/viper"
	"regexp"
)

func InitRegex() {
	viper.Set("tag.regex.comment", regexp.MustCompile(fmt.Sprintf(`^//.*?@(?i:%s?):\s*(.*)$`, viper.GetString("tag.comment-prefix"))))
	viper.Set("tag.regex.inject", regexp.MustCompile("`.+`$"))
	viper.Set("tag.regex.tags", regexp.MustCompile(`[\w_]+:"[^"]+"`))
	viper.Set("tag.regex.all", regexp.MustCompile(".*"))
}

func GetRegex(name string) *regexp.Regexp {
	regex := viper.Get(name)

	return regex.(*regexp.Regexp)
}
