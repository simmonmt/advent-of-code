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
	"strconv"
	"strings"
	"testing"

	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	//go:embed sample.txt
	rawSample  string
	sampleLine string

	//go:embed input.txt
	rawInput  string
	inputLine string
)

func TestValidateMath(t *testing.T) {
	for i, dirs := range []string{sampleLine, inputLine} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			first, second := findRepeat(dirs)

			prologueLen := first.lastPartIdx + 1
			prologueHeight := first.height

			repLen := second.lastPartIdx - first.lastPartIdx
			repHeight := second.height - first.height

			trialLen := prologueLen + repLen*7
			wantHeight := prologueHeight + repHeight*7

			if got := measureHeight(dirs, trialLen); got != wantHeight {
				t.Errorf("got %v, want %v", got, wantHeight)
			}
		})
	}
}

func TestMeasureTallHeight(t *testing.T) {
	for i, dirs := range []string{sampleLine, inputLine} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			height := 100000
			want := int64(measureHeight(dirs, height))
			got := measureTallHeight(dirs, int64(height))

			if got != want {
				t.Errorf("want %v, got %v", want, got)
			}
		})
	}
}

func TestSolveA(t *testing.T) {
	if got, want := solveA(sampleLine), 3068; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
	if got, want := solveA(inputLine), 3133; got != want {
		t.Errorf("solveA(input) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	if got, want := solveB(sampleLine), int64(1514285714288); got != want {
		t.Errorf("solveB(sample) = %v, want %v", got, want)
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)

	sampleLine = strings.Split(rawSample, "\n")[0]
	inputLine = strings.Split(rawInput, "\n")[0]

	os.Exit(m.Run())
}
