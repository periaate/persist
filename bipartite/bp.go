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

func Make[K comparable, V any]() Bipartite[K, V] {
	bg := Bipartite[K, V]{
		R: make(map[K]*Part[K, V], 0),
		L: make(map[K]*Part[K, V], 0),
	}

	return bg
}

type Bipartite[K comparable, V any] struct {
	R map[K]*Part[K, V]
	L map[K]*Part[K, V]
}

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

func (bg *Bipartite[K, V]) Get(side string, k K) (map[K]*Part[K, V], error) {
	if _, ok := rSide[side]; ok {
		if part, ok := bg.R[k]; ok {
			return part.Targets, nil
		}
		return nil, ErrNotExist
	}
	if _, ok := lSide[side]; ok {
		if part, ok := bg.L[k]; ok {
			return part.Targets, nil
		}
		return nil, ErrNotExist
	}
	return nil, ErrBadSide
}

func (bg *Bipartite[K, V]) Add(s string, k K, v V) error {
	if _, ok := rSide[s]; ok {
		if _, ok := bg.R[k]; ok {
			return fmt.Errorf("vertex with name %v already exists in right", k)
		}
		bg.R[k] = &Part[K, V]{
			Targets: map[K]*Part[K, V]{},
			Value:   v,
		}
		return nil
	}
	if _, ok := lSide[s]; ok {
		if _, ok := bg.L[k]; ok {
			return fmt.Errorf("vertex with name %v already exists in left", k)
		}
		bg.L[k] = &Part[K, V]{
			Targets: map[K]*Part[K, V]{},
			Value:   v,
		}
		return nil
	}
	return fmt.Errorf("no such side")
}

func (bg *Bipartite[K, V]) Edge(rName, lName K) error {
	var rVert *Part[K, V]
	var lVert *Part[K, V]
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

type Part[K comparable, V any] struct {
	Targets map[K]*Part[K, V]
	Value   V
}
