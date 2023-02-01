package gspriority_queue

// priorityQueueItem stores the queue item which has a `priority` attribute to sort itself in heap.
type priorityQueueItem[T any] struct {
	value    T
	priority int64
}

// priorityQueueHeap is a heap manager, of which the underlying `array` is a array implementing a heap structure.
type priorityQueueHeap[T any] struct {
	array []priorityQueueItem[T]
	isBig bool
}

func newPriorityQueueHeap[T any](isBig bool) *priorityQueueHeap[T] {
	return &priorityQueueHeap[T]{
		array: make([]priorityQueueItem[T], 0),
		isBig: isBig,
	}
}

// Len is used to implement the interface of sort.Interface.
func (h *priorityQueueHeap[T]) Len() int {
	return len(h.array)
}

// Less is used to implement the interface of sort.Interface.
// The least one is placed to the top of the heap.
func (h *priorityQueueHeap[T]) Less(i, j int) bool {
	if h.isBig {
		return h.array[i].priority > h.array[j].priority
	}
	return h.array[i].priority < h.array[j].priority
}

// Swap is used to implement the interface of sort.Interface.
func (h *priorityQueueHeap[T]) Swap(i, j int) {
	if len(h.array) == 0 {
		return
	}
	h.array[i], h.array[j] = h.array[j], h.array[i]
}

// Push pushes an item to the heap.
func (h *priorityQueueHeap[T]) Push(x any) {
	h.array = append(h.array, (x).(priorityQueueItem[T]))
}

// Pop retrieves, removes and returns the most high priority item from the heap.
func (h *priorityQueueHeap[T]) Pop() any {
	length := len(h.array)
	if length == 0 {
		return nil
	}
	item := h.array[length-1]
	h.array = h.array[0 : length-1]
	return item
}
