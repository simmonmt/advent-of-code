package main

import (
	_ "embed"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/simmonmt/aoc/2024/common/logger"
	"github.com/simmonmt/aoc/2024/common/testutils"
)

var (
	//go:embed combined_samples.txt
	rawSample       string
	sampleTestCases = []testutils.SampleTestCase{
		testutils.SampleTestCase{
			WantA: 126384, WantB: -1,
		},
	}
)

// 789
// 456
// 123
// _0A
func TestNumPad(t *testing.T) {
	type TestCase struct {
		Cur, Dest rune
		Want      []string
	}

	testCases := []TestCase{
		TestCase{Cur: 'A', Dest: '3', Want: []string{"^A"}},
		TestCase{Cur: 'A', Dest: '2', Want: []string{"<^A", "^<A"}},
		TestCase{Cur: 'A', Dest: '1', Want: []string{"<^<A", "^<<A"}},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := NewNumPad().PressFrom(tc.Cur, tc.Dest)
			slices.Sort(got)

			if !reflect.DeepEqual(got, tc.Want) {
				t.Errorf("Press(%v) from %v = %v, want %v", tc.Dest, tc.Cur, got, tc.Want)
			}
		})
	}
}

// _^A
// <v>
func TestDirPad(t *testing.T) {
	type TestCase struct {
		Cur, Dest rune
		Want      []string
	}

	testCases := []TestCase{
		TestCase{Cur: 'A', Dest: '>', Want: []string{"vA"}},
		TestCase{Cur: 'A', Dest: 'v', Want: []string{"<vA", "v<A"}},
		TestCase{Cur: 'A', Dest: '<', Want: []string{"<v<A", "v<<A"}},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := NewDirPad().PressFrom(tc.Cur, tc.Dest)
			slices.Sort(got)

			if !reflect.DeepEqual(got, tc.Want) {
				t.Errorf("Press(%v) from %v = %v, want %v", tc.Dest, tc.Cur, got, tc.Want)
			}
		})
	}
}

func TestMinCostsBase(t *testing.T) {
	n1 := &Node{To: '0', Level: 1}
	n2 := &Node{To: '2', Level: 1}
	n1.Next = n2

	want := []*Cost{&Cost{To: "2", Cost: 2}}
	got := findMinCosts(n1, "A")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findMinCosts mismatch; -want,+got:\n%s\n", diff)
	}
}

func TestMinCostsSimple(t *testing.T) {
	n21 := &Node{To: '<', Level: 2}
	n22 := &Node{To: 'A', Level: 2}
	n21.Next = n22

	n1 := &Node{To: '0', Level: 1, Paths: map[rune][]*Node{'A': []*Node{n21}}}
	n21.Parent = n1
	n22.Parent = n1

	want := []*Cost{&Cost{To: "0A", Cost: 2}}
	got := findMinCosts(n1, "AA")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findMinCosts mismatch; -want,+got:\n%s\n", diff)
	}
}

func TestMinCosts(t *testing.T) {
	type TestCase struct {
		Input string
		Stack []Keypad
		Want  []*Cost
	}

	testCases := []TestCase{
		TestCase{
			Input: "<",
			Stack: []Keypad{NewDirPad()},
			Want:  []*Cost{&Cost{To: "<A", Cost: 4}},
		},
		TestCase{
			Input: "0",
			Stack: []Keypad{NewNumPad()},
			Want:  []*Cost{&Cost{To: "0A", Cost: 2}},
		},
		TestCase{
			Input: "029A",
			Stack: []Keypad{NewNumPad(), NewDirPad(), NewDirPad()},
			Want:  []*Cost{&Cost{To: "AAAA", Cost: 68}},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			start := buildPathNodes(tc.Input, 1, nil, tc.Stack)

			froms := strings.Repeat("A", len(tc.Stack)+1)

			got := findMinCosts(start, froms)
			if diff := cmp.Diff(tc.Want, got); diff != "" {
				t.Errorf("findMinCosts mismatch; -want,+got:\n%s\n", diff)
			}
		})
	}
}

func TestSolveA(t *testing.T) {
	for _, tc := range sampleTestCases {
		if tc.WantA == -1 {
			continue
		}

		t.Run(tc.File, func(t *testing.T) {
			input, err := parseInput(tc.Body)
			if err != nil {
				t.Fatal(err)
			}

			if got := solveA(input); got != tc.WantA {
				t.Errorf("solveA(sample) = %v, want %v", got, tc.WantA)
			}
		})
	}
}

func TestSolveB(t *testing.T) {
	for _, tc := range sampleTestCases {
		if tc.WantB == -1 {
			continue
		}

		t.Run(tc.File, func(t *testing.T) {
			input, err := parseInput(tc.Body)
			if err != nil {
				t.Fatal(err)
			}

			if got := solveB(input); got != tc.WantB {
				t.Errorf("solveB(sample) = %v, want %v", got, tc.WantB)
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	testutils.PopulateTestCases(rawSample, sampleTestCases)
	os.Exit(m.Run())
}
