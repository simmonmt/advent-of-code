package main

import (
	_ "embed"
	"os"
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
			WantA: -1, WantB: -1,
		},
		testutils.SampleTestCase{
			WantA: -1, WantB: 117440,
		},
	}

	sampleAWants = []string{
		"4,6,3,5,6,3,5,2,1,0", "",
	}
)

func TestSolveA(t *testing.T) {
	for i, tc := range sampleTestCases {
		t.Run(tc.File, func(t *testing.T) {
			input, err := parseInput(tc.Body)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := solveA(input), sampleAWants[i]; want != "" && got != want {
				t.Errorf("solveA(sample) = %v, want %v", got, want)
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

func TestRunProgram(t *testing.T) {
	type TestCase struct {
		in   uint64
		want byte
	}

	testCases := []TestCase{
		TestCase{in: 0b000, want: 0b110}, TestCase{in: 0b111000, want: 0b001},
		TestCase{in: 0b001, want: 0b111},
		TestCase{in: 0b010, want: 0b101},
		TestCase{in: 0b011, want: 0b110},
		TestCase{in: 0b100, want: 0b010},
		TestCase{in: 0b101, want: 0b011},
		TestCase{in: 0b110, want: 0b000}, TestCase{in: 0b11100110, want: 0b111},
		TestCase{in: 0b111, want: 0b001}, TestCase{in: 0b01110111, want: 0b110},
	}

	program := []byte{2, 4, 1, 3, 7, 5, 0, 3, 1, 5, 4, 4, 5, 5}

	for _, tc := range testCases {
		t.Run(strconv.FormatUint(uint64(tc.in), 2), func(t *testing.T) {
			regs := map[string]int{"A": int(tc.in), "B": 0, "C": 0}
			if got := byte(runProgram(program, regs)[0]); got != tc.want {
				t.Errorf("got %b %d want %b %d", got, got, tc.want, tc.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	testutils.PopulateTestCases(rawSample, sampleTestCases)
	os.Exit(m.Run())
}
