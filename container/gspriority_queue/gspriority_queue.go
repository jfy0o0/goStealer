package gspriority_queue

import (
	"container/heap"
	"github.com/jfy0o0/goStealer/container/gstype"
	"math"
	"sync"
)

// priorityQueue is an abstract data type similar to a regular queue or stack data structure in which
// each element additionally has a "priority" associated with it. In a priority queue, an element with
// high priority is served before an element with low priority.
// priorityQueue is based on heap structure.
type priorityQueue[T any] struct {
	mu           sync.Mutex
	heap         *priorityQueueHeap[T]      // the underlying queue items manager using heap.
	nextPriority *gstype.AtomicValue[int64] // nextPriority stores the next priority value of the heap, which is used to check if necessary to call the Pop of heap by Timer.
}

// newPriorityQueue creates and returns a priority queue.
func New[T any](isBig bool) *priorityQueue[T] {
	queue := &priorityQueue[T]{
		heap:         newPriorityQueueHeap[T](isBig),
		nextPriority: gstype.NewAtomicValue[int64](math.MaxInt64),
	}
	heap.Init(queue.heap)
	return queue
}

// NextPriority retrieves and returns the minimum and the most priority value of the queue.
func (q *priorityQueue[T]) NextPriority() int64 {
	return q.nextPriority.Val()
}

// Push pushes a value to the queue.
// The `priority` specifies the priority of the value.
// The lesser the `priority` value the higher priority of the `value`.
func (q *priorityQueue[T]) Push(value T, priority int64) {
	q.mu.Lock()
	defer q.mu.Unlock()
	heap.Push(q.heap, priorityQueueItem[T]{
		value:    value,
		priority: priority,
	})
	// Update the minimum priority using atomic operation.
	nextPriority := q.nextPriority.Val()
	if priority >= nextPriority {
		return
	}
	q.nextPriority.Set(priority)
}

// Pop retrieves, removes and returns the most high priority value from the queue.
func (q *priorityQueue[T]) Pop() any {
	q.mu.Lock()
	defer q.mu.Unlock()
	if v := heap.Pop(q.heap); v != nil {
		var nextPriority int64 = math.MaxInt64
		if len(q.heap.array) > 0 {
			nextPriority = q.heap.array[0].priority
		}
		q.nextPriority.Set(nextPriority)
		return v.(priorityQueueItem[T]).value
	}
	return nil
}
