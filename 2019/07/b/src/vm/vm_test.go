// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
				0,  // scratch
			},
			input:          []int{15},
			expectedOutput: []int{16},
		},
		TestCase{ // branching
			ramVals: []int{
				11105, 1, 4, // 0: jit 1, 4
				99,          // 3: hlt catches jit fail
				11106, 0, 8, // 4: jif 0, 8
				99,       // 7: hlt catches jif fail
				11104, 1, // out 1
				99, // hlt
			},
			input:          []int{},
			expectedOutput: []int{1},
		},
		TestCase{ // less than and equals
			ramVals: []int{
				1107, 0, 1, 13, // 0: lt 0,1 => 13
				1108, 0, 1, 14, // 4: eq 0,1 => 14
				4, 13, // 8: out *13
				4, 14, // 10: out *14
				99,   // 12: hlt
				0, 0, // 13: x, y
			},
			input:          []int{},
			expectedOutput: []int{1, 0},
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
				if got := io.(*ioImpl).Written(); !reflect.DeepEqual(got, tc.expectedOutput) {
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
