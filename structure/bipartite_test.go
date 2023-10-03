package structure_test

import (
	. "partdb/structure"
	"testing"
)

func TestMake(t *testing.T) {
	bg := Make[Side, string]()
	if len(bg.R) != 0 || len(bg.L) != 0 {
		t.Errorf("Failed to create an empty Bipartite graph")
	}
}

func TestAddVertices(t *testing.T) {
	bg := Make[Side, string]()
	bg.Add(Right, 1, "Right")
	bg.Add(Left, 2, "Left")

	if _, ok := bg.R[1]; !ok {
		t.Errorf("Failed to add vertex to the right side")
	}

	if _, ok := bg.L[2]; !ok {
		t.Errorf("Failed to add vertex to the left side")
	}
}

func TestAddEdges(t *testing.T) {
	bg := Make[Side, string]()
	bg.Add(Right, 1, "Right")
	bg.Add(Left, 2, "Left")
	bg.Edge(1, 2)

	if _, ok := bg.R[1]; !ok || bg.R[1].Targets[2] != bg.L[2] {
		t.Errorf("Failed to add edge from right vertex to left vertex")
	}

	if _, ok := bg.L[2]; !ok || bg.L[2].Targets[1] != bg.R[1] {
		t.Errorf("Failed to add edge from left vertex to right vertex")
	}
}
