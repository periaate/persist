package partdb

import (
	"encoding/gob"
	"os"
)

func Serialize[K comparable, V any](hm *HMap[K, V], path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(hm)
	if err != nil {
		return err
	}
	return nil
}

func Deserialize[K comparable, V any](path string, hashFn func(K) uint64) (*HMap[K, V], error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	hm, err := New[K, V](hashFn, 8)
	if err != nil {
		return nil, err
	}

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(hm)
	if err != nil {
		return nil, err
	}

	return hm, nil
}
