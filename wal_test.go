package partdb

import (
	"os"
	"testing"
)

// TestDeserialize tests the Deserialize function
func TestPersistWal(t *testing.T) {

	pm, err := NewPersist[uint64, []byte](Hash_u64(), 32, "testlog.lgob")
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(2, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))
	_ = pm.Set(1, []byte("Hello, World!"))

}

func TestLoadPersistMap(t *testing.T) {
	persistPath := "test.wal"

	originalPM, err := NewPersist[uint64, string](Hash_u64(), 16, persistPath)
	if err != nil {
		t.Fatalf("Failed to create original PersistMap: %v", err)
	}

	originalPM.Set(1, "one")
	originalPM.Set(2, "two")
	originalPM.Set(3, "three")
	originalPM.fm.wal.close()

	loadedPM, err := Rebuild[uint64, string](nil, Hash_u64(), persistPath)
	if err != nil {
		t.Fatalf("Failed to load PersistMap: %v", err)
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

	if err := os.Remove(persistPath); err != nil {
		t.Errorf("Failed to remove test WAL file: %v", err)
	}
}
