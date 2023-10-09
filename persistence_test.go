package partdb

import (
	"os"
	"testing"
)

func def() *Bipartite[string, string] {
	aKeys, bKeys := generateKeys(3, 12)
	bp := NewBipartite[string, string]()
	for _, aKey := range aKeys {
		bp.Index(A, aKey, bKeys...)
	}
	return bp
}

func TestSerialize(t *testing.T) {
	bp := def()
	err := serializeData[string, string](*bp, "testFile.gob")
	if err != nil {
		t.Error(err)
	}
}

func TestDeserialize(t *testing.T) {
	_, bKeys := generateKeys(3, 12)
	bp, err := deserializeData[string, string]("testFile.gob")
	if err != nil {
		t.Error(err)
	}

	aKeys, err := bp.ListPart(A)
	if err != nil {
		t.Fatal(err)
	}

	for aKey := range aKeys {
		targets, err := bp.ListKey(A, aKey)
		if err != nil {
			t.Error(err)
		}
		for _, bKey := range bKeys {
			if _, ok := targets[bKey]; !ok {
				t.Error(bKey, "not found in targets")
			}
		}
	}

	err = os.Remove("testFile.gob")
	if err != nil {
		t.Fatal("unable to delete file", err)
	}
}
