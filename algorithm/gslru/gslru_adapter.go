package gslru

type LruAdapter[T comparable] interface {
	PushNormal(key T)
	PushRaiseWithPop(key T)
	PushRaise(key T)
}
