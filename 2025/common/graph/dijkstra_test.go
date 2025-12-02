package graph

import (
	"fmt"
	"reflect"
	"testing"
)

type testGraph struct {
	// Map of node IDs to edge+costs
	nodes map[NodeID]map[NodeID]int
}

func (g *testGraph) Neighbors(id NodeID) []NodeID {
	edges, found := g.nodes[id]
	if !found {
		panic("bad from")
	}

	out := []NodeID{}
	for dest, _ := range edges {
		out = append(out, dest)
	}
	return out
}

func (g *testGraph) NeighborDistance(from, to NodeID) int {
	edges, found := g.nodes[from]
	if !found {
		panic("bad from")
	}

	for dest, cost := range edges {
		if dest == to {
			return cost
		}
	}

	panic("bad to")
}

func TestShortestPath(t *testing.T) {
	g := &testGraph{
		nodes: map[NodeID]map[NodeID]int{
			"1": map[NodeID]int{"2": 5, "3": 15},
			"2": map[NodeID]int{"1": 5, "3": 6},
			"3": map[NodeID]int{"1": 15, "2": 6, "4": 2},
			"4": map[NodeID]int{"3": 2},
			"5": map[NodeID]int{},
		},
	}

	type TestCase struct {
		start, end NodeID
		want       []NodeID
	}

	testCases := []TestCase{
		TestCase{"1", "2", []NodeID{"2"}},
		TestCase{"1", "3", []NodeID{"2", "3"}},
		TestCase{"1", "4", []NodeID{"2", "3", "4"}},
		TestCase{"1", "5", nil},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v-%v", tc.start, tc.end), func(t *testing.T) {
			if got := ShortestPath(tc.start, tc.end, g); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ShortestPath(%v,%v,g) = %v, want %v",
					tc.start, tc.end, got, tc.want)
			}

		})
	}
}
