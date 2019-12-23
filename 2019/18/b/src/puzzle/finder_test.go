package puzzle

import (
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestFindAllPaths(t *testing.T) {
	type TestCase struct {
		Map           []string
		ExpectedPaths map[string][]Path
	}

	testCases := []TestCase{
		TestCase{
			Map: map1,
			ExpectedPaths: map[string][]Path{
				"@": []Path{
					Path{"a", 2, []string{}},
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
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b, start := NewBoard(tc.Map)
			if got := FindAllPaths(b, start); !reflect.DeepEqual(got, tc.ExpectedPaths) {
				t.Errorf("FindAllPaths(_, %v) = %v, %v, want %v", start, got, tc.ExpectedPaths)
			}
		})
	}
}

var (
	box = []string{
		"#####",
		"#b.A#",
		"#C#.#",
		"#..@#",
		"#####",
	}
)

func TestFindAllPathsToKey(t *testing.T) {
	b, start := NewBoard(box)

	expected := []Path{
		Path{Dest: "b", Dist: 4, Doors: []string{"A"}},
		Path{Dest: "b", Dist: 4, Doors: []string{"C"}},
	}

	paths := FindAllPathsToKey(b, start, "b")
	sort.Sort(PathsByDest(paths))

	if !reflect.DeepEqual(paths, expected) {
		t.Errorf("FindAllPathsToKey b = %v, want %v", paths, expected)
	}
}
