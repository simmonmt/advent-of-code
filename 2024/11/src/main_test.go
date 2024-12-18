package main

import (
	_ "embed"
	"os"
	"strconv"
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
			WantA: 55312, WantB: -1,
		},
	}
)

func TestSplitDigits(t *testing.T) {
	type TestCase struct {
		in               int64
		num, left, right int64
	}

	testCases := []TestCase{
		TestCase{1, 1, -1, -1},
		TestCase{10, 2, 1, 0},
		TestCase{123, 3, -1, -1},
		TestCase{1234, 4, 12, 34},
	}

	for _, tc := range testCases {
		t.Run(strconv.FormatInt(tc.in, 10), func(t *testing.T) {
			num, left, right := splitDigits(tc.in)
			if num != tc.num || (tc.left != -1 && left != tc.left) || (tc.right != -1 && right != tc.right) {
				wantLeft := strconv.FormatInt(tc.left, 10)
				if tc.left == -1 {
					wantLeft = "_"
				}
				wantRight := strconv.FormatInt(tc.right, 10)
				if tc.right == -1 {
					wantRight = "_"
				}

				t.Errorf("splitDigits(%v) = %v, %v, %v; want %v, %v, %v",
					tc.in, num, left, right, tc.num, wantLeft, wantRight)
			}
		})
	}
}

func TestTransform(t *testing.T) {
	type TestCase struct {
		in, want []int64
	}

	seq := [][]int64{
		[]int64{125, 17},
		[]int64{253000, 1, 7},
		[]int64{253, 0, 2024, 14168},
		[]int64{512072, 1, 20, 24, 28676032},
		[]int64{512, 72, 2024, 2, 0, 2, 4, 2867, 6032},
		[]int64{1036288, 7, 2, 20, 24, 4048, 1, 4048, 8096, 28, 67, 60, 32},
		[]int64{2097446912, 14168, 4048, 2, 0, 2, 4, 40, 48, 2024, 40, 48, 80, 96, 2, 8, 6, 7, 6, 0, 3, 2},
	}

	testCases := []TestCase{}
	for i := 1; i < len(seq); i++ {
		testCases = append(testCases, TestCase{in: seq[i-1], want: seq[i]})
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := transform(tc.in)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("transform mismatch; -want,+got:\n%s\n", diff)
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
