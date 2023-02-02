package gstimer

// DefaultOptions creates and returns a default options object for Timer creation.
func DefaultOptions() TimerOptions {
	return TimerOptions{
		Interval: defaultTimerInterval,
	}
}
