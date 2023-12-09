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

	//go:embed sample2.txt
	rawSample2   string
	sample2Lines []string

	//go:embed sample3.txt
	rawSample3   string
	sample3Lines []string
)

func TestParseInput(t *testing.T) {
	instructions, nodeMap, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if want := "RL"; instructions != want {
		t.Errorf("instructions got %v, want %v", instructions, want)
	}

	wantNodeMap := map[string]Node{
		"AAA": {Left: "BBB", Right: "CCC"},
		"BBB": {Left: "DDD", Right: "EEE"},
		"CCC": {Left: "ZZZ", Right: "GGG"},
		"DDD": {Left: "DDD", Right: "DDD"},
		"EEE": {Left: "EEE", Right: "EEE"},
		"GGG": {Left: "GGG", Right: "GGG"},
		"ZZZ": {Left: "ZZZ", Right: "ZZZ"},
	}

	if diff := cmp.Diff(wantNodeMap, nodeMap); diff != "" {
		t.Errorf("parseInput(sampleLines) mismatch; -want,+got:\n%s\n", diff)
	}
}

func TestSolveA(t *testing.T) {
	instructions, nodeMap, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(instructions, nodeMap), 2; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveA2(t *testing.T) {
	instructions, nodeMap, err := parseInput(sample2Lines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(instructions, nodeMap), 6; got != want {
		t.Errorf("solveA(sample2) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	instructions, nodeMap, err := parseInput(sample3Lines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(instructions, nodeMap), int64(6); got != want {
		t.Errorf("solveB(sample) = %v, want %v", got, want)
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)

	sampleLines = strings.Split(rawSample, "\n")
	if len(sampleLines) > 0 && sampleLines[len(sampleLines)-1] == "" {
		sampleLines = sampleLines[0 : len(sampleLines)-1]
	}

	sample2Lines = strings.Split(rawSample2, "\n")
	if len(sample2Lines) > 0 && sample2Lines[len(sample2Lines)-1] == "" {
		sample2Lines = sample2Lines[0 : len(sample2Lines)-1]
	}

	sample3Lines = strings.Split(rawSample3, "\n")
	if len(sample3Lines) > 0 && sample3Lines[len(sample3Lines)-1] == "" {
		sample3Lines = sample3Lines[0 : len(sample3Lines)-1]
	}

	os.Exit(m.Run())
}
