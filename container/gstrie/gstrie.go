package gstrie

import (
	"github.com/jfy0o0/goStealer/internal/rwmutex"
)

type trieCallback[T any] interface {
	// AddClash
	// @param 1: old
	// @param 2: new
	// @ret : insert ok ?
	AddClash(T, T) bool

	// Del
	// del callback
	// @param : v
	// @ret : del ok?
	Del(T, interface{}) (needSoftDel bool, happenDel bool)

	// Find
	// find callback
	// @param: v
	// @ret : find ok ?
	Find([]T) ([]T, bool)
}

type Trie[T any] struct {
	Root     *Node[T]
	mu       *rwmutex.RWMutex
	callback trieCallback[T]
}

func NewTrie[T any](isSafe bool, callback ...trieCallback[T]) *Trie[T] {
	t := &Trie[T]{
		Root: NewRootNode[T](0),
		mu:   rwmutex.New(isSafe),
	}
	if len(callback) > 0 {
		t.callback = callback[0]
	} else {
		t.callback = &trieDefaultCallback[T]{}
	}

	return t
}
func (tree *Trie[T]) ReloadFromTrie(newTrie *Trie[T]) {
	tree.mu.Lock()
	newTrie.mu.Lock()
	defer newTrie.mu.Unlock()
	defer tree.mu.Unlock()
	tree.Root = newTrie.Root
}

func (tree *Trie[T]) Add(word string, value T) bool {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	var current = tree.Root
	var runes = []rune(word)
	var length = len(runes)

	if length == 0 {
		return false
	}

	for position := len(runes) - 1; position >= 0; position-- {
		r := runes[position]
		if next, ok := current.Children[r]; ok {
			current = next
		} else {
			newNode := NewNode[T](r)
			current.Children[r] = newNode
			current = newNode
		}
		if position == 0 {
			if current.IsPathEnd() {
				return tree.callback.AddClash(current.v, value)
			}
			current.v = value
			current.isPathEnd = true
		}
	}
	return true
}

func (tree *Trie[T]) Del(word string, value interface{}) bool {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	var current = tree.Root
	var runes = []rune(word)
	var length = len(runes)

	if length == 0 {
		return false
	}
	for position := length - 1; position >= 0; position-- {
		r := runes[position]
		if next, ok := current.Children[r]; !ok {
			return false
		} else {
			current = next
		}
	}
	needSoftDel, happenDel := tree.callback.Del(current.v, value)
	//if !tree.callback.Del(current.v, value) {
	//	return false
	//}
	if needSoftDel {
		current.SoftDel()
	}
	return happenDel
}

func (tree *Trie[T]) Find(text string) (arr []T, ret bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	var (
		parent  = tree.Root
		current *Node[T]
		runes   = []rune(text)
		found   bool
		//hitNode *Node[T]
		length = len(runes)
	)
	if length == 0 {
		return arr, false
	}

	for position := len(runes) - 1; position >= 0; position-- {
		current, found = parent.Children[runes[position]]
		if !found {
			break
		}
		parent = current
		if !current.IsPathEnd() {
			continue
		}
		//hitNode = current
		arr = append(arr, current.v)
	}
	if len(arr) == 0 {
		return arr, false
	}
	return tree.callback.Find(arr)
}
