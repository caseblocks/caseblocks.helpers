package helpers

import (
	"fmt"

	"github.com/martini-contrib/render"
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

func RespondWithError(err error, statusCode int, r render.Render) {
	var msg string
	switch statusCode {
	case 422:
		msg = "Bad request."
	case 500:
		msg = "Unable to process request."
	}
	errorMsg := fmt.Sprintf("%s %s", msg, err)
	fmt.Println(errorMsg)
	r.JSON(statusCode, map[string]interface{}{"error": errorMsg})
}
