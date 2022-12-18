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
	"reflect"
	"strings"
	"testing"

	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	//go:embed sample.txt
	rawSample   string
	sampleLines []string
)

func TestAllNeighbors(t *testing.T) {
	type TestCase struct {
		start string
		want  []string
	}

	testCases := []TestCase{
		TestCase{
			start: "0,0",
			want:  []string{"1,0", "0,1"},
		},
	}

	g, _, _, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}
	ac := &astarClient{g}

	for _, tc := range testCases {
		t.Run(tc.start, func(t *testing.T) {
			neighbors := ac.AllNeighbors(tc.start)
			if !reflect.DeepEqual(neighbors, tc.want) {
				t.Errorf("AllNeighbors(%s) = %v, want %v",
					tc.start, neighbors, tc.want)
			}
		})
	}
}

func TestSolveA(t *testing.T) {
	g, start, end, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(g, start, end), 31; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	g, _, end, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(g, end), 29; got != want {
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