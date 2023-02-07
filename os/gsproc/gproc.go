package gsproc

import (
	"os"
	"time"
)

var (
	processPid       = os.Getpid() // processPid is the pid of current process.
	processStartTime = time.Now()  // processStartTime is the start time of current process.
)

// Pid returns the pid of current process.
func Pid() int {
	return processPid
}

//// PPid returns the custom parent pid if exists, or else it returns the system parent pid.
//func PPid() int {
//	if !IsChild() {
//		return Pid()
//	}
//	ppidValue := os.Getenv(envKeyPPid)
//	if ppidValue != "" && ppidValue != "0" {
//		return gconv.Int(ppidValue)
//	}
//	return PPidOS()
//}

// PPidOS returns the system parent pid of current process.
// Note that the difference between PPidOS and PPid function is that the PPidOS returns
// the system ppid, but the PPid functions may return the custom pid by gproc if the custom
// ppid exists.
func PPidOS() int {
	return os.Getppid()
}

// IsChild checks and returns whether current process is a child process.
// A child process is forked by another gproc process.
func IsChild() bool {
	ppidValue := os.Getenv(envKeyPPid)
	return ppidValue != "" && ppidValue != "0"
}

//// SetPPid sets custom parent pid for current process.
//func SetPPid(ppid int) error {
//	if ppid > 0 {
//		return os.Setenv(envKeyPPid, gconv.String(ppid))
//	} else {
//		return os.Unsetenv(envKeyPPid)
//	}
//}

// StartTime returns the start time of current process.
func StartTime() time.Time {
	return processStartTime
}

// Uptime returns the duration which current process has been running
func Uptime() time.Duration {
	return time.Now().Sub(processStartTime)
}
