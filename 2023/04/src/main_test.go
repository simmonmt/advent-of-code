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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/simmonmt/aoc/2023/common/logger"
)

var (
	//go:embed sample.txt
	rawSample   string
	sampleLines []string
)

func TestParseInput(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	want := []*Card{
		&Card{
			Num:     1,
			Winning: []int{41, 48, 83, 86, 17},
			Have:    []int{83, 86, 6, 31, 17, 9, 48, 53},
		},
		&Card{
			Num:     2,
			Winning: []int{13, 32, 20, 16, 61},
			Have:    []int{61, 30, 68, 82, 17, 32, 24, 19},
		},
		&Card{
			Num:     3,
			Winning: []int{1, 21, 53, 59, 44},
			Have:    []int{69, 82, 63, 72, 16, 21, 14, 1},
		},
		&Card{
			Num:     4,
			Winning: []int{41, 92, 73, 84, 69},
			Have:    []int{59, 84, 76, 51, 58, 5, 54, 83},
		},
		&Card{
			Num:     5,
			Winning: []int{87, 83, 26, 28, 32},
			Have:    []int{88, 30, 70, 12, 93, 22, 82, 36},
		},
		&Card{
			Num:     6,
			Winning: []int{31, 18, 13, 56, 72},
			Have:    []int{74, 77, 10, 23, 35, 67, 36, 11},
		},
	}

	if diff := cmp.Diff(want, input); diff != "" {
		t.Errorf("parseInput(sampleLines) mismatch (-want +got):\n%s", diff)
	}
}

func TestSolveA(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(input), 13; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(input), 30; got != want {
		t.Errorf("solveB(sample) = %v, want %v", got, want)
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
