package gsmap

import (
	"encoding/json"
	"github.com/jfy0o0/goStealer/internal/rwmutex"
)

type AnyAnyMap[K comparable, V any] struct {
	mu   rwmutex.RWMutex
	data map[K]V
}

// NewIntAnyMap returns an empty IntAnyMap object.
// The parameter `safe` is used to specify whether using map in concurrent-safety,
// which is false in default.
func NewAnyAnyMap[K comparable, V any](safe ...bool) *AnyAnyMap[K, V] {
	return &AnyAnyMap[K, V]{
		mu:   rwmutex.Create(safe...),
		data: make(map[K]V),
	}
}

// NewIntAnyMapFrom creates and returns a hash map from given map `data`.
// Note that, the param `data` map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
func NewAnyAnyMapFrom[K comparable, V any](data map[K]V, safe ...bool) *AnyAnyMap[K, V] {
	return &AnyAnyMap[K, V]{
		mu:   rwmutex.Create(safe...),
		data: data,
	}
}

// Iterator iterates the hash map readonly with custom callback function `f`.
// If `f` returns true, then it continues iterating; or false to stop.
func (m *AnyAnyMap[K, V]) Iterator(f func(k K, v V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

// Clone returns a new hash map with copy of current map data.
func (m *AnyAnyMap[K, V]) Clone() *AnyAnyMap[K, V] {
	return NewAnyAnyMapFrom[K, V](m.MapCopy(), m.mu.IsSafe())
}

// Map returns the underlying data map.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (m *AnyAnyMap[K, V]) Map() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.mu.IsSafe() {
		return m.data
	}
	data := make(map[K]V, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

// MapCopy returns a copy of the underlying data of the hash map.
func (m *AnyAnyMap[K, V]) MapCopy() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data := make(map[K]V, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

// Set sets key-value to the hash map.
func (m *AnyAnyMap[K, V]) Set(key K, val V) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[K]V)
	}
	m.data[key] = val
	m.mu.Unlock()
}

// Sets batch sets key-values to the hash map.
func (m *AnyAnyMap[K, V]) Sets(data map[K]V) {
	m.mu.Lock()
	if m.data == nil {
		m.data = data
	} else {
		for k, v := range data {
			m.data[k] = v
		}
	}
	m.mu.Unlock()
}

// Search searches the map with given `key`.
// Second return parameter `found` is true if key was found, otherwise false.
func (m *AnyAnyMap[K, V]) Search(key K) (value V, found bool) {
	m.mu.RLock()
	if m.data != nil {
		value, found = m.data[key]
	}
	m.mu.RUnlock()
	return
}

// Get returns the value by given `key`.
func (m *AnyAnyMap[K, V]) Get(key K) (value V) {
	m.mu.RLock()
	if m.data != nil {
		value, _ = m.data[key]
	}
	m.mu.RUnlock()
	return
}

// Pop retrieves and deletes an item from the map.
func (m *AnyAnyMap[K, V]) Pop() (key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for key, value = range m.data {
		delete(m.data, key)
		return
	}
	return
}

// Pops retrieves and deletes `size` items from the map.
// It returns all items if size == -1.
func (m *AnyAnyMap[K, V]) Pops(size int) map[K]V {
	m.mu.Lock()
	defer m.mu.Unlock()
	if size > len(m.data) || size == -1 {
		size = len(m.data)
	}
	if size == 0 {
		return nil
	}
	var (
		index  = 0
		newMap = make(map[K]V, size)
	)
	for k, v := range m.data {
		delete(m.data, k)
		newMap[k] = v
		index++
		if index == size {
			break
		}
	}
	return newMap
}

//// doSetWithLockCheck checks whether value of the key exists with mutex.Lock,
//// if not exists, set value to the map with given `key`,
//// or else just return the existing value.
////
//// When setting value, if `value` is type of `func() interface {}`,
//// it will be executed with mutex.Lock of the hash map,
//// and its return value will be set to the map with `key`.
////
//// It returns value with given `key`.
//func (m *AnyAnyMap[K,V]) doSetWithLockCheck(key int, value T) T {
//	m.mu.Lock()
//	defer m.mu.Unlock()
//	if m.data == nil {
//		m.data = make(map[int]T)
//	}
//	if v, ok := m.data[key]; ok {
//		return v
//	}
//	var v interface{} = value
//	if f, ok := v.(func() T); ok {
//		value = f()
//	}
//	if value != nil {
//		m.data[key] = value
//	}
//	return value
//}

//// GetOrSet returns the value by key,
//// or sets value with given `value` if it does not exist and then returns this value.
//func (m *AnyAnyMap[K,V]) GetOrSet(key int, value T) T {
//	if v, ok := m.Search(key); !ok {
//		return m.doSetWithLockCheck(key, value)
//	} else {
//		return v
//	}
//}
//
//// GetOrSetFunc returns the value by key,
//// or sets value with returned value of callback function `f` if it does not exist and returns this value.
//func (m *AnyAnyMap[K,V]) GetOrSetFunc(key int, f func() T) T {
//	if v, ok := m.Search(key); !ok {
//		return m.doSetWithLockCheck(key, f())
//	} else {
//		return v
//	}
//}

//// GetOrSetFuncLock returns the value by key,
//// or sets value with returned value of callback function `f` if it does not exist and returns this value.
////
//// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function `f`
//// with mutex.Lock of the hash map.
//func (m *AnyAnyMap[K,V]) GetOrSetFuncLock(key int, f func() T) T {
//	if v, ok := m.Search(key); !ok {
//		return m.doSetWithLockCheck(key, f())
//	} else {
//		return v
//	}
//}

//// SetIfNotExist sets `value` to the map if the `key` does not exist, and then returns true.
//// It returns false if `key` exists, and `value` would be ignored.
//func (m *AnyAnyMap[K,V]) SetIfNotExist(key int, value T) bool {
//	if !m.Contains(key) {
//		m.doSetWithLockCheck(key, value)
//		return true
//	}
//	return false
//}
//
//// SetIfNotExistFunc sets value with return value of callback function `f`, and then returns true.
//// It returns false if `key` exists, and `value` would be ignored.
//func (m *AnyAnyMap[K,V]) SetIfNotExistFunc(key int, f func() T) bool {
//	if !m.Contains(key) {
//		m.doSetWithLockCheck(key, f())
//		return true
//	}
//	return false
//}

//// SetIfNotExistFuncLock sets value with return value of callback function `f`, and then returns true.
//// It returns false if `key` exists, and `value` would be ignored.
////
//// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
//// it executes function `f` with mutex.Lock of the hash map.
//func (m *AnyAnyMap[K,V]) SetIfNotExistFuncLock(key int, f func() T) bool {
//	if !m.Contains(key) {
//		m.doSetWithLockCheck(key, f)
//		return true
//	}
//	return false
//}

// Removes batch deletes values of the map by keys.
func (m *AnyAnyMap[K, V]) Removes(keys []K) {
	m.mu.Lock()
	if m.data != nil {
		for _, key := range keys {
			delete(m.data, key)
		}
	}
	m.mu.Unlock()
}

// Remove deletes value from map by given `key`, and return this deleted value.
func (m *AnyAnyMap[K, V]) Remove(key K) (value V) {
	m.mu.Lock()
	if m.data != nil {
		var ok bool
		if value, ok = m.data[key]; ok {
			delete(m.data, key)
		}
	}
	m.mu.Unlock()
	return
}

// Keys returns all keys of the map as a slice.
func (m *AnyAnyMap[K, V]) Keys() []K {
	m.mu.RLock()
	var (
		keys  = make([]K, len(m.data))
		index = 0
	)
	for key := range m.data {
		keys[index] = key
		index++
	}
	m.mu.RUnlock()
	return keys
}

// Values returns all values of the map as a slice.
func (m *AnyAnyMap[K, V]) Values() []V {
	m.mu.RLock()
	var (
		values = make([]V, len(m.data))
		index  = 0
	)
	for _, value := range m.data {
		values[index] = value
		index++
	}
	m.mu.RUnlock()
	return values
}

// Contains checks whether a key exists.
// It returns true if the `key` exists, or else false.
func (m *AnyAnyMap[K, V]) Contains(key K) bool {
	var ok bool
	m.mu.RLock()
	if m.data != nil {
		_, ok = m.data[key]
	}
	m.mu.RUnlock()
	return ok
}

// Size returns the size of the map.
func (m *AnyAnyMap[K, V]) Size() int {
	m.mu.RLock()
	length := len(m.data)
	m.mu.RUnlock()
	return length
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (m *AnyAnyMap[K, V]) IsEmpty() bool {
	return m.Size() == 0
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (m *AnyAnyMap[K, V]) Clear() {
	m.mu.Lock()
	m.data = make(map[K]V)
	m.mu.Unlock()
}

// Replace the data of the map with given `data`.
func (m *AnyAnyMap[K, V]) Replace(data map[K]V) {
	m.mu.Lock()
	m.data = data
	m.mu.Unlock()
}

// LockFunc locks writing with given callback function `f` within RWMutex.Lock.
func (m *AnyAnyMap[K, V]) LockFunc(f func(m map[K]V)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f(m.data)
}

// RLockFunc locks reading with given callback function `f` within RWMutex.RLock.
func (m *AnyAnyMap[K, V]) RLockFunc(f func(m map[K]V)) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f(m.data)
}

// Merge merges two hash maps.
// The `other` map will be merged into the map `m`.
func (m *AnyAnyMap[K, V]) Merge(other *AnyAnyMap[K, V]) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = other.MapCopy()
		return
	}
	if other != m {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	for k, v := range other.data {
		m.data[k] = v
	}
}

// String returns the map as a string.
func (m *AnyAnyMap[K, V]) String() string {
	b, _ := m.MarshalJSON()
	return string(b)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (m AnyAnyMap[K, V]) MarshalJSON() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(m.data)
}

func (m *AnyAnyMap[K, V]) GetOrNew(key K, f func() V) (value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[K]V)
	}

	var ok bool
	value, ok = m.data[key]
	if ok {
		return
	}
	value = f()
	m.data[key] = value
	return
}
