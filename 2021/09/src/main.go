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
	"container/list"
	"flag"
	"fmt"
	"log"
	"sort"

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/grid"
	"github.com/simmonmt/aoc/2021/common/logger"
	"github.com/simmonmt/aoc/2021/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func solveA(g *grid.Grid) {
	risk := 0

	g.Walk(func(ctr pos.P2, v interface{}) {
		d := v.(int)
		for _, n := range g.AllNeighbors(ctr, false) {
			if g.Get(n).(int) <= d {
				return
			}

		}
		logger.LogF("%v=%d lowest", ctr, d)
		risk += d + 1
	})

	fmt.Println("A", risk)
}

func bfs(start pos.P2, getNeighbors func(p pos.P2) []pos.P2) {
	seen := map[pos.P2]bool{}

	queue := list.New()
	queue.PushBack(start)

	for queue.Front() != nil {
		p := queue.Remove(queue.Front()).(pos.P2)
		if _, found := seen[p]; found {
			continue
		}

		seen[p] = true

		for _, n := range getNeighbors(p) {
			if _, found := seen[n]; found {
				continue
			}
			queue.PushBack(n)
		}
	}
}

func basinDump(g *grid.Grid, basinNums map[pos.P2]int) {
	for y := 0; y < g.Height(); y++ {
		fmt.Print("|")
		for x := 0; x < g.Width(); x++ {
			p := pos.P2{X: x, Y: y}
			num, found := basinNums[p]
			if found {
				fmt.Printf("%d", num)
			} else if g.Get(p).(int) == 9 {
				fmt.Print("^")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println("|")
	}
}

func solveB(g *grid.Grid) {
	thisBasinNum := 0
	basinNums := map[pos.P2]int{}
	basinSizes := map[int]int{}

	g.Walk(func(ctr pos.P2, d interface{}) {
		if d == 9 {
			return
		}

		if _, found := basinNums[ctr]; found {
			return
		}

		thisBasinNum++
		bfs(ctr, func(p pos.P2) []pos.P2 {
			if g.Get(p) == 9 {
				return nil
			}

			if _, found := basinNums[p]; found {
				basinDump(g, basinNums)
				panic(fmt.Sprintf(
					"basin collision; want to put basin num %v at %v; has %v",
					thisBasinNum, p, basinNums[p]))
			}

			basinNums[p] = thisBasinNum
			basinSizes[thisBasinNum]++
			return g.AllNeighbors(p, false)
		})
	})

	if logger.Enabled() {
		basinDump(g, basinNums)
	}

	basinsBySize := []int{}
	for num := range basinSizes {
		basinsBySize = append(basinsBySize, num)
	}
	sort.Slice(basinsBySize, func(i, j int) bool {
		return basinSizes[basinsBySize[i]] <
			basinSizes[basinsBySize[j]]
	})

	topThree := basinsBySize[len(basinsBySize)-3 : len(basinsBySize)]
	out := 1
	for _, num := range topThree {
		out *= basinSizes[num]
	}

	fmt.Println("B", out)
}

func readInput(path string) ([]string, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, err
	}

	return lines, err
}

func newGrid(lines []string) *grid.Grid {
	w, h := len(lines[0]), len(lines)
	g := grid.New(w, h)

	for y, line := range lines {
		for x, r := range line {
			d := int(r - '0')
			g.Set(pos.P2{X: x, Y: y}, d)
		}
	}

	return g
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

	g := newGrid(lines)
	solveA(g)
	solveB(g)
}
