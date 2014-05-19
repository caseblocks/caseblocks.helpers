package helpers

import (
	"log"
	"os"
)

func PanicIf(err error, logger *log.Logger) {
	if err != nil {
		logger.Panicln(err)
	}
}

func Logger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags)
}
