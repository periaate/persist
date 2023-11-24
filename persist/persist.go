package persist

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/periaate/persist"
)

type Wrap[T any] struct {
	Obj *T

	LogEncoder *gob.Encoder
	File       *os.File

	Name     string
	MainPath string
	GobPath  string
	TempPath string
}

func (wr *Wrap[T]) Close() error {
	return wr.Dump()
}

func getLogName(prefix string, ext string) string {
	return fmt.Sprintf("%s-%s%s", prefix, time.Now().Format("20060102-150405"), ext)
}

func (wr *Wrap[T]) Dump() error {
	f, err := os.Create(wr.GobPath)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	if err != nil {
		return err
	}
	return enc.Encode(wr.Obj)
}

func (wr *Wrap[T]) Append(el any) error {
	return nil
	//return wr.LogEncoder.Encode(el)
}

func New[T any](src, name string, t *T) (wr *Wrap[T], err error) {
	if t == nil {
		return nil, fmt.Errorf("t can not be nil")
	}

	err = persist.EnsureDir(src)
	if err != nil {
		return nil, err
	}

	mainPath := filepath.Join(src, name)
	gobPath := filepath.Join(mainPath, fmt.Sprint(name, ".gob"))

	err = persist.EnsureDir(mainPath)
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(gobPath)
	if !os.IsNotExist(err) {
		file, err := os.Open(gobPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(t)
		if err != nil {
			return nil, err
		}

	}

	tempPath := filepath.Join(mainPath, "temp")
	err = persist.EnsureDir(tempPath)
	if err != nil {
		return nil, err
	}

	wr = &Wrap[T]{
		Obj:      t,
		Name:     name,
		MainPath: mainPath,
		GobPath:  gobPath,
		TempPath: tempPath,
	}
	//err = wr.Rotate(false)
	//if err != nil {
	//	return nil, err
	//}

	return wr, nil
}

func (wr *Wrap[T]) Rotate(check bool) error {
	//if wr.File != nil {
	//	if err := wr.File.Close(); err != nil {
	//		return err
	//	}
	//}
	//logName := getLogName(wr.Name, ".lgob")
	//logPath := filepath.Join(wr.TempPath, logName)

	if check {
		checkName := getLogName(wr.TempPath, ".gob")
		checkPath := filepath.Join(wr.TempPath, checkName)

		checkFile, err := os.Create(checkPath)
		if err != nil {
			return err
		}
		defer checkFile.Close()

		encoder := gob.NewEncoder(checkFile)
		encoder.Encode(wr.Obj)
	}

	//logFile, err := os.Create(logPath)
	//if err != nil {
	//	return err
	//}

	//wr.LogEncoder = gob.NewEncoder(logFile)
	//wr.File.Close()
	//wr.File = logFile
	return nil
}
