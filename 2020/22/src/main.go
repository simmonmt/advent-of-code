package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"

	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/intmath"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	playerPattern = regexp.MustCompile(`^Player ([0-9]+):$`)
)

func ScoreResults(decks *[2]*Deck) int {
	logger.LogF("== Post-game results ==")
	for _, deck := range decks {
		logger.LogF("Player %v's deck: %v", deck.Name(), deck.Cards())
	}

	winner := decks[0]
	if winner.Empty() {
		winner = decks[1]
	}

	score := 0
	cards := winner.Cards()
	for i := len(cards) - 1; i >= 0; i-- {
		score += cards[i] * (len(cards) - i)
	}

	return score
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		log.Fatal(err)
	}

	originalDecks := []*Deck{}
	for len(lines) > 0 {
		parts := playerPattern.FindStringSubmatch(lines[0])
		if parts == nil {
			log.Fatalf("bad player match: %v", lines[0])
		}

		name := parts[1]

		blank := -1
		for blank = 1; blank < len(lines); blank++ {
			if lines[blank] == "" {
				break
			}
		}

		deckLines := lines[1:blank]
		deckNums := []int{}
		for _, line := range deckLines {
			deckNums = append(deckNums, intmath.AtoiOrDie(line))
		}
		lines = lines[intmath.IntMin(len(lines), blank+1):]

		originalDecks = append(originalDecks,
			newDeck(name, deckNums, len(deckNums)*2))
	}

	if len(originalDecks) != 2 {
		panic("bad num")
	}

	{
		decks := [2]*Deck{originalDecks[0].Clone(), originalDecks[1].Clone()}
		PlayNormalGame(&decks)
		fmt.Printf("A: %v\n", ScoreResults(&decks))
	}

	{
		decks := [2]*Deck{originalDecks[0].Clone(), originalDecks[1].Clone()}
		PlayRecursiveCards(&decks)
		fmt.Printf("B: %v\n", ScoreResults(&decks))
	}
}
