package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/simmonmt/aoc/2019/22/a/src/puzzle"
	"github.com/simmonmt/aoc/2019/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func readInput(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := []string{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lines, nil
}

func findValue(cards []int, want int) int {
	for i, card := range cards {
		if card == want {
			return i
		}
	}
	return -1
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	cmds, err := puzzle.ParseCommands(lines)
	if err != nil {
		log.Fatal(err)
	}

	cards := make([]int, 10007)
	for i := range cards {
		cards[i] = i
	}

	cards = puzzle.RunCommands(cards, cmds)

	fmt.Printf("card 2019 location: %v\n", findValue(cards, 2019))
}
