// Copyright 2022 Google LLC
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

package main

import (
	_ "embed"
	"os"
	"strings"
	"testing"

	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	//go:embed sample.txt
	rawSample   string
	sampleLines []string
)

func TestToFromSnafu(t *testing.T) {
	type TestCase struct {
		snafu string
		num   int
	}

	testCases := []TestCase{
		TestCase{"1=-0-2", 1747},
		TestCase{"12111", 906},
		TestCase{"2=0=", 198},
		TestCase{"21", 11},
		TestCase{"2=01", 201},
		TestCase{"111", 31},
		TestCase{"20012", 1257},
		TestCase{"112", 32},
		TestCase{"1=-1=", 353},
		TestCase{"1-12", 107},
		TestCase{"12", 7},
		TestCase{"1=", 3},
		TestCase{"122", 37},
		TestCase{"1", 1},
		TestCase{"2", 2},
		TestCase{"1=", 3},
		TestCase{"1-", 4},
		TestCase{"10", 5},
		TestCase{"11", 6},
		TestCase{"12", 7},
		TestCase{"2=", 8},
		TestCase{"2-", 9},
		TestCase{"20", 10},
		TestCase{"1=0", 15},
		TestCase{"1-0", 20},
		TestCase{"1=11-2", 2022},
		TestCase{"1-0---0", 12345},
		TestCase{"1121-1110-1=0", 314159265},
		TestCase{"2=-01", 976},
	}

	for _, tc := range testCases {
		t.Run(tc.snafu, func(t *testing.T) {
			if got := fromSnafu(tc.snafu); got != tc.num {
				t.Errorf("FromSnafu('%v') = %v, want %v",
					tc.snafu, got, tc.num)
			}
			if got := toSnafu(tc.num); got != tc.snafu {
				t.Errorf("ToSnafu(%v) = '%v', want '%v'",
					tc.num, got, tc.snafu)
			}
		})
	}
}

func TestSolveA(t *testing.T) {
	if got, want := solveA(sampleLines), "2=-1=0"; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)

	sampleLines = strings.Split(rawSample, "\n")
	if len(sampleLines) > 0 && sampleLines[len(sampleLines)-1] == "" {
		sampleLines = sampleLines[0 : len(sampleLines)-1]
	}

	os.Exit(m.Run())
}
