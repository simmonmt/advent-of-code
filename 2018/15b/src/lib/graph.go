package lib

import (
	"github.com/soniakeys/graph"
)

type Graph struct {
	al graph.LabeledAdjacencyList
}

func NewGraph(sz int) *Graph {
	return &Graph{
		al: make([][]graph.Half, sz),
	}
}

type GraphEdge struct {
	From, To graph.NI
}

func (g *Graph) AllEdges() []GraphEdge {
	outMap := map[GraphEdge]bool{}

	addEdge := func(from, to graph.NI) {
		if from < to {
			outMap[GraphEdge{from, to}] = true
		} else {
			outMap[GraphEdge{to, from}] = true
		}
	}

	for i := range g.al {
		from := graph.NI(i)
		for _, h := range g.al[from] {
			addEdge(from, h.To)
		}
	}

	out := make([]GraphEdge, len(outMap))
	i := 0
	for e := range outMap {
		out[i] = e
		i++
	}
	return out
}

func (g *Graph) HasEdge(from, to graph.NI) bool {
	return g.hasDirectedEdge(from, to) &&
		g.hasDirectedEdge(to, from)
}

func (g *Graph) hasDirectedEdge(from, to graph.NI) bool {
	for _, h := range g.al[from] {
		if h.To == to {
			return true
		}
	}
	return false
}

func (g *Graph) AddEdge(from, to graph.NI) {
	g.addDirectedEdge(from, to)
	g.addDirectedEdge(to, from)
}

func (g *Graph) addDirectedEdge(from, to graph.NI) {
	if g.al[int(from)] == nil {
		g.al[int(from)] = []graph.Half{}
	} else {
		for _, h := range g.al[int(from)] {
			if h.To == to {
				panic("edge already exists")
			}
		}
	}

	g.al[int(from)] = append(g.al[int(from)], graph.Half{To: to})
}

func (g *Graph) RemoveEdge(from, to graph.NI) {
	g.removeDirectedEdge(from, to)
	g.removeDirectedEdge(to, from)
}

func (g *Graph) removeDirectedEdge(from, to graph.NI) {
	row := g.al[int(from)]
	if row == nil {
		panic("no source for edge to remove")
	}

	filtered := []graph.Half{}
	found := false
	for _, h := range row {
		if h.To == to {
			if found {
				panic("double entry")
			}
			found = true
		} else {
			filtered = append(filtered, h)
		}
	}
	g.al[int(from)] = filtered

	if !found {
		panic("remove unfound")
	}
}

func (g *Graph) ShortestPath(from, to graph.NI, dist func(from, to graph.NI) int) []graph.NI {
	heuristic := func(cand graph.NI) float64 {
		return float64(dist(cand, to))
	}

	weight := func(cand graph.LI) float64 { return 1 }

	fullPath, _ := g.al.AStarAPath(from, to, heuristic, weight)
	if len(fullPath.Path) == 0 {
		return nil
	}

	path := []graph.NI{}
	for _, n := range fullPath.Path {
		path = append(path, n.To)
	}
	return path
}
