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
	"fmt"
	"os"
	"testing"

	"github.com/simmonmt/aoc/2020/common/logger"
	"github.com/simmonmt/aoc/2020/common/pos"
)

func TestNeighborCounterB(t *testing.T) {
	type TestCase struct {
		lines []string
		p     []pos.P2
		want  []int
	}

	testCases := []TestCase{
		TestCase{
			lines: []string{ //
				".......#.", //
				"...#.....", //
				".#.......", //
				".........", //
				"..#L....#", //
				"....#....", //
				".........", //
				"#........", //
				"...#.....", //
			},
			p:    []pos.P2{pos.P2{X: 3, Y: 4}},
			want: []int{8},
		},
		TestCase{
			lines: []string{ //
				".............", //
				".L.L.#.#.#.#.", //
				".............", //
			},
			p: []pos.P2{
				pos.P2{X: 1, Y: 1},
				pos.P2{X: 3, Y: 1},
			},
			want: []int{0, 1},
		},
		TestCase{
			lines: []string{ //
				".##.##.", //
				"#.#.#.#", //
				"##...##", //
				"...L...", //
				"##...##", //
				"#.#.#.#", //
				".##.##.", //
			},
			p:    []pos.P2{pos.P2{X: 3, Y: 3}},
			want: []int{0},
		},
	}

	for tcNum, tc := range testCases {
		for i := range tc.p {
			t.Run(fmt.Sprintf("%d/%d", tcNum, i), func(t *testing.T) {
				logger.LogF("test %d/%d\n", tcNum, i)
				b := newBoard(tc.lines)
				got := neighborCounterB(b, tc.p[i])
				if got != tc.want[i] {
					t.Errorf("occupiedNeighbors(_, %v) = %v, want %v",
						tc.p[i], got, tc.want[i])
				}
			})
		}
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
