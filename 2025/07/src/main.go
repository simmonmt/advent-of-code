package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/simmonmt/aoc/2025/common/dir"
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
	Start pos.P2
	Grid  *grid.Grid[rune]
}

func parseInput(lines []string) (*Input, error) {
	var start pos.P2
	g, err := grid.NewFromLines(lines, func(p pos.P2, r rune) (rune, error) {
		if r == 'S' {
			start = p
		}
		return r, nil
	})
	if err != nil {
		return nil, err
	}

	return &Input{Start: start, Grid: g}, nil
}

func solveA(input *Input) int {
	g := input.Grid.Clone()
	todo := []pos.P2{input.Start}

	numSplits := 0
	for len(todo) != 0 {
		next := []pos.P2{}

		for _, p := range todo {
			p2 := dir.DIR_SOUTH.From(p)

			switch g.GetOr(p2, 'X') {
			case '.':
				next = append(next, p2)
				g.Set(p2, '|')
			case '^':
				newSplit := false
				splits := []pos.P2{dir.DIR_EAST.From(p2), dir.DIR_WEST.From(p2)}
				for _, split := range splits {
					if g.GetOr(split, 'X') == '.' {
						newSplit = true
						next = append(next, split)
						g.Set(split, '|')
					}
				}
				if newSplit {
					numSplits++
				}
			}
		}

		//g.Dump(true, grid.RuneDumper)
		//fmt.Println(numSplits)

		todo = next
	}

	return numSplits
}

func findScore(p pos.P2, nextSplitter map[pos.P2]pos.P2, scores map[pos.P2]int) int {
	if score, found := scores[p]; found {
		//fmt.Println(p, "found in cache", score)
		return score
	}

	sp, found := nextSplitter[p]
	if !found {
		//fmt.Println(p, "fell off", 1)
		scores[p] = 1
		return 1
	}

	left := findScore(dir.DIR_WEST.From(sp), nextSplitter, scores)
	right := findScore(dir.DIR_EAST.From(sp), nextSplitter, scores)
	//fmt.Println(p, "splitter", left, right)
	tot := left + right
	scores[p] = tot
	return tot
}

func solveB(input *Input) int {
	splitters := []pos.P2{}
	nextSplitter := map[pos.P2]pos.P2{}
	scores := map[pos.P2]int{}

	g := input.Grid
	for x := 0; x < g.Width(); x++ {
		q := []pos.P2{}
		for y := 0; y < g.Height(); y++ {
			p := pos.P2{X: x, Y: y}
			if r := g.GetOr(p, 'X'); r == '.' || r == 'S' {
				q = append(q, p)
			} else if r == '^' {
				splitters = append(splitters, p)
				for _, qp := range q {
					nextSplitter[qp] = p
				}
				q = []pos.P2{}
			}
		}
	}

	score := findScore(input.Start, nextSplitter, scores)

	return score
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
