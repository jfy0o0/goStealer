package gstimer

import (
	"github.com/jfy0o0/goStealer/container/gspriority_queue"
	"github.com/jfy0o0/goStealer/container/gstype"
	"sync"
	"time"
)

// Timer is the timer manager, which uses ticks to calculate the timing interval.
type Timer struct {
	mu      sync.RWMutex
	queue   *gspriority_queue.PriorityQueue[*Entry] // queue is a priority queue based on heap structure.
	status  *gstype.Int                             // status is the current timer status.
	ticks   *gstype.Int64                           // ticks is the proceeded interval number by the timer.
	options TimerOptions                            // timer options is used for timer configuration.
}

// TimerOptions is the configuration object for Timer.
type TimerOptions struct {
	Interval time.Duration // Interval is the interval escaped of the timer.
}

func New(options ...TimerOptions) *Timer {
	t := &Timer{
		queue:  gspriority_queue.New[*Entry](false),
		status: gstype.NewInt(StatusRunning),
		ticks:  gstype.NewInt64(),
	}
	if len(options) > 0 {
		t.options = options[0]
	} else {
		t.options = DefaultOptions()
	}
	go t.loop()
	return t
}

// Add adds a timing job to the timer, which runs in interval of <interval>.
func (t *Timer) Add(interval time.Duration, job JobFunc) *Entry {
	return t.createEntry(interval, job, false, -1, StatusReady)
}

// AddEntry adds a timing job to the timer with detailed parameters.
//
// The parameter <interval> specifies the running interval of the job.
//
// The parameter <singleton> specifies whether the job running in singleton mode.
// There's only one of the same job is allowed running when its a singleton mode job.
//
// The parameter <times> specifies limit for the job running times, which means the job
// exits if its run times exceeds the <times>.
//
// The parameter <status> specifies the job status when it's firstly added to the timer.
func (t *Timer) AddEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
	return t.createEntry(interval, job, singleton, times, status)
}

// AddSingleton is a convenience function for add singleton mode job.
func (t *Timer) AddSingleton(interval time.Duration, job JobFunc) *Entry {
	return t.createEntry(interval, job, true, -1, StatusReady)
}

// AddOnce is a convenience function for adding a job which only runs once and then exits.
func (t *Timer) AddOnce(interval time.Duration, job JobFunc) *Entry {
	return t.createEntry(interval, job, true, 1, StatusReady)
}

// AddTimes is a convenience function for adding a job which is limited running times.
func (t *Timer) AddTimes(interval time.Duration, times int, job JobFunc) *Entry {
	return t.createEntry(interval, job, true, times, StatusReady)
}

// DelayAdd adds a timing job after delay of <interval> duration.
// Also see Add.
func (t *Timer) DelayAdd(delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(delay, func() {
		t.Add(interval, job)
	})
}

// DelayAddEntry adds a timing job after delay of <interval> duration.
// Also see AddEntry.
func (t *Timer) DelayAddEntry(delay time.Duration, interval time.Duration, job JobFunc, singleton bool, times int, status int) {
	t.AddOnce(delay, func() {
		t.AddEntry(interval, job, singleton, times, status)
	})
}

// DelayAddSingleton adds a timing job after delay of <interval> duration.
// Also see AddSingleton.
func (t *Timer) DelayAddSingleton(delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(delay, func() {
		t.AddSingleton(interval, job)
	})
}

// DelayAddOnce adds a timing job after delay of <interval> duration.
// Also see AddOnce.
func (t *Timer) DelayAddOnce(delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(delay, func() {
		t.AddOnce(interval, job)
	})
}

// DelayAddTimes adds a timing job after delay of <interval> duration.
// Also see AddTimes.
func (t *Timer) DelayAddTimes(delay time.Duration, interval time.Duration, times int, job JobFunc) {
	t.AddOnce(delay, func() {
		t.AddTimes(interval, times, job)
	})
}

// Start starts the timer.
func (t *Timer) Start() {
	t.status.Set(StatusRunning)
}

// Stop stops the timer.
func (t *Timer) Stop() {
	t.status.Set(StatusStopped)
}

// Close closes the timer.
func (t *Timer) Close() {
	t.status.Set(StatusClosed)
}

// loop starts the ticker using a standalone goroutine.
func (t *Timer) loop() {
	go func() {
		var (
			currentTimerTicks   int64
			timerIntervalTicker = time.NewTicker(t.options.Interval)
		)
		defer timerIntervalTicker.Stop()
		for {
			select {
			case <-timerIntervalTicker.C:
				// Check the timer status.
				switch t.status.Val() {
				case StatusRunning:
					// Timer proceeding.
					if currentTimerTicks = t.ticks.Add(1); currentTimerTicks >= t.queue.NextPriority() {
						t.proceed(currentTimerTicks)
					}

				case StatusStopped:
					// Do nothing.

				case StatusClosed:
					// Timer exits.
					return
				}
			}
		}
	}()
}

// proceed function proceeds the timer job checking and running logic.
func (t *Timer) proceed(currentTimerTicks int64) {
	var (
		entry *Entry
		ok    bool
	)
	for {
		entry, ok = t.queue.Pop()
		if !ok {
			break
		}
		// It checks if it meets the ticks' requirement.
		if jobNextTicks := entry.nextTicks.Val(); currentTimerTicks < jobNextTicks {
			// It pushes the job back if current ticks does not meet its running ticks requirement.
			t.queue.Push(entry, entry.nextTicks.Val())
			break
		}
		// It checks the job running requirements and then does asynchronous running.
		entry.doCheckAndRunByTicks(currentTimerTicks)
		// Status check: push back or ignore it.
		if entry.Status() != StatusClosed {
			// It pushes the job back to queue for next running.
			t.queue.Push(entry, entry.nextTicks.Val())
		}
	}
}

// createEntry creates and adds a timing job to the timer.
func (t *Timer) createEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
	var (
		infinite = false
	)
	if times <= 0 {
		infinite = true
	}
	var (
		intervalTicksOfJob = int64(interval / t.options.Interval)
	)
	if intervalTicksOfJob == 0 {
		// If the given interval is lesser than the one of the wheel,
		// then sets it to one tick, which means it will be run in one interval.
		intervalTicksOfJob = 1
	}
	var (
		nextTicks = t.ticks.Val() + intervalTicksOfJob
		entry     = &Entry{
			job:       job,
			timer:     t,
			ticks:     intervalTicksOfJob,
			times:     gstype.NewInt(times),
			status:    gstype.NewInt(status),
			singleton: gstype.NewBool(singleton),
			nextTicks: gstype.NewInt64(nextTicks),
			infinite:  gstype.NewBool(infinite),
		}
	)
	t.queue.Push(entry, nextTicks)
	return entry
}
