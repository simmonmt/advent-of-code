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
	"log"

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

func solveA(b *Board) {
	risk := 0

	b.Walk(func(ctr pos.P2, d int) {
		for _, n := range ctr.AllNeighbors(false) {
			if n.X < 0 || n.Y < 0 || n.X >= b.Width() || n.Y >= b.Height() {
				continue
			}

			if b.Get(n) <= d {
				return
			}

		}
		logger.LogF("%v=%d lowest", ctr, d)
		risk += d + 1
	})

	fmt.Println("A", risk)
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
}
