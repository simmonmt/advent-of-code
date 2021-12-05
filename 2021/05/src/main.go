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
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/logger"
	"github.com/simmonmt/aoc/2021/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	inputPattern = regexp.MustCompile(`^(\d+,\d+) -> (\d+,\d+)$`)
)

type Path struct {
	From, To pos.P2
}

func readInput(path string) ([]Path, error) {
	lines, err := filereader.Lines(path)
	if err != nil {
		return nil, err
	}

	out := []Path{}
	for i, line := range lines {
		parts := inputPattern.FindStringSubmatch(line)
		if parts == nil {
			return nil, fmt.Errorf("%d: parse failure", i)
		}

		from, err := pos.P2FromString(parts[1])
		if err != nil {
			return nil, fmt.Errorf("%d: bad from: %v", i, err)
		}

		to, err := pos.P2FromString(parts[2])
		if err != nil {
			return nil, fmt.Errorf("%d: bad to: %v", i, err)
		}

		out = append(out, Path{From: from, To: to})
	}

	return out, nil

}

func hvOnly(in []Path) []Path {
	out := []Path{}
	for _, path := range in {
		if path.From.X == path.To.X || path.From.Y == path.To.Y {
			out = append(out, path)
		}
	}
	return out
}

func fillGrid(paths []Path) map[pos.P2]int {
	grid := map[pos.P2]int{}

	for _, path := range paths {
		from, to := path.From, path.To

		inc := pos.P2{0, 0}

		if from.X < to.X {
			inc.X = 1
		} else if from.X > to.X {
			inc.X = -1
		}

		if from.Y < to.Y {
			inc.Y = 1
		} else if from.Y > to.Y {
			inc.Y = -1
		}

		if inc.X == 0 && inc.Y == 0 {
			panic(fmt.Sprintf(
				"bad inc: from %+v to %+v inc %+v\n",
				from, to, inc))
		}

		p := from
		for {
			grid[p]++

			if p.Equals(to) {
				break
			}

			p.Add(inc)
		}
	}

	return grid
}

func printGridTo(w io.Writer, grid map[pos.P2]int) {
	maxX, maxY := 0, 0
	for p := range grid {
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {
			p := pos.P2{X: x, Y: y}
			n := grid[p]

			if n == 0 {
				fmt.Fprint(w, ".")
			} else if n < 10 {
				fmt.Fprintf(w, "%d", n)
			} else {
				fmt.Fprint(w, "X")
			}
		}
		fmt.Fprint(w, "\n")
	}
}

func printGrid(grid map[pos.P2]int) {
	printGridTo(os.Stdout, grid)
}

func solve(paths []Path) int {
	grid := fillGrid(paths)
	if logger.Enabled() {
		printGrid(grid)
	}

	count := 0
	for _, num := range grid {
		if num > 1 {
			count++
		}
	}
	return count
}

func solveA(paths []Path) {
	paths = hvOnly(paths)
	logger.LogLn("#hvpaths:", len(paths))
	logger.LogLn(paths)

	fmt.Println("A", solve(paths))
}

func solveB(paths []Path) {
	fmt.Println("B", solve(paths))
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	paths, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}
	logger.LogLn("#paths:", len(paths))
	logger.LogLn(paths)

	solveA(paths)
	solveB(paths)
}
