package partdb

import (
	"encoding/gob"
	"log/slog"
	"os"
)

// NewPersistentInstance returns a persistent instance.
func NewPersistentInstance[K comparable, V any]() PersistentInstance[K, V] {
	bp := NewBipartite[K, V]()
	pbp := &persistentBp[K, V]{bp}
	return pbp

}

type persistentBp[K comparable, V any] struct {
	*Bipartite[K, V]
}

// Serializes the instance into a file at the path argument. Overwrites or creates file.
func (p *persistentBp[K, V]) Serialize(path string) error {
	return SerializeBipartite[K, V](p.Bipartite, path)
}

// Deserializes the instance from a file at the path argument.
func (p *persistentBp[K, V]) Deserialize(path string) error {
	bp, err := DeserializeBipartite[K, V](path)
	if err != nil {
		return err
	}
	p.Bipartite = bp
	return nil
}

type flatPartite[K comparable, V any] struct {
	Akeys      map[K]flatPartiteNode[K, V]
	Bkeyvalues map[K]V
}

type flatPartiteNode[K comparable, V any] struct {
	Avalue V
	Bkeys  []K
}

func flatten[K comparable, V any](bp *Bipartite[K, V]) *flatPartite[K, V] {
	slog.Info("flattening bipartite")
	fp := &flatPartite[K, V]{
		Akeys:      map[K]flatPartiteNode[K, V]{},
		Bkeyvalues: map[K]V{},
	}
	for aKey, aNode := range bp.A {
		fpn := flatPartiteNode[K, V]{
			Avalue: aNode.Value,
			Bkeys:  make([]K, len(aNode.Targets)),
		}

		var i int
		for bKey := range aNode.Targets {
			fpn.Bkeys[i] = bKey
			i++
		}

		fp.Akeys[aKey] = fpn
	}

	for bKey, bNode := range bp.B {
		fp.Bkeyvalues[bKey] = bNode.Value
	}

	return fp
}

func rebuild[K comparable, V any](fp *flatPartite[K, V]) *Bipartite[K, V] {
	slog.Info("rebuilding bipartite from gob binary")
	bp := NewBipartite[K, V]()
	for aKey, aNode := range fp.Akeys {
		bp.Index(false, aKey, aNode.Bkeys...)
		if v, ok := bp.A[aKey]; ok {
			v.Value = aNode.Avalue
		}
	}

	for bKey, bValue := range fp.Bkeyvalues {
		if v, ok := bp.B[bKey]; ok {
			v.Value = bValue
		}
	}

	return bp
}

func SerializeBipartite[K comparable, V any](bp *Bipartite[K, V], path string) error {
	fp := flatten(bp)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(fp)
	if err != nil {
		return err
	}
	return nil
}

func DeserializeBipartite[K comparable, V any](path string) (*Bipartite[K, V], error) {
	var data *flatPartite[K, V]
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	bp := rebuild[K, V](data)

	return bp, nil
}
