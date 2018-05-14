// +build !release

package utils

import (
	"io"
	"log"
	"os"
)

var Debug = true

func Log(msg ...interface{}) {
	log.Println(msg...)
}

func GetLogWriter() io.Writer {
	return os.Stdout
}

func CloseLogWriter() {

}
