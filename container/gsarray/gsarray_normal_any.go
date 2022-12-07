package gsarray

import (
	"fmt"
	"github.com/jfy0o0/goStealer/errors/gscode"
	"github.com/jfy0o0/goStealer/errors/gserror"
	"github.com/jfy0o0/goStealer/internal/rwmutex"
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

// NewArrayRange creates and returns a array by a range from `start` to `end`
// with step value `step`.
func NewArrayRange[T int](start, end, step int, safe ...bool) *Array[int] {
	if step == 0 {
		panic(fmt.Sprintf(`invalid step value: %d`, step))
	}
	slice := make([]int, (end-start+1)/step)
	index := 0
	for i := start; i <= end; i += step {
		slice[index] = i
		index++
	}
	return NewArrayFrom[int](slice, safe...)
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
		return nil, false
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
		return nil, false
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
		return nil, false
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
		return nil, false
	}
	value = a.array[index]
	a.array = a.array[:index]
	return value, true
}
