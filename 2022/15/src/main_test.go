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
	"github.com/simmonmt/aoc/2022/common/pos"
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

	want := []Sensor{
		Sensor{pos.P2{2, 18}, pos.P2{-2, 15}},
		Sensor{pos.P2{9, 16}, pos.P2{10, 16}},
	}

	if !reflect.DeepEqual(input, want) {
		t.Errorf("parseInput(sampleLines) = %v, want %v",
			input, want)
	}
}

func TestFurthestIntersections(t *testing.T) {
	type TestCase struct {
		sensor      Sensor
		row         int
		left, right int
	}

	testCases := []TestCase{
		TestCase{
			Sensor{pos.P2{8, 7}, pos.P2{2, 10}},
			10, 2, 14,
		},
		TestCase{
			Sensor{pos.P2{8, 7}, pos.P2{8, 10}},
			10, 8, 8,
		},
		TestCase{
			Sensor{pos.P2{8, 7}, pos.P2{8, 10}},
			9, 7, 9,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			left, right := furthestIntersections(tc.sensor, tc.row)
			if left != tc.left || right != tc.right {
				t.Errorf("furthestIntersections(%v, %v) = %v, %v, want %v, %v",
					tc.sensor, tc.row, left, right, tc.left, tc.right)
			}
		})
	}
}

func TestSolveA(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(input, 10), 26; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(input, 0, 20), 56000011; got != want {
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
