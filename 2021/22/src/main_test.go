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

package main

import (
	"os"
	"reflect"
	"sort"
	"strconv"
	"testing"

	"github.com/simmonmt/aoc/2021/common/logger"
)

func lessArray(a1, a2 []int) bool {
	for i := 0; i < len(a1); i++ {
		if a1[i] < a2[i] {
			return true
		} else if a1[i] > a2[i] {
			return false
		}
	}
	return false
}

func sortCubes(cubes []*Cube) {
	sort.Slice(cubes, func(i, j int) bool {
		c1, c2 := cubes[i], cubes[j]
		return lessArray(
			[]int{c1.xLo, c1.xHi, c1.yLo, c1.yHi, c1.zLo, c1.zHi},
			[]int{c2.xLo, c2.xHi, c2.yLo, c2.yHi, c2.zLo, c2.zHi},
		)
	})
}

func TestCubeOverlaps(t *testing.T) {
	type TestCase struct {
		a, b Cube
		want bool
	}

	testCases := []TestCase{
		TestCase{
			a:    Cube{-5, -3, 2, 0, 5, -3},
			b:    Cube{10, 12, 10, 12, 10, 12},
			want: false,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := tc.a.Overlaps(&tc.b); got != tc.want {
				t.Errorf("overlap %v %v got %v, want %v",
					tc.a, tc.b, got, tc.want)
			}
			if got := tc.b.Overlaps(&tc.a); got != tc.want {
				t.Errorf("overlap %v %v got %v, want %v",
					tc.b, tc.a, got, tc.want)
			}
		})
	}

}

func TestCubeSize(t *testing.T) {
	type TestCase struct {
		in   *Cube
		want int64
	}

	testCases := []TestCase{
		TestCase{
			in:   &Cube{5, 9, 2, 8, 1, 7},
			want: 5 * 7 * 7,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := tc.in.Size(); got != tc.want {
				t.Errorf("Size(%v) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

func TestCubeSub(t *testing.T) {
	type TestCase struct {
		in   *Cube
		sub  *Cube
		want []*Cube
	}

	testCases := []TestCase{
		TestCase{
			in:  &Cube{5, 9, 2, 8, 1, 7},
			sub: &Cube{6, 7, 3, 5, 4, 6},
			want: []*Cube{
				&Cube{5, 9, 2, 8, 7, 7},
				&Cube{5, 9, 2, 8, 1, 3},
				&Cube{5, 5, 2, 8, 4, 6},
				&Cube{8, 9, 2, 8, 4, 6},
				&Cube{6, 7, 2, 2, 4, 6},
				&Cube{6, 7, 6, 8, 4, 6},
			},
		},
		TestCase{
			in:  &Cube{5, 9, 2, 8, 1, 7},
			sub: &Cube{5, 7, 2, 5, 1, 6},
			want: []*Cube{
				&Cube{5, 9, 2, 8, 7, 7},
				&Cube{8, 9, 2, 8, 1, 6},
				&Cube{5, 7, 6, 8, 1, 6},
			},
		},
		TestCase{
			in:  &Cube{5, 9, 2, 8, 1, 7},
			sub: &Cube{6, 9, 3, 8, 4, 7},
			want: []*Cube{
				&Cube{5, 9, 2, 8, 1, 3}, // bottom
				&Cube{5, 5, 2, 8, 4, 7}, // left
				&Cube{6, 9, 2, 2, 4, 7}, // front
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := tc.in.Sub(tc.sub)
			sortCubes(got)
			sortCubes(tc.want)

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Sub = %v, want %v", got, tc.want)
			}

			gotSz := int64(0)
			for _, c := range got {
				gotSz += c.Size()
			}

			wantSz := tc.in.Size() - tc.sub.Size()
			if gotSz != wantSz {
				t.Errorf("got sz = %v, want sz %v", gotSz, wantSz)
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
