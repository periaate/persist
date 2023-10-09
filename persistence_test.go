package partdb

import (
	"os"
	"testing"
)

func TestSerialization(t *testing.T) {
	basekeys, bKeys := generateKeys(3, 12)
	bp := NewBipartite[string, string]()
	for _, aKey := range basekeys {
		bp.Index(A, aKey, bKeys...)
	}

	err := SerializeBipartite[string, string](bp, "testFile.gob")
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = os.Remove("testFile.gob")
		if err != nil {
			t.Fatal("unable to delete file", err)
		}
	}()

	bpdeserialize, err := DeserializeBipartite[string, string]("testFile.gob")
	if err != nil {
		t.Error(err)
	}

	aKeys, err := bpdeserialize.ListPart(A)
	if err != nil {
		t.Fatal(err)
	}

	for aKey := range aKeys {
		targets, err := bpdeserialize.ListKey(A, aKey)
		if err != nil {
			t.Error(err)
		}
		for _, bKey := range bKeys {
			if _, ok := targets[bKey]; !ok {
				t.Error(bKey, "not found in targets")
			}
		}
	}
}

func TestPersistentInstance(t *testing.T) {
	aKeys, bKeys := generateKeys(3, 12)

	instance := NewPersistentInstance[string, string]()
	for _, aKey := range aKeys {
		instance.Index(A, aKey, bKeys...)
	}

	testfileName := "persistentTestFile.gob"

	err := instance.Serialize(testfileName)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = os.Remove(testfileName)
		if err != nil {
			t.Fatal("unable to delete file", err)
		}
	}()

	err = instance.Deserialize(testfileName)
	if err != nil {
		t.Error(err)
	}

	aKeyMap, err := instance.ListPart(A)
	if err != nil {
		t.Fatal(err)
	}

	for aKey := range aKeyMap {
		targets, err := instance.ListKey(A, aKey)
		if err != nil {
			t.Error(err)
		}
		for _, bKey := range bKeys {
			if _, ok := targets[bKey]; !ok {
				t.Error(bKey, "not found in targets")
			}
		}
	}
}
