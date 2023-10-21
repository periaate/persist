package unordered

import (
	"fmt"
	"testing"

	"github.com/periaate/partdb"
)

const (
	times uint64 = 10_000
)

// TestNew tests the New function
func TestNew(t *testing.T) {
	_, err := New[uint64, string](partdb.NewHashU64(), 32)
	if err != nil {
		t.Fatalf("Failed to create HMap: %v", err)
	}
}

// TestGetSet tests the Get and Set methods
func TestGetSet(t *testing.T) {
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
