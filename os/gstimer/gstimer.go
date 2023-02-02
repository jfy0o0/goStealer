package gstimer

import "time"

const (
	StatusReady          = 0                      // Job or Timer is ready for running.
	StatusRunning        = 1                      // Job or Timer is already running.
	StatusStopped        = 2                      // Job or Timer is stopped.
	StatusClosed         = -1                     // Job or Timer is closed and waiting to be deleted.
	panicExit            = "exit"                 // panicExit is used for custom job exit with panic.
	defaultTimerInterval = 100 * time.Millisecond // defaultTimerInterval is the default timer interval in milliseconds.
)
