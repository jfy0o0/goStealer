package gslru

import (
	"fmt"
	"github.com/jfy0o0/goStealer/container/gslist"
	"github.com/jfy0o0/goStealer/container/gsmap"
	"github.com/jfy0o0/goStealer/container/gstype"
	"sync"
)

// LRU cache object.
// It uses list.List from stdlib for its underlying doubly linked list.
type Lru[T comparable] struct {
	mutex   sync.RWMutex
	data    *gsmap.AnyAnyMap[T, *gslist.Element] // Key mapping to the item of the list.
	list    *gslist.List[T]                      // Key list.
	closed  *gstype.Bool                         // Closed or not.
	adapter LruAdapter[T]
	cap     int
	size    int
}

// newMemCacheLru creates and returns a new LRU object.
func New[T comparable](cap int, adapter LruAdapter[T]) *Lru[T] {
	lru := &Lru[T]{
		data:    gsmap.NewAnyAnyMap[T, *gslist.Element](false),
		list:    gslist.New[T](false),
		closed:  gstype.NewBool(),
		cap:     cap,
		adapter: adapter,
	}
	return lru
}

func (lru *Lru[T]) Cap() int {
	return lru.cap
}

// Close closes the LRU object.
func (lru *Lru[T]) Close() {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()
	lru.closed.Set(true)
}

// Remove deletes the `key` FROM `lru`.
func (lru *Lru[T]) Remove(key T) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()
	if !lru.closed.Val() {
		return
	}
	if v := lru.data.Get(key); v != nil {
		lru.size--
		lru.data.Remove(key)
		lru.list.Remove(v)
	}
}

// Size returns the size of `lru`.
func (lru *Lru[T]) Size() int {
	lru.mutex.RLock()
	defer lru.mutex.RUnlock()
	return lru.size
}

// Push pushes `key` to the tail of `lru`.
func (lru *Lru[T]) Push(key T) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()
	if !lru.closed.Val() {
		return
	}
	if v := lru.data.Get(key); v != nil {
		lru.list.Remove(v)
		lru.data.Set(key, lru.list.PushFront(key))
		lru.adapter.PushRaise(key)
		return
	}

	if lru.size < lru.cap {
		lru.data.Set(key, lru.list.PushFront(key))
		lru.size++
		lru.adapter.PushNormal(key)
		return
	}

	v, ok := lru.list.PopBack()
	if !ok {
		return
	}

	lru.data.Remove(v)
	lru.data.Set(key, lru.list.PushFront(key))
	lru.adapter.PushRaiseWithPop(key)
}

// Pop deletes and returns the key from tail of `lru`.
func (lru *Lru[T]) Pop() (t T, ok bool) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()
	if !lru.closed.Val() {
		return
	}
	if v, ok := lru.list.PopBack(); ok {
		lru.data.Remove(v)
		return v, true
	}
	return t, false
}

func (lru *Lru[T]) Dump() {
	fmt.Println(lru.data.String())
}
