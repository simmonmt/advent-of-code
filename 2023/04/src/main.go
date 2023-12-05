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
	"regexp"
	"strconv"

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/strutil"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	cardPattern = regexp.MustCompile(`^Card +(\d+): ([0-9 ]+) \| ([0-9 ]+)$`)
)

type Card struct {
	Num     int
	Winning []int
	Have    []int
}

func parseCard(line string) (*Card, error) {
	parts := cardPattern.FindStringSubmatch(line)
	if parts == nil {
		return nil, fmt.Errorf("regex failed match")
	}

	card := &Card{}

	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("bad card number: %v", err)
	}
	card.Num = num

	winning, err := strutil.ListOfNumbers(parts[2])
	if err != nil {
		return nil, fmt.Errorf("bad winning nums: %v", err)
	}
	card.Winning = winning

	have, err := strutil.ListOfNumbers(parts[3])
	if err != nil {
		return nil, fmt.Errorf("bad have nums: %v", err)
	}
	card.Have = have

	return card, nil
}

func parseInput(lines []string) ([]*Card, error) {
	cards := []*Card{}
	for i, line := range lines {
		card, err := parseCard(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line %d: %v", i, err)
		}

		cards = append(cards, card)
	}
	return cards, nil
}

func findWinningNums(card *Card) []int {
	have := map[int]bool{}
	for _, num := range card.Have {
		have[num] = true
	}

	out := []int{}
	for _, num := range card.Winning {
		if _, ok := have[num]; ok {
			out = append(out, num)
		}
	}
	return out
}

func solveA(input []*Card) int {
	out := 0

	for _, card := range input {
		score := 0
		for _ = range findWinningNums(card) {
			if score == 0 {
				score = 1
			} else {
				score *= 2
			}
		}
		out += score
	}

	return out
}

func solveB(cards []*Card) int {
	deck := map[int]int{}
	for _, card := range cards {
		deck[card.Num] += 1
		mult := deck[card.Num]

		numWinning := len(findWinningNums(card))
		for i := 0; i < numWinning; i++ {
			cardNum := card.Num + 1 + i
			deck[cardNum] += mult
		}
	}

	out := 0
	for _, num := range deck {
		out += num
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
		logger.Fatalf("read failed: %v", err)
	}

	input, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("parse failed: %v", err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
