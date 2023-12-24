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
	"github.com/simmonmt/aoc/2023/common/pos"
)

var (
	//go:embed sample.txt
	rawSample   string
	sampleLines []string
)

func TestParseInput(t *testing.T) {
	ins := []string{
		"1,0,1~1,2,1",
		"0,2,3~2,2,3",
		"0,0,1~0,0,4",
	}

	wants := []*Block{
		&Block{Num: 1, R: Range{From: pos.P3{X: 1, Y: 2, Z: 1}, To: pos.P3{X: 1, Y: 0, Z: 1}}},
		&Block{Num: 2, R: Range{From: pos.P3{X: 2, Y: 2, Z: 3}, To: pos.P3{X: 0, Y: 2, Z: 3}}},
		&Block{Num: 3, R: Range{From: pos.P3{X: 0, Y: 0, Z: 4}, To: pos.P3{X: 0, Y: 0, Z: 1}}},
	}

	gots, err := parseInput(ins)
	if err != nil {
		t.Errorf(`parseInput(ins) = _, %v, want _, nil`, err)
		return
	}

	if diff := cmp.Diff(gots, wants); diff != "" {
		t.Errorf("parseInput mismatch; -want,+got:\n%s\n", diff)
	}
}

func TestSolveA(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(input), 5; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(input), 7; got != want {
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
