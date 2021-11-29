// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"intmath"
	"log"
	"math"
	"os"
	"strings"

	"logger"
)

var (
	verbose   = flag.Bool("verbose", false, "verbose")
	dumpFinal = flag.Bool("dump_final", true, "dump final board")
	gtThresh  = flag.Int("gt_thresh", -1, "gt thresh")
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

func uniquify(posns []Pos) []Pos {
	seen := make(map[Pos]struct{}, len(posns))

	i := 0
	for _, p := range posns {
		if _, found := seen[p]; found {
			continue
		}
		posns[i] = p
		seen[p] = struct{}{}
		i++
	}

	logger.LogF("reduced %d to %d", len(posns), i)
	return posns[:i]
}

type Extent struct {
	S, E int
}

func traverse(starts []Pos, board *Board, path string) []Pos {
	logger.LogF("traverse posns %v path %v", starts, path)

	posns := make([]Pos, len(starts))
	copy(posns, starts)

	i := 0
	for i < len(path) {
		logger.LogF("i=%d of %v", i, path)
		nextI := i + 1

		switch path[i] {
		case 'N', 'S', 'E', 'W':
			dir := rune(path[i])
			for j := range posns {
				logger.LogF("moving %v from %v", string(dir), posns[j])
				posns[j] = board.Move(posns[j], dir)
				if *verbose {
					board.Dump()
				}
			}

		case '(':
			extents, endIdx := parseGroup(path[i:])
			logger.LogF("found group with %d extents end %v, starting at %v",
				len(extents), endIdx, posns)
			ends := []Pos{}
			for _, e := range extents {
				sub := path[e.S+i : e.E+i]
				newEnds := traverse(posns, board, sub)
				logger.LogF("extent %v done: began %v, returned %v", sub, posns, newEnds)
				ends = append(ends, newEnds...)
			}
			posns = uniquify(ends)
			nextI = i + endIdx

		case '|', ')':
			break
		default:
			panic("unknown " + string(path[i]))
		}

		i = nextI
	}

	logger.LogF("traverse of %v path %v done", starts, path)
	return posns
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

func findDistances(board *Board, start Pos) map[Pos]int {
	distances := map[Pos]int{start: 0}

	board.BFS(start, func(pos Pos, neighbors []Pos) {
		dist, found := distances[pos]
		if !found {
			panic("pos unfound")
		}

		for _, n := range neighbors {
			existingDist := math.MaxInt32
			if _, found := distances[n]; found {
				existingDist = distances[n]
			}
			distances[n] = intmath.IntMin(existingDist, dist+1)
		}
	})

	logger.LogLn(distances)
	return distances
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	lines, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	origin := Pos{0, 0}
	board := NewBoard(origin)

	for _, line := range lines {
		logger.LogF("input is %v", line)
		line = strings.TrimPrefix(line, "^")
		line = strings.TrimSuffix(line, "$")
		traverse([]Pos{origin}, board, line)
	}

	if *dumpFinal {
		board.Dump()
	}

	distances := findDistances(board, origin)

	greaterThan := 0
	maxDist := 0
	for _, d := range distances {
		if d > maxDist {
			maxDist = d
		}

		if d > *gtThresh {
			greaterThan++
		}
	}

	fmt.Printf("max dist %v\n", maxDist)
	if *gtThresh > 0 {
		fmt.Printf("nodes with path >%d: %d\n", *gtThresh, greaterThan)
	}
}
