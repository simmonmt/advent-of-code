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
		testutils.SampleTestCase{ // sample_1.txt
			WantA: 2, WantB: -1,
		},
		testutils.SampleTestCase{ // sample_2.txt
			WantA: 4, WantB: -1,
		},
		testutils.SampleTestCase{ // sample_3.txt
			WantA: 3, WantB: -1,
		},
		testutils.SampleTestCase{ // sample_4.txt
			WantA: 1, WantB: -1,
		},
		testutils.SampleTestCase{ // sample_5.txt
			WantA: 36, WantB: 81,
		},
		testutils.SampleTestCase{ // sample_6.txt
			WantA: -1, WantB: 3,
		},
		testutils.SampleTestCase{ // sample_7.txt
			WantA: -1, WantB: 13,
		},
		testutils.SampleTestCase{ // sample_8.txt
			WantA: -1, WantB: 227,
		},
	}
)

func TestParseInput(t *testing.T) {
	for _, tc := range sampleTestCases {
		if tc.WantInput == nil {
			continue
		}

		t.Run(tc.File, func(t *testing.T) {
			input, err := parseInput(tc.Body)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.WantInput, input); diff != "" {
				t.Errorf("parseInput mismatch; -want,+got:\n%s\n", diff)
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
