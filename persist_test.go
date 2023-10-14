package partdb

import (
	"os"
	"testing"
)

// TestNew tests the New function
func TestNew(t *testing.T) {
	t.Cleanup(CleanupTestFiles)
	_, err := New[uint64, string](Hash_u64(), 32)
	if err != nil {
		t.Fatalf("Failed to create HMap: %v", err)
	}
}

// TestGetSet tests the Get and Set methods
func TestGetSet(t *testing.T) {
	t.Cleanup(CleanupTestFiles)
	hm, _ := New[uint64, string](Hash_u64(), 32)
	hm.Set(1, "one")

	el, ok := hm.Get(1)
	if !ok || el.Value != "one" {
		t.Fatalf("Get or Set method failed")
	}
}

// TestSerialize tests the Serialize function
func TestSerialize(t *testing.T) {
	t.Cleanup(CleanupTestFiles)
	hm, _ := New[uint64, string](Hash_u64(), 32)
	hm.Set(1, "one")
	Serialize(hm, "serialize_test.gob")

	if _, err := os.Stat("serialize_test.gob"); os.IsNotExist(err) {
		t.Fatalf("Serialization failed, file not found")
	}
}

// TestDeserialize tests the Deserialize function
func TestDeserialize(t *testing.T) {
	t.Cleanup(CleanupTestFiles)
	hm1, _ := New[uint64, string](Hash_u64(), 32)
	hm1.Set(1, "one")
	err := Serialize(hm1, "deserialize_test.gob")
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	hm2, err := Deserialize[uint64, string]("deserialize_test.gob", Hash_u64())
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}
	if hm2 == nil {
		t.Fatal("the deserialized hashmap is nil")
	}

	for _, v := range hm1.Elements {
		if v.HashedKey == 0 {
			continue
		}
		if _, ok := hm2.Get(v.Key); !ok {
			t.Fatalf("Original and deserialized HMaps are not the same")
		}
	}
}
