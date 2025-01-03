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
			WantInput: &Input{
				Wires: map[string]State{
					"x00": ST_ON,
					"x01": ST_ON,
					"x02": ST_ON,
					"y00": ST_OFF,
					"y01": ST_ON,
					"y02": ST_OFF,
				},
				Gates: []Gate{
					Gate{ID: 0, In1: "x00", In2: "y00", Out: "z00", Type: GT_AND},
					Gate{ID: 1, In1: "x01", In2: "y01", Out: "z01", Type: GT_XOR},
					Gate{ID: 2, In1: "x02", In2: "y02", Out: "z02", Type: GT_OR},
				},
			},
			WantA: 4, WantB: -1,
		},
		testutils.SampleTestCase{
			WantA: 2024, WantB: -1,
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

func TestMain(m *testing.M) {
	logger.Init(true)
	testutils.PopulateTestCases(rawSample, sampleTestCases)
	os.Exit(m.Run())
}
