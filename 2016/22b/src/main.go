package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"astar"
	"intmath"
	"logger"
)

var (
	sizePattern = regexp.MustCompile(`^/dev/grid/node-x([0-9]+)-y([0-9]+) +([0-9]+)T +([0-9]+)T +[0-9]+T +[0-9]+%$`)

	immovableThresh = flag.Int("immovable_thresh", -1, "Immovable threshold")
	verbose         = flag.Bool("verbose", false, "verbose")
	start           = flag.String("start", "", "start node")
)

func readInput(r io.Reader) (*Board, *PlayState, error) {
	nodes := map[Pos]bool{}

	maxX, maxY := 0, 0
	var empty Pos

	reader := bufio.NewReader(r)
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if !strings.HasPrefix(line, "/") {
			continue
		}

		line = strings.TrimSpace(line)
		matches := sizePattern.FindStringSubmatch(line)
		if matches == nil {
			return nil, nil, fmt.Errorf("%d: failed to parse", lineNum)
		}

		x := intmath.AtoiOrDie(matches[1])
		y := intmath.AtoiOrDie(matches[2])
		//size := intmath.AtoiOrDie(matches[3])
		used := intmath.AtoiOrDie(matches[4])

		maxX = intmath.IntMax(maxX, x)
		maxY = intmath.IntMax(maxY, y)

		pos := Pos{x, y}

		switch {
		case used > *immovableThresh:
			nodes[pos] = false
		case used == 0:
			empty = pos
			nodes[pos] = true
			fmt.Printf("empty is %v\n", empty)
		default:
			nodes[pos] = true
		}
	}

	width := maxX + 1
	height := maxY + 1

	board := NewBoard(width, height)
	for p, s := range nodes {
		board.Set(p, s)
	}

	playState := &PlayState{
		Empty: empty,
		Goal:  Pos{width - 1, 0},
	}

	return board, playState, nil
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *immovableThresh == -1 {
		log.Fatal("--immovable_thresh is required")
	}

	board, playState, err := readInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	goal := &PlayState{Goal: Pos{0, 0}}

	helper := NewAStarHelper(board)
	steps := astar.AStar(playState.Encode(), goal.Encode(), helper)

	fmt.Printf("num steps %v\n", len(steps))

	fmt.Println()
	for i := len(steps) - 1; i >= 0; i-- {
		board.Dump(Decode(steps[i]))
		fmt.Println()
	}
}
