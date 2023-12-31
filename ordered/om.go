package ordered

import (
	"fmt"
	"sync"

	"github.com/periaate/persist"
)

type OrdinalMap[K comparable, V any] struct {
	Elements OrdinalElements[K, V]
	hashFn   func(K) uint64
	resizer  func(uint64) uint64

	Threshold float64
	Len       uint64
	Max       uint64

	mutex sync.RWMutex
}

func New[K comparable, V any](hfn func(K) uint64, maxSize uint64) (*OrdinalMap[K, V], error) {
	hm := &OrdinalMap[K, V]{}
	if maxSize < 8 {
		maxSize = 8
	}
	hm.Max = maxSize

	if hfn == nil {
		return nil, fmt.Errorf("hash function can not be nil")
	}
	hm.hashFn = hfn

	hm.Elements = make([]Orderable[K, V], hm.Max)

	hm.resizer = persist.DefaultInterpolate()
	hm.Threshold = 0.6

	return hm, nil
}

func (hm *OrdinalMap[K, V]) resize() {
	oldEls := hm.Elements
	hm.Max = hm.Max * hm.resizer(hm.Max)
	hm.Elements = make([]Orderable[K, V], hm.Max)

	for _, oldEl := range oldEls {
		if oldEl.HashedKey != 0 {
			hm.Elements[oldEl.HashedKey%hm.Max] = oldEl
		}
	}
}

func (hm *OrdinalMap[K, V]) Get(key K) (el Orderable[K, V], ok bool) {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()
	i := hm.hashFn(key)
	el = hm.Elements[i%hm.Max]
	for el.HashedKey != 0 {
		if el.Key == key {
			return el, true
		}

		i++
		el = hm.Elements[i%hm.Max]
	}

	return Orderable[K, V]{}, false
}

func (hm *OrdinalMap[K, V]) Set(key K, value V) error {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	hash := hm.hashFn(key)
	i := hash
	el := hm.Elements[i%hm.Max]
	for el.HashedKey != 0 {
		if el.Key == key {
			hm.Elements[i%hm.Max] = Orderable[K, V]{HashedKey: hash, Key: key, Value: value}
			return nil
		}

		i++
		el = hm.Elements[i%hm.Max]
	}

	hm.Len++
	hm.Elements[i%hm.Max] = Orderable[K, V]{HashedKey: hash, Key: key, Value: value}
	if float64(hm.Len)/float64(hm.Max) > hm.Threshold {
		hm.resize()
	}
	return nil
}
