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
