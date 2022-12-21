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
	"strconv"
	"strings"
	"testing"

	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	//go:embed sample.txt
	rawSample   string
	sampleLines []string
)

func TestParseInput(t *testing.T) {
	input, err := parseInput(sampleLines[0:2])
	if err != nil {
		t.Fatal(err)
	}

	want := map[string]*Monkey{
		"root": &Monkey{name: "root", op: ADD, a: "pppw", b: "sjmn"},
		"dbpl": &Monkey{name: "dbpl", op: IMM, imm: 5},
	}

	if !reflect.DeepEqual(input, want) {
		t.Errorf("parseInput(sampleLines) = %v, want %v",
			input, sampleLines)
	}
}

func TestUnsolveOp(t *testing.T) {
	type TestCase struct {
		op        Operation
		result    int64
		known     int64
		knownLeft bool
		want      int64
	}

	testCases := []TestCase{
		TestCase{ADD, 43, 42, true, 1},
		TestCase{ADD, 43, 42, false, 1},
		TestCase{SUB, 43, 50, true, 7},
		TestCase{SUB, 43, 50, false, 93},
		TestCase{MUL, 10, 5, true, 2},
		TestCase{MUL, 10, 5, false, 2},
		TestCase{DIV, 10, 20, true, 2},
		TestCase{DIV, 10, 20, false, 200},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := unsolveOp(tc.op, tc.result, tc.known, tc.knownLeft); got != tc.want {
				t.Errorf("unsolveOp(%v,%v,%v,%v) = %v, want %v",
					tc.op, tc.result, tc.known, tc.knownLeft, got, tc.want)
			}
		})
	}
}

func TestSolveA(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(input), int64(152); got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(input), int64(301); got != want {
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
