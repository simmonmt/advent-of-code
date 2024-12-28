// Copyright 2024 Google LLC
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

package pos

import (
	"reflect"
	"strconv"
	"testing"
)

func TestP2FromStringOK(t *testing.T) {
	want := P2{1, -2}
	if got, err := P2FromString("1,-2"); err != nil || !got.Equals(want) {
		t.Errorf(`P2FromString("1,-2") = %v, %v, want %v, nil`,
			got, err, want)
	}
}

func TestP2FromStringBad(t *testing.T) {
	if _, err := P2FromString("1,bob"); err == nil {
		t.Errorf(`P2FromString("1,bob") = _, %v, want _, non-nil`, err)
	}
}

func TestP2(t *testing.T) {
	type TestCase struct {
		p1, p2            P2
		lessThan          bool
		manhattanDistance int
	}

	testCases := []TestCase{
		TestCase{P2{1, 1}, P2{4, 2}, true, 4},
		TestCase{P2{1, 1}, P2{1, 2}, true, 1},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := tc.p1.LessThan(tc.p2); got != tc.lessThan {
				t.Errorf("want %v < %v, got %v", tc.p1, tc.p2, got)
			}
			if got := tc.p2.LessThan(tc.p1); got == tc.lessThan {
				t.Errorf("want %v < %v, got %v", tc.p2, tc.p1, got)
			}
			if got := tc.p1.Cmp(tc.p1); got != 0 {
				t.Errorf("want %v <=> %v = 0, got %v", tc.p1, tc.p1, got)
			}

			if got := tc.p1.ManhattanDistance(tc.p2); got != tc.manhattanDistance {
				t.Errorf("manhattan distance %v, %v = %v, want %v",
					tc.p1, tc.p2, got, tc.manhattanDistance)
			}
			if got := tc.p2.ManhattanDistance(tc.p1); got != tc.manhattanDistance {
				t.Errorf("manhattan distance %v, %v = %v, want %v",
					tc.p2, tc.p1, got, tc.manhattanDistance)
			}

			p := tc.p1
			p.Add(P2{2, 1})
			want := P2{tc.p1.X + 2, tc.p1.Y + 1}
			if !p.Equals(want) {
				t.Errorf("%v.Add({2,1}) = %v, want %v",
					tc.p1, p, want)
			}
		})
	}
}

func TestWalkP2(t *testing.T) {
	saw := []P2{}

	WalkP2(2, 3, func(p P2) { saw = append(saw, p) })

	expected := []P2{
		P2{Y: 0, X: 0},
		P2{Y: 0, X: 1},
		P2{Y: 0, X: 2},
		P2{Y: 1, X: 0},
		P2{Y: 1, X: 1},
		P2{Y: 1, X: 2},
	}

	if !reflect.DeepEqual(saw, expected) {
		t.Errorf("WalkP2 saw %+v, want %+v", saw, expected)
	}
}

func TestP3FromStringOK(t *testing.T) {
	want := P3{1, -2, 3}
	if got, err := P3FromString("1,-2,3"); err != nil || !got.Equals(want) {
		t.Errorf(`P3FromString("1,-2,3") = %v, %v, want %v, nil`,
			got, err, want)
	}
}

func TestP3FromStringBad(t *testing.T) {
	if _, err := P3FromString("1,bob,3"); err == nil {
		t.Errorf(`P3FromString("1,bob,3") = _, %v, want _, non-nil`, err)
	}
}

func TestP3AllNeighbors(t *testing.T) {
	p := P3{0, 0, 0}

	want := map[P3]bool{
		// x=-1
		P3{-1, -1, -1}: true,
		P3{-1, -1, 0}:  true,
		P3{-1, -1, 1}:  true,

		P3{-1, 0, -1}: true,
		P3{-1, 0, 0}:  true,
		P3{-1, 0, 1}:  true,

		P3{-1, 1, -1}: true,
		P3{-1, 1, 0}:  true,
		P3{-1, 1, 1}:  true,

		// x=0
		P3{0, -1, -1}: true,
		P3{0, -1, 0}:  true,
		P3{0, -1, 1}:  true,

		P3{0, 0, -1}: true,
		//P3{0,0,0}: true,
		P3{0, 0, 1}: true,

		P3{0, 1, -1}: true,
		P3{0, 1, 0}:  true,
		P3{0, 1, 1}:  true,

		// x=1
		P3{1, -1, -1}: true,
		P3{1, -1, 0}:  true,
		P3{1, -1, 1}:  true,

		P3{1, 0, -1}: true,
		P3{1, 0, 0}:  true,
		P3{1, 0, 1}:  true,

		P3{1, 1, -1}: true,
		P3{1, 1, 0}:  true,
		P3{1, 1, 1}:  true,
	}

	got := map[P3]bool{}
	for _, n := range p.AllNeighbors(true) {
		got[n] = true
	}

	if reflect.DeepEqual(got, want) {
		return
	}

	for n := range got {
		if _, found := want[n]; !found {
			t.Errorf("got has %v, not in want", n)
		}
	}

	for n := range want {
		if _, found := got[n]; !found {
			t.Errorf("want has %v, not in got", n)
		}
	}

	t.Errorf("got want mismatch, can't figure out why")
}

func TestP4AllNeighbors(t *testing.T) {
	p := P4{0, 0, 0, 0}

	want := map[P4]bool{
		P4{-1, -1, -1, -1}: true,
		P4{-1, -1, -1, 0}:  true,
		P4{-1, -1, -1, 1}:  true,
		P4{-1, -1, 0, -1}:  true,
		P4{-1, -1, 0, 0}:   true,
		P4{-1, -1, 0, 1}:   true,
		P4{-1, -1, 1, -1}:  true,
		P4{-1, -1, 1, 0}:   true,
		P4{-1, -1, 1, 1}:   true,
		P4{-1, 0, -1, -1}:  true,
		P4{-1, 0, -1, 0}:   true,
		P4{-1, 0, -1, 1}:   true,
		P4{-1, 0, 0, -1}:   true,
		P4{-1, 0, 0, 0}:    true,
		P4{-1, 0, 0, 1}:    true,
		P4{-1, 0, 1, -1}:   true,
		P4{-1, 0, 1, 0}:    true,
		P4{-1, 0, 1, 1}:    true,
		P4{-1, 1, -1, -1}:  true,
		P4{-1, 1, -1, 0}:   true,
		P4{-1, 1, -1, 1}:   true,
		P4{-1, 1, 0, -1}:   true,
		P4{-1, 1, 0, 0}:    true,
		P4{-1, 1, 0, 1}:    true,
		P4{-1, 1, 1, -1}:   true,
		P4{-1, 1, 1, 0}:    true,
		P4{-1, 1, 1, 1}:    true,

		P4{0, -1, -1, -1}: true,
		P4{0, -1, -1, 0}:  true,
		P4{0, -1, -1, 1}:  true,
		P4{0, -1, 0, -1}:  true,
		P4{0, -1, 0, 0}:   true,
		P4{0, -1, 0, 1}:   true,
		P4{0, -1, 1, -1}:  true,
		P4{0, -1, 1, 0}:   true,
		P4{0, -1, 1, 1}:   true,
		P4{0, 0, -1, -1}:  true,
		P4{0, 0, -1, 0}:   true,
		P4{0, 0, -1, 1}:   true,
		P4{0, 0, 0, -1}:   true,
		//P4{0, 0, 0, 0}:    true,
		P4{0, 0, 0, 1}:   true,
		P4{0, 0, 1, -1}:  true,
		P4{0, 0, 1, 0}:   true,
		P4{0, 0, 1, 1}:   true,
		P4{0, 1, -1, -1}: true,
		P4{0, 1, -1, 0}:  true,
		P4{0, 1, -1, 1}:  true,
		P4{0, 1, 0, -1}:  true,
		P4{0, 1, 0, 0}:   true,
		P4{0, 1, 0, 1}:   true,
		P4{0, 1, 1, -1}:  true,
		P4{0, 1, 1, 0}:   true,
		P4{0, 1, 1, 1}:   true,

		P4{1, -1, -1, -1}: true,
		P4{1, -1, -1, 0}:  true,
		P4{1, -1, -1, 1}:  true,
		P4{1, -1, 0, -1}:  true,
		P4{1, -1, 0, 0}:   true,
		P4{1, -1, 0, 1}:   true,
		P4{1, -1, 1, -1}:  true,
		P4{1, -1, 1, 0}:   true,
		P4{1, -1, 1, 1}:   true,
		P4{1, 0, -1, -1}:  true,
		P4{1, 0, -1, 0}:   true,
		P4{1, 0, -1, 1}:   true,
		P4{1, 0, 0, -1}:   true,
		P4{1, 0, 0, 0}:    true,
		P4{1, 0, 0, 1}:    true,
		P4{1, 0, 1, -1}:   true,
		P4{1, 0, 1, 0}:    true,
		P4{1, 0, 1, 1}:    true,
		P4{1, 1, -1, -1}:  true,
		P4{1, 1, -1, 0}:   true,
		P4{1, 1, -1, 1}:   true,
		P4{1, 1, 0, -1}:   true,
		P4{1, 1, 0, 0}:    true,
		P4{1, 1, 0, 1}:    true,
		P4{1, 1, 1, -1}:   true,
		P4{1, 1, 1, 0}:    true,
		P4{1, 1, 1, 1}:    true,
	}

	got := map[P4]bool{}
	for _, n := range p.AllNeighbors() {
		got[n] = true
	}

	if reflect.DeepEqual(got, want) {
		return
	}

	numMismatches := 0
	for n := range got {
		if _, found := want[n]; !found {
			t.Errorf("got has %v, not in want", n)
			numMismatches++
		}
	}

	for n := range want {
		if _, found := got[n]; !found {
			t.Errorf("want has %v, not in got", n)
			numMismatches++
		}
	}

	t.Errorf("got want mismatchs: %v", numMismatches)
}
