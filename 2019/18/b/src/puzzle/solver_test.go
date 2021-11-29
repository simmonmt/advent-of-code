// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package puzzle

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/simmonmt/aoc/2019/common/pos"
)

func makeKeyMap(keys ...string) map[string]bool {
	out := map[string]bool{}
	for _, key := range keys {
		out[key] = true
	}
	return out
}

func TestNode(t *testing.T) {
	type TestCase struct {
		posns []pos.P2
		keys  map[string]bool
		str   string
	}

	testCases := []TestCase{
		TestCase{
			posns: []pos.P2{pos.P2{1, 1}},
			keys:  makeKeyMap(),
			str:   "1,1_",
		},
		TestCase{
			posns: []pos.P2{pos.P2{1, 1}, pos.P2{2, 2}},
			keys:  makeKeyMap("a", "b"),
			str:   "1,1_2,2_a,b",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := nodeToString(tc.posns, tc.keys); got != tc.str {
				t.Errorf(`nodeToString(%v, %v) = "%v", want "%v"`,
					tc.posns, tc.keys, got, tc.str)
			}

			if posns, keys := parseNode(tc.str); !reflect.DeepEqual(posns, tc.posns) || !reflect.DeepEqual(keys, tc.keys) {
				t.Errorf(`parseNode("%v") = %v, %v, want %v, %v`,
					tc.str, posns, keys, tc.posns, tc.keys)
			}
		})
	}
}

type AllNeighborsTestCase struct {
	state     *astarState
	initial   string
	neighbors []string
}

func testAllNeighbors(t *testing.T, testCases []AllNeighborsTestCase) {
	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := tc.state.AllNeighbors(tc.initial)
			if !reflect.DeepEqual(got, tc.neighbors) {
				t.Errorf("AllNeighbors(%v) = %v, want %v",
					tc.initial, got, tc.neighbors)
			}
		})
	}
}

func TestAllNeighborsSingle(t *testing.T) {
	board, starts := NewBoard(map1)
	if len(starts) != 1 {
		t.Errorf("wanted %d starts, got %d", 1, len(starts))
		return
	}
	start := starts[0]

	state := &astarState{
		board: board,
		graphs: map[pos.P2]map[string][]Path{
			start: FindAllPaths(board, start),
		},
	}

	testCases := []AllNeighborsTestCase{
		AllNeighborsTestCase{
			state:   state,
			initial: nodeToString([]pos.P2{start}, makeKeyMap()),
			neighbors: []string{
				nodeToString([]pos.P2{board.KeyLoc("a")}, makeKeyMap()),
			},
		},
		AllNeighborsTestCase{
			state:   state,
			initial: nodeToString([]pos.P2{start}, makeKeyMap("a")),
			neighbors: []string{
				nodeToString([]pos.P2{board.KeyLoc("a")}, makeKeyMap("a")),
				nodeToString([]pos.P2{board.KeyLoc("b")}, makeKeyMap("a")),
			},
		},
	}

	testAllNeighbors(t, testCases)
}

func TestAllNeighborsMulti(t *testing.T) {
	board, starts := NewBoard(fourStartMap1)
	graphs := FindAllPathsMulti(board, starts)

	state := &astarState{
		board:  board,
		graphs: graphs,
	}

	testCases := []AllNeighborsTestCase{
		AllNeighborsTestCase{
			state:   state,
			initial: nodeToString(starts, makeKeyMap()),
			neighbors: []string{
				nodeToString([]pos.P2{
					board.KeyLoc("a"), starts[1], starts[2], starts[3],
				}, makeKeyMap()),
			},
		},
		AllNeighborsTestCase{
			state:   state,
			initial: nodeToString(starts, makeKeyMap("c")),
			neighbors: []string{
				nodeToString([]pos.P2{
					board.KeyLoc("a"), starts[1], starts[2], starts[3],
				}, makeKeyMap("c")),
				nodeToString([]pos.P2{
					starts[0], board.KeyLoc("d"), starts[2], starts[3],
				}, makeKeyMap("c")),
			},
		},
		AllNeighborsTestCase{
			// Start with all initial positions except the first
			// robot, which starts on 'a'. This means it'll
			// immediately get 'a', which will allow the fourth
			// robot to move to 'b'.
			state: state,
			initial: nodeToString([]pos.P2{
				board.KeyLoc("a"), starts[1], starts[2], starts[3],
			}, makeKeyMap()),
			neighbors: []string{
				nodeToString([]pos.P2{
					board.KeyLoc("a"), starts[1], starts[2],
					board.KeyLoc("b"),
				}, makeKeyMap("a")),
			},
		},
	}

	testAllNeighbors(t, testCases)
}

func TestFindShortestPathOneStart(t *testing.T) {
	type TestCase struct {
		board        *Board
		start        pos.P2
		graph        map[string][]Path
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
				tc.board, tc.graph, tc.start)
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
		expectedCost int
	}

	makeTestCase := func(lines []string, expectedCost int) TestCase {
		board, starts := NewBoard(lines)
		graphs := FindAllPathsMulti(board, starts)

		return TestCase{
			board:        board,
			starts:       starts,
			graphs:       graphs,
			expectedCost: expectedCost,
		}
	}

	testCases := []TestCase{
		makeTestCase(fourStartMap1, 8),
		makeTestCase(fourStartMap2, 24),
		makeTestCase(fourStartMap3, 32),
		makeTestCase(fourStartMap4, 72),
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			path, cost := FindShortestPathMulti(
				tc.board, tc.graphs, tc.starts)
			if cost != tc.expectedCost {
				t.Errorf("FindShortestPathMultiStart() = %v, %v, want cost %v",
					path, cost, tc.expectedCost)
			}
		})
	}
}
