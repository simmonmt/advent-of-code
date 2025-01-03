package main

import (
	_ "embed"
	"os"
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
			WantInput: []string{""},
			WantA:     37327623, WantB: -1,
		},
		testutils.SampleTestCase{
			WantInput: []string{""},
			WantA:     -1, WantB: 23,
		},
	}
)

func TestRound(t *testing.T) {
	num := 123
	got := make([]int, 10)
	for i := range 10 {
		num = round(num)
		got[i] = num
	}

	want := []int{
		15887950, 16495136, 527345, 704524, 1553684,
		12683156, 11100544, 12249484, 7753432, 5908254,
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("round mismatch; -want,+got:\n%s\n", diff)
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
