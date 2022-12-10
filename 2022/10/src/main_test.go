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
	"fmt"
	"log"
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

func parseInstructionsOrDie(lines []string) []Inst {
	insts, err := parseInstructions(lines)
	if err != nil {
		log.Fatal(err)
	}
	return insts
}

func TestParseInstructions(t *testing.T) {
	in := []string{"addx 15", "addx -11", "noop"}
	want := []Inst{Inst{"addx", 15}, Inst{"addx", -11}, Inst{"noop", 0}}

	if got, err := parseInstructions(in); err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("parseInstructions(_) = %v, %v, want %v, nil", got, err, want)
	}
}

func TestOverSprite(t *testing.T) {
	type TestCase struct {
		sp, x int
		want  bool
	}

	testCases := []TestCase{
		TestCase{1, 0, true},
		TestCase{1, 1, true},
		TestCase{1, 2, true},
		TestCase{1, 3, false},

		TestCase{5, 3, false},
		TestCase{5, 4, true},
		TestCase{5, 5, true},
		TestCase{5, 6, true},
		TestCase{5, 7, false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d_%d", tc.sp, tc.x), func(t *testing.T) {
			if got := overSprite(tc.sp, tc.x); got != tc.want {
				t.Errorf("overSprite(%d, %d) = %v, want %v",
					tc.sp, tc.x, got, tc.want)
			}
		})
	}
}

func TestSolveA(t *testing.T) {
	if got, want := solveA(parseInstructionsOrDie(sampleLines)), 13140; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	want := []string{
		"##..##..##..##..##..##..##..##..##..##..",
		"###...###...###...###...###...###...###.",
		"####....####....####....####....####....",
		"#####.....#####.....#####.....#####.....",
		"######......######......######......####",
		"#######.......#######.......#######.....",
	}

	got := solveB(parseInstructionsOrDie(sampleLines))

	for i := 0; i < len(want); i++ {
		if i >= len(got) {
			t.Errorf("missing line %d in got", i)
			continue
		}

		if want[i] != got[i] {
			fmt.Printf("want %d: %v\n", i, want[i])
			fmt.Printf("got  %d: %v\n", i, got[i])

			t.Errorf("line %d mismatch", i)
		}
	}

	if t.Failed() {
		fmt.Println("got dump")
		for i, line := range got {
			fmt.Printf("%d: %v\n", i, line)
		}
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
