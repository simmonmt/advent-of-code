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

	"github.com/google/go-cmp/cmp"
	"github.com/simmonmt/aoc/2023/common/logger"
)

var (
	//go:embed sample.txt
	rawSample   string
	sampleLines []string
)

func oneElemRanges(nums ...int64) []Range {
	out := []Range{}
	for _, num := range nums {
		out = append(out, Range{num, num})
	}
	return out
}

func parseGardenMapOrDie(lines []string) *GardenMap {
	gm, err := parseGardenMap(lines)
	if err != nil {
		logger.Fatalf("parseGardenMap fail: %v", err)
	}
	return gm
}

func TestParseInput(t *testing.T) {
	input := []string{
		"seeds: 79 14 55 13",
		"",
		"seed-to-soil map:",
		"50 98 2",
		"52 50 48",
		"",
		"soil-to-nothing map:",
		"1 2 3",
	}

	wantSeedNums := []int64{79, 14, 55, 13}

	wantGardenMaps := []*GardenMap{
		&GardenMap{
			From: "seed",
			To:   "soil",
			Ranges: []*RangeMap{
				&RangeMap{Src: Range{Lo: 98, Hi: 99}, Dest: Range{Lo: 50, Hi: 51}},
				&RangeMap{Src: Range{Lo: 50, Hi: 97}, Dest: Range{Lo: 52, Hi: 99}},
			},
		},
		&GardenMap{
			From: "soil",
			To:   "nothing",
			Ranges: []*RangeMap{
				&RangeMap{Src: Range{Lo: 2, Hi: 4}, Dest: Range{Lo: 1, Hi: 3}},
			},
		},
	}

	seedNums, gardenMaps, err := parseInput(input)
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Errorf("parseInput(sampleLines) = _, _, %v, want nil", err)
	}

	if diff := cmp.Diff(wantSeedNums, seedNums); diff != "" {
		t.Errorf("parseInput seedNums mismatch +want, -got:\n%s\n", diff)
	}
	if diff := cmp.Diff(wantGardenMaps, gardenMaps); diff != "" {
		t.Errorf("parseInput gardenMaps mismatch +want, -got:\n%s\n", diff)
	}
}

func TestSolveLevel(t *testing.T) {
	type TestCase struct {
		in   []Range
		gm   *GardenMap
		want []Range
	}

	testCases := []TestCase{
		TestCase{
			in: oneElemRanges(79, 14, 55, 13),
			gm: parseGardenMapOrDie([]string{
				"seed-to-soil map:",
				"50 98 2",
				"52 50 48",
			}),
			want: oneElemRanges(13, 14, 57, 81),
		},
		TestCase{
			in: oneElemRanges(13, 14, 57, 81),
			gm: parseGardenMapOrDie([]string{
				"soil-to-fertilizer map:",
				"0 15 37",
				"37 52 2",
				"39 0 15",
			}),
			want: oneElemRanges(52, 53, 57, 81),
		},
		TestCase{
			in: oneElemRanges(52, 53, 57, 81),
			gm: parseGardenMapOrDie([]string{
				"fertilizer-to-water map:",
				"49 53 8",
				"0 11 42",
				"42 0 7",
				"57 7 4",
			}),
			want: oneElemRanges(41, 49, 53, 81),
		},
		TestCase{
			in: oneElemRanges(41, 49, 53, 81),
			gm: parseGardenMapOrDie([]string{
				"water-to-light map:",
				"88 18 7",
				"18 25 70",
			}),
			want: oneElemRanges(34, 42, 46, 74),
		},
		TestCase{
			in: oneElemRanges(34, 42, 46, 74),
			gm: parseGardenMapOrDie([]string{
				"light-to-temperature map:",
				"45 77 23",
				"81 45 19",
				"68 64 13",
			}),
			want: oneElemRanges(34, 42, 82, 78),
		},
		TestCase{
			in: oneElemRanges(34, 42, 78, 82),
			gm: parseGardenMapOrDie([]string{
				"temperature-to-humidity map:",
				"0 69 1",
				"1 0 69",
			}),
			want: oneElemRanges(35, 43, 78, 82),
		},
		TestCase{
			in: oneElemRanges(35, 43, 78, 82),
			gm: parseGardenMapOrDie([]string{
				"humidity-to-location map:",
				"60 56 37",
				"56 93 4",
			}),
			want: oneElemRanges(35, 43, 82, 86),
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sortRanges(tc.in)
			sortGardenMap(tc.gm)

			got := solveLevel(tc.in, tc.gm)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("solveLevel(%v, %v) mismatch -want, +got\n%s\n",
					tc.in, tc.gm, diff)
			}
		})
	}
}

func TestSolveA(t *testing.T) {
	seedNums, gardenMaps, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(seedNums, gardenMaps), int64(35); got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	seedNums, gardenMaps, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(seedNums, gardenMaps), int64(46); got != want {
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
