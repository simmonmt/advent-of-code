// Copyright 2022 Google LLC
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

package grid

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/simmonmt/aoc/2022/common/pos"
)

func TestGrid(t *testing.T) {
	g := New[string](5, 6)

	if got := g.Width(); got != 5 {
		t.Errorf("g.Width() = %v, want 5", got)
	}

	if got := g.Height(); got != 6 {
		t.Errorf("g.Height() = %v, want 6", got)
	}

	p := pos.P2{1, 2}
	rp := pos.P2{2, 1}
	op := pos.P2{3, 4}
	value1 := "testvalue1"
	value2 := "testvalue2"

	g.Set(p, value1)
	g.Set(op, value2)

	if got := g.Get(p); got != value1 {
		t.Errorf("g.Get(%v) = %v, want %v", p, got, value1)
	}
	if got := g.Get(rp); got != "" {
		t.Errorf("g.Get(%v) = '%v', want ''", rp, got)
	}
	if got := g.Get(op); got != value2 {
		t.Errorf("g.Get(%v) = %v, want %v", op, got, value2)
	}

	if got := g.IsValid(pos.P2{1, 2}); !got {
		t.Errorf("g.IsValid({1,2}) = %v, want true", got)
	}
	if got := g.IsValid(pos.P2{99, 99}); got {
		t.Errorf("g.IsValid({99,99}) = %v, want false", got)
	}
}

func TestGridWalk(t *testing.T) {
	g := New[string](3, 2)

	for y := 0; y < 2; y++ {
		for x := 0; x < 3; x++ {
			p := pos.P2{X: x, Y: y}
			g.Set(p, p.String())
		}
	}

	wantPos := []pos.P2{
		pos.P2{0, 0}, pos.P2{1, 0}, pos.P2{2, 0},
		pos.P2{0, 1}, pos.P2{1, 1}, pos.P2{2, 1},
	}

	wantVals := []string{
		"0,0", "1,0", "2,0",
		"0,1", "1,1", "2,1",
	}

	gotPos := []pos.P2{}
	gotVals := []string{}

	g.Walk(func(p pos.P2, v string) {
		gotPos = append(gotPos, p)
		gotVals = append(gotVals, v)
	})

	if !reflect.DeepEqual(gotPos, wantPos) {
		t.Errorf("Walk gotPos = %v, want %v", gotPos, wantPos)
	}
	if !reflect.DeepEqual(gotVals, wantVals) {
		t.Errorf("Walk gotVals = %v, want %v", gotVals, wantVals)
	}
}

func TestGridAllNeighbors(t *testing.T) {
	type TestCase struct {
		p           pos.P2
		includeDiag bool
		want        []pos.P2
	}

	testCases := []TestCase{
		TestCase{
			pos.P2{0, 0},
			false,
			[]pos.P2{pos.P2{1, 0}, pos.P2{0, 1}},
		},
		TestCase{
			pos.P2{0, 0},
			true,
			[]pos.P2{pos.P2{1, 0}, pos.P2{0, 1}, pos.P2{1, 1}},
		},
		TestCase{
			pos.P2{1, 1},
			false,
			[]pos.P2{
				pos.P2{0, 1}, pos.P2{2, 1}, pos.P2{1, 0},
				pos.P2{1, 2},
			},
		},
		TestCase{
			pos.P2{1, 1},
			true,
			[]pos.P2{
				pos.P2{0, 1}, pos.P2{2, 1}, pos.P2{1, 0},
				pos.P2{1, 2},
				pos.P2{0, 0}, pos.P2{2, 0},
				pos.P2{0, 2}, pos.P2{2, 2},
			},
		},
	}

	b := New[int](3, 3)

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := b.AllNeighbors(tc.p, tc.includeDiag)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("AllNeighbors(%v, %v) = %v, want %v",
					tc.p, tc.includeDiag,
					got, tc.want)
			}
		})
	}
}

func TestSparseGrid(t *testing.T) {
	g := NewSparseGrid()

	p := pos.P2{1, 2}
	rp := pos.P2{2, 1}
	op := pos.P2{3, 4}
	value1 := "testvalue1"
	value2 := "testvalue2"

	g.Set(p, value1)
	g.Set(op, value2)

	if got, found := g.Get(p); !found || !reflect.DeepEqual(got, value1) {
		t.Errorf("g.Get(%v) = %v, %v, want %v, true", p, got, found, value1)
	}
	if got, found := g.Get(rp); found || got != nil {
		t.Errorf("g.Get(%v) = %v, %v, want nil, false", rp, got, found)
	}
	if got, found := g.Get(op); !found || !reflect.DeepEqual(got, value2) {
		t.Errorf("g.Get(%v) = %v, %v, want %v, true", op, got, found, value2)
	}

	if got, want := g.Start(), (pos.P2{1, 2}); !got.Equals(want) {
		t.Errorf("g.Start() = %v, want %v", got, want)
	}
	if got, want := g.End(), (pos.P2{3, 4}); !got.Equals(want) {
		t.Errorf("g.End() = %v, want %v", got, want)
	}
}
