package gstype

import (
	"encoding/json"
	"sync/atomic"
)

// Interface is a struct for concurrent-safe operation for type interface{}.
type AtomicValue[T any] struct {
	value atomic.Value
}

// NewInterface creates and returns a concurrent-safe object for interface{} type,
// with given initial value `value`.
func NewAtomicValue[T any](value ...T) *AtomicValue[T] {
	t := &AtomicValue[T]{}
	if len(value) > 0 && interface{}(value[0]) != nil {
		t.value.Store(value[0])
	}
	return t
}

// Clone clones and returns a new concurrent-safe object for interface{} type.
func (v *AtomicValue[T]) Clone() *AtomicValue[T] {
	return NewAtomicValue[T](v.Val())
}

// Set atomically stores `value` into t.value and returns the previous value of t.value.
// Note: The parameter `value` cannot be nil.
func (v *AtomicValue[T]) Set(value T) (old T) {
	old = v.Val()
	v.value.Store(value)
	return
}

// Val atomically loads and returns t.value.
func (v *AtomicValue[T]) Val() T {
	t := v.value.Load()
	t2, _ := t.(T)
	return t2
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (v AtomicValue[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Val())
}
