package gstimer

// DefaultOptions creates and returns a default options object for Timer creation.
func DefaultOptions() TimerOptions {
	return TimerOptions{
		Interval: defaultTimerInterval,
	}
}

// Exit is used in timing job internally, which exits and marks it closed from timer.
// The timing job will be automatically removed from timer later. It uses "panic-recover"
// mechanism internally implementing this feature, which is designed for simplification
// and convenience.
func Exit() {
	panic(panicExit)
}
