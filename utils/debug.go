// +build !release

package utils

import (
	"fmt"
	"io"
	"log"
	"os"
)

var Debug = true

func Log(msg ...interface{}) {
	log.Output(2, fmt.Sprintln(msg...))
}

func Logf(format string, msg ...interface{}) {
	log.Output(2, fmt.Sprintf(format, msg...))
}

func GetLogWriter() io.Writer {
	return os.Stdout
}

func CloseLogWriter() {

}
