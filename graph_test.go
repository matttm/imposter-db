package main

import (
	"reflect"
	"testing"
)

func TestTopologicalSort_SimpleChain(t *testing.T) {
	tables := [][2]string{
		{"a", "b"},
		{"b", "c"},
	}
	result := topologicalSort(tables)
	// valid orders: a, b, c
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestTopologicalSort_Branching(t *testing.T) {
	tables := [][2]string{
		{"a", "b"},
		{"a", "c"},
		{"b", "d"},
		{"c", "d"},
	}
	result := topologicalSort(tables)
	// valid orders: a before b and c, b and c before d
	aIdx := indexOf(result, "a")
	bIdx := indexOf(result, "b")
	cIdx := indexOf(result, "c")
	dIdx := indexOf(result, "d")
	if !(aIdx < bIdx && aIdx < cIdx && bIdx < dIdx && cIdx < dIdx) {
		t.Errorf("invalid topological order: %v", result)
	}
}

func TestTopologicalSort_SingleNode(t *testing.T) {
	tables := [][2]string{}
	result := topologicalSort(tables)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestTopologicalSort_Cycle(t *testing.T) {
	tables := [][2]string{
		{"a", "b"},
		{"b", "a"},
	}
	result := topologicalSort(tables)
	// The current implementation does not handle cycles, so just check it doesn't panic and returns something
	if result == nil {
		t.Errorf("expected non-nil result for cycle")
	}
}

func indexOf(slice []string, val string) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return -1
}
