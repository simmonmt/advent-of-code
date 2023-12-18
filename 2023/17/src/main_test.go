// Copyright 2023 Google LLC
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
	"strconv"
	"strings"
	"testing"

	"github.com/simmonmt/aoc/2023/common/logger"
)

var (
	//go:embed sample.txt
	rawSample   string
	sampleLines []string

	//go:embed sample2.txt
	rawSample2   string
	sample2Lines []string
)

func TestSolveA(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(input), 102; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	type TestCase struct {
		lines []string
		want  int
	}

	testCases := []TestCase{
		TestCase{lines: sampleLines, want: 94},
		TestCase{lines: sample2Lines, want: 71},
	}

	*dumpGrid = true

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			input, err := parseInput(tc.lines)
			if err != nil {
				t.Fatal(err)
			}

			if got := solveB(input); got != tc.want {
				t.Errorf("solveB(sample) = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(*verbose)

	sampleLines = strings.Split(rawSample, "\n")
	if len(sampleLines) > 0 && sampleLines[len(sampleLines)-1] == "" {
		sampleLines = sampleLines[0 : len(sampleLines)-1]
	}

	sample2Lines = strings.Split(rawSample2, "\n")
	if len(sample2Lines) > 0 && sample2Lines[len(sample2Lines)-1] == "" {
		sample2Lines = sample2Lines[0 : len(sample2Lines)-1]
	}

	os.Exit(m.Run())
}
