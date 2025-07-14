package main

import (
	"fmt"
	"slices"
)

var PERM uint8 = 7
var TEMP uint8 = 1

func topologicalSort(tables [][2]string) ([]string, error) {
	var adj map[string][]string = make(map[string][]string)
	var isChild map[string]bool = make(map[string]bool) // indicates whether node has an incoming edge(s)
	var visited map[string]uint8 = make(map[string]uint8)
	for _, edge := range tables {
		p := edge[0]
		c := edge[1]
		isChild[c] = true
		adj[p] = append(adj[p], c)
	}
	topoSorted := []string{}
	for _, curr := range tables {
		p := curr[0]
		// if p is never a child, it must have no incoming eedges
		if !isChild[p] {
			if err := dfs(adj, visited, &topoSorted, p); err != nil {
				return nil, err
			}
		}
	}
	slices.Reverse(topoSorted)
	fmt.Printf("topoSorted: %v\n", topoSorted)
	return topoSorted, nil
}

func dfs(adj map[string][]string, visited map[string]uint8, topo *[]string, curr string) error {
	if visited[curr] == PERM {
		return nil
	}
	if visited[curr] == TEMP {
		return fmt.Errorf("Error: cycle detected in graph")
	}
	visited[curr] = TEMP
	neighbors := adj[curr]
	for _, neighbor := range neighbors {
		// visited[neighbor] = TEMP
		if err := dfs(adj, visited, topo, neighbor); err != nil {
			return err
		}
		// visited[neighbor] = 0
	}
	visited[curr] = PERM
	*topo = append(*topo, curr)
	return nil
}
