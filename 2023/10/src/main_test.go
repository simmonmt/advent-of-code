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

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
)

var (
	//go:embed samples.txt
	rawSamples string
	samples    [][]string
)

func TestSolveA(t *testing.T) {
	wants := []int{4, 4, 8, 11, 23, 70, 8, 80}

	for i, sample := range samples {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			g, start, err := NewBoard(sample)
			if err != nil {
				t.Fatal(err)
			}

			if got := solveA(g, start); got != wants[i] {
				t.Errorf("solveA(sample) = %v, want %v", got, wants[i])
			}
		})
	}
}

func TestSolveB(t *testing.T) {
	wants := []int{1, 1, 1, 8, 4, 8, 9, 10}

	for i, sample := range samples {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			g, start, err := NewBoard(sample)
			if err != nil {
				t.Fatal(err)
			}

			if got := solveB(g, start); got != wants[i] {
				t.Errorf("solveB(sample) = %v, want %v", got, wants[i])
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)

	var err error
	samples, err = filereader.BlankSeparatedGroupsFromLines(
		strings.Split(rawSamples, "\n"))
	if err != nil {
		panic("bad sample")
	}

	os.Exit(m.Run())
}
