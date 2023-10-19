package partdb

import (
	"math"
)

// by using days as the smallest value instead of unix time
// we can modify counting sort to work as "range sort"
// allowing for range selections to occur very easily on results
// as the sorting will only be done on elements which are in the
// correct range of time

type Orderable[K comparable, V any] struct {
	Lex string
	Num uint16

	HashedKey uint64
	Key       K
	Value     V
}

type OrdinalElements[K comparable, V any] []Orderable[K, V]

func (o OrdinalElements[K, V]) Len() int           { return len(o) }
func (o OrdinalElements[K, V]) Less(i, j int) bool { return o[i].Lex < o[j].Lex }
func (o OrdinalElements[K, V]) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

func (o OrdinalElements[K, V]) Order(from, to uint16, origin *K) OrdinalElements[K, V] {
	if len(o) == 0 {
		return nil
	}
	var k int
	var r bool
	if to < from {
		k = int(from - to)
		r = true
	} else {
		k = math.MaxUint16
	}
	count := make([]int, k+1)
	var c int

	if r {
		// Count occurrences of each value.
		for _, v := range o {
			if v.Num >= to && v.Num <= from {
				c++
				count[v.Num-to]++
			}
		}
	} else {
		for _, v := range o {
			count[v.Num]++
		}
	}

	// Build and apply offset by summing the counts of previous values.
	for i := 1; i <= k; i++ {
		count[i] += count[i-1]
	}

	var result OrdinalElements[K, V]
	if r {
		result = make([]Orderable[K, V], c)
		for _, v := range o {
			if v.Num >= to && v.Num <= from {
				result[len(result)-count[v.Num-to]] = v
				count[v.Num-to]--
			}
		}
	} else {
		result = make([]Orderable[K, V], len(o))
		for _, v := range o {

			result[len(o)-count[v.Num]] = v
			count[v.Num]--
		}
	}

	if origin != nil {
		orgK := *origin

		ret := 0
		for i, v := range result {
			if v.Key == orgK {
				ret = i + 1
			}
		}

		return result[ret:]
	}

	return result
}
