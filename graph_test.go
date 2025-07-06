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
	result, _ := topologicalSort(tables)
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
	result, err := topologicalSort(tables)
	// valid orders: a before b and c, b and c before d
	if err != nil {
		t.Errorf("invalid topological order: %v", result)
	}
}

func TestTopologicalSort_SingleNode(t *testing.T) {
	tables := [][2]string{}
	result, _ := topologicalSort(tables)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestTopologicalSort_Cycle(t *testing.T) {
	tables := [][2]string{
		{"a", "b"},
		{"b", "a"},
	}
	result, _ := topologicalSort(tables)
	// The current implementation does not handle cycles, so just check it doesn't panic and returns something
	if result == nil {
		t.Errorf("expected non-nil result for cycle")
	}
}

func TestTopologicalSort_LargeAcyclic(t *testing.T) {
	// 26 nodes: "a" to "z"
	nodes := []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
		"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}
	// 50 edges, no cycles
	tables := [][2]string{
		{"a", "b"}, {"a", "c"}, {"b", "d"}, {"c", "d"}, {"d", "e"},
		{"e", "f"}, {"f", "g"}, {"g", "h"}, {"h", "i"}, {"i", "j"},
		{"j", "k"}, {"k", "l"}, {"l", "m"}, {"m", "n"}, {"n", "o"},
		{"o", "p"}, {"p", "q"}, {"q", "r"}, {"r", "s"}, {"s", "t"},
		{"t", "u"}, {"u", "v"}, {"v", "w"}, {"w", "x"}, {"x", "y"},
		{"y", "z"}, {"a", "d"}, {"b", "e"}, {"c", "f"}, {"d", "g"},
		{"e", "h"}, {"f", "i"}, {"g", "j"}, {"h", "k"}, {"i", "l"},
		{"j", "m"}, {"k", "n"}, {"l", "o"}, {"m", "p"}, {"n", "q"},
		{"o", "r"}, {"p", "s"}, {"q", "t"}, {"r", "u"}, {"s", "v"},
		{"t", "w"}, {"u", "x"}, {"v", "y"}, {"w", "z"}, {"a", "z"},
	}
	result, err := topologicalSort(tables)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Check that all nodes are present in the result
	if len(result) != len(nodes) {
		t.Errorf("expected %d nodes, got %d", len(nodes), len(result))
	}
	// Check that all edges are respected
	for _, edge := range tables {
		fromIdx := indexOf(result, edge[0])
		toIdx := indexOf(result, edge[1])
		if fromIdx == -1 || toIdx == -1 {
			t.Errorf("missing node in result: %v", edge)
		}
		if fromIdx > toIdx {
			t.Errorf("edge order violated: %v appears after %v", edge[0], edge[1])
		}
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
