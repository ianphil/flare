package errors

import "github.com/iphilpot/flare/apis/logger"

// HandleError - prints and logs error
func HandleError(err error) {
	if err != nil {
		logger.Log.Fatalln(err)
	}
}
