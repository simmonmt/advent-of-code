package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/simmonmt/aoc/2025/common/filereader"
	"github.com/simmonmt/aoc/2025/common/grid"
	"github.com/simmonmt/aoc/2025/common/logger"
	"github.com/simmonmt/aoc/2025/common/pos"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

type Input struct {
	Grid *grid.Grid[rune]
	Locs map[pos.P2]bool
}

func parseInput(lines []string) (*Input, error) {
	locs := map[pos.P2]bool{}
	g, _ := grid.NewFromLines(lines, func(p pos.P2, r rune) (rune, error) {
		if r == '@' {
			locs[p] = true
		}
		return r, nil
	})

	return &Input{g, locs}, nil
}

func canRemove(g *grid.Grid[rune], p pos.P2) bool {
	num := 0
	for _, n := range g.AllNeighbors(p, true) {
		if g.GetOr(n, '.') == '@' {
			num++
		}
	}
	return num < 4
}

func solveA(input *Input) int {
	seen := map[pos.P2]bool{}

	input.Grid.Walk(func(p pos.P2, r rune) {
		if r == '.' {
			return
		}

		if canRemove(input.Grid, p) {
			seen[p] = true
		}
	})

	return len(seen)
}

func oneRound(g *grid.Grid[rune], locs map[pos.P2]bool) int {
	toRemove := []pos.P2{}
	for l := range locs {
		if canRemove(g, l) {
			toRemove = append(toRemove, l)
		}
	}

	for _, p := range toRemove {
		delete(locs, p)
		g.Set(p, '.')
	}

	return len(toRemove)
}

func solveB(input *Input) int {
	g := input.Grid.Clone()
	locs := map[pos.P2]bool{}
	for l := range input.Locs {
		locs[l] = true
	}

	tot := 0
	for {
		n := oneRound(g, locs)
		if n == 0 {
			break
		}
		tot += n
	}

	return tot
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
