package unordered

import (
	"github.com/periaate/persist/persist"
)

type Map[K comparable, V any] struct {
	persist.Wrap[HMap[K, V]]
}

func Initialize[K comparable, V any](src, name string, hfn func(K) uint64, size uint64) (db *Map[K, V], err error) {
	hm, err := New[K, V](hfn, size)
	if err != nil {
		return nil, err
	}

	wr, err := persist.New(src, name, hm)
	if err != nil {
		return nil, err
	}

	db = &Map[K, V]{*wr}

	return db, nil
}

func (db *Map[K, V]) Get(key K) (el Element[K, V], ok bool) { return db.Obj.Get(key) }

func (db *Map[K, V]) Set(key K, value V) error {
	el := Element[K, V]{
		HashedKey: db.Obj.hashFn(key),
		Key:       key,
		Value:     value,
	}
	if err := db.Append(el); err != nil {
		return err
	}
	return db.Obj.Set(key, value)
}
