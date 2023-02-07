package gspool

import (
	"github.com/jfy0o0/goStealer/container/gslist"
	"github.com/jfy0o0/goStealer/container/gstype"
	"github.com/jfy0o0/goStealer/errors/gscode"
	"github.com/jfy0o0/goStealer/errors/gserror"
	"github.com/jfy0o0/goStealer/os/gstimer"
	"time"
)

// Pool is an Object-Reusable Pool.
type Pool[T any] struct {
	list    *gslist.List[*poolItem[T]] // Available/idle items list.
	closed  *gstype.Bool               // Whether the pool is closed.
	TTL     time.Duration              // Time To Live for pool items.
	NewFunc func() (T, error)          // Callback function to create pool item.
	// ExpireFunc is the for expired items destruction.
	// This function needs to be defined when the pool items
	// need to perform additional destruction operations.
	// Eg: net.Conn, os.File, etc.
	ExpireFunc func(T)
	timer      *gstimer.Timer
}

// Pool item.
type poolItem[T any] struct {
	value    T     // Item value.
	expireAt int64 // Expire timestamp in milliseconds.
}

// Creation function for object.
type NewFunc[T any] func() (T, error)

// Destruction function for object.
type ExpireFunc[T any] func(T)

// New creates and returns a new object pool.
// To ensure execution efficiency, the expiration time cannot be modified once it is set.
//
// Note the expiration logic:
// ttl = 0 : not expired;
// ttl < 0 : immediate expired after use;
// ttl > 0 : timeout expired;
func New[T any](ttl time.Duration, newFunc NewFunc[T], expireFunc ...ExpireFunc[T]) *Pool[T] {
	r := &Pool[T]{
		list:    gslist.New[*poolItem[T]](true),
		closed:  gstype.NewBool(),
		TTL:     ttl,
		NewFunc: newFunc,
		timer:   gstimer.New(),
	}
	if len(expireFunc) > 0 {
		r.ExpireFunc = expireFunc[0]
	}
	r.timer.AddSingleton(time.Second, r.checkExpireItems)
	return r
}

// Put puts an item to pool.
func (p *Pool[T]) Put(value T) error {
	if p.closed.Val() {
		return gserror.NewCode(gscode.CodeInvalidOperation, "pool is closed")
	}
	item := &poolItem[T]{
		value: value,
	}
	if p.TTL == 0 {
		item.expireAt = 0
	} else {
		// As for Golang version < 1.13, there's no method Milliseconds for time.Duration.
		// So we need calculate the milliseconds using its nanoseconds value.
		item.expireAt = time.Now().UnixMilli() + p.TTL.Nanoseconds()/1000000
	}
	p.list.PushBack(item)
	return nil
}

// Clear clears pool, which means it will remove all items from pool.
func (p *Pool[T]) Clear() {
	if p.ExpireFunc != nil {
		for {
			if r, ok := p.list.PopFront(); ok {
				p.ExpireFunc(r.value)
			} else {
				break
			}
		}
	} else {
		p.list.RemoveAll()
	}

}

// Get picks and returns an item from pool. If the pool is empty and NewFunc is defined,
// it creates and returns one from NewFunc.
func (p *Pool[T]) Get() (T, error) {
	for !p.closed.Val() {
		if r, ok := p.list.PopFront(); ok {
			f := r
			if f.expireAt == 0 || f.expireAt > time.Now().UnixMilli() {
				return f.value, nil
			} else if p.ExpireFunc != nil {
				// TODO: move expire function calling asynchronously from `Get` operation.
				p.ExpireFunc(f.value)
			}
		} else {
			break
		}
	}
	if p.NewFunc != nil {
		return p.NewFunc()
	}
	return nil, gserror.NewCode(gscode.CodeInvalidOperation, "pool is empty")
}

// Size returns the count of available items of pool.
func (p *Pool[T]) Size() int {
	return p.list.Len()
}

// Close closes the pool. If <p> has ExpireFunc,
// then it automatically closes all items using this function before it's closed.
// Commonly you do not need call this function manually.
func (p *Pool[T]) Close() {
	p.closed.Set(true)
}

// checkExpire removes expired items from pool in every second.
func (p *Pool[T]) checkExpireItems() {
	if p.closed.Val() {
		// If p has ExpireFunc,
		// then it must close all items using this function.
		if p.ExpireFunc != nil {
			for {
				if r, ok := p.list.PopFront(); ok {
					p.ExpireFunc(r.value)
				} else {
					break
				}
			}
		}
		p.timer.Stop()
	}
	// All items do not expire.
	if p.TTL == 0 {
		return
	}
	// The latest item expire timestamp in milliseconds.
	var latestExpire int64 = -1
	// Retrieve the current timestamp in milliseconds, it expires the items
	// by comparing with this timestamp. It is not accurate comparison for
	// every items expired, but high performance.
	var timestampMilli = time.Now().UnixNano() / 1e6
	for {
		if latestExpire > timestampMilli {
			break
		}
		if r, ok := p.list.PopFront(); ok {
			item := r
			latestExpire = item.expireAt
			// TODO improve the auto-expiration mechanism of the pool.
			if item.expireAt > timestampMilli {
				p.list.PushFront(item)
				break
			}
			if p.ExpireFunc != nil {
				p.ExpireFunc(item.value)
			}
		} else {
			break
		}
	}
}
