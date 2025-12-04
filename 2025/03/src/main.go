package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/simmonmt/aoc/2025/common/filereader"
	"github.com/simmonmt/aoc/2025/common/logger"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

type Input struct {
	Lines []string
}

func parseInput(lines []string) (*Input, error) {
	return &Input{Lines: lines}, nil
}

func findJoltage(line string) int {
	h1 := '0'
	h1i := -1
	for i, r := range line[0 : len(line)-1] {
		if r > h1 {
			h1 = r
			h1i = i
		}
	}

	h2 := '0'
	for _, r := range line[h1i+1:] {
		if r > h2 {
			h2 = r
		}
	}

	return int((byte(h1)-'0')*10 + (byte(h2) - '0'))
}

func solveA(input *Input) int {
	out := 0
	for _, line := range input.Lines {
		j := findJoltage(line)
		out += j
	}
	return out
}

func solveB(input *Input) int {
	return -1
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *input == "" {
		logger.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		logger.Fatalf("failed to read input: %v", err)
	}

	input, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
