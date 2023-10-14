package partdb

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"time"
)

func (w *Wal[K, V]) Append(el Element[K, V]) error {
	if w.encoder == nil {
		file, err := os.Create(w.path)
		if err != nil {
			return err
		}
		if err != nil {
			if err != os.ErrExist {
				return err
			}
			file, err = os.Open(w.path)
			if err != nil {
				return err
			}
		}
		w.encoder = gob.NewEncoder(file)
		w.file = file
	}

	return w.encoder.Encode(&el)
}

// Rebuild rebuilds a new PersistMap from a wal or replays a wal to the passed PersistMap.
func Rebuild[K comparable, V any](pm *PersistMap[K, V], hfn func(K) uint64, persistPath string) (_ *PersistMap[K, V], err error) {
	if pm == nil {
		pm, err = NewPersist[K, V](hfn, 32, persistPath)
		if err != nil {
			return nil, err
		}
	}

	file, err := os.Open(persistPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	for {
		var el Element[K, V]
		err := decoder.Decode(&el)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		pm.hm.Set(el.Key, el.Value)
	}

	return pm, nil
}

type Wal[K comparable, V any] struct {
	path    string
	encoder *gob.Encoder
	file    *os.File
}

type FileManager[K comparable, V any] struct {
	basePath string
	wal      *Wal[K, V]
}

func (w *Wal[K, V]) close() error {
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

func OpenFileManager[K comparable, V any](prefix string) *FileManager[K, V] {
	fm := &FileManager[K, V]{
		wal:      &Wal[K, V]{path: GenerateFileName(prefix)},
		basePath: prefix,
	}
	return fm
}

func GenerateFileName(prefix string) string {
	return fmt.Sprintf("%s-%s.lgob", prefix, time.Now().Format("20060102-150405"))
}

func (p *PersistMap[K, V]) Checkpoint(path string) error { return nil } // Unimplemented

func (p *PersistMap[K, V]) Dump(path string) error {
	p.hm.mutex.Lock()
	defer p.hm.mutex.Unlock()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	err = enc.Encode(p.hm)
	return err
}

func (fm *FileManager[K, V]) Append(el Element[K, V]) error { return fm.wal.Append(el) }
