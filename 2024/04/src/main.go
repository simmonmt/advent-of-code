package main

import (
	"flag"
	"fmt"

	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/grid"
	"github.com/simmonmt/aoc/2024/common/logger"
	"github.com/simmonmt/aoc/2024/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	xmas = []rune{'X', 'M', 'A', 'S'}
	mas  = []rune{'M', 'A', 'S'}

	allDirs = []pos.P2{
		pos.P2{X: 1, Y: 0},
		pos.P2{X: 1, Y: 1},
		pos.P2{X: 0, Y: 1},
		pos.P2{X: -1, Y: 1},
		pos.P2{X: -1, Y: 0},
		pos.P2{X: -1, Y: -1},
		pos.P2{X: 0, Y: -1},
		pos.P2{X: 1, Y: -1},
	}

	diagDirs = []pos.P2{
		pos.P2{X: 1, Y: 1},
		pos.P2{X: -1, Y: 1},
		pos.P2{X: -1, Y: -1},
		pos.P2{X: 1, Y: -1},
	}
)

func parseInput(lines []string) (*grid.Grid[rune], error) {
	return grid.NewFromLines[rune](lines,
		func(p pos.P2, r rune) (rune, error) {
			return r, nil
		})
}

func findFrom(g *grid.Grid[rune], p, rel pos.P2, cur rune, toFind []rune) bool {
	if cur != toFind[0] {
		return false
	}

	if len(toFind) == 1 {
		return true
	}

	toFind = toFind[1:]
	np := p
	np.Add(rel)

	cur, ok := g.Get(np)
	if !ok {
		return false
	}

	return findFrom(g, np, rel, cur, toFind)
}

func solveA(g *grid.Grid[rune]) int64 {
	total := 0
	g.Walk(func(p pos.P2, r rune) {
		for _, rel := range allDirs {
			if findFrom(g, p, rel, r, xmas) {
				total++
			}
		}
	})

	return int64(total)
}

func solveB(g *grid.Grid[rune]) int64 {
	as := map[pos.P2]int{}
	total := 0

	g.Walk(func(p pos.P2, r rune) {
		for _, rel := range diagDirs {
			if !findFrom(g, p, rel, r, mas) {
				continue
			}

			aPos := p
			aPos.Add(rel)
			as[aPos]++

			if n := as[aPos]; n == 2 {
				total++
			} else if n > 2 {
				panic("too many")
			}
		}
	})

	return int64(total)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

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
