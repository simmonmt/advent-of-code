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
	"testing"

	"github.com/simmonmt/aoc/2019/common/pos"
)

func parseBoardFromStrings(strs []string) *Board {
	b := NewBoard()
	for y := range strs {
		for x, r := range strs[y] {
			p := pos.P2{x, y}
			b.Set(p, r)
		}
	}
	return b
}

func TestFindIntersections(t *testing.T) {
	boardStrs := []string{
		"..#..........",
		"..#..........",
		"#######...###",
		"#.#...#...#.#",
		"#############",
		"..#...#...#..",
		"..#####...#..",
	}

	b := parseBoardFromStrings(boardStrs)

	want := map[pos.P2]bool{
		pos.P2{2, 2}:  true,
		pos.P2{2, 4}:  true,
		pos.P2{6, 4}:  true,
		pos.P2{10, 4}: true,
	}

	if got := FindIntersections(b); !reflect.DeepEqual(want, got) {
		t.Errorf("FindIntersections = %v, want %v", got, want)
	}
}

func TestSumAlignmentParams(t *testing.T) {
	ps := map[pos.P2]bool{
		pos.P2{2, 2}:  true,
		pos.P2{2, 4}:  true,
		pos.P2{6, 4}:  true,
		pos.P2{10, 4}: true,
	}

	if want, got := 76, SumAlignmentParams(ps); want != got {
		t.Errorf("SumAlignmentParams(%v) = %v, want %v", ps, got, want)
	}
}
