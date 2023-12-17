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
)

func TestIsOneBit(t *testing.T) {
	type TestCase struct {
		a, b  int64
		is    bool
		which int
	}

	testCases := []TestCase{
		TestCase{0x13ef, 0x13f0, false, -1},
		TestCase{0x13ef, 0x13ee, true, 0},
		TestCase{0x13ef, 0x13cf, true, 5},
		TestCase{0x4c, 0xc, true, 6},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if is, which := isOneBit(tc.a, tc.b); is != tc.is || (tc.is && which != tc.which) {
				t.Errorf("isOneBit(%v,%v) = %v, %v; want %v, %v",
					tc.a, tc.b, is, which, tc.is, tc.which)
			}
		})
	}
}

func TestSolveA(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(input), 405; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(input), 400; got != want {
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
