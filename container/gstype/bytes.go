package gstype

import (
	"encoding/base64"
	"sync/atomic"
)

// Bytes is a struct for concurrent-safe operation for type []byte.
type Bytes struct {
	value atomic.Value
}

// NewBytes creates and returns a concurrent-safe object for []byte type,
// with given initial value <value>.
func NewBytes(value ...[]byte) *Bytes {
	t := &Bytes{}
	if len(value) > 0 {
		t.value.Store(value[0])
	}
	return t
}

// Clone clones and returns a new concurrent-safe object for []byte type.
func (v *Bytes) Clone() *Bytes {
	return NewBytes(v.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
// Note: The parameter <value> cannot be nil.
func (v *Bytes) Set(value []byte) (old []byte) {
	old = v.Val()
	v.value.Store(value)
	return
}

// Val atomically loads and returns t.value.
func (v *Bytes) Val() []byte {
	if s := v.value.Load(); s != nil {
		return s.([]byte)
	}
	return nil
}

// String implements String interface for string printing.
func (v *Bytes) String() string {
	return string(v.Val())
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (v *Bytes) MarshalJSON() ([]byte, error) {
	val := v.Val()
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(val)))
	base64.StdEncoding.Encode(dst, val)
	return []byte(`"` + string(dst) + `"`), nil
}
