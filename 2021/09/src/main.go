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
	"github.com/simmonmt/aoc/2021/common/logger"
	"github.com/simmonmt/aoc/2021/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Board struct {
	w, h int
	m    []int
}

func NewBoard(lines []string) *Board {
	w, h := len(lines[0]), len(lines)

	b := &Board{
		w: w,
		h: h,
		m: make([]int, w*h),
	}

	for y, line := range lines {
		for x, r := range line {
			d := int(r - '0')
			b.Set(pos.P2{X: x, Y: y}, d)
		}
	}

	return b
}

func (b *Board) Width() int {
	return b.w
}

func (b *Board) Height() int {
	return b.h
}

func (b *Board) Set(p pos.P2, d int) {
	b.m[b.w*p.Y+p.X] = d
}

func (b *Board) Get(p pos.P2) int {
	return b.m[b.w*p.Y+p.X]
}

func (b *Board) Walk(walker func(p pos.P2, d int)) {
	for y := 0; y < b.h; y++ {
		for x := 0; x < b.w; x++ {
			p := pos.P2{X: x, Y: y}
			walker(p, b.Get(p))
		}
	}
}

func (b *Board) AllNeighbors(p pos.P2) []pos.P2 {
	out := []pos.P2{}
	for _, n := range p.AllNeighbors(false) {
		if n.X < 0 || n.Y < 0 {
			continue
		}
		if n.X >= b.Width() || n.Y >= b.Height() {
			continue
		}
		out = append(out, n)
	}
	return out
}

func solveA(b *Board) {
	risk := 0

	b.Walk(func(ctr pos.P2, d int) {
		for _, n := range b.AllNeighbors(ctr) {
			if b.Get(n) <= d {
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

func basinDump(b *Board, basinNums map[pos.P2]int) {
	for y := 0; y < b.Height(); y++ {
		fmt.Print("|")
		for x := 0; x < b.Width(); x++ {
			p := pos.P2{X: x, Y: y}
			num, found := basinNums[p]
			if found {
				fmt.Printf("%d", num)
			} else if b.Get(p) == 9 {
				fmt.Print("^")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println("|")
	}
}

func solveB(b *Board) {
	thisBasinNum := 0
	basinNums := map[pos.P2]int{}
	basinSizes := map[int]int{}

	b.Walk(func(ctr pos.P2, d int) {
		if d == 9 {
			return
		}

		if _, found := basinNums[ctr]; found {
			return
		}

		thisBasinNum++
		bfs(ctr, func(p pos.P2) []pos.P2 {
			if b.Get(p) == 9 {
				return nil
			}

			if _, found := basinNums[p]; found {
				basinDump(b, basinNums)
				panic(fmt.Sprintf(
					"basin collision; want to put basin num %v at %v; has %v",
					thisBasinNum, p, basinNums[p]))
			}

			basinNums[p] = thisBasinNum
			basinSizes[thisBasinNum]++
			return b.AllNeighbors(p)
		})
	})

	if logger.Enabled() {
		basinDump(b, basinNums)
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

	b := NewBoard(lines)
	solveA(b)
	solveB(b)
}
