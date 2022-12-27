package gstrie

type Node[T any] struct {
	isRootNode bool
	isPathEnd  bool
	Character  rune
	Children   map[rune]*Node[T]
	v          T
}

func NewNode[T any](character rune) *Node[T] {
	return &Node[T]{
		Character: character,
		Children:  make(map[rune]*Node[T], 0),
	}
}

func NewRootNode[T any](character rune) *Node[T] {
	return &Node[T]{
		isRootNode: true,
		Character:  character,
		Children:   make(map[rune]*Node[T], 0),
	}
}

func (node *Node[T]) IsLeafNode() bool {
	return len(node.Children) == 0
}

func (node *Node[T]) IsRootNode() bool {
	return node.isRootNode
}

func (node *Node[T]) IsPathEnd() bool {
	return node.isPathEnd
}

func (node *Node[T]) SoftDel() {
	node.isPathEnd = false
}
