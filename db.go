package partdb

import (
	"fmt"
	"log/slog"
)

func NewInstance[K comparable, V any]() Instance[K, V] {
	return NewBipartite[K, V]()
}

func NewBipartite[K comparable, V any]() *Bipartite[K, V] {
	return &Bipartite[K, V]{
		A: map[K]*Node[K, V]{},
		B: map[K]*Node[K, V]{},
	}
}

type Bipartite[K comparable, V any] struct {
	A map[K]*Node[K, V]
	B map[K]*Node[K, V]
}

func newNode[K comparable, V any]() *Node[K, V] {
	return &Node[K, V]{Targets: map[K]*Node[K, V]{}}

}

type Node[K comparable, V any] struct {
	Targets map[K]*Node[K, V]
	Value   V
}

func partName(part bool) string {
	if !part {
		return "A"
	}
	return "B"
}
func (bp *Bipartite[K, V]) getPart(part bool) map[K]*Node[K, V] {
	if !part {
		return bp.A
	}
	return bp.B
}
func (bp *Bipartite[K, V]) getNode(part bool, key K) *Node[K, V] {
	s := bp.getPart(part)
	if v, ok := s[key]; ok {
		return v
	}
	return nil
}

func (bp *Bipartite[K, V]) getOrMake(part bool, key K) *Node[K, V] {
	s := bp.getPart(part)
	if v, ok := s[key]; ok {
		return v
	}

	s[key] = newNode[K, V]()
	return s[key]
}

// Index takes a side, a key, and a key slice. The key is added to the side if it does not exist.
// If the keys in the slice do not exist on the opposite side, they are created.
// All keys are two way indexed. Existing keys are modified.
func (bp Bipartite[K, V]) Index(part bool, key K, keys ...K) error {
	slog.Info("indexing",
		"side", partName(part),
		"key", key,
		"length", len(keys),
	)

	keyNode := bp.getOrMake(part, key)

	for _, target := range keys {
		targetNode := bp.getOrMake(!part, target)
		bp.edge(key, target, keyNode, targetNode)
	}

	return nil
}

// Remove takes a side and key, deleting the key from the side.
// All edges to this key are removed before deletion of the key.
func (bp Bipartite[K, V]) Remove(part bool, key K) error {
	slog.Info("removing",
		"side", partName(part),
		"key", key,
	)

	keyNode := bp.getOrMake(part, key)

	for _, target := range keyNode.Targets {
		delete(target.Targets, key)
	}
	s := bp.getPart(part)
	delete(s, key)

	return nil
}

// GetKey takes a side and a key, returning the value of the key, or returning an error if key does not exist.
func (bp Bipartite[K, V]) GetKey(part bool, key K) (*V, error) {
	slog.Info("getting",
		"side", partName(part),
		"key", key,
	)

	node := bp.getNode(part, key)
	if node == nil {
		slog.Info("key asked which was not found", key)
		return nil, fmt.Errorf("node not found")
	}

	return &node.Value, nil
}

// GetKey takes a side and a key, returning the value of the key, or returning an error if key does not exist.
func (bp Bipartite[K, V]) SetKey(part bool, key K, value V) error {
	slog.Info("setting key with value",
		"side", partName(part),
		"key", key,
		"value", value,
	)

	node := bp.getOrMake(part, key)
	node.Value = value
	return nil
}

// LiseKeys takes a side and a key, returning the index map of that key.
func (bp Bipartite[K, V]) ListKey(part bool, key K) (map[K]V, error) {
	slog.Info("listing key",
		"side", partName(part),
		"key", key,
	)

	node := bp.getNode(part, key)
	if node == nil {
		slog.Info("key asked which was not found", key)
		return nil, fmt.Errorf("node not found")
	}

	res := make(map[K]V, len(node.Targets))
	for key, targetNode := range node.Targets {
		res[key] = targetNode.Value
	}

	return res, nil
}

// ListPart takes a side and returns a map of all keys and their values on that side.
func (bp Bipartite[K, V]) ListPart(part bool) (map[K]V, error) {
	slog.Info("listing side",
		"side", partName(part),
	)

	s := bp.getPart(part)

	res := make(map[K]V, len(s))
	for key, targetNode := range s {
		res[key] = targetNode.Value
	}

	return res, nil
}

func (bp *Bipartite[K, V]) edge(aKey, bKey K, aNode, bNode *Node[K, V]) {
	aNode.Targets[bKey] = bNode
	bNode.Targets[aKey] = aNode
}
