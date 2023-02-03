package gslist

import (
	"container/list"
	"encoding/json"
	"github.com/jfy0o0/goStealer/internal/rwmutex"
)

type (
	// List is a doubly linked list containing a concurrent-safe/unsafe switch.
	// The switch should be set when its initialization and cannot be changed then.
	List[T any] struct {
		mu   rwmutex.RWMutex
		list *list.List
	}
	// Element the item type of the list.
	Element = list.Element
)

// New creates and returns a new empty doubly linked list.
func New[T any](safe ...bool) *List[T] {
	return &List[T]{
		mu:   rwmutex.Create(safe...),
		list: list.New(),
	}
}

// NewFrom creates and returns a list from a copy of given slice `array`.
// The parameter `safe` is used to specify whether using list in concurrent-safety,
// which is false in default.
func NewFrom[T any](array []T, safe ...bool) *List[T] {
	l := list.New()
	for _, v := range array {
		l.PushBack(v)
	}
	return &List[T]{
		mu:   rwmutex.Create(safe...),
		list: l,
	}
}

// PushFront inserts a new element `e` with value `v` at the front of list `l` and returns `e`.
func (l *List[T]) PushFront(v T) (e *Element) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.PushFront(v)
	l.mu.Unlock()
	return
}

// PushBack inserts a new element `e` with value `v` at the back of list `l` and returns `e`.
func (l *List[T]) PushBack(v T) (e *Element) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.PushBack(v)
	l.mu.Unlock()
	return
}

// PushFronts inserts multiple new elements with values `values` at the front of list `l`.
func (l *List[T]) PushFronts(values []T) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	for _, v := range values {
		l.list.PushFront(v)
	}
	l.mu.Unlock()
}

// PushBacks inserts multiple new elements with values `values` at the back of list `l`.
func (l *List[T]) PushBacks(values []T) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	for _, v := range values {
		l.list.PushBack(v)
	}
	l.mu.Unlock()
}

// PopBack removes the element from back of `l` and returns the value of the element.
func (l *List[T]) PopBack() (value T, ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return value, false
	}
	if e := l.list.Back(); e != nil {
		value = l.list.Remove(e).(T)
		ok = true
	}
	return
}

// PopFront removes the element from front of `l` and returns the value of the element.
func (l *List[T]) PopFront() (value T, ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return value, false
	}
	if e := l.list.Front(); e != nil {
		value = l.list.Remove(e).(T)
		ok = true
	}
	return
}

// PopBacks removes `max` elements from back of `l`
// and returns values of the removed elements as slice.
func (l *List[T]) PopBacks(max int) (values []T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	length := l.list.Len()
	if length > 0 {
		if max > 0 && max < length {
			length = max
		}
		values = make([]T, length)
		for i := 0; i < length; i++ {
			v := l.list.Remove(l.list.Back())
			values[i] = v.(T)
		}
	}
	return
}

// PopFronts removes `max` elements from front of `l`
// and returns values of the removed elements as slice.
func (l *List[T]) PopFronts(max int) (values []T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	length := l.list.Len()
	if length > 0 {
		if max > 0 && max < length {
			length = max
		}
		values = make([]T, length)
		for i := 0; i < length; i++ {
			v := l.list.Remove(l.list.Front())
			values[i] = v.(T)
		}
	}
	return
}

// PopBackAll removes all elements from back of `l`
// and returns values of the removed elements as slice.
func (l *List[T]) PopBackAll() []T {
	return l.PopBacks(-1)
}

// PopFrontAll removes all elements from front of `l`
// and returns values of the removed elements as slice.
func (l *List[T]) PopFrontAll() []T {
	return l.PopFronts(-1)
}

// FrontAll copies and returns values of all elements from front of `l` as slice.
func (l *List[T]) FrontAll() (values []T) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		values = make([]T, length)
		for i, e := 0, l.list.Front(); i < length; i, e = i+1, e.Next() {
			values[i] = e.Value.(T)
		}
	}
	return
}

// BackAll copies and returns values of all elements from back of `l` as slice.
func (l *List[T]) BackAll() (values []T) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		values = make([]T, length)
		for i, e := 0, l.list.Back(); i < length; i, e = i+1, e.Prev() {
			values[i] = e.Value.(T)
		}
	}
	return
}

// FrontValue returns value of the first element of `l` or nil if the list is empty.
func (l *List[T]) FrontValue() (value T) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	if e := l.list.Front(); e != nil {
		value = e.Value.(T)
	}
	return
}

// BackValue returns value of the last element of `l` or nil if the list is empty.
func (l *List[T]) BackValue() (value interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	if e := l.list.Back(); e != nil {
		value = e.Value
	}
	return
}

// Front returns the first element of list `l` or nil if the list is empty.
func (l *List[T]) Front() (e *Element) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	e = l.list.Front()
	return
}

// Back returns the last element of list `l` or nil if the list is empty.
func (l *List[T]) Back() (e *Element) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	e = l.list.Back()
	return
}

// Len returns the number of elements of list `l`.
// The complexity is O(1).
func (l *List[T]) Len() (length int) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length = l.list.Len()
	return
}

// Size is alias of Len.
func (l *List[T]) Size() int {
	return l.Len()
}

// MoveBefore moves element `e` to its new position before `p`.
// If `e` or `p` is not an element of `l`, or `e` == `p`, the list is not modified.
// The element and `p` must not be nil.
func (l *List[T]) MoveBefore(e, p *Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	l.list.MoveBefore(e, p)
}

// MoveAfter moves element `e` to its new position after `p`.
// If `e` or `p` is not an element of `l`, or `e` == `p`, the list is not modified.
// The element and `p` must not be nil.
func (l *List[T]) MoveAfter(e, p *Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	l.list.MoveAfter(e, p)
}

// MoveToFront moves element `e` to the front of list `l`.
// If `e` is not an element of `l`, the list is not modified.
// The element must not be nil.
func (l *List[T]) MoveToFront(e *Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	l.list.MoveToFront(e)
}

// MoveToBack moves element `e` to the back of list `l`.
// If `e` is not an element of `l`, the list is not modified.
// The element must not be nil.
func (l *List[T]) MoveToBack(e *Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	l.list.MoveToBack(e)
}

// PushBackList inserts a copy of an other list at the back of list `l`.
// The lists `l` and `other` may be the same, but they must not be nil.
func (l *List[T]) PushBackList(other *List[T]) {
	if l != other {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	l.list.PushBackList(other.list)
}

// PushFrontList inserts a copy of an other list at the front of list `l`.
// The lists `l` and `other` may be the same, but they must not be nil.
func (l *List[T]) PushFrontList(other *List[T]) {
	if l != other {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	l.list.PushFrontList(other.list)
}

// InsertAfter inserts a new element `e` with value `v` immediately after `p` and returns `e`.
// If `p` is not an element of `l`, the list is not modified.
// The `p` must not be nil.
func (l *List[T]) InsertAfter(p *Element, v T) (e *Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.InsertAfter(v, p)
	return
}

// InsertBefore inserts a new element `e` with value `v` immediately before `p` and returns `e`.
// If `p` is not an element of `l`, the list is not modified.
// The `p` must not be nil.
func (l *List[T]) InsertBefore(p *Element, v T) (e *Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.InsertBefore(v, p)
	return
}

// Remove removes `e` from `l` if `e` is an element of list `l`.
// It returns the element value e.Value.
// The element must not be nil.
func (l *List[T]) Remove(e *Element) (value T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	value = l.list.Remove(e).(T)
	return
}

// Removes removes multiple elements `es` from `l` if `es` are elements of list `l`.
func (l *List[T]) Removes(es []*Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	for _, e := range es {
		l.list.Remove(e)
	}
	return
}

// RemoveAll removes all elements from list `l`.
func (l *List[T]) RemoveAll() {
	l.mu.Lock()
	l.list = list.New()
	l.mu.Unlock()
}

// Clear is alias of RemoveAll.
func (l *List[T]) Clear() {
	l.RemoveAll()
}

// RLockFunc locks reading with given callback function `f` within RWMutex.RLock.
func (l *List[T]) RLockFunc(f func(list *list.List)) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list != nil {
		f(l.list)
	}
}

// LockFunc locks writing with given callback function `f` within RWMutex.Lock.
func (l *List[T]) LockFunc(f func(list *list.List)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	f(l.list)
}

// Iterator is alias of IteratorAsc.
func (l *List[T]) Iterator(f func(e *Element) bool) {
	l.IteratorAsc(f)
}

// IteratorAsc iterates the list readonly in ascending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (l *List[T]) IteratorAsc(f func(e *Element) bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		for i, e := 0, l.list.Front(); i < length; i, e = i+1, e.Next() {
			if !f(e) {
				break
			}
		}
	}
}

// IteratorDesc iterates the list readonly in descending order with given callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (l *List[T]) IteratorDesc(f func(e *Element) bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length := l.list.Len()
	if length > 0 {
		for i, e := 0, l.list.Back(); i < length; i, e = i+1, e.Prev() {
			if !f(e) {
				break
			}
		}
	}
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (l *List[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.FrontAll())
}
