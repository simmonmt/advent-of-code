package main

import (
	_ "embed"
	"os"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/simmonmt/aoc/2025/common/logger"
	"github.com/simmonmt/aoc/2025/common/testutils"
)

var (
	//go:embed combined_samples.txt
	rawSample       string
	sampleTestCases = []testutils.SampleTestCase{
		testutils.SampleTestCase{
			WantInput: &Input{Ranges: []Range{
				{From: 11, To: 22},
				{From: 95, To: 115},
				{From: 998, To: 1012},
				{From: 1188511880, To: 1188511890},
				{From: 222220, To: 222224},
				{From: 1698522, To: 1698528},
				{From: 446443, To: 446449},
				{From: 38593856, To: 38593862},
				{From: 565653, To: 565659},
				{From: 824824821, To: 824824827},
				{From: 2121212118, To: 2121212124},
			}},
			WantA: 1227775554, WantB: 4174379265,
		},
	}
)

func TestIsAValid(t *testing.T) {
	testCases := map[int]bool{
		11: false,
		12: true,
		22: false,
	}

	for in, want := range testCases {
		t.Run(strconv.Itoa(in), func(t *testing.T) {
			if got := isAValid(in); got != want {
				t.Errorf("isAValid(%v) = %v, want %v", in, got, want)
			}
		})
	}
}

func TestIsBValid(t *testing.T) {
	testCases := map[int]bool{
		11:         false,
		12:         true,
		22:         false,
		999:        false,
		1010:       false,
		1188511885: false,
	}

	for in, want := range testCases {
		t.Run(strconv.Itoa(in), func(t *testing.T) {
			if got := isBValid(in); got != want {
				t.Errorf("isBValid(%v) = %v, want %v", in, got, want)
			}
		})
	}
}

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
		if tc.WantA == nil {
			continue
		}

		t.Run(tc.File, func(t *testing.T) {
			input, err := parseInput(tc.Body)
			if err != nil {
				t.Fatal(err)
			}

			got := solveA(input)
			if diff := cmp.Diff(tc.WantA, got); diff != "" {
				t.Errorf("solveA(sample) = %v, want %v", got, tc.WantA)
			}
		})
	}
}

func TestSolveB(t *testing.T) {
	for _, tc := range sampleTestCases {
		if tc.WantB == nil {
			continue
		}

		t.Run(tc.File, func(t *testing.T) {
			input, err := parseInput(tc.Body)
			if err != nil {
				t.Fatal(err)
			}

			got := solveB(input)
			if diff := cmp.Diff(tc.WantB, got); diff != "" {
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
