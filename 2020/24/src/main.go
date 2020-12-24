package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/logger"
	"github.com/simmonmt/aoc/2020/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func move(p pos.P2, dir string) pos.P2 {
	switch dir {
	case "e":
		return pos.P2{X: p.X + 1, Y: p.Y}
	case "w":
		return pos.P2{X: p.X - 1, Y: p.Y}
	case "nw":
		return pos.P2{X: p.X, Y: p.Y - 1}
	case "ne":
		return pos.P2{X: p.X + 1, Y: p.Y - 1}
	case "sw":
		return pos.P2{X: p.X - 1, Y: p.Y + 1}
	case "se":
		return pos.P2{X: p.X, Y: p.Y + 1}
	default:
		panic("bad dir")
	}
}

func parseLine(in string) ([]string, error) {
	out := []string{}
	for len(in) > 0 {
		cmdLen := 2
		if strings.HasPrefix(in, "e") || strings.HasPrefix(in, "w") {
			cmdLen = 1
		}

		if cmdLen > len(in) {
			return nil, fmt.Errorf("bad command %v", in)
		}

		out = append(out, in[0:cmdLen])
		in = in[cmdLen:]
	}
	return out, nil
}

type Color bool

const (
	BLACK Color = true
	WHITE Color = false
)

func RunCmds(lines [][]string) map[pos.P2]Color {
	tiles := map[pos.P2]Color{}

	for _, cmds := range lines {
		p := pos.P2{X: 0, Y: 0}
		for _, cmd := range cmds {
			p = move(p, cmd)
		}

		logger.LogF("%s: %v now %v", cmds, p, !tiles[p])
		tiles[p] = !tiles[p]
	}

	return tiles
}

func CountBlack(tiles map[pos.P2]Color) int {
	numBlack := 0
	for _, c := range tiles {
		if c == BLACK {
			numBlack++
		}
	}

	return numBlack
}

func neighbors(p pos.P2) []pos.P2 {
	out := []pos.P2{}
	for _, d := range []string{"ne", "nw", "e", "w", "se", "sw"} {
		out = append(out, move(p, d))
	}
	return out
}

func numBlackNeighbors(tiles map[pos.P2]Color, p pos.P2) int {
	num := 0
	for _, n := range neighbors(p) {
		if tiles[n] == BLACK {
			num++
		}
	}
	return num
}

func evolve(old map[pos.P2]Color) map[pos.P2]Color {
	toEval := map[pos.P2]bool{}
	for p, c := range old {
		if c == BLACK {
			toEval[p] = true
		}

		for _, n := range neighbors(p) {
			toEval[n] = true
		}
	}

	out := map[pos.P2]Color{}
	for p := range toEval {
		if old[p] == BLACK {
			if num := numBlackNeighbors(old, p); num == 1 || num == 2 {
				out[p] = BLACK
			}
		} else {
			if numBlackNeighbors(old, p) == 2 {
				out[p] = BLACK
			}
		}
	}

	return out
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		log.Fatal(err)
	}

	cmds := [][]string{}
	for _, line := range lines {
		lineCmds, err := parseLine(line)
		if err != nil {
			log.Fatal(err)
		}

		cmds = append(cmds, lineCmds)
	}

	tiles := RunCmds(cmds)
	fmt.Printf("A: %d\n", CountBlack(tiles))

	for stepNum := 1; stepNum <= 100; stepNum++ {
		tiles = evolve(tiles)
		logger.LogF("Day %v: %v", stepNum, CountBlack(tiles))
	}
	fmt.Printf("B: %d\n", CountBlack(tiles))

}
