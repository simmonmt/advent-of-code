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
	"strconv"
	"testing"

	"github.com/simmonmt/aoc/2021/common/logger"
	"github.com/simmonmt/aoc/2021/common/pos"
)

func TestRotatePos(t *testing.T) {
	type TestCase struct {
		in                  pos.P3
		xRots, yRots, zRots []int
		wants               []pos.P3
	}

	testCases := []TestCase{
		TestCase{
			in:    pos.P3{X: 1, Y: 0, Z: 0},
			yRots: []int{0, 90, 180, 270},
			wants: []pos.P3{
				pos.P3{1, 0, 0},
				pos.P3{0, 0, -1},
				pos.P3{-1, 0, 0},
				pos.P3{0, 0, 1},
			},
		},
		TestCase{
			in:    pos.P3{X: 0, Y: 1, Z: 0},
			xRots: []int{0, 90, 180, 270},
			wants: []pos.P3{
				pos.P3{0, 1, 0},
				pos.P3{0, 0, 1},
				pos.P3{0, -1, 0},
				pos.P3{0, 0, -1},
			},
		},
		TestCase{
			in:    pos.P3{X: 1, Y: 0, Z: 0},
			zRots: []int{0, 90, 180, 270},
			wants: []pos.P3{
				pos.P3{1, 0, 0},
				pos.P3{0, 1, 0},
				pos.P3{-1, 0, 0},
				pos.P3{0, -1, 0},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var mat []int
			var rots []int
			if tc.xRots != nil {
				mat = RX
				rots = tc.xRots
			} else if tc.yRots != nil {
				mat = RY
				rots = tc.yRots
			} else if tc.zRots != nil {
				mat = RZ
				rots = tc.zRots
			} else {
				panic("bad case")
			}

			for i, deg := range rots {
				want := tc.wants[i]
				if got := rotatePos(mat, deg, tc.in); !reflect.DeepEqual(got, want) {
					t.Errorf("rotatePos(%v, %v, %v) = %v, want %v",
						mat, deg, tc.in, got, want)
				}
			}
		})
	}
}

func TestAllOrientationsAreUnique(t *testing.T) {
	f := NewOneScannerField(0, []pos.P3{pos.P3{1, 2, 3}})

	ofs := map[string]int{}
	for i, of := range allOrientations(f) {
		logger.LogF("orientation %d: %v", i, of.Beacons[0])

		ps := of.Beacons[0].String()
		if dup, found := ofs[ps]; found {
			t.Errorf("orientation duplicate %v and %v", i, dup)
		}
		ofs[ps] = i
	}
}

func TestAllOrientations(t *testing.T) {
	type TestCase struct {
		in    []pos.P3
		wants [][]pos.P3
	}

	testCases := []TestCase{
		TestCase{
			in: []pos.P3{
				pos.P3{-1, -1, 1},
				pos.P3{-2, -2, 2},
				pos.P3{-3, -3, 3},
				pos.P3{-2, -3, 1},
				pos.P3{5, 6, -4},
				pos.P3{8, 0, 7},
			},
			wants: [][]pos.P3{
				[]pos.P3{
					pos.P3{1, -1, 1},
					pos.P3{2, -2, 2},
					pos.P3{3, -3, 3},
					pos.P3{2, -1, 3},
					pos.P3{-5, 4, -6},
					pos.P3{-8, -7, 0},
				},
				[]pos.P3{
					pos.P3{-1, -1, -1},
					pos.P3{-2, -2, -2},
					pos.P3{-3, -3, -3},
					pos.P3{-1, -3, -2},
					pos.P3{4, 6, 5},
					pos.P3{-7, 0, 8},
				},
				[]pos.P3{
					pos.P3{1, 1, -1},
					pos.P3{2, 2, -2},
					pos.P3{3, 3, -3},
					pos.P3{1, 3, -2},
					pos.P3{-4, -6, 5},
					pos.P3{7, 0, 8},
				},
				[]pos.P3{
					pos.P3{1, 1, 1},
					pos.P3{2, 2, 2},
					pos.P3{3, 3, 3},
					pos.P3{3, 1, 2},
					pos.P3{-6, -4, -5},
					pos.P3{0, 7, -8},
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			f := NewOneScannerField(0, tc.in)

			matches := map[int]int{}
			for ofn, of := range allOrientations(f) {
				for wn, want := range tc.wants {
					if _, found := matches[wn]; found {
						continue
					}

					if reflect.DeepEqual(of.Beacons, want) {
						matches[wn] = ofn
						break
					}
				}
			}

			for i := range tc.wants {
				if _, found := matches[i]; !found {
					t.Errorf("missing match for want %v", i)
				}
			}
		})
	}
}

func TestMatchFields(t *testing.T) {
	type TestCase struct {
		ref, cand []pos.P3
	}

	testCases := []TestCase{
		TestCase{
			ref: []pos.P3{
				pos.P3{404, -588, -901},
				pos.P3{528, -643, 409},
				pos.P3{-838, 591, 734},
				pos.P3{390, -675, -793},
				pos.P3{-537, -823, -458},
				pos.P3{-485, -357, 347},
				pos.P3{-345, -311, 381},
				pos.P3{-661, -816, -575},
				pos.P3{-876, 649, 763},
				pos.P3{-618, -824, -621},
				pos.P3{553, 345, -567},
				pos.P3{474, 580, 667},
				pos.P3{-447, -329, 318},
				pos.P3{-584, 868, -557},
				pos.P3{544, -627, -890},
				pos.P3{564, 392, -477},
				pos.P3{455, 729, 728},
				pos.P3{-892, 524, 684},
				pos.P3{-689, 845, -530},
				pos.P3{423, -701, 434},
				pos.P3{7, -33, -71},
				pos.P3{630, 319, -379},
				pos.P3{443, 580, 662},
				pos.P3{-789, 900, -551},
				pos.P3{459, -707, 401},
			},
			cand: []pos.P3{
				pos.P3{686, 422, 578},
				pos.P3{605, 423, 415},
				pos.P3{515, 917, -361},
				pos.P3{-336, 658, 858},
				pos.P3{95, 138, 22},
				pos.P3{-476, 619, 847},
				pos.P3{-340, -569, -846},
				pos.P3{567, -361, 727},
				pos.P3{-460, 603, -452},
				pos.P3{669, -402, 600},
				pos.P3{729, 430, 532},
				pos.P3{-500, -761, 534},
				pos.P3{-322, 571, 750},
				pos.P3{-466, -666, -811},
				pos.P3{-429, -592, 574},
				pos.P3{-355, 545, -477},
				pos.P3{703, -491, -529},
				pos.P3{-328, -685, 520},
				pos.P3{413, 935, -424},
				pos.P3{-391, 539, -444},
				pos.P3{586, -435, 557},
				pos.P3{-364, -763, -893},
				pos.P3{807, -499, -711},
				pos.P3{755, -354, -619},
				pos.P3{553, 889, -390},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ref := NewOneScannerField(0, tc.ref)
			cand := NewOneScannerField(1, tc.cand)

			if matched := matchFields(ref, cand); matched == nil {
				t.Errorf("no match")
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
