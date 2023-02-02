package gstype

import (
	"encoding/json"
	"strconv"
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

func (v AtomicValue[T]) String() string {
	t := v.value.Load()
	switch t.(type) {
	case bool:
		t2, _ := t.(bool)
		if t2 {
			return "true"
		} else {
			return "false"
		}
	case byte:
		t2, _ := t.(byte)
		return strconv.FormatUint(uint64(t2), 10)
	case []byte:
		t2, _ := t.([]byte)
		return string(t2)
	case float32:
		t2, _ := t.(float32)
		return strconv.FormatFloat(float64(t2), 'g', -1, 32)
	case float64:
		t2, _ := t.(float64)
		return strconv.FormatFloat(t2, 'g', -1, 64)
	case int:
		t2, _ := t.(int)
		return strconv.Itoa(t2)
	case int32:
		t2, _ := t.(int32)
		return strconv.Itoa(int(t2))
	case int64:
		t2, _ := t.(int64)
		return strconv.FormatInt(t2, 10)
	case string:
		t2, _ := t.(string)
		return t2
	case uint:
		t2, _ := t.(uint)
		return strconv.FormatUint(uint64(t2), 10)
	case uint32:
		t2, _ := t.(uint32)
		return strconv.FormatUint(uint64(t2), 10)
	case uint64:
		t2, _ := t.(uint64)
		return strconv.FormatUint(t2, 10)
	}
	return "nil"
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (v AtomicValue[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Val())
}
