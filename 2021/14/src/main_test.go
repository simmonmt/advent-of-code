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

	"github.com/simmonmt/aoc/2021/common/logger"
)

func TestSolve(t *testing.T) {
	type TestCase struct {
		in         string
		numSteps   int
		rewriteMap map[string]rune
		want       map[string]int
	}

	nncbRewriteMap := map[string]rune{
		"CH": 'B',
		"HH": 'N',
		"CB": 'H',
		"NH": 'C',
		"HB": 'C',
		"HC": 'B',
		"HN": 'C',
		"NN": 'C',
		"BH": 'H',
		"NC": 'B',
		"NB": 'B',
		"BN": 'B',
		"BB": 'N',
		"BC": 'B',
		"CC": 'N',
		"CN": 'C',
	}

	testCases := []TestCase{
		TestCase{
			in:         "NNCB",
			numSteps:   0,
			rewriteMap: nncbRewriteMap,
			want: map[string]int{
				"B": 1,
				"C": 1,
				"N": 2,
			},
		},
		TestCase{
			in:         "NNCB",
			numSteps:   1,
			rewriteMap: nncbRewriteMap,
			want: map[string]int{
				"B": 2,
				"C": 2,
				"H": 1,
				"N": 2,
			},
		},
		TestCase{
			in:         "NNCB",
			numSteps:   2,
			rewriteMap: nncbRewriteMap,
			want: map[string]int{
				"B": 6,
				"C": 4,
				"H": 1,
				"N": 2,
			},
		},
		TestCase{
			in:         "NNCB",
			numSteps:   3,
			rewriteMap: nncbRewriteMap,
			want: map[string]int{
				"B": 11,
				"C": 5,
				"H": 4,
				"N": 5,
			},
		},
		TestCase{
			in:         "NNCB",
			numSteps:   10,
			rewriteMap: nncbRewriteMap,
			want: map[string]int{
				"B": 1749,
				"C": 298,
				"H": 161,
				"N": 865,
			},
		},
		TestCase{
			in:         "NNCB",
			numSteps:   40,
			rewriteMap: nncbRewriteMap,
			want: map[string]int{
				"B": 2192039569602,
				"H": 3849876073,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v_%d", tc.in, tc.numSteps),
			func(t *testing.T) {
				got := solveSeq(tc.in, tc.numSteps, tc.rewriteMap)

				for s, wantN := range tc.want {
					if gotN := got[s]; gotN != wantN {
						t.Errorf("totals[%v] = %v, want %v",
							s, gotN, wantN)
					}
				}
			})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
