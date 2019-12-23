package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/simmonmt/aoc/2019/20/b/src/puzzle"
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

	board := puzzle.NewBoard(lines)
	allPaths := puzzle.FindAllPathsFromAllPortals(board)

	start := board.Gate("AA").PortalOut()
	end := board.Gate("ZZ").GateOut()

	cost, found := puzzle.Solve(board, allPaths, start, end)
	if !found {
		log.Fatal("no path")
	}

	// Solve overcounts by 1 -- it ends in ZZ's gate rather than
	// ZZ's portal.
	fmt.Printf("result %d\n", cost-1)
}
