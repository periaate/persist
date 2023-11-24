package unordered

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/periaate/persist"
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
	t.Cleanup(CleanupTestFiles)
	name := "testing"
	persistPath := "test"

	originalPM, err := Initialize[uint64, string](persistPath, name, persist.NewHashU64(), 32)
	if err != nil {
		t.Fatalf("Failed to create original PersistMap: %v", err)
	}

	for i, v := range []string{"one", "two", "three"} {
		if err = originalPM.Set(uint64(i), v); err != nil {
			t.Fatal(err)
		}
	}

	if err = originalPM.Close(); err != nil {
		t.Fatal(err)
	}

	loadedPM, err := Initialize[uint64, string](persistPath, name, persist.NewHashU64(), 32)
	if err != nil {
		t.Fatalf("Failed to create original PersistMap: %v", err)
	}

	if originalPM == nil || loadedPM == nil {
		t.Fatalf("One of the PersistMaps is nil")
	}

	for _, el := range originalPM.Obj.Elements {
		if el.HashedKey != 0 {
			loadedEl, ok := loadedPM.Obj.Get(el.Key)
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
