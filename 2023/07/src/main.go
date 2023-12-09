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
	"flag"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type HandType int

const (
	HT_FIVE HandType = iota
	HT_FOUR
	HT_FULL
	HT_THREE
	HT_TWO
	HT_ONE
	HT_HIGH
)

func calcType(hand string, jokersWild bool) HandType {
	m := map[rune]int{}
	maxRune := ' '
	maxRuneNum := -1
	for _, r := range hand {
		m[r]++

		if r != 'J' && (maxRuneNum == -1 || m[r] > maxRuneNum) {
			maxRune = r
			maxRuneNum = m[r]
		}
	}

	if jokersWild {
		if numJokers := m['J']; numJokers > 0 {
			m[maxRune] += numJokers
			delete(m, 'J')
		}
	}

	sizes := []int{}
	for _, num := range m {
		sizes = append(sizes, num)
	}
	sort.Ints(sizes)

	if len(sizes) == 1 {
		return HT_FIVE
	}
	if len(sizes) == 2 && sizes[0] == 1 {
		return HT_FOUR
	}
	if len(sizes) == 2 && sizes[0] == 2 {
		return HT_FULL
	}
	if len(sizes) == 3 && sizes[len(sizes)-1] == 3 {
		return HT_THREE
	}
	if len(sizes) == 3 && sizes[len(sizes)-1] == 2 && sizes[len(sizes)-2] == 2 {
		return HT_TWO
	}
	if len(sizes) == 4 {
		return HT_ONE
	}
	if len(sizes) == 5 {
		return HT_HIGH
	}
	panic("unknown " + hand)
}

type Hand struct {
	Val  string
	Type HandType
	Bid  int
}

func NewHand(hand string, bid int, jokersWild bool) Hand {
	return Hand{hand, calcType(hand, jokersWild), bid}
}

var (
	normalFaceStrength = map[rune]int{
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
		'T': 10,
		'J': 11,
		'Q': 12,
		'K': 13,
		'A': 14,
	}
	wildFaceStrength = map[rune]int{
		'J': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
		'T': 10,
		'Q': 12,
		'K': 13,
		'A': 14,
	}
)

func (h Hand) StrongerThan(oh Hand, jokersWild bool) bool {
	faceStrength := normalFaceStrength
	if jokersWild {
		faceStrength = wildFaceStrength
	}

	if h.Type != oh.Type {
		return h.Type < oh.Type
	}

	for i, r := range h.Val {
		or := rune(oh.Val[i])

		if s, os := faceStrength[r], faceStrength[or]; s != os {
			result := "stronger"
			if s < os {
				result = "weaker"
			}
			logger.Infof("%s %s than %s because %d vs %d for %c %c", h.Val, result, oh.Val, s, os, r, or)
			return s > os
		}
	}

	panic("same")
}

func parseInput(lines []string, jokersWild bool) ([]Hand, error) {
	hands := []Hand{}
	for _, line := range lines {
		a, b, found := strings.Cut(line, " ")
		if !found {
			return nil, fmt.Errorf("bad hand")
		}

		bid, err := strconv.Atoi(b)
		if err != nil {
			return nil, fmt.Errorf("bad bid")
		}

		hands = append(hands, NewHand(a, bid, jokersWild))
	}

	return hands, nil
}

func solveA(lines []string) int64 {
	hands, err := parseInput(lines, false)
	if err != nil {
		panic("bad input")
	}

	sort.Slice(hands, func(i, j int) bool {
		return hands[i].StrongerThan(hands[j], false)
	})

	logger.Infof("got hands %v", hands)

	out := int64(0)
	for i, hand := range hands {
		out += int64(len(hands)-i) * int64(hand.Bid)
	}
	return out
}

func solveB(lines []string) int64 {
	hands, err := parseInput(lines, true)
	if err != nil {
		panic("bad input")
	}

	sort.Slice(hands, func(i, j int) bool {
		return hands[i].StrongerThan(hands[j], true)
	})

	logger.Infof("got hands %v", hands)

	out := int64(0)
	for i, hand := range hands {
		out += int64(len(hands)-i) * int64(hand.Bid)
	}
	return out
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		logger.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		logger.Fatalf("failed to read input: %v", err)
	}

	fmt.Println("A", solveA(lines))
	fmt.Println("B", solveB(lines))
}
