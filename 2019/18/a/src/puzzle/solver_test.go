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
	"reflect"
	"strconv"
	"testing"
)

func TestAllNeighbors(t *testing.T) {
	state := &astarState{
		graph: map[string][]Path{
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
		numKeys: 2,
	}

	cases := map[string][]string{
		"@_":  []string{"a_"},
		"@_a": []string{"a_a", "b_a"},
	}

	for start, want := range cases {
		if got := state.AllNeighbors(start); !reflect.DeepEqual(got, want) {
			t.Errorf("AllNeighbors(%v) = %v, want %v", start, got, want)
		}
	}

}

func TestFindShortestPath(t *testing.T) {
	type TestCase struct {
		graph        map[string][]Path
		numKeys      int
		expectedCost int
	}

	makeTestCase := func(lines []string, expectedCost int) TestCase {
		b, start := NewBoard(lines)
		graph := FindAllPaths(b, start)

		return TestCase{
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

			path, cost := FindShortestPath(tc.graph, tc.numKeys, "@")
			if cost != tc.expectedCost {
				t.Errorf("FindShortestPath() = %v, %v, want cost %v",
					path, cost, tc.expectedCost)
			}
		})
	}
}
