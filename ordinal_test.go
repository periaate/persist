package partdb

import (
	"reflect"
	"testing"
)

func TestOrderWithEmptySlice(t *testing.T) {
	var oe OrdinalElements[int, string]
	result := oe.Order(10, 0, nil)
	if len(result) != 0 {
		t.Errorf("Expected empty slice, got: %v", result)
	}
}

func TestOrderWithRangedFromTo(t *testing.T) {
	oe := OrdinalElements[int, string]{
		{Num: 5, Value: "c"}, // 2
		{Num: 1, Value: "a"}, // 0
		{Num: 2, Value: "d"}, // 3
		{Num: 2, Value: "d"}, // 3
		{Num: 3, Value: "b"}, // 1
	}

	result := oe.Order(5, 2, nil)
	expected := OrdinalElements[int, string]{
		{Num: 5, Value: "c"}, // 0
		{Num: 3, Value: "b"}, // 1
		{Num: 2, Value: "d"}, // 2
		{Num: 2, Value: "d"}, // 2
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got: %v", expected, result)
	}
}

func TestOrderWithAllElements(t *testing.T) {
	oe := OrdinalElements[int, string]{
		{Num: 1, Value: "a"},
		{Num: 3, Value: "b"},
		{Num: 5, Value: "c"},
		{Num: 2, Value: "d"},
	}

	result := oe.Order(0, 0, nil)
	expected := OrdinalElements[int, string]{
		{Num: 5, Value: "c"},
		{Num: 3, Value: "b"},
		{Num: 2, Value: "d"},
		{Num: 1, Value: "a"},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got: %v", expected, result)
	}
}

func TestOrderWithOrigin(t *testing.T) {
	oe := OrdinalElements[int, string]{
		{Num: 1, Key: 1, Value: "a"},
		{Num: 3, Key: 3, Value: "b"},
		{Num: 3, Key: 3, Value: "b"},
		{Num: 5, Key: 5, Value: "c"},
		{Num: 5, Key: 5, Value: "c"},
		{Num: 2, Key: 2, Value: "d"},
	}
	origin := 3
	result := oe.Order(5, 1, &origin)
	expected := OrdinalElements[int, string]{
		{Num: 2, Key: 2, Value: "d"},
		{Num: 1, Key: 1, Value: "a"},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got: %v", expected, result)
	}
}

func TestOrderWithNonExistentOrigin(t *testing.T) {
	oe := OrdinalElements[int, string]{
		{Num: 1, Value: "a"},
		{Num: 3, Value: "b"},
		{Num: 5, Value: "c"},
		{Num: 2, Value: "d"},
		{Num: 1, Value: "a"},
		{Num: 1, Value: "a"},
	}
	expected := OrdinalElements[int, string]{
		{Num: 5, Value: "c"},
		{Num: 3, Value: "b"},
		{Num: 2, Value: "d"},
		{Num: 1, Value: "a"},
		{Num: 1, Value: "a"},
		{Num: 1, Value: "a"},
	}
	origin := 4
	result := oe.Order(5, 1, &origin)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got: %v", expected, result)
	}
}
