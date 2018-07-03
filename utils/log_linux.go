package utils

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// RedirectStderr to the file passed in
func RedirectStderr() (err error) {
	logFilename := filepath.Join(LogDir(), strings.ToLower(EXEName())+"-error.log")
	logFile, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0644)
	if err != nil {
		return
	}
	err = syscall.Dup2(int(logFile.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		return
	}
	return
}
