package main

import (
	"flag"
	"fmt"
	"iter"
	"log"
	"os"
	"runtime/pprof"

	"github.com/simmonmt/aoc/2025/common/filereader"
	"github.com/simmonmt/aoc/2025/common/logger"
	"github.com/simmonmt/aoc/2025/common/mtsmath"
	"github.com/simmonmt/aoc/2025/common/pos"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

type Input struct {
	Marks []pos.P2
}

func parseInput(lines []string) (*Input, error) {
	marks := []pos.P2{}
	for i, line := range lines {
		mark, err := pos.P2FromString(line)
		if err != nil {
			return nil, fmt.Errorf("%d: bad pos: %v", i+1, err)
		}
		marks = append(marks, mark)
	}
	return &Input{Marks: marks}, nil
}

type Pair struct {
	A, B pos.P2
}

func MakePairs(ps []pos.P2) iter.Seq[Pair] {
	return func(yield func(Pair) bool) {
		for i, a := range ps[0 : len(ps)-1] {
			for _, b := range ps[i+1:] {
				if !yield(Pair{a, b}) {
					return
				}
			}
		}
	}
}

func solveA(input *Input) int {
	maxArea := 0
	for pair := range MakePairs(input.Marks) {
		area := (mtsmath.Abs(pair.A.X-pair.B.X) + 1) *
			(mtsmath.Abs(pair.A.Y-pair.B.Y) + 1)
		//fmt.Println(pair, area)
		if area > maxArea {
			maxArea = area
		}
	}
	return maxArea
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
