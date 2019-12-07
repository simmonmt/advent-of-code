package vm

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/simmonmt/aoc/2019/common/logger"
)

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
		TestCase{ // input and output
			ramVals: []int{
				3, 9, // in => *9
				1001, 9, 1, 9, // add *9, 1 => *9
				4, 9, // out *9
				99, // hlt
				0,  // scratch},
			},
			input:          []int{15},
			expectedOutput: []int{16},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			logger.LogF("running case %d", i)

			ram := NewRam(tc.ramVals...)
			io := NewIO(tc.input...)

			if err := Run(ram, io, 0); err != nil {
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
