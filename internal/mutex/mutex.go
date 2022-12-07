package mutex

import "sync"

// Mutex is a sync.Mutex with a switch for concurrent safe feature.
type Mutex struct {
	sync.Mutex
	safe bool
}

// New creates and returns a new *Mutex.
// The parameter `safe` is used to specify whether using this mutex in concurrent-safety,
// which is false in default.
func New(safe ...bool) *Mutex {
	mu := new(Mutex)
	if len(safe) > 0 {
		mu.safe = safe[0]
	} else {
		mu.safe = false
	}
	return mu
}

func (mu *Mutex) IsSafe() bool {
	return mu.safe
}

func (mu *Mutex) Lock() {
	if mu.safe {
		mu.Mutex.Lock()
	}
}

func (mu *Mutex) Unlock() {
	if mu.safe {
		mu.Mutex.Unlock()
	}
}
