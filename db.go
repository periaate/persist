package partdb

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

type Map[K comparable, V any] struct {
	name     string
	path     string
	mainPath string
	tempPath string

	hm *HMap[K, V]

	encoder *gob.Encoder
	file    *os.File
}

func (db *Map[K, V]) Close() error {
	defer db.file.Close()
	return db.Dump()
}

func logName(prefix string, ext string) string {
	return fmt.Sprintf("%s-%s%s", prefix, time.Now().Format("20060102-150405"), ext)
}

func (db *Map[K, V]) Dump() error {
	db.hm.mutex.Lock()
	defer db.hm.mutex.Unlock()

	f, err := os.Create(db.mainPath)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	return enc.Encode(db.hm)
}

func (m *Map[K, V]) Append(el Element[K, V]) error {
	return m.encoder.Encode(&el)
}

func Initialize[K comparable, V any](name string, path string, hfn func(K) uint64) (*Map[K, V], error) {

	err := EnsureDir(path)
	if err != nil {
		return nil, err
	}

	folderPath := filepath.Join(path, name)
	filePath := filepath.Join(folderPath, fmt.Sprint(name, ".gob"))

	err = EnsureDir(folderPath)
	if err != nil {
		return nil, err
	}

	ohm, err := New[K, V](hfn, 8)
	if err != nil {
		return nil, err
	}

	db := &Map[K, V]{
		name:     name,
		mainPath: filePath,
		path:     folderPath,
		hm:       ohm,
	}

	_, err = os.Stat(filePath)
	if !os.IsNotExist(err) {
		hm, _ := New[K, V](hfn, 8)
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(hm)
		if err != nil {
			return nil, err
		}

		hm.hashFn = hfn
		hm.resizer = DefaultInterpolate()

		db.hm = hm
	}

	tempPath := filepath.Join(folderPath, "temp")
	err = EnsureDir(tempPath)
	if err != nil {
		return nil, err
	}

	db.tempPath = tempPath
	err = db.rotate(false)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Map[K, V]) rotate(check bool) error {
	db.hm.mutex.Lock()
	defer db.hm.mutex.Unlock()
	if db.file != nil {
		db.file.Close()
	}
	logname := logName(db.name, ".lgob")
	logpath := filepath.Join(db.tempPath, logname)

	if check {
		checkname := logName(db.tempPath, ".gob")
		checkpath := filepath.Join(db.tempPath, checkname)

		checkf, err := os.Create(checkpath)
		if err != nil {
			return err
		}
		defer checkf.Close()

		encoder := gob.NewEncoder(checkf)
		encoder.Encode(db.hm)
	}

	logf, err := os.Create(logpath)
	if err != nil {
		return err
	}

	db.encoder = gob.NewEncoder(logf)
	db.file = logf
	return nil
}

func (db *Map[K, V]) Get(key K) (el Element[K, V], ok bool) { return db.hm.Get(key) }

func (db *Map[K, V]) Set(key K, value V) error {
	if el, n, ok, diff := db.diff(key, value); diff {
		if !ok {
			el = Element[K, V]{
				HashedKey: db.hm.hashFn(key),
				Key:       key,
				Value:     value,
			}
		}
		if err := db.Append(el); err != nil {
			return err
		}
		if ok {
			db.hm.Elements[n].Value = value
			return nil
		}

		return db.hm.Set(key, value)
	}
	return fmt.Errorf("key not found")
}

func (db *Map[K, V]) diff(key K, value V) (el Element[K, V], n uint64, ok bool, diff bool) {
	el, ok, n = db.hm.get(key)
	if ok {
		if reflect.DeepEqual(el.Value, value) {
			fmt.Println("Deeply equal")
			return el, 0, ok, false
		}

		fmt.Println("Not deeply equal")
		return el, n, ok, true
	}

	return el, 0, ok, true
}
