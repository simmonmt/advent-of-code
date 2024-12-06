package main

import (
	"flag"
	"fmt"

	"github.com/simmonmt/aoc/2024/common/dir"
	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/grid"
	"github.com/simmonmt/aoc/2024/common/logger"
	"github.com/simmonmt/aoc/2024/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Input struct {
	G        *grid.Grid[bool]
	StartPos pos.P2
	StartDir dir.Dir
}

func parseInput(lines []string) (*Input, error) {
	input := &Input{}

	var err error
	input.G, err = grid.NewFromLines(lines, func(p pos.P2, r rune) (bool, error) {
		if r == '#' {
			return true, nil
		} else if r == '.' {
			return false, nil
		}

		input.StartPos = p
		var ok bool
		if input.StartDir, ok = dir.ParseIcon(r); !ok {
			return false, fmt.Errorf(`bad rune %d`, int(r))
		}
		return false, nil
	})

	return input, err
}

type State struct {
	p pos.P2
	d dir.Dir
}

type Result int

const (
	TERMINATED Result = iota
	LOOP
)

func runMaze(input *Input, obs pos.P2) (map[pos.P2]bool, Result) {
	states := map[State]bool{}
	result := TERMINATED

	for p, d := input.StartPos, input.StartDir; input.G.IsValid(p); {
		s := State{p, d}
		if _, found := states[s]; found {
			result = LOOP
			break
		}
		states[s] = true

		nd := d
		np := d.From(p)

		for {
			atObstacle := false
			if np.Equals(obs) {
				atObstacle = true
			} else if v, ok := input.G.Get(np); ok && v == true {
				atObstacle = true
			}

			if !atObstacle {
				break
			}

			nd = nd.Right()
			np = nd.From(p)

			if nd == d {
				panic("turned around")
			}
		}

		p, d = np, nd
	}

	visited := map[pos.P2]bool{}
	for s := range states {
		visited[s.p] = true
	}

	return visited, result
}

func solveA(input *Input) int64 {
	visited, result := runMaze(input, pos.P2{X: -99, Y: -99})
	if result != TERMINATED {
		panic("not terminated")
	}

	return int64(len(visited))
}

func solveB(input *Input) int64 {
	visited, result := runMaze(input, pos.P2{X: -99, Y: -99})
	if result != TERMINATED {
		panic("not terminated")
	}

	// We only need to put obstacles on positions the guard visited as those
	// are the only ones that'll cause him to change direction.
	num := int64(0)
	for p := range visited {
		if p.Equals(input.StartPos) {
			continue
		}

		if _, result := runMaze(input, p); result == LOOP {
			num++
		}
	}

	return num
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
