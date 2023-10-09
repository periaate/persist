package partdb

import (
	"fmt"
	"testing"
)

const (
	A = false
	B = true
)

// NewTestInstance constructs a Bipartite graph for testing with some preloaded data
func NewTestInstance(aKeys, bKeys []string) Instance[string, string] {
	instance := NewInstance[string, string]()
	for _, aKey := range aKeys {
		instance.Index(A, aKey, bKeys...)
	}
	return instance
}

func generateKeys(aSize, bSize int) ([]string, []string) {
	aKeys := make([]string, aSize)
	bKeys := make([]string, bSize)
	for i := range aKeys {
		aKeys[i] = fmt.Sprint("A KEY", i)
	}
	for i := range bKeys {
		bKeys[i] = fmt.Sprint("B KEY", i)
	}
	return aKeys, bKeys
}

func TestIndex(t *testing.T) {
	aKeys, bKeys := generateKeys(3, 12)
	instance := NewTestInstance(aKeys, bKeys)

	for _, aKey := range aKeys {
		res, err := instance.ListKey(false, aKey)
		if err != nil {
			t.Fatal("fatal error encountered", err)
		}
		for _, bKey := range bKeys {
			if _, ok := res[bKey]; !ok {
				t.Error("bKey", bKey, "was not found in", aKey)
			}
		}
	}
}

func TestRemove(t *testing.T) {
	aKeys, bKeys := generateKeys(1, 4)
	instance := NewTestInstance(aKeys, bKeys)

	instance.Remove(B, bKeys[0])
	instance.Remove(B, bKeys[1])

	aRes, err := instance.ListKey(false, aKeys[0])
	if err != nil {
		t.Fatal("fatal error encountered", err)
	}
	if len(aRes) != 2 {
		t.Error("wrong length after removal, expected 2, received", len(aRes))
	}
}
