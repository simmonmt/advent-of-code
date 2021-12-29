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

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/grid"
	"github.com/simmonmt/aoc/2021/common/logger"
	"github.com/simmonmt/aoc/2021/common/pos"
)

var (
	verbose  = flag.Bool("verbose", false, "verbose")
	input    = flag.String("input", "", "input file")
	numSteps = flag.Int("num_steps", 10, "number of steps")
)

type Cuke int

const (
	EMPTY Cuke = iota
	DOWN
	RIGHT
)

func (c Cuke) String() string {
	switch c {
	case EMPTY:
		return "."
	case DOWN:
		return "v"
	case RIGHT:
		return ">"
	default:
		panic("bad cuke")
	}
}

type Board struct {
	g *grid.Grid
}

func NewBoard(lines []string) (*Board, error) {
	g := grid.New(len(lines[0]), len(lines))

	for y := range lines {
		for x := range lines[0] {
			p := pos.P2{x, y}

			switch c := lines[y][x]; c {
			case '.':
				g.Set(p, EMPTY)
			case 'v':
				g.Set(p, DOWN)
			case '>':
				g.Set(p, RIGHT)
			default:
				return nil, fmt.Errorf("unknown char %v at %v",
					string(c), p)
			}
		}
	}

	return &Board{g: g}, nil
}

func (b *Board) DumpTo(o io.Writer) {
	b.g.Walk(func(p pos.P2, v interface{}) {
		c := v.(Cuke)
		fmt.Fprint(o, c)
		if p.X == b.g.Width()-1 {
			fmt.Fprintln(o)
		}
	})
}

func (b *Board) Move() int {
	moves := map[pos.P2]pos.P2{}
	numMoves := 0

	b.g.Walk(func(p pos.P2, v interface{}) {
		if c := v.(Cuke); c == RIGHT {
			right := pos.P2{(p.X + 1) % b.g.Width(), p.Y}
			if b.get(right) == EMPTY {
				moves[p] = right
			}
		}
	})
	numMoves += len(moves)
	for from, to := range moves {
		b.set(from, EMPTY)
		b.set(to, RIGHT)
	}

	moves = map[pos.P2]pos.P2{}
	b.g.Walk(func(p pos.P2, v interface{}) {
		if c := v.(Cuke); c == DOWN {
			down := pos.P2{p.X, (p.Y + 1) % b.g.Height()}
			if b.get(down) == EMPTY {
				moves[p] = down
			}
		}
	})
	numMoves += len(moves)
	for from, to := range moves {
		b.set(from, EMPTY)
		b.set(to, DOWN)
	}

	return numMoves
}

func (b *Board) get(p pos.P2) Cuke {
	if p.X >= b.g.Width() || p.Y >= b.g.Height() {
		return EMPTY
	}
	return b.g.Get(p).(Cuke)
}

func (b *Board) set(p pos.P2, c Cuke) {
	b.g.Set(p, c)
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

	b, err := NewBoard(lines)
	if err != nil {
		log.Fatalf("failed to build board: %v", err)
	}

	if logger.Enabled() {
		fmt.Println("initial:")
		b.DumpTo(os.Stdout)
	}

	stepNum := 1
	for ; stepNum <= *numSteps; stepNum++ {
		n := b.Move()

		if logger.Enabled() {
			logger.LogF("\nafter step %v (%v moves):", stepNum, n)
			b.DumpTo(os.Stdout)
		}

		if n == 0 {
			break
		}
	}

	fmt.Printf("stopped after %v steps\n", stepNum)
}
