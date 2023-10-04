package bipartite_test

import (
	"fmt"
	. "partdb/bipartite"
	"testing"
)

const (
	Right = "RIGHT"
	Left  = "LEFT"
)

func TestMake(t *testing.T) {
	bg := Make[string, string]()
	if len(bg.R) != 0 || len(bg.L) != 0 {
		t.Errorf("Failed to create an empty Bipartite graph")
	}
}

func TestAdd(t *testing.T) {
	bg := Make[string, string]()
	bg.Add(Right, "Hello", ",")
	bg.Add(Left, "World", "!")

	if v, ok := bg.R["Hello"]; ok {
		if v.Value != "," {
			t.Errorf("right vertex has wrong value")
		}

	} else {
		t.Errorf("failed to add vertex to the right side")
	}

	if v, ok := bg.L["World"]; ok {

		if v.Value != "!" {
			t.Errorf("left vertex has wrong value")
		}

	} else {
		t.Errorf("failed to add vertex to the left side")
	}
}
func TestAddR(t *testing.T) {
	bg := Make[string, string]()
	bg.AddR("Hello", ",")

	if v, ok := bg.R["Hello"]; ok {
		if v.Value != "," {
			t.Errorf("right vertex has wrong value")
		}

	} else {
		t.Errorf("failed to add vertex to the right side")
	}
}
func TestAddL(t *testing.T) {
	bg := Make[string, string]()
	bg.AddL("World", "!")

	if v, ok := bg.L["World"]; ok {

		if v.Value != "!" {
			t.Errorf("left vertex has wrong value")
		}

	} else {
		t.Errorf("failed to add vertex to the left side")
	}
}

func TestAddEdges(t *testing.T) {
	rKey := "Hello"
	lKey := "World"
	bg := Make[string, string]()
	bg.Add(Right, rKey, ",")
	bg.Add(Left, lKey, "!")
	bg.Edge(lKey, rKey)

	if _, ok := bg.R[rKey]; !ok || bg.R[rKey].Targets[lKey] != bg.L[lKey] {
		t.Errorf("Failed to add edge from right vertex to left vertex")
	}

	if _, ok := bg.L[lKey]; !ok || bg.L[lKey].Targets[rKey] != bg.R[rKey] {
		t.Errorf("Failed to add edge from left vertex to right vertex")
	}
}

func TestGet(t *testing.T) {
	rKey := "Hello"
	lKey := "World"
	bg := Make[string, string]()
	bg.AddR(rKey, ",")
	bg.AddL(lKey, "!")
	bg.Edge(rKey, lKey)

	rNode, err := bg.Get(Right, "Hello")
	if err != nil {
		t.Fatal(err)
	}
	if rNode.Value != "," {
		t.Fatal(fmt.Errorf("received unexpected value, wanted \",\", received: \"%s\"", rNode.Value))
	}

	lNode, err := bg.Get(Left, "World")
	if err != nil {
		t.Fatal(err)
	}
	if lNode.Value != "!" {
		t.Fatal(fmt.Errorf("received unexpected value, wanted \"!\", received: \"%s\"", lNode.Value))
	}
}

func TestList(t *testing.T) {
	keyArr := []string{"one", "two", "three"}
	bg := Make[string, string]()
	bg.AddValueless(Left, keyArr)
	bg.AddR("Hello", "World")

	bg.Edge("one", "Hello")
	bg.Edge("two", "Hello")
	bg.Edge("three", "Hello")
	keyList, err := bg.List(Right, "Hello")
	if err != nil {
		t.Fatal(err)
	}
	for i, key := range keyList {
		expect := keyArr[i]
		if key != expect {
			t.Fatal(fmt.Errorf("received unexpected value, wanted \"%s\", received: \"%s\"", expect, key))
		}
	}
}
