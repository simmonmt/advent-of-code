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

package puzzle

import (
	"fmt"
	"testing"

	"github.com/simmonmt/aoc/2019/common/intmath"
	"github.com/simmonmt/aoc/2019/common/pos"
)

var (
	refBoard1Str = []string{
		"      ",
		"   ## ",
		"  #.S#",
		" #o.# ",
		" #.#  ",
		"  #   ",
	}

	refBoard2Str = []string{
		" ##   ",
		"#..## ",
		"#.#..#",
		"#.S.# ",
		" ###  ",
	}
)

func runeToTile(r rune) Tile {
	switch r {
	case '#':
		return TILE_WALL
	case '.':
		return TILE_OPEN
	case 'o':
		return TILE_GOAL
	default:
		panic(fmt.Sprintf("bad rune '%v'", string(r)))
	}
}

func parseStringBoard(sb []string) (*Board, pos.P2) {
	b := NewBoard()
	var startPos pos.P2

	for y, line := range sb {
		for x, r := range line {
			p := pos.P2{x, y}

			if r == ' ' {
				continue
			} else if r == 'S' {
				startPos = p
				r = '.'
			}

			b.Set(p, runeToTile(r))
		}
	}

	return b, startPos
}

func moveTo(ref *Board, pos pos.P2, dir Dir) (newPos pos.P2, t Tile) {
	newPos = dir.From(pos)
	t = ref.Get(newPos)
	if t == TILE_UNKNOWN {
		panic(fmt.Sprintf("reached unreachable at %v", newPos))
	} else if t == TILE_WALL {
		newPos = pos
	}

	fmt.Printf("exp %v asked %s => new %v (%s)\n", pos, dir, newPos, t)
	return
}

func TestExplore(t *testing.T) {
	refBoard, refPos := parseStringBoard(refBoard1Str)
	refBoard = refBoard.CenterAt(refPos)

	expStart := pos.P2{0, 0}
	expBoard := NewBoard()
	expBoard.Set(expStart, TILE_OPEN)

	Explore(expBoard, expStart,
		func(pos pos.P2, dir Dir) (newPos pos.P2, t Tile) {
			newPos, t = moveTo(refBoard, pos, dir)
			return
		})

	min := pos.P2{
		X: intmath.IntMin(refBoard.min.X, expBoard.min.X),
		Y: intmath.IntMin(refBoard.min.Y, expBoard.min.Y),
	}
	max := pos.P2{
		X: intmath.IntMax(refBoard.max.X, expBoard.max.X),
		Y: intmath.IntMax(refBoard.max.Y, expBoard.max.Y),
	}

	diff := false
	for y := min.Y; y <= max.Y; y++ {
		for x := min.X; x <= max.X; x++ {
			p := pos.P2{x, y}
			rt := refBoard.Get(p)
			if et := expBoard.Get(p); et != rt {
				t.Errorf("diff at %v: ref %s exp %s", p, rt, et)
				diff = true
			}
		}
	}

	if diff {
		PrintBoard(refBoard, pos.P2{-1, -1})
		fmt.Println()
		PrintBoard(refBoard, pos.P2{-1, -1})
	}
}

func TestFill(t *testing.T) {
	b, refStart := parseStringBoard(refBoard2Str)
	b = b.CenterAt(refStart)

	start := pos.P2{0, 0}
	if max := Fill(b, start); max != 4 {
		t.Errorf("Fill(b, start) = %d, want 4", max)
	}
}
