package gitconfig

type node[T any] struct {
	val        T
	next, prev *node[T]
}

type list[T any] struct {
	root node[T]
}

func (l *list[T]) front() *node[T] {
	return l.root.next
}

func (l *list[T]) pushBack(val T) *node[T] {
	node := &node[T]{
		val: val,
	}

	if l.root.next == nil {
		l.root.next = node
		l.root.prev = node
		return node
	}

	node.prev = l.root.prev
	l.root.prev.next = node
	l.root.prev = node

	return node
}

func (l *list[T]) remove(node *node[T]) {
	if node.prev == nil {
		l.root.next = node.next
	} else {
		node.prev.next = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	}
}
