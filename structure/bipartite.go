package structure

import (
	"fmt"
)

// Filler types r and l for type safety in struct.
type rType uint8
type lType byte
type Side uint8

func (s Side) String() string {
	switch s {
	case Right:
		return "RIGHT"
	case Left:
		return "LEFT"
	}
	return ""
}

const (
	Right Side = iota + 1
	Left
)

func Make[K comparable, V any]() Bipartite[K, V] {
	bg := Bipartite[K, V]{
		R: make(map[K]*Graph[rType, lType, K, V], 0),
		L: make(map[K]*Graph[lType, rType, K, V], 0),
	}

	return bg
}

type Bipartite[K comparable, V any] struct {
	// The inverse ordering of the types enforces the bipartite structure.
	R map[K]*Graph[rType, lType, K, V]
	L map[K]*Graph[lType, rType, K, V]
}

type Graph[A, B any, K comparable, V any] struct {
	Targets map[K]*Graph[B, A, K, V]
	Values  map[K]V
	Name    K
	Side    Side
}

func (g Graph[A, B, K, V]) String() string {
	return fmt.Sprintf("Vertex %v in %s", g.Name, g.Side)
}

func (g *Graph[A, B, K, V]) Add(k K, target *Graph[B, A, K, V]) error {
	if _, ok := g.Targets[k]; ok {
		return fmt.Errorf("key with name %v already exists in right %s", k, g)
	}
	g.Targets[k] = target
	return nil
}

func (bg *Bipartite[K, V]) Add(s Side, k K, v V) error {
	switch s {
	case Right:
		if _, ok := bg.R[k]; ok {
			return fmt.Errorf("vertex with name %v already exists in right", k)
		}
		bg.R[k] = &Graph[rType, lType, K, V]{
			Targets: map[K]*Graph[lType, rType, K, V]{},
			Values:  make(map[K]V),
			Name:    k,
			Side:    Right,
		}
	case Left:
		if _, ok := bg.L[k]; ok {
			return fmt.Errorf("vertex with name %v already exists in left", k)
		}
		bg.L[k] = &Graph[lType, rType, K, V]{
			Targets: map[K]*Graph[rType, lType, K, V]{},
			Values:  make(map[K]V),
			Name:    k,
			Side:    Left,
		}
	}
	return nil
}

func (bg *Bipartite[K, V]) Edge(rName, lName K) error {
	var rVert *Graph[rType, lType, K, V]
	var lVert *Graph[lType, rType, K, V]
	var ok bool
	if rVert, ok = bg.R[rName]; !ok {
		return fmt.Errorf("vertex with name %v doesn't exist in in %s", rName, Right)
	}
	if lVert, ok = bg.L[lName]; !ok {
		return fmt.Errorf("vertex with name %v doesn't exist in in %s", lName, Left)
	}

	err := rVert.Add(lName, lVert)
	if err != nil {
		return err
	}
	lVert.Add(rName, rVert)
	if err != nil {
		return err
	}
	return nil
}
