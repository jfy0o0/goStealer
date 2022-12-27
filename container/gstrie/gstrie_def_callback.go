package gstrie

type trieDefaultCallback[T any] struct {
}

func (a *trieDefaultCallback[T]) AddConflict(T, T) bool {
	return true
}

func (a *trieDefaultCallback[T]) Del(T, interface{}) bool {
	return true
}

func (a *trieDefaultCallback[T]) Find(v T) (T, bool) {
	return v, true
}
