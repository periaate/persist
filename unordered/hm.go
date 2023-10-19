package partdb

import (
	"fmt"
	"sync"
)

type Element[K comparable, V any] struct {
	HashedKey uint64
	Key       K
	Value     V
}

type HMap[K comparable, V any] struct {
	Elements []Element[K, V]
	hashFn   func(K) uint64
	resizer  func(uint64) uint64

	Threshold float64
	Len       uint64
	Max       uint64

	mutex sync.RWMutex
}

func New[K comparable, V any](hfn func(K) uint64, maxSize uint64) (*HMap[K, V], error) {
	hm := &HMap[K, V]{}
	if maxSize < 8 {
		maxSize = 8
	}
	hm.Max = maxSize

	if hfn == nil {
		return nil, fmt.Errorf("hash function can not be nil")
	}
	hm.hashFn = hfn

	hm.Elements = make([]Element[K, V], hm.Max)

	hm.resizer = DefaultInterpolate()

	return hm, nil
}

func (hm *HMap[K, V]) resize() {
	oldEls := hm.Elements
	hm.Max = hm.Max * hm.resizer(hm.Max)
	hm.Elements = make([]Element[K, V], hm.Max)

	for _, oldEl := range oldEls {
		if oldEl.HashedKey != 0 {
			hm.Elements[oldEl.HashedKey%hm.Max] = oldEl
		}
	}
}

func (hm *HMap[K, V]) Get(key K) (el Element[K, V], ok bool) {
	// Race
	i := hm.hashFn(key)
	el = hm.Elements[i%hm.Max]
	for el.HashedKey != 0 {
		if el.Key == key {
			return el, true
		}

		i++
		el = hm.Elements[i%hm.Max]
	}

	return Element[K, V]{}, false
}

func (hm *HMap[K, V]) get(key K) (el Element[K, V], ok bool, n uint64) {
	// Race
	i := hm.hashFn(key)
	el = hm.Elements[i%hm.Max]
	for el.HashedKey != 0 {
		if el.Key == key {
			return el, true, i
		}

		i++
		el = hm.Elements[i%hm.Max]
	}

	return Element[K, V]{}, false, 0
}

func (hm *HMap[K, V]) Set(key K, value V) error {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	hash := hm.hashFn(key)
	i := hash
	el := hm.Elements[i%hm.Max]
	for el.HashedKey != 0 {
		if el.Key == key {
			hm.Elements[i%hm.Max] = Element[K, V]{hash, key, value}
			return nil
		}

		i++
		el = hm.Elements[i%hm.Max]
	}

	hm.Elements[i%hm.Max] = Element[K, V]{hash, key, value}
	if float64(hm.Len)/float64(hm.Max) > hm.Threshold {
		hm.resize()
	}
	return nil
}
