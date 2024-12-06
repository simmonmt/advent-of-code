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

type Result int

const (
	TERMINATED Result = iota
	LOOP
)

func runMaze(input *Input, obs pos.P2, keepVisited bool) (map[pos.P2]bool, Result) {
	// The run time of this function is dominated by access time to
	// 'states'. If it's a map, depending on how complicated the key is, the
	// runtime (A+B) can be close to 15s. Simpler keys and preallocated maps
	// reduce the time somewhat, but hashing is still expensive. Nothing
	// beats a preallocated array, which takes it down to 480ms (A+B).
	w, h := input.G.Width(), input.G.Height()
	states := make([]bool, w*h*5)
	result := TERMINATED
	visited := map[pos.P2]bool{}

	for p, d := input.StartPos, input.StartDir; input.G.IsValid(p); {
		i := int(d)*w*h + p.X*h + p.Y
		if states[i] {
			result = LOOP
			break
		}
		states[i] = true

		if keepVisited {
			visited[p] = true
		}

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

	return visited, result
}

func solveA(input *Input) int64 {
	visited, result := runMaze(input, pos.P2{X: -99, Y: -99}, true)
	if result != TERMINATED {
		panic("not terminated")
	}

	return int64(len(visited))
}

func solveB(input *Input) int64 {
	visited, result := runMaze(input, pos.P2{X: -99, Y: -99}, true)
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

		if _, result := runMaze(input, p, false); result == LOOP {
			num++
		}
	}

	return num
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
