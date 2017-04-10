package logger

import (
	"log"
	"os"
)

var Log *log.Logger

func FailOnError(err error, msg string) {
	if err != nil {
		Log.Fatalf("%s: %s", msg, err.Error())
	}
}

func LogOnError(err error, msg string) {
	if err != nil {
		Log.Printf("%s: %s", msg, err.Error())
	}
}

func CreateLog(file *os.File, FilePath string, Debug bool) {
	if !Debug {
		var err error
		file, err = os.OpenFile(FilePath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		FailOnError(err, "Cannot open log")

		Log = log.New(file, "", log.LstdFlags | log.Lshortfile)
	} else {
		Log = log.New(os.Stdout, "", log.LstdFlags | log.Ltime | log.Ldate)
	}
}