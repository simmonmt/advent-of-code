package main

import (
	_ "embed"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/simmonmt/aoc/2024/common/logger"
	"github.com/simmonmt/aoc/2024/common/testutils"
)

type AAndB struct {
	As, Bs []int
}

var (
	//go:embed combined_samples.txt
	rawSample       string
	sampleTestCases = []testutils.SampleTestCase{
		testutils.SampleTestCase{
			WantInput: &AAndB{
				As: []int{3, 4, 2, 1, 3, 3},
				Bs: []int{4, 3, 5, 3, 9, 3},
			},
			WantA: 11, WantB: 31,
		},
	}
)

func TestParseInput(t *testing.T) {
	for _, tc := range sampleTestCases {
		if tc.WantInput == nil {
			continue
		}

		t.Run(tc.File, func(t *testing.T) {
			as, bs, err := parseInput(tc.Body)
			if err != nil {
				t.Fatal(err)
			}

			got := &AAndB{as, bs}
			if diff := cmp.Diff(tc.WantInput, got); diff != "" {
				t.Errorf("parseInput mismatch; -want,+got:\n%s\n", diff)
			}
		})
	}
}

func TestSolveA(t *testing.T) {
	for _, tc := range sampleTestCases {
		t.Run(tc.File, func(t *testing.T) {
			as, bs, err := parseInput(tc.Body)
			if err != nil {
				t.Fatal(err)
			}

			if got := solveA(as, bs); got != tc.WantA {
				t.Errorf("solveA(sample) = %v, want %v", got, tc.WantA)
			}
		})
	}
}

func TestSolveB(t *testing.T) {
	for _, tc := range sampleTestCases {
		t.Run(tc.File, func(t *testing.T) {
			as, bs, err := parseInput(tc.Body)
			if err != nil {
				t.Fatal(err)
			}

			if got := solveB(as, bs); got != tc.WantB {
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
