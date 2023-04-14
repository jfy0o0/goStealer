package gstrie

type trieDefaultCallback[T any] struct {
}

func (a *trieDefaultCallback[T]) AddClash(T, T) bool {
	return true
}

func (a *trieDefaultCallback[T]) Del(T, interface{}) (bool, bool) {
	return true, true
}

func (a *trieDefaultCallback[T]) Find(arr []T) ([]T, bool) {
	return arr, true
}
