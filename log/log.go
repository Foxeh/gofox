package log

import (
	"io"
	"log"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {

	Info = log.New(infoHandle,
		"INFO - ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING - ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR - ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func ErrCheck(msg string, err error) {
	if err != nil {
		Error.Printf("%s: %+v", msg, err)
	}
}
