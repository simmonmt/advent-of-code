package main

import (
	_ "embed"
	"fmt"
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
			WantA: 0, WantB: -1,
		},
	}

	sampleASolutions = []map[int]int{
		map[int]int{2: 14, 4: 14, 6: 2, 8: 4, 10: 2, 12: 3, 20: 1, 36: 1, 38: 1, 40: 1, 64: 1},
	}
	sampleBSolutions = []map[int]int{
		map[int]int{
			50: 32, 52: 31, 54: 29, 56: 39, 58: 25, 60: 23, 62: 20, 64: 19, 66: 12,
			68: 14, 70: 12, 72: 22, 74: 4, 76: 3,
		},
	}
)

func parseOrDie(body []string) *Input {
	input, err := parseInput(body)
	if err != nil {
		panic(fmt.Sprintf("bad body: %v", err))
	}
	return input
}

func TestFindSolutions(t *testing.T) {
	type TestCase struct {
		maxDist, minWant int
		want             map[int]int
	}

	testCases := []TestCase{
		TestCase{
			maxDist: 2,
			minWant: 0,
			want:    map[int]int{2: 14, 4: 14, 6: 2, 8: 4, 10: 2, 12: 3, 20: 1, 36: 1, 38: 1, 40: 1, 64: 1},
		},
		TestCase{
			maxDist: 20,
			minWant: 50,
			want: map[int]int{
				50: 32, 52: 31, 54: 29, 56: 39, 58: 25, 60: 23, 62: 20, 64: 19, 66: 12,
				68: 14, 70: 12, 72: 22, 74: 4, 76: 3,
			},
		},
	}

	input := parseOrDie(sampleTestCases[0].Body)
	for _, tc := range testCases {
		t.Run(strconv.Itoa(tc.maxDist), func(t *testing.T) {
			got := map[int]int{}
			for d, n := range findSolutions(input, tc.maxDist) {
				if d >= tc.minWant {
					got[d] = n
				}
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("findSolutions mismatch; -want,+got:\n%s\n", diff)
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	testutils.PopulateTestCases(rawSample, sampleTestCases)
	os.Exit(m.Run())
}
