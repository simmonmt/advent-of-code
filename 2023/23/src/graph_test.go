package main

import (
	"sort"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/simmonmt/aoc/2023/common/pos"
)

var smallGraph = []string{
	"#.###", // 0
	"#...#", // 1
	"#.#.#", // 2
	"#.>.#", // 3
	"#.#.#", // 4
	"#.<.#", // 5
	"###.#", // 6
}

func p2(x, y int) pos.P2 {
	return pos.P2{X: x, Y: y}
}

func edge(x, y, dist int) Edge {
	return Edge{End: pos.P2{X: x, Y: y}, Dist: dist}
}

func normalizeGraph(g *Graph) {
	for _, dests := range g.Edges {
		sort.Slice(dests, func(i, j int) bool {
			return dests[i].End.LessThan(dests[j].End)
		})
	}
}

func TestGraph(t *testing.T) {
	type TestCase struct {
		in                 []string
		restrict           bool
		want, wantRestrict *Graph
	}

	testCases := []TestCase{
		TestCase{
			in:       smallGraph,
			restrict: true,
			want: &Graph{
				Start: p2(1, 0),
				End:   p2(3, 6),
				Intersections: map[pos.P2]bool{
					p2(1, 0): true,
					p2(1, 1): true,
					p2(1, 3): true,
					p2(3, 3): true,
					p2(3, 5): true,
					p2(3, 6): true,
				},
				Edges: map[pos.P2][]Edge{
					p2(1, 0): []Edge{edge(1, 1, 1)},
					p2(1, 1): []Edge{edge(1, 0, 1), edge(1, 3, 2), edge(3, 3, 4)},
					p2(1, 3): []Edge{edge(1, 1, 2), edge(3, 3, 2)},
					p2(3, 3): []Edge{edge(1, 1, 4), edge(3, 5, 2)},
					p2(3, 5): []Edge{edge(1, 3, 4), edge(3, 3, 2), edge(3, 6, 1)},
					p2(3, 6): []Edge{edge(3, 5, 1)},
				},
			},
		},
		TestCase{
			in:       smallGraph,
			restrict: false,
			want: &Graph{
				Start: p2(1, 0),
				End:   p2(3, 6),
				Intersections: map[pos.P2]bool{
					p2(1, 0): true,
					p2(1, 1): true,
					p2(1, 3): true,
					p2(3, 3): true,
					p2(3, 5): true,
					p2(3, 6): true,
				},
				Edges: map[pos.P2][]Edge{
					p2(1, 0): []Edge{edge(1, 1, 1)},
					p2(1, 1): []Edge{edge(1, 0, 1), edge(1, 3, 2), edge(3, 3, 4)},
					p2(1, 3): []Edge{edge(1, 1, 2), edge(3, 3, 2), edge(3, 5, 4)},
					p2(3, 3): []Edge{edge(1, 1, 4), edge(1, 3, 2), edge(3, 5, 2)},
					p2(3, 5): []Edge{edge(1, 3, 4), edge(3, 3, 2), edge(3, 6, 1)},
					p2(3, 6): []Edge{edge(3, 5, 1)},
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			board, err := parseInput(tc.in)
			if err != nil {
				t.Fatalf("bad graph: %v", err)
			}

			got, err := BuildGraph(board, tc.restrict)
			if err != nil {
				t.Fatalf("failed to build graph: %v", err)
			}

			normalizeGraph(got)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("graphs mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
