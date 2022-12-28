package gstype

import (
	"encoding/json"
	"sync/atomic"
)

// Interface is a struct for concurrent-safe operation for type interface{}.
type Interface[T any] struct {
	value atomic.Value
}

// NewInterface creates and returns a concurrent-safe object for interface{} type,
// with given initial value `value`.
func NewInterface[T any](value ...T) *Interface[T] {
	t := &Interface[T]{}
	if len(value) > 0 && value[0] != nil {
		t.value.Store(value[0])
	}
	return t
}

// Clone clones and returns a new concurrent-safe object for interface{} type.
func (v *Interface[T]) Clone() *Interface[T] {
	return NewInterface[T](v.Val())
}

// Set atomically stores `value` into t.value and returns the previous value of t.value.
// Note: The parameter `value` cannot be nil.
func (v *Interface[T]) Set(value T) (old T) {
	old = v.Val()
	v.value.Store(value)
	return
}

// Val atomically loads and returns t.value.
func (v *Interface[T]) Val() T {
	return v.value.Load()
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (v Interface[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Val())
}
