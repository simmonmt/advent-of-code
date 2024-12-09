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
	"github.com/simmonmt/aoc/2024/common/mtsmath"
	"github.com/simmonmt/aoc/2024/common/pos"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func parseInput(lines []string) (*grid.Grid[rune], error) {
	return grid.NewFromLines[rune](lines, func(p pos.P2, r rune) (rune, error) {
		return r, nil
	})
}

func solveA(g *grid.Grid[rune]) int64 {
	allAnts := map[rune][]pos.P2{}
	g.Walk(func(p pos.P2, r rune) {
		if r != '.' {
			allAnts[r] = append(allAnts[r], p)
		}
	})

	allPodes := map[pos.P2]bool{}
	for _, ants := range allAnts {
		for i := 0; i < len(ants)-1; i++ {
			for j := i + 1; j < len(ants); j++ {
				a, b := ants[i], ants[j]

				slx, sly := b.X-a.X, b.Y-a.Y

				podes := []pos.P2{
					pos.P2{X: b.X + slx, Y: b.Y + sly},
					pos.P2{X: a.X - slx, Y: a.Y - sly},
				}

				for _, pode := range podes {
					if g.IsValid(pode) {
						allPodes[pode] = true
					}
				}
			}
		}
	}

	return int64(len(allPodes))
}

func simplify(n, d int) (int, int) {
	for {
		gcd := int(mtsmath.GCD(int64(n), int64(d)))
		if gcd == 1 {
			return n, d
		}
		n /= gcd
		d /= gcd
	}
}

func extend(g *grid.Grid[rune], p pos.P2, slx, sly int, podes map[pos.P2]bool) {
	for {
		p.X += slx
		p.Y += sly
		if !g.IsValid(p) {
			return
		}

		podes[p] = true
	}
}

func solveB(g *grid.Grid[rune]) int64 {
	allAnts := map[rune][]pos.P2{}
	g.Walk(func(p pos.P2, r rune) {
		if r != '.' {
			allAnts[r] = append(allAnts[r], p)
		}
	})

	allPodes := map[pos.P2]bool{}
	for _, ants := range allAnts {
		if len(ants) == 1 {
			continue
		}

		for _, ant := range ants {
			allPodes[ant] = true
		}

		for i := 0; i < len(ants)-1; i++ {
			for j := i + 1; j < len(ants); j++ {
				a, b := ants[i], ants[j]

				slx, sly := simplify(b.X-a.X, b.Y-a.Y)

				extend(g, b, slx, sly, allPodes)
				extend(g, a, -slx, -sly, allPodes)
			}
		}
	}

	return int64(len(allPodes))
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
