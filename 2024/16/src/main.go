package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/simmonmt/aoc/2024/common/dir"
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
	Start, End pos.P2
	Grid       *grid.Grid[rune]
}

func parseInput(lines []string) (*Input, error) {
	var start, end pos.P2

	g, _ := grid.NewFromLines[rune](lines, func(p pos.P2, r rune) (rune, error) {
		if r == 'S' {
			start = p
			r = '.'
		} else if r == 'E' {
			end = p
			r = '.'
		}

		return r, nil
	})

	return &Input{Start: start, End: end, Grid: g}, nil
}

func get[T any](g *grid.Grid[T], p pos.P2) T {
	v, _ := g.Get(p)
	return v
}

type State struct {
	D dir.Dir
	P pos.P2
}

func (s *State) Equals(o State) bool {
	return s.D == o.D && s.P.Equals(o.P)
}

type Path struct {
	Cost int
	Prev []State
}

func findBestPaths(g *grid.Grid[rune], start State, end pos.P2) (best map[State]*Path, bestEnd State) {
	best = map[State]*Path{}

	todo := map[State]int{start: 0}

	for len(todo) > 0 {
		next := map[State]int{}

		for cur, curCost := range todo {
			for _, d := range []dir.Dir{cur.D, cur.D.Left(), cur.D.Right()} {
				n := d.From(cur.P)
				if get(g, n) != '.' {
					continue
				}

				ns := State{P: n, D: d}

				cost := curCost + 1
				if d != cur.D {
					cost += 1000
				}

				curPath := best[ns]
				if curPath == nil || curPath.Cost > cost {
					curPath = &Path{Cost: cost, Prev: []State{cur}}
					best[ns] = curPath
				} else if curPath.Cost == cost {
					curPath.Prev = append(curPath.Prev, cur)
				} else { // curPath.Cost < cost
					continue
				}

				next[ns] = cost
			}
		}

		todo = next
	}

	bestCost := -1
	for s, p := range best {
		if s.P.Equals(end) && (bestCost == -1 || p.Cost < bestCost) {
			bestEnd = s
			bestCost = p.Cost
		}
	}

	return best, bestEnd
}

func solveA(input *Input) int {
	start := State{D: dir.DIR_EAST, P: input.Start}
	best, bestEnd := findBestPaths(input.Grid, start, input.End)
	return best[bestEnd].Cost
}

func solveB(input *Input) int {
	start := State{P: input.Start, D: dir.DIR_EAST}
	best, bestEnd := findBestPaths(input.Grid, start, input.End)

	todo := map[State]int{bestEnd: 1}
	touched := map[pos.P2]bool{}

	for len(todo) > 0 {
		next := map[State]int{}

		for cur := range todo {
			touched[cur.P] = true

			if cur.Equals(start) {
				continue
			}
			if _, found := best[cur]; !found {
				panic(fmt.Sprintf("no best for %v", cur))
			}
			for _, ps := range best[cur].Prev {
				next[ps] = 1
			}
		}

		todo = next
	}

	// input.Grid.Dump(true, func(p pos.P2, r rune, _ bool) string {
	// 	if _, found := touched[p]; found {
	// 		return "O"
	// 	} else {
	// 		return string(r)
	// 	}
	// })

	return len(touched)
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
