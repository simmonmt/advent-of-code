package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
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
	dimsFlag   = flag.String("dims", "101,103", "board dimensions")

	robotPattern = regexp.MustCompile(`^p=([0-9]+,[0-9]+) v=([-0-9]+,[-0-9]+)$`)
)

type Robot struct {
	P, V pos.P2
}

func parseInput(lines []string) ([]*Robot, error) {
	out := []*Robot{}

	for i, line := range lines {
		parts := robotPattern.FindStringSubmatch(line)
		if len(parts) != 3 {
			return nil, fmt.Errorf("%d: bad match", i+1)
		}

		p, err := pos.P2FromString(parts[1])
		if err != nil {
			return nil, fmt.Errorf("%d: bad pos: %v", i+1, err)
		}

		v, err := pos.P2FromString(parts[2])
		if err != nil {
			return nil, fmt.Errorf("%d: bad velocity: %v", i+1, err)
		}

		out = append(out, &Robot{P: p, V: v})
	}

	return out, nil
}

func moveRobot(r *Robot, dims pos.P2) {
	r.P.X = (r.P.X + dims.X + r.V.X) % dims.X
	r.P.Y = (r.P.Y + dims.Y + r.V.Y) % dims.Y
}

func solveA(input []*Robot, dims pos.P2) int64 {
	robots := make([]Robot, len(input))
	for i, in := range input {
		robots[i] = *in
	}

	for i := 0; i < 100; i++ {
		for j := range robots {
			moveRobot(&robots[j], dims)
		}
	}

	whichQuad := func(i, sz int) int {
		if i < sz/2 {
			return 0
		}
		if i > sz/2 {
			return 1
		}
		return -1
	}

	quads := [2][2]int{}
	for x := 0; x < dims.X; x++ {
		xq := whichQuad(x, dims.X)
		for y := 0; y < dims.Y; y++ {
			yq := whichQuad(y, dims.Y)
			p := pos.P2{X: x, Y: y}

			if xq == -1 || yq == -1 {
				continue
			}

			for i := range robots {
				if robots[i].P.Equals(p) {
					quads[yq][xq]++
				}
			}
		}
	}

	score := quads[0][0] * quads[0][1] * quads[1][0] * quads[1][1]
	return int64(score)
}

func dumpRobots(robots []Robot, dims pos.P2) {
	g := grid.New[bool](dims.X, dims.Y)
	for i := range robots {
		g.Set(robots[i].P, true)
	}

	g.Dump(false, func(_ pos.P2, v bool, _ bool) string {
		if v {
			return "#"
		}
		return "."
	})
	fmt.Println()
}

// 7051
func solveB(input []*Robot, dims pos.P2) int64 {
	robots := make([]Robot, len(input))
	for i, in := range input {
		robots[i] = *in
	}

	for i := 0; i < 10000; i++ {
		for j := range robots {
			moveRobot(&robots[j], dims)
		}

		//if (i-46)%103 == 0 {
		fmt.Println(i)
		dumpRobots(robots, dims)
		//}
	}

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

	dims, err := pos.P2FromString(*dimsFlag)
	if err != nil {
		logger.Fatalf("failed to process dims: %v", err)
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

	fmt.Println("A", solveA(input, dims))
	fmt.Println("B", solveB(input, dims))
}
