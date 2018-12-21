package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
)

func readInput() ([]string, error) {
	lines := []string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lines, nil
}

type Extent struct {
	S, E int
}

func traverse(pos Pos, board *Board, path string) {
	logger.LogF("traverse pos %v path %v\n", pos, path)
	i := 0
	for i < len(path) {
		nextI := i + 1

		switch path[i] {
		case 'N', 'S', 'E', 'W':
			pos = board.Move(pos, string(path[i]))
		case '(':
			extents, endIdx := parseGroup(path[i:])
			for _, e := range extents {
				// traverse the extent and the stuff beyond it
				traverse(pos, board, path[e.S+i:e.E+i]+path[endIdx+i:])
				// fmt.Println("---")
			}
			return
		case '|', ')':
			break
		default:
			panic("unknown " + string(path[i]))
		}

		i = nextI
	}
}

func parseGroup(str string) ([]Extent, int) {
	level := 0

	extents := []Extent{}
	curExtent := Extent{}
	for i, c := range str {
		switch c {
		case '(':
			level++
			if level == 1 {
				curExtent.S = i + 1
				curExtent.E = -1
			}
		case '|':
			if level == 1 {
				curExtent.E = i
				extents = append(extents, curExtent)

				curExtent.S = i + 1
				curExtent.E = -1
			}
		case ')':
			if level == 1 {
				curExtent.E = i
				extents = append(extents, curExtent)
			}

			level--
			if level == 0 {
				return extents, i
			}
		}
	}

	panic("ran out of string")
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	lines, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	board := NewBoard()

	for _, line := range lines {
		line = strings.TrimPrefix(line, "^")
		line = strings.TrimSuffix(line, "$")
		traverse(Pos{0, 0}, board, line)
	}

	board.Dump(Pos{0, 0})
}
