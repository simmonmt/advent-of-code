package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2019/common/intmath"
	"github.com/simmonmt/aoc/2019/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func readInput(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := []string{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lines, nil
}

type Pos struct {
	X, Y int
}

func draw(grid, ref map[Pos]int, cmds []string, intersect func(Pos, int, int)) {
	pos := Pos{0, 0}
	numSteps := 0
	for _, cmd := range cmds {
		dir := cmd[0]
		dist, err := strconv.Atoi(cmd[1:])
		if err != nil {
			panic(fmt.Sprintf("bad cmd %v", cmd))
		}

		dest := pos
		var inc Pos
		switch dir {
		case 'U':
			dest.Y -= dist
			inc = Pos{0, -1}
			break
		case 'D':
			dest.Y += dist
			inc = Pos{0, 1}
			break
		case 'L':
			dest.X -= dist
			inc = Pos{-1, 0}
			break
		case 'R':
			dest.X += dist
			inc = Pos{1, 0}
			break
		default:
			panic(fmt.Sprintf("bad dir in cmd %v", cmd))
		}

		for pos != dest {
			numSteps++
			pos.X += inc.X
			pos.Y += inc.Y
			grid[pos] = numSteps
			if ref[pos] != 0 {
				intersect(pos, numSteps, ref[pos])
			}
		}

		logger.LogF("%c %v", dir, dist)
	}
}

func dump(grid map[Pos]int) {
	var minX, maxX, minY, maxY int
	first := true
	for pos := range grid {
		if first {
			minX, maxX = pos.X, pos.X
			minY, maxY = pos.Y, pos.Y
			first = false
		} else {
			minX = intmath.IntMin(minX, pos.X)
			maxX = intmath.IntMax(maxX, pos.X)
			minY = intmath.IntMin(minY, pos.Y)
			maxY = intmath.IntMax(maxY, pos.Y)
		}
	}

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if grid[Pos{x, y}] != 0 {
				fmt.Print("+")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func solve(first, second []string) {
	grid1 := map[Pos]int{}
	grid2 := map[Pos]int{}

	draw(grid1, map[Pos]int{}, first, func(Pos, int, int) {})
	if logger.Enabled() {
		dump(grid1)
	}

	closest := map[Pos]int{}
	draw(grid2, grid1, second, func(pos Pos, aSteps, bSteps int) {
		fmt.Printf("intersect at %v a %v b %v\n", pos, aSteps, bSteps)
		if _, found := closest[pos]; !found {
			closest[pos] = aSteps + bSteps
		}
	})
	if logger.Enabled() {
		dump(grid2)
	}

	lowest := -1
	for _, dist := range closest {
		if lowest == -1 || dist < lowest {
			lowest = dist
		}
	}

	fmt.Printf("closest dist %d\n", lowest)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	solve(strings.Split(lines[0], ","), strings.Split(lines[1], ","))
}
