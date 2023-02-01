package gsqueue

import (
	"github.com/jfy0o0/goStealer/container/gslist"
	"github.com/jfy0o0/goStealer/container/gstype"
	"math"
)

// Queue is a concurrent-safe queue built on doubly linked list and channel.
type Queue[T any] struct {
	limit  int                       // Limit for queue size.
	list   *gslist.List[T]           // Underlying list structure for data maintaining.
	closed *gstype.AtomicValue[bool] // Whether queue is closed.
	events chan struct{}             // Events for data writing.
	C      chan T                    // Underlying channel for data reading.
}

const (
	defaultQueueSize = 10000 // Size for queue buffer.
	defaultBatchSize = 10    // Max batch size per-fetching from list.
)

// New returns an empty queue object.
// Optional parameter <limit> is used to limit the size of the queue, which is unlimited in default.
// When <limit> is given, the queue will be static and high performance which is comparable with stdlib channel.
func New[T any](limit ...int) *Queue[T] {
	q := &Queue[T]{
		closed: gstype.NewAtomicValue[bool](),
	}
	if len(limit) > 0 && limit[0] > 0 {
		q.limit = limit[0]
		q.C = make(chan T, limit[0])
	} else {
		q.list = gslist.New[T](true)
		q.events = make(chan struct{}, math.MaxInt32)
		q.C = make(chan T, defaultQueueSize)
		go q.asyncLoopFromListToChannel()
	}
	return q
}

// asyncLoopFromListToChannel starts an asynchronous goroutine,
// which handles the data synchronization from list <q.list> to channel <q.C>.
func (q *Queue[T]) asyncLoopFromListToChannel() {
	defer func() {
		if q.closed.Val() {
			_ = recover()
		}
	}()
	for !q.closed.Val() {
		<-q.events
		for !q.closed.Val() {
			if length := q.list.Len(); length > 0 {
				if length > defaultBatchSize {
					length = defaultBatchSize
				}
				for _, v := range q.list.PopFronts(length) {
					// When q.C is closed, it will panic here, especially q.C is being blocked for writing.
					// If any error occurs here, it will be caught by recover and be ignored.
					q.C <- v
				}
			} else {
				break
			}
		}
		// Clear q.events to remain just one event to do the next synchronization check.
		for i := 0; i < len(q.events)-1; i++ {
			<-q.events
		}
	}
	// It should be here to close q.C if <q> is unlimited size.
	// It's the sender's responsibility to close channel when it should be closed.
	close(q.C)
}

// Push pushes the data <v> into the queue.
// Note that it would panics if Push is called after the queue is closed.
func (q *Queue[T]) Push(v T) {
	if q.limit > 0 {
		q.C <- v
	} else {
		q.list.PushBack(v)
		if len(q.events) < defaultQueueSize {
			q.events <- struct{}{}
		}
	}
}

// Pop pops an item from the queue in FIFO way.
// Note that it would return nil immediately if Pop is called after the queue is closed.
func (q *Queue[T]) Pop() T {
	return <-q.C
}

// Close closes the queue.
// Notice: It would notify all goroutines return immediately,
// which are being blocked reading using Pop method.
func (q *Queue[T]) Close() {
	q.closed.Set(true)
	if q.events != nil {
		close(q.events)
	}
	if q.limit > 0 {
		close(q.C)
	}
	for i := 0; i < defaultBatchSize; i++ {
		q.Pop()
	}
}

// Len returns the length of the queue.
// Note that the result might not be accurate as there's a
// asynchronous channel reading the list constantly.
func (q *Queue[T]) Len() (length int) {
	if q.list != nil {
		length += q.list.Len()
	}
	length += len(q.C)
	return
}

// Size is alias of Len.
func (q *Queue[T]) Size() int {
	return q.Len()
}
