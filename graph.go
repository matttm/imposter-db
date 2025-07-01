package main

func topologicalSort(tables [][2]string) []string {
	var adj map[string][]string
	for _, edge := range tables {
		p := edge[0]
		c := edge[1]
		adj[p] = append(adj[p], c)
	}
	topoSorted := []string{}
	return topoSorted
}

func dfs(adj map[string][]string, topo []string, curr string) {
}
