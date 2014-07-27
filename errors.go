package helpers

import (
	"fmt"
)

func PanicToLogIf(err error, logger Log) {
	if err != nil {
		logger.Panic(err.Error())
	}
}

func PanicIf(err error) {
	if err != nil {
		NewConsoleLogger().Panic(err.Error())
	}
}

func Ap(something interface{}) {
	fmt.Printf("\n%#v\n", something)
}
