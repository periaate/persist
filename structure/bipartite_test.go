package structure_test

import (
	"fmt"
	. "partdb/structure"
	"testing"
)

func TestMake(t *testing.T) {
	bg := Make[string]()
	if len(bg.R) != 0 || len(bg.L) != 0 {
		t.Errorf("Failed to create an empty Bipartite graph")
	}
}

func TestAddVertices(t *testing.T) {
	bg := Make[int]()
	bg.Add(Right, 1)
	bg.Add(Left, 2)

	if _, ok := bg.R[1]; !ok {
		t.Errorf("Failed to add vertex to the right side")
	}

	if _, ok := bg.L[2]; !ok {
		t.Errorf("Failed to add vertex to the left side")
	}
}

func TestAddEdges(t *testing.T) {
	bg := Make[int]()
	bg.Add(Right, 1)
	bg.Add(Left, 2)
	bg.Edge(1, 2)

	if _, ok := bg.R[1]; !ok || bg.R[1].Targets[2] != bg.L[2] {
		t.Errorf("Failed to add edge from right vertex to left vertex")
	}

	if _, ok := bg.L[2]; !ok || bg.L[2].Targets[1] != bg.R[1] {
		t.Errorf("Failed to add edge from left vertex to right vertex")
	}
}

func TestGet(t *testing.T) {

	bg := Make[string]()
	bg.Add(Right, "Hello")
	bg.Add(Left, "World")
	bg.Edge("Hello", "World")

	r, err := bg.Get(Right, "Hello")
	if err != nil {
		t.Fatal(err)
	}

	l, err := bg.Get(Left, "World")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(r, l)
}
