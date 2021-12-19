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
	"strconv"
	"testing"

	"github.com/simmonmt/aoc/2021/common/logger"
)

func TestSNumber(t *testing.T) {
	n := &SNumber{
		Left:  &SNumber{Lit: 3},
		Right: &SNumber{Left: &SNumber{Lit: 1}, Right: &SNumber{Lit: 2}},
	}

	want := "[3,[1,2]]"
	if got := n.String(); got != want {
		t.Errorf("SNumber.String = %v, want %v", got, want)
	}
}

func TestParseSNumber(t *testing.T) {
	testCases := []string{
		"[1,2]",
		"[[1,2],3]",
		"[9,[8,7]]",
		"[[1,9],[8,5]]",
		"[[[[0,7],4],[15,[0,13]]],[1,1]]",
		"[[[[1,2],[3,4]],[[5,6],[7,8]]],9]",
		"[[[9,[3,8]],[[0,9],6]],[[[3,7],[4,9]],3]]",
		"[[[[1,3],[5,3]],[[1,3],[8,7]]],[[[4,9],[6,9]],[[8,2],[7,3]]]]",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			if n, err := parseSNumber(tc); err != nil || n == nil || n.String() != tc {
				t.Errorf(`parseSNumber("%v") = %v, %v, want %v, nil`,
					tc, n, err, tc)
			}
		})
	}
}

type OpTestCase struct {
	in, want string
}

func doOpTests(t *testing.T, testCases []OpTestCase, opName string, op func(*SNumber)) {
	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			sn, err := parseSNumber(tc.in)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}

			op(sn)
			if sn.String() != tc.want {
				t.Errorf("%v = %v, want %v",
					opName, sn.String(), tc.want)
			}
		})
	}
}

func TestSplitSNumber(t *testing.T) {
	testCases := []OpTestCase{
		OpTestCase{
			in:   "[[[[0,7],4],[15,[0,13]]],[1,1]]",
			want: "[[[[0,7],4],[[7,8],[0,13]]],[1,1]]",
		},
		OpTestCase{
			in:   "[[[[0,7],4],[[7,8],[0,13]]],[1,1]]",
			want: "[[[[0,7],4],[[7,8],[0,[6,7]]]],[1,1]]",
		},
	}

	doOpTests(t, testCases, "split",
		func(sn *SNumber) { splitSNumber(sn) })
}

func TestExplodeSNumber(t *testing.T) {
	testCases := []OpTestCase{
		OpTestCase{
			in:   "[[[[[9,8],1],2],3],4]",
			want: "[[[[0,9],2],3],4]",
		},
		OpTestCase{
			in:   "[7,[6,[5,[4,[3,2]]]]]",
			want: "[7,[6,[5,[7,0]]]]",
		},
		OpTestCase{
			in:   "[[6,[5,[4,[3,2]]]],1]",
			want: "[[6,[5,[7,0]]],3]",
		},
		OpTestCase{
			in:   "[[3,[2,[1,[7,3]]]],[6,[5,[4,[3,2]]]]]",
			want: "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]",
		},
		OpTestCase{
			in:   "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]",
			want: "[[3,[2,[8,0]]],[9,[5,[7,0]]]]",
		},
	}

	doOpTests(t, testCases, "explode",
		func(sn *SNumber) { explodeSNumber(sn) })
}

func TestReduceSNumber(t *testing.T) {
	testCases := []OpTestCase{
		OpTestCase{
			in:   "[[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]",
			want: "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]",
		},
	}

	doOpTests(t, testCases, "reduce", reduceSNumber)
}

func TestAddSNumbers(t *testing.T) {
	type TestCase struct {
		ins  []string
		want string
	}

	testCases := []TestCase{
		TestCase{
			ins: []string{
				"[1,1]",
				"[2,2]",
				"[3,3]",
				"[4,4]",
			},
			want: "[[[[1,1],[2,2]],[3,3]],[4,4]]",
		},
		TestCase{
			ins: []string{
				"[1,1]",
				"[2,2]",
				"[3,3]",
				"[4,4]",
				"[5,5]",
			},
			want: "[[[[3,0],[5,3]],[4,4]],[5,5]]",
		},
		TestCase{
			ins: []string{
				"[1,1]",
				"[2,2]",
				"[3,3]",
				"[4,4]",
				"[5,5]",
				"[6,6]",
			},
			want: "[[[[5,0],[7,4]],[5,5]],[6,6]]",
		},
		TestCase{
			ins: []string{
				"[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]",
				"[7,[[[3,7],[4,3]],[[6,3],[8,8]]]]",
				"[[2,[[0,8],[3,4]]],[[[6,7],1],[7,[1,6]]]]",
				"[[[[2,4],7],[6,[0,5]]],[[[6,8],[2,8]],[[2,1],[4,5]]]]",
				"[7,[5,[[3,8],[1,4]]]]",
				"[[2,[2,2]],[8,[8,1]]]",
				"[2,9]",
				"[1,[[[9,3],9],[[9,0],[0,7]]]]",
				"[[[5,[7,4]],7],1]",
				"[[[[4,2],2],6],[8,7]]",
			},
			want: "[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ins := []*SNumber{}
			for _, in := range tc.ins {
				sn, err := parseSNumber(in)
				if err != nil {
					t.Fatalf("parse failure: %v", err)
				}
				ins = append(ins, sn)
			}

			got := addSNumbers(ins)
			if got.String() != tc.want {
				t.Errorf("add = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestMagnitude(t *testing.T) {
	type TestCase struct {
		in   string
		want int
	}

	testCases := []TestCase{
		TestCase{
			in:   "[[1,2],[[3,4],5]]",
			want: 143,
		},
		TestCase{
			in:   "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]",
			want: 1384,
		},
		TestCase{
			in:   "[[[[1,1],[2,2]],[3,3]],[4,4]]",
			want: 445,
		},
		TestCase{
			in:   "[[[[3,0],[5,3]],[4,4]],[5,5]]",
			want: 791,
		},
		TestCase{
			in:   "[[[[5,0],[7,4]],[5,5]],[6,6]]",
			want: 1137,
		},
		TestCase{
			in:   "[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]",
			want: 3488,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			sn, err := parseSNumber(tc.in)
			if err != nil {
				t.Fatalf("parse failure: %v", err)
			}

			if got := magnitude(sn); got != tc.want {
				t.Errorf(`magnitude("%v") = %v, want %v`, sn, got, tc.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
