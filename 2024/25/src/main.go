package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/grid"
	"github.com/simmonmt/aoc/2024/common/logger"
	"github.com/simmonmt/aoc/2024/common/pos"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

type Input struct {
	Locks  [][]int
	Keys   [][]int
	Height int
}

func parseItem(lines []string) []int {
	g, _ := grid.NewFromLines[bool](lines, func(p pos.P2, r rune) (bool, error) {
		return r == '#', nil
	})

	nums := make([]int, g.Width())
	g.Walk(func(p pos.P2, v bool) {
		if v {
			nums[p.X]++
		}
	})

	return nums
}

func parseInput(lines []string) (*Input, error) {
	groups := filereader.BlankSeparatedGroupsFromLines(lines)

	locks := [][]int{}
	keys := [][]int{}
	for _, group := range groups {
		if group[0][0] == '#' {
			locks = append(locks, parseItem(group[1:]))
		} else {
			keys = append(keys, parseItem(group[0:len(group)-1]))
		}
	}

	return &Input{Locks: locks, Keys: keys, Height: len(groups[0]) - 2}, nil
}

func solveA(input *Input) int64 {
	num := 0
	for _, lock := range input.Locks {
		for _, key := range input.Keys {
			fits := true
			for i := range lock {
				if lock[i]+key[i] > input.Height {
					fits = false
					break
				}
			}

			if fits {
				num++
			}
		}
	}

	return int64(num)
}

func solveB(input *Input) int64 {
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
