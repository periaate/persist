package bipartite

import (
	"errors"
	"fmt"
)

var (
	ErrBadSide  = errors.New("requested side does not exist")
	ErrExist    = errors.New("vertex already exists")
	ErrNotExist = errors.New("vertex does not exist")
)

var (
	rSide = map[string]any{
		"r":     "",
		"R":     "",
		"right": "",
		"Right": "",
		"RIGHT": "",
	}

	lSide = map[string]any{
		"l":    "",
		"L":    "",
		"left": "",
		"Left": "",
		"LEFT": "",
	}
)

// Make constructs a generic KV bipartite graph.
func Make[K comparable, V any]() Bipartite[K, V] {
	bg := Bipartite[K, V]{
		R: make(map[K]*Node[K, V], 0),
		L: make(map[K]*Node[K, V], 0),
	}

	return bg
}

// Bipartite is a generic KV bipartite graph.
type Bipartite[K comparable, V any] struct {
	R map[K]*Node[K, V]
	L map[K]*Node[K, V]
}

// Find looks for a key in a given side, returning an error if not found.
func (bg *Bipartite[K, V]) Find(side string, k K) error {
	if _, ok := rSide[side]; ok {
		if _, ok := bg.R[k]; ok {
			return nil
		}
		return ErrNotExist
	}
	if _, ok := lSide[side]; ok {
		if _, ok := bg.L[k]; ok {
			return nil
		}
		return ErrNotExist
	}
	return ErrBadSide
}

// Get looks for a key in the given side, returning it if found.
func (bg *Bipartite[K, V]) Get(side string, k K) (*Node[K, V], error) {
	if _, ok := rSide[side]; ok {
		if part, ok := bg.R[k]; ok {
			return part, nil
		}
		return nil, ErrNotExist
	}
	if _, ok := lSide[side]; ok {
		if part, ok := bg.L[k]; ok {
			return part, nil
		}
		return nil, ErrNotExist
	}
	return nil, ErrBadSide
}

// Add adds a key value pair to the given side.
func (bg *Bipartite[K, V]) Add(s string, k K, v V) error {
	if _, ok := rSide[s]; ok {
		if _, ok := bg.R[k]; ok {
			return fmt.Errorf("vertex with name %v already exists in right", k)
		}
		bg.R[k] = &Node[K, V]{
			Targets: map[K]*Node[K, V]{},
			Value:   v,
		}
		return nil
	}
	if _, ok := lSide[s]; ok {
		if _, ok := bg.L[k]; ok {
			return fmt.Errorf("vertex with name %v already exists in left", k)
		}
		bg.L[k] = &Node[K, V]{
			Targets: map[K]*Node[K, V]{},
			Value:   v,
		}
		return nil
	}
	return fmt.Errorf("no such side")
}

// AddR adds makes a new key-value pair onto the right part.
func (bg *Bipartite[K, V]) AddR(k K, v V) error {
	if _, ok := bg.R[k]; ok {
		return fmt.Errorf("vertex with name %v already exists in right", k)
	}
	bg.R[k] = &Node[K, V]{
		Targets: map[K]*Node[K, V]{},
		Value:   v,
	}
	return nil
}

// AddL adds makes a new key-value pair onto the left part.
func (bg *Bipartite[K, V]) AddL(k K, v V) error {
	if _, ok := bg.L[k]; ok {
		return fmt.Errorf("vertex with name %v already exists in right", k)
	}
	bg.L[k] = &Node[K, V]{
		Targets: map[K]*Node[K, V]{},
		Value:   v,
	}
	return nil
}

// AddValueless adds an array of valueless keys to the given side.
func (bg *Bipartite[K, V]) AddValueless(side string, keys []K) error {
	if _, ok := rSide[side]; ok {
		for _, key := range keys {
			if _, ok := bg.R[key]; ok {
				continue
			}
			bg.R[key] = &Node[K, V]{
				Targets: map[K]*Node[K, V]{},
			}
		}
		return nil
	}
	if _, ok := lSide[side]; ok {
		for _, key := range keys {
			if _, ok := bg.L[key]; ok {
				continue
			}
			bg.L[key] = &Node[K, V]{
				Targets: map[K]*Node[K, V]{},
			}
		}
		return nil
	}
	return fmt.Errorf("no such side")
}

// Edge adds edges from a left key to a right key and from a right key to a left key.
func (bg *Bipartite[K, V]) Edge(lName, rName K) error {
	var rVert *Node[K, V]
	var lVert *Node[K, V]
	var ok bool
	if rVert, ok = bg.R[rName]; !ok {
		return ErrNotExist
	}
	if lVert, ok = bg.L[lName]; !ok {
		return ErrNotExist
	}

	if _, ok := rVert.Targets[lName]; ok {
		return ErrExist
	}
	if _, ok := lVert.Targets[rName]; ok {
		return ErrExist
	}

	rVert.Targets[lName] = lVert
	lVert.Targets[rName] = rVert
	return nil
}

// List looks for a key in the given side, returning all of the keys of its edges.
func (bg *Bipartite[K, V]) List(side string, k K) ([]K, error) {
	if _, ok := rSide[side]; ok {
		if node, ok := bg.R[k]; ok {
			list := []K{}
			for key := range node.Targets {
				list = append(list, key)
			}
			return list, nil
		}
		return nil, ErrNotExist
	}
	if _, ok := lSide[side]; ok {
		if node, ok := bg.L[k]; ok {
			list := []K{}
			for key := range node.Targets {
				list = append(list, key)
			}
			return list, nil
		}
		return nil, ErrNotExist
	}
	return nil, ErrBadSide
}

// Node is a node used in partite graphs.
type Node[K comparable, V any] struct {
	// Targets are the edges this node has to other nodes.
	Targets map[K]*Node[K, V]
	Value   V
}
