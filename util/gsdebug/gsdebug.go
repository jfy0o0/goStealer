package gsdebug

import (
	"fmt"
	"github.com/jfy0o0/goStealer/util/gstool"
	"os"
	"syscall"
	"time"
)

func RedirectStderr() {
	filename := fmt.Sprintf("%v_%v_%v.stderr", gstool.GetCurrentProcessName(), os.Getpid(), time.Now().UnixNano())
	logFile, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0666)
	syscall.Dup2(int(logFile.Fd()), 2)
}

func RedirectStdout() {
	filename := fmt.Sprintf("%v_%v_%v.stdout", gstool.GetCurrentProcessName(), os.Getpid(), time.Now().UnixNano())
	logFile, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0666)
	syscall.Dup2(int(logFile.Fd()), 1)
}
