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
	"fmt"
	"os"
	"sort"
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

func TestParseInput(t *testing.T) {
	input, err := parseInput(sampleLines, false)
	if err != nil {
		t.Fatal(err)
	}

	want := []Hand{
		Hand{"32T3K", HT_ONE, 765},
		Hand{"T55J5", HT_THREE, 684},
		Hand{"KK677", HT_TWO, 28},
		Hand{"KTJJT", HT_TWO, 220},
		Hand{"QQQJA", HT_THREE, 483},
	}

	if diff := cmp.Diff(want, input); diff != "" {
		t.Errorf("parseInput(sampleLines) mismatch; -want,+got\n%s\n", diff)
	}
}

func TestHandType(t *testing.T) {
	type TestCase struct {
		in   string
		wild bool
		want HandType
	}

	testCases := []TestCase{
		TestCase{"AAAAA", false, HT_FIVE},
		TestCase{"AA8AA", false, HT_FOUR},
		TestCase{"23332", false, HT_FULL},
		TestCase{"TTT98", false, HT_THREE},
		TestCase{"23432", false, HT_TWO},
		TestCase{"A23A4", false, HT_ONE},
		TestCase{"23456", false, HT_HIGH},

		TestCase{"32T3K", false, HT_ONE},
		TestCase{"T55J5", false, HT_THREE},
		TestCase{"KK677", false, HT_TWO},
		TestCase{"KTJJT", false, HT_TWO},
		TestCase{"QQQJA", false, HT_THREE},

		TestCase{"32T3K", true, HT_ONE},
		TestCase{"T55J5", true, HT_FOUR},
		TestCase{"KK677", true, HT_TWO},
		TestCase{"KTJJT", true, HT_FOUR},
		TestCase{"QQQJA", true, HT_FOUR},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s%v", tc.in, tc.wild), func(t *testing.T) {
			if got := NewHand(tc.in, 0, tc.wild).Type; got != tc.want {
				t.Errorf("%v got %v want %v", tc.in, got, tc.want)
			}
		})
	}
}

func TestHandStrength(t *testing.T) {
	type TestCase struct {
		in   []string
		wild bool
	}

	testCases := []TestCase{
		TestCase{
			in:   []string{"A", "K", "Q", "J", "T", "9", "8", "7", "6", "5", "4", "3", "2"},
			wild: false,
		},
		TestCase{
			in:   []string{"A", "K", "Q", "T", "9", "8", "7", "6", "5", "4", "3", "2", "J"},
			wild: true,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			hands := []Hand{}
			for _, val := range tc.in {
				hands = append(hands, NewHand(val, 0, tc.wild))
			}

			sort.Slice(hands, func(i, j int) bool {
				return hands[i].StrongerThan(hands[j], tc.wild)
			})

			got := []string{}
			for _, hand := range hands {
				got = append(got, hand.Val)
			}

			if diff := cmp.Diff(tc.in, got); diff != "" {
				t.Errorf("strength sort mismatch; -want,+got\n%s\n", diff)
			}
		})
	}
}

func TestSolveA(t *testing.T) {
	if got, want := solveA(sampleLines), int64(6440); got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	if got, want := solveB(sampleLines), int64(5905); got != want {
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
