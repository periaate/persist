package partdb

import (
	"encoding/gob"
	"os"
)

type flatPartite[K comparable, V any] struct {
	Akeys      map[K]flatPartiteNode[K, V]
	Bkeyvalues map[K]V
}

type flatPartiteNode[K comparable, V any] struct {
	Avalue V
	Bkeys  []K
}

func flatten[K comparable, V any](bp Bipartite[K, V]) *flatPartite[K, V] {
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

func rebuild[K comparable, V any](fp flatPartite[K, V]) *Bipartite[K, V] {
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

func serializeData[K comparable, V any](bp Bipartite[K, V], fileName string) error {
	fp := flatten(bp)

	file, err := os.Create(fileName)
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

func deserializeData[K comparable, V any](fileName string) (*Bipartite[K, V], error) {
	var data flatPartite[K, V]
	file, err := os.Open(fileName)
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
