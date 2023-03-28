package gsarray

import (
	"encoding/json"
	"github.com/jfy0o0/goStealer/errors/gscode"
	"github.com/jfy0o0/goStealer/errors/gserror"
	"github.com/jfy0o0/goStealer/internal/rwmutex"
	"math"
	"sort"
)

// Array is a golang array with rich features.
// It contains a concurrent-safe/unsafe switch, which should be set
// when its initialization and cannot be changed then.
type Array[T any] struct {
	mu    rwmutex.RWMutex
	array []T
}

// New creates and returns an empty array.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func New[T any](safe ...bool) *Array[T] {
	return NewArraySize[T](0, 0, safe...)
}

// NewArray is alias of New, please see New.
func NewArray[T any](safe ...bool) *Array[T] {
	return NewArraySize[T](0, 0, safe...)
}

// NewArraySize create and returns an array with given size and cap.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewArraySize[T any](size int, cap int, safe ...bool) *Array[T] {
	return &Array[T]{
		mu:    rwmutex.Create(safe...),
		array: make([]T, size, cap),
	}
}

// NewFrom is alias of NewArrayFrom.
// See NewArrayFrom.
func NewFrom[T int](array []T, safe ...bool) *Array[T] {
	return NewArrayFrom(array, safe...)
}

// NewFromCopy is alias of NewArrayFromCopy.
// See NewArrayFromCopy.
func NewFromCopy[T any](array []T, safe ...bool) *Array[T] {
	return NewArrayFromCopy(array, safe...)
}

// NewArrayFrom creates and returns an array with given slice `array`.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewArrayFrom[T any](array []T, safe ...bool) *Array[T] {
	return &Array[T]{
		mu:    rwmutex.Create(safe...),
		array: array,
	}
}

// NewArrayFromCopy creates and returns an array from a copy of given slice `array`.
// The parameter `safe` is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewArrayFromCopy[T any](array []T, safe ...bool) *Array[T] {
	newArray := make([]T, len(array))
	copy(newArray, array)
	return &Array[T]{
		mu:    rwmutex.Create(safe...),
		array: newArray,
	}
}

// At returns the value by the specified index.
// If the given `index` is out of range of the array, it returns `nil`.
func (a *Array[T]) At(index int) (value T) {
	value, _ = a.Get(index)
	return
}

// Get returns the value by the specified index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *Array[T]) Get(index int) (value T, found bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if index < 0 || index >= len(a.array) {
		return value, false
	}
	return a.array[index], true
}

// Set sets value to specified index.
func (a *Array[T]) Set(index int, value T) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return gserror.NewCodef(gscode.CodeInvalidParameter, "index %d out of array range %d", index, len(a.array))
	}
	a.array[index] = value
	return nil
}

// SetArray sets the underlying slice array with the given `array`.
func (a *Array[T]) SetArray(array []T) *Array[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.array = array
	return a
}

// Replace replaces the array items by given `array` from the beginning of array.
func (a *Array[T]) Replace(array []T) *Array[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	max := len(array)
	if max > len(a.array) {
		max = len(a.array)
	}
	for i := 0; i < max; i++ {
		a.array[i] = array[i]
	}
	return a
}

// SortFunc sorts the array by custom function `less`.
func (a *Array[T]) SortFunc(less func(v1, v2 T) bool) *Array[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	sort.Slice(a.array, func(i, j int) bool {
		return less(a.array[i], a.array[j])
	})
	return a
}

// InsertBefore inserts the `value` to the front of `index`.
func (a *Array[T]) InsertBefore(index int, value T) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return gserror.NewCodef(gscode.CodeInvalidParameter, "index %d out of array range %d", index, len(a.array))
	}
	rear := append([]T{}, a.array[index:]...)
	a.array = append(a.array[0:index], value)
	a.array = append(a.array, rear...)
	return nil
}

// InsertAfter inserts the `value` to the back of `index`.
func (a *Array[T]) InsertAfter(index int, value T) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if index < 0 || index >= len(a.array) {
		return gserror.NewCodef(gscode.CodeInvalidParameter, "index %d out of array range %d", index, len(a.array))
	}
	rear := append([]T{}, a.array[index+1:]...)
	a.array = append(a.array[0:index+1], value)
	a.array = append(a.array, rear...)
	return nil
}

// Remove removes an item by index.
// If the given `index` is out of range of the array, the `found` is false.
func (a *Array[T]) Remove(index int) (value T, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.doRemoveWithoutLock(index)
}

// doRemoveWithoutLock removes an item by index without lock.
func (a *Array[T]) doRemoveWithoutLock(index int) (value T, found bool) {
	if index < 0 || index >= len(a.array) {
		return value, false
	}
	// Determine array boundaries when deleting to improve deletion efficiency.
	if index == 0 {
		value := a.array[0]
		a.array = a.array[1:]
		return value, true
	} else if index == len(a.array)-1 {
		value := a.array[index]
		a.array = a.array[:index]
		return value, true
	}
	// If it is a non-boundary delete,
	// it will involve the creation of an array,
	// then the deletion is less efficient.
	value = a.array[index]
	a.array = append(a.array[:index], a.array[index+1:]...)
	return value, true
}

//// RemoveValue removes an item by value.
//// It returns true if value is found in the array, or else false if not found.
//func (a *Array[T]) RemoveValue(value T) bool {
//	if i := a.Search(value); i != -1 {
//		a.Remove(i)
//		return true
//	}
//	return false
//}

// PushLeft pushes one or multiple items to the beginning of array.
func (a *Array[T]) PushLeft(value ...T) *Array[T] {
	a.mu.Lock()
	a.array = append(value, a.array...)
	a.mu.Unlock()
	return a
}

// PushRight pushes one or multiple items to the end of array.
// It equals to Append.
func (a *Array[T]) PushRight(value ...T) *Array[T] {
	a.mu.Lock()
	a.array = append(a.array, value...)
	a.mu.Unlock()
	return a
}

//// PopRand randomly pops and return an item out of array.
//// Note that if the array is empty, the `found` is false.
//func (a *Array[T]) PopRand() (value T, found bool) {
//	a.mu.Lock()
//	defer a.mu.Unlock()
//	return a.doRemoveWithoutLock(grand.Intn(len(a.array)))
//}

//// PopRands randomly pops and returns `size` items out of array.
//func (a *Array[T]) PopRands(size int) []T {
//	a.mu.Lock()
//	defer a.mu.Unlock()
//	if size <= 0 || len(a.array) == 0 {
//		return nil
//	}
//	if size >= len(a.array) {
//		size = len(a.array)
//	}
//	array := make([]interface{}, size)
//	for i := 0; i < size; i++ {
//		array[i], _ = a.doRemoveWithoutLock(grand.Intn(len(a.array)))
//	}
//	return array
//}

// PopLeft pops and returns an item from the beginning of array.
// Note that if the array is empty, the `found` is false.
func (a *Array[T]) PopLeft() (value T, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.array) == 0 {
		return value, false
	}
	value = a.array[0]
	a.array = a.array[1:]
	return value, true
}

// PopRight pops and returns an item from the end of array.
// Note that if the array is empty, the `found` is false.
func (a *Array[T]) PopRight() (value T, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	index := len(a.array) - 1
	if index < 0 {
		return value, false
	}
	value = a.array[index]
	a.array = a.array[:index]
	return value, true
}

// PopLefts pops and returns `size` items from the beginning of array.
func (a *Array[T]) PopLefts(size int) (v []T) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size <= 0 || len(a.array) == 0 {
		return v
	}
	if size >= len(a.array) {
		array := a.array
		a.array = a.array[:0]
		return array
	}
	value := a.array[0:size]
	a.array = a.array[size:]
	return value
}

// PopRights pops and returns `size` items from the end of array.
func (a *Array[T]) PopRights(size int) (v []T) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size <= 0 || len(a.array) == 0 {
		return v
	}
	index := len(a.array) - size
	if index <= 0 {
		array := a.array
		a.array = a.array[:0]
		return array
	}
	value := a.array[index:]
	a.array = a.array[:index]
	return value
}

// Range picks and returns items by range, like array[start:end].
// Notice, if in concurrent-safe usage, it returns a copy of slice;
// else a pointer to the underlying data.
//
// If `end` is negative, then the offset will start from the end of array.
// If `end` is omitted, then the sequence will have everything from start up
// until the end of the array.
func (a *Array[T]) Range(start int, end ...int) (v []T) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	offsetEnd := len(a.array)
	if len(end) > 0 && end[0] < offsetEnd {
		offsetEnd = end[0]
	}
	if start > offsetEnd {
		return v
	}
	if start < 0 {
		start = 0
	}
	array := v
	if a.mu.IsSafe() {
		array = make([]T, offsetEnd-start)
		copy(array, a.array[start:offsetEnd])
	} else {
		array = a.array[start:offsetEnd]
	}
	return array
}

// SubSlice returns a slice of elements from the array as specified
// by the `offset` and `size` parameters.
// If in concurrent safe usage, it returns a copy of the slice; else a pointer.
//
// If offset is non-negative, the sequence will start at that offset in the array.
// If offset is negative, the sequence will start that far from the end of the array.
//
// If length is given and is positive, then the sequence will have up to that many elements in it.
// If the array is shorter than the length, then only the available array elements will be present.
// If length is given and is negative then the sequence will stop that many elements from the end of the array.
// If it is omitted, then the sequence will have everything from offset up until the end of the array.
//
// Any possibility crossing the left border of array, it will fail.
func (a *Array[T]) SubSlice(offset int, length ...int) (v []T) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	size := len(a.array)
	if len(length) > 0 {
		size = length[0]
	}
	if offset > len(a.array) {
		return v
	}
	if offset < 0 {
		offset = len(a.array) + offset
		if offset < 0 {
			return v
		}
	}
	if size < 0 {
		offset += size
		size = -size
		if offset < 0 {
			return v
		}
	}
	end := offset + size
	if end > len(a.array) {
		end = len(a.array)
		size = len(a.array) - offset
	}
	if a.mu.IsSafe() {
		s := make([]T, size)
		copy(s, a.array[offset:])
		return s
	} else {
		return a.array[offset:end]
	}
}

// Append is alias of PushRight, please See PushRight.
func (a *Array[T]) Append(value ...T) *Array[T] {
	a.PushRight(value...)
	return a
}

// Len returns the length of array.
func (a *Array[T]) Len() int {
	a.mu.RLock()
	length := len(a.array)
	a.mu.RUnlock()
	return length
}

// Slice returns the underlying data of array.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (a *Array[T]) Slice() []T {
	if a.mu.IsSafe() {
		a.mu.RLock()
		defer a.mu.RUnlock()
		array := make([]T, len(a.array))
		copy(array, a.array)
		return array
	} else {
		return a.array
	}
}

// Clone returns a new array, which is a copy of current array.
func (a *Array[T]) Clone() (newArray *Array[T]) {
	a.mu.RLock()
	array := make([]T, len(a.array))
	copy(array, a.array)
	a.mu.RUnlock()
	return NewArrayFrom(array, a.mu.IsSafe())
}

// Clear deletes all items of current array.
func (a *Array[T]) Clear() *Array[T] {
	a.mu.Lock()
	if len(a.array) > 0 {
		a.array = make([]T, 0)
	}
	a.mu.Unlock()
	return a
}

// Contains checks whether a value exists in the array.
func (a *Array[T]) Contains(value T) bool {
	return a.Search(value) != -1
}

// Search searches array by `value`, returns the index of `value`,
// or returns -1 if not exists.
func (a *Array[T]) Search(value T) int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.array) == 0 {
		return -1
	}
	result := -1
	for index, v := range a.array {
		var v1 interface{} = v
		var v2 interface{} = value
		if v1 == v2 {
			result = index
			break
		}
	}
	return result
}

// Unique uniques the array, clear repeated items.
// Example: [1,1,2,3,2] -> [1,2,3]
func (a *Array[T]) Unique() *Array[T] {
	a.mu.Lock()
	for i := 0; i < len(a.array)-1; i++ {
		for j := i + 1; j < len(a.array); {
			var v1 interface{} = a.array[i]
			var v2 interface{} = a.array[j]
			if v1 == v2 {
				a.array = append(a.array[:j], a.array[j+1:]...)
			} else {
				j++
			}
		}
	}
	a.mu.Unlock()
	return a
}

// LockFunc locks writing by callback function `f`.
func (a *Array[T]) LockFunc(f func(array []T)) *Array[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	f(a.array)
	return a
}

// RLockFunc locks reading by callback function `f`.
func (a *Array[T]) RLockFunc(f func(array []T)) *Array[T] {
	a.mu.RLock()
	defer a.mu.RUnlock()
	f(a.array)
	return a
}

// Fill fills an array with num entries of the value `value`,
// keys starting at the `startIndex` parameter.
func (a *Array[T]) Fill(startIndex int, num int, value T) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if startIndex < 0 || startIndex > len(a.array) {
		return gserror.NewCodef(gscode.CodeInvalidParameter, "index %d out of array range %d", startIndex, len(a.array))
	}
	for i := startIndex; i < startIndex+num; i++ {
		if i > len(a.array)-1 {
			a.array = append(a.array, value)
		} else {
			a.array[i] = value
		}
	}
	return nil
}

// Chunk splits an array into multiple arrays,
// the size of each array is determined by `size`.
// The last chunk may contain less than size elements.
func (a *Array[T]) Chunk(size int) (v [][]T) {
	if size < 1 {
		return v
	}
	a.mu.RLock()
	defer a.mu.RUnlock()
	length := len(a.array)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]T
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, a.array[i*size:end])
		i++
	}
	return n
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
// Note that do not use pointer as its receiver here.
func (a *Array[T]) MarshalJSON() ([]byte, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return json.Marshal(a.array)
}

// Pad pads array to the specified length with `value`.
// If size is positive then the array is padded on the right, or negative on the left.
// If the absolute value of `size` is less than or equal to the length of the array
// then no padding takes place.
func (a *Array[T]) Pad(size int, val T) *Array[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size == 0 || (size > 0 && size < len(a.array)) || (size < 0 && size > -len(a.array)) {
		return a
	}
	n := size
	if size < 0 {
		n = -size
	}
	n -= len(a.array)
	tmp := make([]T, n)
	for i := 0; i < n; i++ {
		tmp[i] = val
	}
	if size > 0 {
		a.array = append(a.array, tmp...)
	} else {
		a.array = append(tmp, a.array...)
	}
	return a
}

// Reverse makes array with elements in reverse order.
func (a *Array[T]) Reverse() *Array[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, j := 0, len(a.array)-1; i < j; i, j = i+1, j-1 {
		a.array[i], a.array[j] = a.array[j], a.array[i]
	}
	return a
}

// Walk applies a user supplied function `f` to every item of array.
func (a *Array[T]) Walk(f func(value T) T) *Array[T] {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, v := range a.array {
		a.array[i] = f(v)
	}
	return a
}

// IsEmpty checks whether the array is empty.
func (a *Array[T]) IsEmpty() bool {
	return a.Len() == 0
}
