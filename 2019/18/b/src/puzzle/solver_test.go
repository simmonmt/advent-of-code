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
		b, starts := NewBoard(lines)
		graph := FindAllPaths(b, starts[0])

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
