// Copyright 2023 Google LLC
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
	_ "embed"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/pos"
)

var (
	//go:embed sample.txt
	rawSample   string
	sampleLines []string
)

func TestCoordTranslator(t *testing.T) {
	type TestCase struct {
		xInflates, yInflates map[int]int
		reals, infs          []pos.P2
	}

	testCases := []TestCase{
		TestCase{
			reals: []pos.P2{pos.P2{X: 1, Y: 2}},
			infs:  []pos.P2{pos.P2{X: 1, Y: 2}},
		},
		TestCase{
			xInflates: map[int]int{2: 2},
			reals: []pos.P2{
				pos.P2{X: 1, Y: 0},
				pos.P2{X: 2, Y: 0},
				pos.P2{X: 3, Y: 0},
			},
			infs: []pos.P2{
				pos.P2{X: 1, Y: 0},
				pos.P2{X: 2, Y: 0},
				pos.P2{X: 4, Y: 0},
			},
		},
		TestCase{
			xInflates: map[int]int{2: 2, 5: 2, 8: 2},
			yInflates: map[int]int{3: 2, 7: 2},
			reals: []pos.P2{
				pos.P2{X: 1, Y: 0},
				pos.P2{X: 3, Y: 0},
				pos.P2{X: 5, Y: 0},
				pos.P2{X: 7, Y: 1},
				pos.P2{X: 9, Y: 0},
			},
			infs: []pos.P2{
				pos.P2{X: 1, Y: 0},
				pos.P2{X: 4, Y: 0},
				pos.P2{X: 6, Y: 0},
				pos.P2{X: 9, Y: 1},
				pos.P2{X: 12, Y: 0},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			xlate := NewCoordTranslator()

			for x, inf := range tc.xInflates {
				xlate.InflateX(x, inf)
			}
			for y, inf := range tc.yInflates {
				xlate.InflateY(y, inf)
			}

			for j := range tc.reals {
				if got := xlate.RealToInflated(tc.reals[j]); !got.Equals(tc.infs[j]) {
					t.Errorf("RealToInflated(%v) = %v, want %v", tc.reals[j], got, tc.infs[j])
				}
			}
		})
	}
}

func TestSolve(t *testing.T) {
	type TestCase struct {
		factor int
		want   int64
	}

	testCases := []TestCase{
		TestCase{2, 374},
		TestCase{10, 1030},
		TestCase{100, 8410},
	}

	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(strconv.Itoa(tc.factor), func(t *testing.T) {
			if got := solve(input, tc.factor); got != tc.want {
				t.Errorf("solve(%d) = %v, want %v", tc.factor, got, tc.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)

	sampleLines = strings.Split(rawSample, "\n")
	if len(sampleLines) > 0 && sampleLines[len(sampleLines)-1] == "" {
		sampleLines = sampleLines[0 : len(sampleLines)-1]
	}

	os.Exit(m.Run())
}
