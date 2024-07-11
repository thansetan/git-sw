package gitconfig

type kv[K comparable, V any] struct {
	key K
	val V
}

type orderedMap[K comparable, V any] struct {
	data map[K]*node[kv[K, V]]
	l    list[kv[K, V]]
}

func newOrderedMap[K comparable, V any]() *orderedMap[K, V] {
	return &orderedMap[K, V]{
		data: make(map[K]*node[kv[K, V]]),
	}
}

func (o *orderedMap[K, V]) put(key K, val V) {
	if _, ok := o.data[key]; !ok {
		node := o.l.pushBack(kv[K, V]{key, val})
		o.data[key] = node
		return
	}
	o.data[key].val.val = val
}

func (o *orderedMap[K, V]) remove(key K) bool {
	node, ok := o.data[key]
	if !ok {
		return false
	}
	o.l.remove(node)
	delete(o.data, key)
	return true
}

func (o *orderedMap[K, V]) getNode(key K) (*node[kv[K, V]], bool) {
	node, ok := o.data[key]
	return node, ok
}

func (o *orderedMap[K, V]) mustGetNode(key K) *node[kv[K, V]] {
	node, ok := o.getNode(key)
	if !ok {
		panic("key can't be found")
	}
	return node
}

func (o *orderedMap[K, V]) get(key K) (V, bool) {
	node, ok := o.getNode(key)
	if !ok {
		var v V
		return v, ok
	}
	return node.val.val, ok
}

func (o *orderedMap[K, V]) mustGet(key K) V {
	val, ok := o.get(key)
	if !ok {
		panic("key can't be found")
	}
	return val
}

func (o *orderedMap[K, V]) len() int {
	return len(o.data)
}

func (o *orderedMap[K, V]) keys() []K {
	keys := make([]K, 0, len(o.data))

	e := o.l.front()

	for e != nil {
		keys = append(keys, e.val.key)
		e = e.next
	}

	return keys
}
