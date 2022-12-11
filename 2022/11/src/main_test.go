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

var (
	monkey1Str = []string{
		"Monkey 14:",
		"  Starting items: 79, 60, 97",
		"  Operation: new = old * old",
		"  Test: divisible by 13",
		"    If true: throw to monkey 1",
		"    If false: throw to monkey 3",
	}
)

func parseMonkeyOrDie(t *testing.T, lines []string) *Monkey {
	m, err := parseMonkey(lines)
	if err != nil {
		t.Fatalf("bad monkey parse: %v", err)
	}
	return m
}

func TestParseMonkey(t *testing.T) {
	want := &Monkey{
		id:          14,
		items:       []int{79, 60, 97},
		op:          &SquareOp{},
		testDivisor: 13,
		trueDest:    1,
		falseDest:   3,
	}

	got, err := parseMonkey(monkey1Str)
	if err != nil {
		t.Fatalf("failed to parse monkey: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want monkey %+v, got %+v", want, got)
	}
}

func TestMonkeyStep(t *testing.T) {
	m := parseMonkeyOrDie(t, monkey1Str)

	wantOut := map[int][]int{
		1: []int{2080},
		3: []int{1200, 3136},
	}

	wantMonkey := &Monkey{
		id:             14,
		items:          []int{},
		numInspections: 3,
		op:             &SquareOp{},
		testDivisor:    13,
		trueDest:       1,
		falseDest:      3,
	}

	if out := m.Step(); !reflect.DeepEqual(out, wantOut) {
		t.Errorf("Step => %v, want %v", out, wantOut)
	}

	if !reflect.DeepEqual(m, wantMonkey) {
		t.Errorf("want result monkey %+v, got %+v", wantMonkey, m)
	}
}

func TestSolveA(t *testing.T) {
	monkeys, err := parseMonkeys(sampleLines)
	if err != nil {
		t.Fatalf("bad parse: %v", err)
	}

	if got, want := solveA(monkeys), 10605; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	monkeys, err := parseMonkeys(sampleLines)
	if err != nil {
		t.Fatalf("bad parse: %v", err)
	}

	if got, want := solveB(monkeys), -1; got != want {
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
