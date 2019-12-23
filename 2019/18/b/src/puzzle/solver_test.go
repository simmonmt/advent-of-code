package puzzle

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/simmonmt/aoc/2019/common/pos"
)

func TestAllNeighbors(t *testing.T) {
	board, starts := NewBoard(map1)
	if len(starts) != 1 {
		t.Errorf("wanted %d starts, got %d", 1, len(starts))
		return
	}
	start := starts[0]

	state := &astarState{
		board: board,

		graphs: map[pos.P2]map[string][]Path{
			start: map[string][]Path{
				"@": []Path{
					Path{"a", 2, nil},
					Path{"b", 4, []string{"A"}},
				},
				"a": []Path{
					Path{"b", 6, []string{"A"}},
				},
				"b": []Path{
					Path{"a", 6, []string{"A"}},
				},
			},
		},
		numKeys: 2,
	}

	type TestCase struct {
		initial   string
		neighbors []string
	}

	makeKeyMap := func(keys ...string) map[string]bool {
		out := map[string]bool{}
		for _, key := range keys {
			out[key] = true
		}
		return out
	}

	testCases := []TestCase{
		TestCase{
			initial: nodeToString([]pos.P2{start}, makeKeyMap()),
			neighbors: []string{
				nodeToString([]pos.P2{board.KeyLoc("a")}, makeKeyMap()),
			},
		},
		TestCase{
			initial: nodeToString([]pos.P2{start}, makeKeyMap("a")),
			neighbors: []string{
				nodeToString([]pos.P2{board.KeyLoc("a")}, makeKeyMap("a")),
				nodeToString([]pos.P2{board.KeyLoc("b")}, makeKeyMap("a")),
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := state.AllNeighbors(tc.initial)
			if !reflect.DeepEqual(got, tc.neighbors) {
				t.Errorf("AllNeighbors(%v) = %v, want %v",
					tc.initial, got, tc.neighbors)
			}
		})
	}
}

func TestFindShortestPathOneStart(t *testing.T) {
	type TestCase struct {
		board        *Board
		start        pos.P2
		graph        map[string][]Path
		numKeys      int
		expectedCost int
	}

	makeTestCase := func(lines []string, expectedCost int) TestCase {
		b, starts := NewBoard(lines)
		if len(starts) != 1 {
			panic(fmt.Sprintf("found %d starts, wanted 1", len(starts)))
		}
		start := starts[0]

		graph := FindAllPaths(b, start)

		return TestCase{
			board:        b,
			start:        start,
			graph:        graph,
			numKeys:      len(b.Keys()),
			expectedCost: expectedCost,
		}
	}

	testCases := []TestCase{
		makeTestCase(map1, 8),
		makeTestCase(map2, 86),
		makeTestCase(map3, 136),
		makeTestCase(map4, 81),
		makeTestCase(map5, 132),
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			path, cost := FindShortestPath(
				tc.board, tc.graph, tc.numKeys, tc.start)
			if cost != tc.expectedCost {
				t.Errorf("FindShortestPath() = %v, %v, want cost %v",
					path, cost, tc.expectedCost)
			}
		})
	}
}

func TestFindShortestPathMultiStart(t *testing.T) {
	type TestCase struct {
		board        *Board
		starts       []pos.P2
		graphs       map[pos.P2]map[string][]Path
		numKeys      int
		expectedCost int
	}

	makeTestCase := func(lines []string, expectedCost int) TestCase {
		b, starts := NewBoard(lines)

		graphs := map[pos.P2]map[string][]Path{}
		for _, start := range starts {
			graphs[start] = FindAllPaths(b, start)
		}

		return TestCase{
			board:        b,
			starts:       starts,
			graphs:       graphs,
			numKeys:      len(b.Keys()),
			expectedCost: expectedCost,
		}
	}

	testCases := []TestCase{
		makeTestCase(fourStartMap1, 0), //8),
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			path, cost := FindShortestPathMultiStart(
				tc.board, tc.graphs, tc.numKeys, tc.starts)
			if cost != tc.expectedCost {
				t.Errorf("FindShortestPathMultiStart() = %v, %v, want cost %v",
					path, cost, tc.expectedCost)
			}
		})
	}
}
