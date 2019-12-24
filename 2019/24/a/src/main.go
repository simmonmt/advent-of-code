package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/simmonmt/aoc/2019/24/a/src/puzzle"
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

	b := puzzle.NewBoard(lines)
	cache := map[string]int{b.Hash(): 0}

	for i := 1; ; i++ {
		b = b.Evolve()
		h := b.Hash()
		if when, found := cache[h]; found {
			fmt.Printf("repeat %d and %d\n", i, when)
			break
		}
		cache[h] = i
	}

	fmt.Printf("biodiversity %v\n", b.Biodiversity())
}
