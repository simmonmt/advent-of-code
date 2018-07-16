package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Dir int

func (d Dir) String() string {
	switch d {
	case DirNorth:
		return "north"
	case DirSouth:
		return "south"
	case DirEast:
		return "east"
	case DirWest:
		return "west"
	default:
		panic(fmt.Sprintf("unknown dir %v", d))
	}
}

const (
	DirNorth Dir = iota
	DirSouth
	DirEast
	DirWest
)

func abs(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

type Step struct {
	left bool
	len  int
}

type Coord struct {
	x, y int
}

func walkNorth(c Coord, l int) Coord {
	c.y -= l
	return c
}

func walkSouth(c Coord, l int) Coord {
	c.y += l
	return c
}

func walkEast(c Coord, l int) Coord {
	c.x += l
	return c
}

func walkWest(c Coord, l int) Coord {
	c.x -= l
	return c
}

var (
	walkers = map[Dir]func(c Coord, l int) Coord{
		DirNorth: walkNorth,
		DirSouth: walkSouth,
		DirEast:  walkEast,
		DirWest:  walkWest,
	}

	nextDirs = map[Dir][]Dir{ // 0=right, 1=left
		DirNorth: []Dir{DirEast, DirWest},
		DirSouth: []Dir{DirWest, DirEast},
		DirEast:  []Dir{DirSouth, DirNorth},
		DirWest:  []Dir{DirNorth, DirSouth},
	}
)

func walk(c Coord, curDir Dir, step Step) (nextCoord Coord, nextDir Dir) {
	idx := 0
	if step.left {
		idx = 1
	}

	nextDir = nextDirs[curDir][idx]
	nextCoord = walkers[nextDir](c, step.len)
	return
}

func readInput(r io.Reader) (steps []Step, err error) {
	steps = []Step{}

	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		parts := strings.Split(line, ", ")
		for _, partStr := range parts {
			part := []byte(partStr)
			if len(part) < 2 {
				return nil, fmt.Errorf("illegal step %v", partStr)
			}

			left := true
			switch part[0] {
			case 'L':
				left = true
				break
			case 'R':
				left = false
				break
			default:
				return nil, fmt.Errorf("unknown step %v in %v", part[0], partStr)
			}

			stepLen, err := strconv.ParseUint(string(part[1:]), 10, 32)
			if err != nil {
				return nil, fmt.Errorf("unknown len %v in %v", string(part[1:]), partStr)
			}

			steps = append(steps, Step{left: left, len: int(stepLen)})
		}
	}

	return
}

func main() {
	steps, err := readInput(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	c := Coord{0, 0}
	dir := DirNorth
	for _, step := range steps {
		fmt.Printf("c %+v d %v, step %+v", c, dir, step)
		c, dir = walk(c, dir, step)
		fmt.Printf(", now c %+v d %v\n", c, dir)
	}

	fmt.Printf("dist: %v\n", abs(c.x)+abs(c.y))
}
