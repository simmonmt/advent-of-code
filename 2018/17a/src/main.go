package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"intmath"
	"lib"
	"logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
)

func parseRange(str string) (min, max int) {
	parts := strings.SplitN(str, "..", 2)
	if len(parts) == 1 {
		val := intmath.AtoiOrDie(parts[0])
		return val, val
	} else {
		return intmath.AtoiOrDie(parts[0]), intmath.AtoiOrDie(parts[1])
	}
}

func readInput() ([]lib.InputLine, error) {
	lines := []string{}
	inputLines := []lib.InputLine{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ", ")
		sort.Strings(parts)

		xpart := parts[0]
		ypart := parts[1]

		xmin, xmax := parseRange(xpart[2:])
		ymin, ymax := parseRange(ypart[2:])

		inputLines = append(inputLines, lib.InputLine{xmin, xmax, ymin, ymax})

		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return inputLines, nil
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	lines, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	board := lib.NewBoard(lines)
	board.Dump()
}
