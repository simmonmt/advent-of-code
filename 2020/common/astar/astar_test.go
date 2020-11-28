package astar

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/simmonmt/aoc/2020/common/logger"
)

type helperNode struct {
	distances map[string]uint
}

type aStarHelper struct {
	nodes map[string]helperNode
}

func (h *aStarHelper) AllNeighbors(start string) []string {
	node, found := h.nodes[start]
	if !found {
		return nil
	}

	neighbors := []string{}
	for neighbor := range node.distances {
		neighbors = append(neighbors, neighbor)
	}
	return neighbors
}

func (h *aStarHelper) EstimateDistance(start, end string) uint {
	if start == end {
		return 0
	} else {
		return 1
	}
}

func (h *aStarHelper) NeighborDistance(n1, n2 string) uint {
	if n1 == n2 {
		return 0
	}

	node, found := h.nodes[n1]
	if !found {
		return 0
	}

	for neighbor, dist := range node.distances {
		if neighbor == n2 {
			return dist
		}
	}

	panic(fmt.Sprintf("no distance for %v to %v", n1, n2))
}

func (h *aStarHelper) GoalReached(cand, goal string) bool {
	return cand == goal
}

func TestAStar(t *testing.T) {
	helper := aStarHelper{
		nodes: map[string]helperNode{
			"start": helperNode{distances: map[string]uint{"a": 15, "d": 20}},
			"a":     helperNode{distances: map[string]uint{"start": 15, "b": 20}},
			"b":     helperNode{distances: map[string]uint{"a": 20, "c": 30}},
			"c":     helperNode{distances: map[string]uint{"b": 30, "end": 40}},
			"d":     helperNode{distances: map[string]uint{"start": 20, "d1": 2, "e": 20}},
			"d1":    helperNode{distances: map[string]uint{"d": 2, "d2": 1}},
			"d2":    helperNode{distances: map[string]uint{"d1": 1, "d3": 1}},
			"d3":    helperNode{distances: map[string]uint{"d2": 1, "e": 1}},
			"e":     helperNode{distances: map[string]uint{"d": 20, "d3": 1, "end": 20}},
			"end":   helperNode{distances: map[string]uint{"c": 40, "e": 20}},
		},
	}

	result := AStar("start", "end", &helper)
	expected := []string{"end", "e", "d3", "d2", "d1", "d", "start"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("start->end, got %v, want %v", result, expected)
	}
}

func TestFScoreMap(t *testing.T) {
	m := newFScoreMap()
	m.Set("a", 6)
	m.Set("b", 2)
	m.Set("c", 1)
	m.Set("d", 1)
	m.Set("e", 1)

	found := []string{}
	m.Walk(func(item *fScoreItem) bool {
		found = append(found, fmt.Sprintf("%v:%v", item.Name, item.Value))
		return true
	})

	expected := []string{"c:1", "d:1", "e:1", "b:2", "a:6"}
	if !reflect.DeepEqual(expected, found) {
		t.Errorf("got %v, wanted %v", found, expected)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	logger.Init(true)

	os.Exit(m.Run())
}
