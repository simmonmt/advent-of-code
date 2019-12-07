package vm

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/simmonmt/aoc/2019/common/logger"
)

func TestSimpleRun(t *testing.T) {
	// This is the example program from day 2
	ram := NewRam(1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50)
	if err := Run(0, ram); err != nil {
		t.Errorf("Run = %v, want nil", err)
		return
	}

	if got := ram.Read(0); got != 3500 {
		t.Errorf("ram[0] = %v, want %v", got, 3500)
	}
}

func TestRun(t *testing.T) {
	type TestCase struct {
		ramVals        []int
		input          []int
		expectedOutput []int
		expectedRam    []int
	}

	testCases := []TestCase{
		TestCase{ // This is the example program from day 2a
			ramVals:     []int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50},
			input:       []int{},
			expectedRam: []int{3500, 9, 10, 70, 2, 3, 11, 0, 99, 30, 40, 50},
		},
		TestCase{ // Example from day 5a
			ramVals:     []int{1002, 4, 3, 4, 33},
			input:       []int{},
			expectedRam: []int{1002, 4, 3, 4, 99},
		},
		TestCase{ // negative numbers
			ramVals:     []int{1101, 100, -1, 4, 0},
			input:       []int{},
			expectedRam: []int{1101, 100, -1, 4, 99},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			logger.LogF("running case %d", i)

			ram := NewRam(tc.ramVals...)
			io := NewIO(tc.input...)

			if err := Run(0, ram); err != nil {
				t.Errorf("run failed: %v", err)
				return
			}

			if tc.expectedOutput != nil {
				if got := io.Written(); !reflect.DeepEqual(got, tc.expectedOutput) {
					t.Errorf("output = %v, want %v", got, tc.expectedOutput)
				}
			}

			if tc.expectedRam != nil {
				for addr, want := range tc.expectedRam {
					if got := ram.Read(addr); got != want {
						t.Errorf("ram verify failed: ram[%d]=%d, want %d", addr, got, want)
					}
				}
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
