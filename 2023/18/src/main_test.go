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

	want := []Command{
		Command{Dir: 'R', Num: 6, Color: 0x70c710},
		Command{Dir: 'D', Num: 5, Color: 0x0dc571},
		Command{Dir: 'L', Num: 2, Color: 0x5713f0},
		Command{Dir: 'D', Num: 2, Color: 0xd2c081},
		Command{Dir: 'R', Num: 2, Color: 0x59c680},
		Command{Dir: 'D', Num: 2, Color: 0x411b91},
		Command{Dir: 'L', Num: 5, Color: 0x8ceee2},
		Command{Dir: 'U', Num: 2, Color: 0xcaa173},
		Command{Dir: 'L', Num: 1, Color: 0x1b58a2},
		Command{Dir: 'U', Num: 2, Color: 0xcaa171},
		Command{Dir: 'R', Num: 2, Color: 0x7807d2},
		Command{Dir: 'U', Num: 3, Color: 0xa77fa3},
		Command{Dir: 'L', Num: 2, Color: 0x015232},
		Command{Dir: 'U', Num: 2, Color: 0x7a21e3},
	}

	if diff := cmp.Diff(want, input); diff != "" {
		t.Errorf("parseInput mismatch; -want,+got:\n%s\n", diff)
	}
}

func TestRowSize(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}
	lcs := makeLocatedCommands(input)

	wantSizes := []int{7, 7, 7, 5, 5, 7, 5, 7, 6, 6}

	for y := 0; y < 10; y++ {
		got := rowSize(y, lcs)
		if want := wantSizes[y]; want != got {
			t.Errorf("row %d want %d got %d", y, want, got)
		}
	}
}

func TestRowRepeats(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}
	lcs := makeLocatedCommands(input)

	wantRepeats := []int{1, 1, 1, 2, 1, 1, 1, 1, 1, 1}

	for y := 0; y < 10; y++ {
		got := rowRepeats(y, lcs)
		if want := wantRepeats[y]; want != got {
			t.Errorf("row %d want %d got %d", y, want, got)
		}
	}
}

func TestSolveA(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(input), int64(62); got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(input), int64(952408144115); got != want {
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
