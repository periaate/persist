package unordered

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/periaate/partdb"
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
	err := db.Dump()
	err2 := db.file.Close()
	return errors.Join(err, err2)
}

func getLogName(prefix string, ext string) string {
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
	if err != nil {
		return err
	}
	return enc.Encode(db.hm)
}

func (db *Map[K, V]) Append(el Element[K, V]) error {
	return db.encoder.Encode(&el)
}

func Initialize[K comparable, V any](name string, path string, hfn func(K) uint64) (*Map[K, V], error) {

	err := partdb.EnsureDir(path)
	if err != nil {
		return nil, err
	}

	folderPath := filepath.Join(path, name)
	filePath := filepath.Join(folderPath, fmt.Sprint(name, ".gob"))

	err = partdb.EnsureDir(folderPath)
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
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(db.hm)
		if err != nil {
			return nil, err
		}

	}

	tempPath := filepath.Join(folderPath, "temp")
	err = partdb.EnsureDir(tempPath)
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
	logName := getLogName(db.name, ".lgob")
	logPath := filepath.Join(db.tempPath, logName)

	if check {
		checkName := getLogName(db.tempPath, ".gob")
		checkPath := filepath.Join(db.tempPath, checkName)

		checkFile, err := os.Create(checkPath)
		if err != nil {
			return err
		}
		defer checkFile.Close()

		encoder := gob.NewEncoder(checkFile)
		encoder.Encode(db.hm)
	}

	logFile, err := os.Create(logPath)
	if err != nil {
		return err
	}

	db.encoder = gob.NewEncoder(logFile)
	db.file = logFile
	return nil
}

func (db *Map[K, V]) Get(key K) (el Element[K, V], ok bool) { return db.hm.Get(key) }

func (db *Map[K, V]) Set(key K, value V) error {
	el := Element[K, V]{
		HashedKey: db.hm.hashFn(key),
		Key:       key,
		Value:     value,
	}
	if err := db.Append(el); err != nil {
		return err
	}
	return db.hm.Set(key, value)
}
