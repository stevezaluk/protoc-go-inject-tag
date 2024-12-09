package verbose

import (
	"github.com/spf13/viper"
	"log"
)

/*
Logf soon to be deprecated. Slog logging will replace this function
*/
func Logf(format string, v ...interface{}) {
	if !viper.GetBool("verbose") {
		return
	}

	log.Printf(format, v...)
}
