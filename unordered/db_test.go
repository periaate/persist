package unordered

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/periaate/partdb"
)

// CleanupTestFiles removes all .gob and .lgob files from the current directory
func CleanupTestFiles() {
	fp := filepath.Join("test", "testing")
	err := os.RemoveAll(fp)
	if err != nil {
		panic(err)
	}
}

func TestLoadPersistMap(t *testing.T) {
	name := "testing"
	persistPath := "test"

	originalPM, err := Initialize[uint64, string](name, persistPath, partdb.NewHashU64())
	if err != nil {
		t.Fatalf("Failed to create original PersistMap: %v", err)
	}

	originalPM.Set(1, "one")
	originalPM.Set(2, "two")
	originalPM.Set(3, "three")
	if err = originalPM.Close(); err != nil {
		t.Fatal(err)
	}

	loadedPM, err := Initialize[uint64, string](name, persistPath, partdb.NewHashU64())
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
	if err = loadedPM.Close(); err != nil {
		t.Fatal(err)
	}

}

const (
	times uint64 = 10_000
)

// TestNew tests the New function
func TestNew(t *testing.T) {
	defer t.Cleanup(CleanupTestFiles)
	_, err := New[uint64, string](partdb.NewHashU64(), 32)
	if err != nil {
		t.Fatalf("Failed to create HMap: %v", err)
	}
}

// TestGetSet tests the Get and Set methods
func TestGetSet(t *testing.T) {
	defer t.Cleanup(CleanupTestFiles)
	hm, _ := New[uint64, string](partdb.NewHashU64(), 32)

	for i := uint64(0); i < times; i++ {
		err := hm.Set(i, fmt.Sprint(i))
		if err != nil {
			t.Error(err)
		}
	}

	for i := uint64(0); i < times; i++ {
		if v, ok := hm.Get(i); !ok {
			t.Error(fmt.Errorf("not found"))
		} else if v.Value != fmt.Sprint(i) {
			t.Error(fmt.Errorf("not found"))
		}
	}
}
