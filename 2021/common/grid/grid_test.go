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

package grid

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/simmonmt/aoc/2021/common/pos"
)

func TestSimple(t *testing.T) {
	g := New(5, 6)

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

	if got := g.Get(p); !reflect.DeepEqual(got, value1) {
		t.Errorf("g.Get(%v) = %v, want %v", p, got, value1)
	}
	if got := g.Get(rp); got != nil {
		t.Errorf("g.Get(%v) = %v, want nil", rp, got)
	}
	if got := g.Get(op); !reflect.DeepEqual(got, value2) {
		t.Errorf("g.Get(%v) = %v, want %v", op, got, value2)
	}
}

func TestWalk(t *testing.T) {
	g := New(3, 2)

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

	g.Walk(func(p pos.P2, v interface{}) {
		gotPos = append(gotPos, p)
		gotVals = append(gotVals, v.(string))
	})

	if !reflect.DeepEqual(gotPos, wantPos) {
		t.Errorf("Walk gotPos = %v, want %v", gotPos, wantPos)
	}
	if !reflect.DeepEqual(gotVals, wantVals) {
		t.Errorf("Walk gotVals = %v, want %v", gotVals, wantVals)
	}
}

func TestAllNeighbors(t *testing.T) {
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

	b := New(3, 3)

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
