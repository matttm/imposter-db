package main

var PERM uint8 = 7
var TEMP uint8 = 1

func topologicalSort(tables [][2]string) []string {
	n := len(tables)
	var adj map[string][]string
	var visited map[string]uint8
	for _, edge := range tables {
		p := edge[0]
		c := edge[1]
		adj[p] = append(adj[p], c)
	}
	topoSorted := []string{}
	dfs(adj, visited, topo, curr)
	return topoSorted
}

func dfs(adj map[string][]string, visited map[string]uint8, topo []string, curr string) {
	if visited[curr] == PERM {
	}
	if visited[curr] == TEMP {
	}
	neighbors := adj[curr]
	topo = append(topo, curr)
	visited[curr] = TEMP
	for _, neighbor := range neighbors {
		visited[neighbor] = TEMP
		dfs(adj, visited, topo, neighbor)
		visited[neighbor] = 0
	}
	visited[curr] = PERM
}
