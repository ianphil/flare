package logger

import (
	"log"
	"os"
)

var (
	// Log - Used to log to file
	Log *log.Logger
)

func init() {
	var logpath = "./info.log"
	var file, err1 = os.Create(logpath)
	if err1 != nil {
		panic(err1)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
	Log.Println("LogFile : " + logpath)
}
