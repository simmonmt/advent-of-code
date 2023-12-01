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
	"reflect"
	"strings"
	"testing"

	"github.com/simmonmt/aoc/2023/common/logger"
)

var (
	//go:embed sample.txt
	rawSampleA   string
	sampleALines []string

	//go:embed sample2.txt
	rawSampleB   string
	sampleBLines []string
)

func TestParseInput(t *testing.T) {
	input, err := parseInput(sampleALines)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(input, sampleALines) {
		t.Errorf("parseInput(sampleALines) = %v, want %v",
			input, sampleALines)
	}
}

func TestSolveA(t *testing.T) {
	input, err := parseInput(sampleALines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(input), 142; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	input, err := parseInput(sampleBLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(input), 281; got != want {
		t.Errorf("solveB(sample) = %v, want %v", got, want)
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)

	sampleALines = strings.Split(rawSampleA, "\n")
	if len(sampleALines) > 0 && sampleALines[len(sampleALines)-1] == "" {
		sampleALines = sampleALines[0 : len(sampleALines)-1]
	}

	sampleBLines = strings.Split(rawSampleB, "\n")
	if len(sampleBLines) > 0 && sampleBLines[len(sampleBLines)-1] == "" {
		sampleBLines = sampleBLines[0 : len(sampleBLines)-1]
	}

	os.Exit(m.Run())
}
