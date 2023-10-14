package partdb

import (
	"os"
	"path/filepath"
	"testing"
)

// CleanupTestFiles removes all .gob and .lgob files from the current directory
func CleanupTestFiles() {
	p := "./"
	files, err := os.ReadDir(p)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext == ".gob" || ext == ".lgob" {
			os.Remove(filepath.Join(p, file.Name()))
		}
	}
}

func TestLoadPersistMap(t *testing.T) {
	t.Cleanup(CleanupTestFiles)
	name := "testing"
	persistPath := "testWal"

	originalPM, err := Initialize[uint64, string](name, persistPath, Hash_u64())
	if err != nil {
		t.Fatalf("Failed to create original PersistMap: %v", err)
	}

	originalPM.Set(1, "one")
	originalPM.Set(2, "two")
	originalPM.Set(3, "three")
	originalPM.Close()

	loadedPM, err := Initialize[uint64, string](name, persistPath, Hash_u64())
	if err != nil {
		t.Fatalf("Failed to create original PersistMap: %v", err)
	}

	if originalPM == nil || loadedPM == nil {
		t.Fatalf("One of the PersistMaps is nil")
	}

	for _, el := range originalPM.hm.Elements {
		if el.HashedKey != 0 {
			loadedEl, ok := loadedPM.hm.Get(el.Key)
			if !ok {
				t.Errorf("Element with key %v not found in loaded PersistMap", el.Key)
			}
			if loadedEl.Value != el.Value {
				t.Errorf("Element values do not match for key %v: original=%v, loaded=%v", el.Key, el.Value, loadedEl.Value)
			}
		}
	}
	loadedPM.Close()
}
