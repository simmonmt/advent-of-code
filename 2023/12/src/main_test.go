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
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/simmonmt/aoc/2023/common/logger"
)

var (
	//go:embed sample.txt
	rawSample   string
	sampleLines []string
)

func rangeOnStates(states []CellState, r Range) string {
	out := ""
	for i, s := range states {
		if i >= r.Left && i <= r.Right+1 {
			out += "_"
		} else {
			out += " "
		}
		out += s.String()
	}
	if r.Right == len(states)-1 {
		out += "_"
	}
	return out
}

func TestRangeOnStates(t *testing.T) {
	inStates := []CellState{CS_YES, CS_NO, CS_YES, CS_NO}
	inRange := Range{1, 2}

	want := " #_._#_."

	if got := rangeOnStates(inStates, inRange); got != want {
		t.Errorf(`rangeOnStates(%v, %v) = "%s", want "%s"`, inStates, inRange, got, want)
	}
}

func TestCellStatesFromString(t *testing.T) {
	in := ".##?."
	want := []CellState{CS_NO, CS_YES, CS_YES, CS_UNKNOWN, CS_NO}

	if got := CellStatesFromString(in); !reflect.DeepEqual(got, want) {
		t.Errorf(`CellStatesFrom("%s") = %v, want %v`, in, got, want)
	}
}

func TestParseInput(t *testing.T) {
	springs, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	want := []*Spring{
		&Spring{States: CellStatesFromString("???.###"), Sizes: []int{1, 1, 3}},
		&Spring{States: CellStatesFromString(".??..??...?##."), Sizes: []int{1, 1, 3}},
		&Spring{States: CellStatesFromString("?#?#?#?#?#?#?#?"), Sizes: []int{1, 3, 1, 6}},
		&Spring{States: CellStatesFromString("????.#...#..."), Sizes: []int{4, 1, 1}},
		&Spring{States: CellStatesFromString("????.######..#####."), Sizes: []int{1, 6, 5}},
		&Spring{States: CellStatesFromString("?###????????"), Sizes: []int{3, 2, 1}},
	}

	if diff := cmp.Diff(want, springs); diff != "" {
		t.Errorf("parseInput mismatch; -want,+got:\n%s\n", diff)
	}
}

type SolveSpringTestCase struct {
	spring              *Spring
	wantStates          []CellState
	wantConstraints     []Range
	wantNumCombinations int
}

func simpleTestCase(statesStart, statesEnd string, size int, wantRange Range, wantNumCombinations int) SolveSpringTestCase {
	return SolveSpringTestCase{
		spring:              &Spring{States: CellStatesFromString(statesStart), Sizes: []int{size}},
		wantStates:          CellStatesFromString(statesEnd),
		wantConstraints:     []Range{wantRange},
		wantNumCombinations: wantNumCombinations,
	}
}

var (
	solveSpringTestCases = []SolveSpringTestCase{
		simpleTestCase("#???", "####", 4, Range{0, 3}, 1),
		simpleTestCase(".#???", ".####", 4, Range{1, 4}, 1),
		simpleTestCase("##???", "####.", 4, Range{0, 3}, 1),
		simpleTestCase("??##??", "??###?", 5, Range{0, 5}, 2),
		SolveSpringTestCase{
			spring:              &Spring{States: CellStatesFromString("???.###"), Sizes: []int{1, 1, 3}},
			wantStates:          CellStatesFromString("#.#.###"),
			wantConstraints:     []Range{Range{0, 0}, Range{2, 2}, Range{4, 6}},
			wantNumCombinations: 1,
		},
		SolveSpringTestCase{
			spring:              &Spring{States: CellStatesFromString(".??..??...?##."), Sizes: []int{1, 1, 3}},
			wantStates:          CellStatesFromString(".??..??...###."),
			wantConstraints:     []Range{Range{1, 2}, Range{5, 6}, Range{10, 12}},
			wantNumCombinations: 4,
		},
		SolveSpringTestCase{
			spring:              &Spring{States: CellStatesFromString("?#?#?#?#?#?#?#?"), Sizes: []int{1, 3, 1, 6}},
			wantStates:          CellStatesFromString(".#.###.#.######"),
			wantConstraints:     []Range{Range{1, 1}, Range{3, 5}, Range{7, 7}, Range{9, 14}},
			wantNumCombinations: 1,
		},
		SolveSpringTestCase{
			spring:              &Spring{States: CellStatesFromString("????.#...#..."), Sizes: []int{4, 1, 1}},
			wantStates:          CellStatesFromString("####.#...#..."),
			wantConstraints:     []Range{Range{0, 3}, Range{5, 5}, Range{9, 9}},
			wantNumCombinations: 1,
		},
		SolveSpringTestCase{
			spring:              &Spring{States: CellStatesFromString("????.######..#####."), Sizes: []int{1, 6, 5}},
			wantStates:          CellStatesFromString("????.######..#####."),
			wantConstraints:     []Range{Range{0, 3}, Range{5, 10}, Range{13, 17}},
			wantNumCombinations: 4,
		},
		SolveSpringTestCase{
			spring:              &Spring{States: CellStatesFromString("?###????????"), Sizes: []int{3, 2, 1}},
			wantStates:          CellStatesFromString(".###.???????"),
			wantConstraints:     []Range{Range{1, 3}, Range{5, 9}, Range{8, 11}},
			wantNumCombinations: 10,
		},
		SolveSpringTestCase{
			spring:              &Spring{States: CellStatesFromString("?????#?##??#??#?"), Sizes: []int{10, 4}},
			wantStates:          CellStatesFromString("##########.####."),
			wantConstraints:     []Range{Range{0, 9}, Range{11, 14}},
			wantNumCombinations: 1,
		},
		SolveSpringTestCase{
			spring:              &Spring{States: CellStatesFromString("?#?.??????"), Sizes: []int{2, 2}},
			wantStates:          CellStatesFromString("?#?.??????"),
			wantConstraints:     []Range{Range{0, 9}, Range{11, 14}},
			wantNumCombinations: 10,
		},
	}
)

// func TestConstrainSpring(t *testing.T) {
// 	for i, tc := range solveSpringTestCases {
// 		t.Run(strconv.Itoa(i), func(t *testing.T) {
// 			t.Logf("constrainSpring(%v)", tc.spring)

// 			states, constraints := constrainSpring(tc.spring)
// 			t.Logf("got states %v, got constraints %v", states, constraints)

// 			if len(constraints) == len(tc.wantConstraints) {
// 				for i, constraint := range constraints {
// 					if !constraint.Equals(tc.wantConstraints[i]) {
// 						t.Logf("constraint %d mismatch; want %v %v, got %v %v for sz %d", i,
// 							rangeOnStates(tc.spring.States, tc.wantConstraints[i]), tc.wantConstraints[i],
// 							rangeOnStates(tc.spring.States, constraints[i]), constraints[i],
// 							tc.spring.Sizes[i])
// 					}
// 				}
// 			}

// 			if !reflect.DeepEqual(constraints, tc.wantConstraints) {
// 				t.Errorf("constraints = %v, want %v", constraints, tc.wantConstraints)
// 			}

// 			// if !reflect.DeepEqual(states, tc.wantStates) {
// 			// 	t.Errorf("states = %v, want %v", states, tc.wantStates)
// 			// }
// 		})
// 	}
// }

func TestRangeIterator(t *testing.T) {
	type TestCase struct {
		ranges []Range
		sizes  []int
		want   [][]Range
	}

	testCases := []TestCase{
		TestCase{
			ranges: []Range{Range{1, 1}},
			sizes:  []int{1},
			want:   [][]Range{[]Range{Range{1, 1}}},
		},
		TestCase{
			ranges: []Range{Range{1, 2}, Range{4, 5}},
			sizes:  []int{1, 1},
			want: [][]Range{
				[]Range{Range{1, 1}, Range{4, 4}},
				[]Range{Range{1, 1}, Range{5, 5}},
				[]Range{Range{2, 2}, Range{4, 4}},
				[]Range{Range{2, 2}, Range{5, 5}},
			},
		},
		TestCase{
			ranges: []Range{Range{1, 4}, Range{4, 5}},
			sizes:  []int{2, 1},
			want: [][]Range{
				[]Range{Range{1, 2}, Range{4, 4}},
				[]Range{Range{1, 2}, Range{5, 5}},
				[]Range{Range{2, 3}, Range{5, 5}},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ri := newRangeIterator(tc.ranges, tc.sizes)
			got := [][]Range{}
			for {
				r, done := ri.Next()
				got = append(got, r)
				if done {
					break
				}
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("iteration mismatch; -want,+got:\n%s", diff)
			}
		})
	}
}

func TestSolveSpring(t *testing.T) {

	for i, tc := range solveSpringTestCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := solveSpring(tc.spring); got != tc.wantNumCombinations {
				t.Errorf("solveSpring(%v) = %v, wantNumCombinations %v", tc.spring, got, tc.wantNumCombinations)
			}
		})
	}
}

func TestSolveA(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(input), 21; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveBSpring(t *testing.T) {
	input := &Spring{States: CellStatesFromString(".??..??...?##."), Sizes: []int{1, 1, 3}}

	if got, want := solveBSpring(input), 1; got != want {
		t.Errorf("solveBSpring = %d, want %d", got, want)
	}
}

func TestSolveB(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(input), 525152; got != want {
		t.Errorf("solveB(sample) = %v, want %v", got, want)
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
