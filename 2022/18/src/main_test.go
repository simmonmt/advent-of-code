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
	"github.com/simmonmt/aoc/2022/common/pos"
)

var (
	//go:embed sample.txt
	rawSample   string
	sampleLines []string
)

func TestFindOutside(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	start, end := pos.P3{1, 0, 0}, pos.P3{4, 4, 7}
	outside := findOutside(start, end, input)

	want := (end.X - start.X + 1) * (end.Y - start.Y + 1) * (end.Z - start.Z + 1)
	want -= (len(input) + 1) // the empty space

	if got := len(outside); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSolveA(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(input), 64; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(input), 58; got != want {
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
