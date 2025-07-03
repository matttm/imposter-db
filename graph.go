package main

import "fmt"

var PERM uint8 = 7
var TEMP uint8 = 1

func topologicalSort(tables [][2]string) []string {
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
			dfs(adj, visited, &topoSorted, p)
		}
	}
	fmt.Printf("topoSorted: %v\n", topoSorted)
	return topoSorted
}

func dfs(adj map[string][]string, visited map[string]uint8, topo *[]string, curr string) {
	if visited[curr] == PERM {
		return
	}
	if visited[curr] == TEMP {
		panic("Error: cycle detected in graph")
	}
	*topo = append(*topo, curr)
	visited[curr] = TEMP
	neighbors := adj[curr]
	for _, neighbor := range neighbors {
		// visited[neighbor] = TEMP
		dfs(adj, visited, topo, neighbor)
		// visited[neighbor] = 0
	}
	visited[curr] = PERM
}
