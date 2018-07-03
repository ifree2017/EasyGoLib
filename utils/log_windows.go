package utils

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

var (
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	procSetStdHandle = kernel32.MustFindProc("SetStdHandle")
)

func setStdHandle(stdhandle int32, handle syscall.Handle) error {
	r0, _, e1 := syscall.Syscall(procSetStdHandle.Addr(), 2, uintptr(stdhandle), uintptr(handle), 0)
	if r0 == 0 {
		if e1 != 0 {
			return error(e1)
		}
		return syscall.EINVAL
	}
	return nil
}

// RedirectStderr to the file passed in
func RedirectStderr() (err error) {
	logFilename := filepath.Join(LogDir(), strings.ToLower(EXEName())+"-error.log")
	logFile, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0644)
	if err != nil {
		return
	}
	err = setStdHandle(syscall.STD_ERROR_HANDLE, syscall.Handle(logFile.Fd()))
	if err != nil {
		return
	}
	// SetStdHandle does not affect prior references to stderr
	os.Stderr = logFile
	return
}
