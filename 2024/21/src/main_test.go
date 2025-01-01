package main

import (
	_ "embed"
	"os"
	"reflect"
	"slices"
	"strconv"
	"testing"

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

// func TestBuildPathNodes(t *testing.T) {
// 	got := buildPathNodes("289A", 1, nil, []Keypad{NewNumPad()})

// 	for n := got; n != nil; n = n.Next {
// 		fmt.Printf("%+v\n", n)

// 		for k, nl := range collections.SortedMapIter(n.Paths) {
// 			for _, n2 := range nl {
// 				out := ""
// 				for n3 := n2; n3 != nil; n3 = n3.Next {
// 					out += fmt.Sprintf(" %d,%s", n3.Level, string(n3.To))
// 				}
// 				fmt.Printf("%c %s\n", k, out)
// 			}
// 		}
// 	}
// 	t.Errorf("no")
// }

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
