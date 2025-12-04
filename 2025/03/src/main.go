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

func findJoltage(line string, n int) int {
	//logger.Infof("in line %v len %d n %d", line, len(line), n)
	digits := make([]rune, n)
	idx := make([]int, n)

	for i := range n {
		start := 0
		if i > 0 {
			start = idx[i-1] + 1
		}
		lim := len(line) - (n - i - 1)
		//logger.Infof("highest among %v", line[start:lim])

		h := '0'
		hi := -1
		for j, r := range line[start:lim] {
			if r > h {
				h = r
				hi = j
			}
		}

		hi += start
		//logger.Infof("found highest %v at %v", string(h), hi)
		digits[i] = h
		idx[i] = hi
	}

	out := 0
	for i := range n {
		out = 10*out + int(digits[i]-'0')
	}
	return out
}

func solveA(input *Input) int {
	out := 0
	for _, line := range input.Lines {
		j := findJoltage(line, 2)
		//logger.Infof("line %v j %v", line, j)
		out += j
	}
	return out
}

func solveB(input *Input) int {
	out := 0
	for _, line := range input.Lines {
		j := findJoltage(line, 12)
		//logger.Infof("line %v j %v", line, j)
		out += j
	}
	return out
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
