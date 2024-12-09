package verbose

import (
	"log"
)

var verbose = false

func Logf(format string, v ...interface{}) {
	if !verbose {
		return
	}
	log.Printf(format, v...)
}
