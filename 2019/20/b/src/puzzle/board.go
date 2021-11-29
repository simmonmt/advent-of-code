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
	"sort"

	"github.com/simmonmt/aoc/2019/common/dir"
	"github.com/simmonmt/aoc/2019/common/pos"
)

type Tile int

const (
	TILE_OPEN Tile = iota
	TILE_PATH
	TILE_WALL
	TILE_GATE
)

func (t Tile) String() string {
	switch t {
	case TILE_OPEN:
		return "_"
	case TILE_PATH:
		return "."
	case TILE_WALL:
		return "#"
	case TILE_GATE:
		return "G"
	default:
		panic("bad tile type")
	}
}

type Gate struct {
	name      string
	pOut, pIn pos.P2 // entrance locations
	gOut, gIn pos.P2 // exit locations
}

func (g Gate) Name() string      { return g.name }
func (g Gate) PortalOut() pos.P2 { return g.pOut }
func (g Gate) PortalIn() pos.P2  { return g.pIn }
func (g Gate) GateOut() pos.P2   { return g.gOut }
func (g Gate) GateIn() pos.P2    { return g.gIn }

func (g Gate) String() string {
	return fmt.Sprintf("[%s: out p:%v,g:%v in p:%v g:%v]",
		g.name, g.pOut, g.gOut, g.pIn, g.gIn)
}

type ByGateName []Gate

func (a ByGateName) Len() int      { return len(a) }
func (a ByGateName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByGateName) Less(i, j int) bool {
	return a[i].name < a[j].name
}

type Board struct {
	w, h int
	c    []Tile

	gatesByName      map[string]*Gate
	gatesByPortalLoc map[pos.P2]*Gate
	gatesByGateLoc   map[pos.P2]*Gate
}

func NewBoard(lines []string) *Board {
	b := &Board{
		w:                len(lines[0]),
		h:                len(lines),
		gatesByName:      map[string]*Gate{},
		gatesByPortalLoc: map[pos.P2]*Gate{},
		gatesByGateLoc:   map[pos.P2]*Gate{},
	}

	b.c = make([]Tile, b.w*b.h)

	gateChars := map[pos.P2]string{}
	for y, line := range lines {
		for x, r := range line {
			b.set(pos.P2{x, y}, r, gateChars)
		}
	}

	b.initGates(gateChars)

	return b
}

func (b *Board) initGates(gateChars map[pos.P2]string) {
	addGate := func(name string, pPos, gPos pos.P2) {
		outside := (gPos.X < 2 || gPos.X >= b.w-2 || gPos.Y < 2 || gPos.Y >= b.h-2)

		gate, found := b.gatesByName[name]
		if !found {
			gate = &Gate{
				name: name,
				pOut: pos.P2{-1, -1},
				gOut: pos.P2{-1, -1},
				pIn:  pos.P2{-1, -1},
				gIn:  pos.P2{-1, -1},
			}
			b.gatesByName[name] = gate
		}

		if outside {
			gate.pOut, gate.gOut = pPos, gPos
		} else {
			gate.pIn, gate.gIn = pPos, gPos
		}

		b.gatesByPortalLoc[pPos] = gate
		b.gatesByGateLoc[gPos] = gate
	}

	for p, ch1 := range gateChars {
		for _, d := range dir.AllDirs {
			gatePos := d.From(p)
			if ch2, found := gateChars[gatePos]; found {
				portalPos := d.From(gatePos)
				if b.Get(portalPos) == TILE_PATH {
					if d == dir.DIR_SOUTH || d == dir.DIR_EAST {
						addGate(ch1+ch2, portalPos, gatePos)
					} else {
						addGate(ch2+ch1, portalPos, gatePos)
					}
				}
			}
		}
	}
}

func (b *Board) set(p pos.P2, r rune, gateChars map[pos.P2]string) {
	i := p.Y*b.w + p.X

	switch {
	case r == '#':
		b.c[i] = TILE_WALL
		break
	case r == '.':
		b.c[i] = TILE_PATH
		break
	case r == ' ':
		b.c[i] = TILE_OPEN
		break
	case r >= 'A' && r <= 'Z':
		b.c[i] = TILE_GATE
		gateChars[p] = string(r)
		break
	default:
		panic(fmt.Sprintf("unexpected char %d %c", r, r))
	}
}

func (b *Board) Get(p pos.P2) Tile {
	if p.X < 0 || p.Y < 0 {
		return TILE_OPEN
	}

	i := p.Y*b.w + p.X
	if i >= len(b.c) {
		return TILE_OPEN
	}

	return b.c[i]
}

func (b *Board) NumGates() int {
	return len(b.gatesByName)
}

func (b *Board) Gates() []Gate {
	gates := make([]Gate, len(b.gatesByName))

	i := 0
	for _, v := range b.gatesByName {
		gates[i] = *v
		i++
	}

	sort.Sort(ByGateName(gates))
	return gates
}

func (b *Board) Gate(name string) Gate {
	return *b.gatesByName[name]
}

func (b *Board) IsPortal(p pos.P2) bool {
	_, found := b.gatesByPortalLoc[p]
	return found
}

func (b *Board) GateByGateLoc(p pos.P2) Gate {
	gate, found := b.gatesByGateLoc[p]
	if !found {
		panic(fmt.Sprintf("no gate at %v", p))
	}
	return *gate
}
