// Copyright 2023 Google LLC
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
	"strconv"

	"github.com/simmonmt/aoc/2023/common/dir"
	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/grid"
	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func parseInput(lines []string) (*grid.Grid[rune], error) {
	return grid.NewFromLines[rune](lines, func(p pos.P2, r rune) (rune, error) {
		if r != '.' && r != '|' && r != '-' && r != '\\' && r != '/' {
			return '?', fmt.Errorf("bad rune %c at %v", r, p)
		}

		return r, nil
	})
}

type Beam struct {
	Loc pos.P2
	Dir dir.Dir
}

type Board struct {
	g        *grid.Grid[rune]
	heads    []*Beam
	presence map[pos.P2][]bool
}

func NewBoard(g *grid.Grid[rune], start *Beam) *Board {
	b := &Board{
		g:        g,
		heads:    []*Beam{start},
		presence: map[pos.P2][]bool{},
	}

	b.setPresence(start)

	return b
}

func (b *Board) hasPresence(beam *Beam) bool {
	dirs, found := b.presence[beam.Loc]
	return found && dirs[int(beam.Dir)]
}

func (b *Board) setPresence(beam *Beam) {
	if _, found := b.presence[beam.Loc]; !found {
		b.presence[beam.Loc] = make([]bool, 5)
	}
	b.presence[beam.Loc][int(beam.Dir)] = true
}

var (
	forwardSlashDirs = map[dir.Dir]dir.Dir{
		dir.DIR_NORTH: dir.DIR_EAST,
		dir.DIR_SOUTH: dir.DIR_WEST,
		dir.DIR_EAST:  dir.DIR_NORTH,
		dir.DIR_WEST:  dir.DIR_SOUTH,
	}

	backSlashDirs = map[dir.Dir]dir.Dir{
		dir.DIR_NORTH: dir.DIR_WEST,
		dir.DIR_SOUTH: dir.DIR_EAST,
		dir.DIR_EAST:  dir.DIR_SOUTH,
		dir.DIR_WEST:  dir.DIR_NORTH,
	}
)

func (b *Board) advanceHead(head *Beam) []*Beam {
	advanceBeam := func(in *Beam) *Beam {
		return &Beam{
			Loc: in.Dir.From(in.Loc),
			Dir: in.Dir,
		}
	}

	switch r, _ := b.g.Get(head.Loc); r {
	case '.':
		return []*Beam{advanceBeam(head)} // keep going
	case '-':
		if head.Dir == dir.DIR_EAST || head.Dir == dir.DIR_WEST {
			return []*Beam{advanceBeam(head)}
		}

		return []*Beam{
			&Beam{Loc: dir.DIR_EAST.From(head.Loc), Dir: dir.DIR_EAST},
			&Beam{Loc: dir.DIR_WEST.From(head.Loc), Dir: dir.DIR_WEST},
		}
	case '|':
		if head.Dir == dir.DIR_NORTH || head.Dir == dir.DIR_SOUTH {
			return []*Beam{advanceBeam(head)}
		}

		return []*Beam{
			&Beam{Loc: dir.DIR_NORTH.From(head.Loc), Dir: dir.DIR_NORTH},
			&Beam{Loc: dir.DIR_SOUTH.From(head.Loc), Dir: dir.DIR_SOUTH},
		}
	case '/':
		newDir := forwardSlashDirs[head.Dir]
		newLoc := newDir.From(head.Loc)

		return []*Beam{&Beam{newLoc, newDir}}
	case '\\':
		newDir := backSlashDirs[head.Dir]
		newLoc := newDir.From(head.Loc)

		return []*Beam{&Beam{newLoc, newDir}}
	default:
		panic("bad rune")
	}
}

func (b *Board) Next() bool {
	if len(b.heads) == 0 {
		return false
	}

	newHeads := []*Beam{}
	for _, head := range b.heads {
		for _, nh := range b.advanceHead(head) {
			if !b.g.IsValid(nh.Loc) {
				continue
			}

			if b.hasPresence(nh) {
				continue
			}
			b.setPresence(nh)

			newHeads = append(newHeads, nh)
		}
	}
	b.heads = newHeads

	return true
}

func (b *Board) Walk(cb func(p pos.P2, r rune, dirs []dir.Dir)) {
	b.g.Walk(func(p pos.P2, r rune) {
		dirs := []dir.Dir{}
		for d, found := range b.presence[p] {
			if found {
				dirs = append(dirs, dir.Dir(d))
			}
		}
		cb(p, r, dirs)
	})
}

func (b *Board) Dump() {
	const presenceToString = " ^V<>"

	b.g.Dump(true, func(p pos.P2, r rune, _ bool) string {
		if r != '.' {
			return string(r)
		}

		char := '.'
		num := 0
		for d, found := range b.presence[p] {
			if !found {
				continue
			}

			num++
			char = rune(presenceToString[d])
		}

		if num <= 1 {
			return string(char)
		}
		return strconv.Itoa(num)
	})
}

func solve(g *grid.Grid[rune], start *Beam) int {
	b := NewBoard(g, start)
	for b.Next() {
	}

	num := 0
	b.Walk(func(p pos.P2, r rune, dirs []dir.Dir) {
		if len(dirs) > 0 {
			num++
		}
	})
	return num
}

func solveA(g *grid.Grid[rune]) int {
	return solve(g, &Beam{Loc: pos.P2{X: 0, Y: 0}, Dir: dir.DIR_EAST})
}

func solveB(input *grid.Grid[rune]) int {
	starts := []*Beam{}

	for x := 0; x < input.Width(); x++ {
		starts = append(starts, &Beam{Loc: pos.P2{X: x, Y: 0}, Dir: dir.DIR_SOUTH})
		starts = append(starts, &Beam{Loc: pos.P2{X: x, Y: input.Height() - 1}, Dir: dir.DIR_NORTH})
	}

	for y := 0; y < input.Height(); y++ {
		starts = append(starts, &Beam{Loc: pos.P2{X: 0, Y: y}, Dir: dir.DIR_EAST})
		starts = append(starts, &Beam{Loc: pos.P2{X: input.Width() - 1, Y: y}, Dir: dir.DIR_WEST})
	}

	high := 0
	for _, start := range starts {
		high = max(high, solve(input, start))
	}
	return high
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		logger.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		logger.Fatalf("failed to read input: %v", err)
	}

	input, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
