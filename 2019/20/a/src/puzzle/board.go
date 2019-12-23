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
	name   string
	p1, p2 pos.P2 // entrance locations
	g1, g2 pos.P2 // exit locations
}

func (g Gate) Name() string    { return g.name }
func (g Gate) Portal1() pos.P2 { return g.p1 }
func (g Gate) Portal2() pos.P2 { return g.p2 }
func (g Gate) Gate1() pos.P2   { return g.g1 }
func (g Gate) Gate2() pos.P2   { return g.g2 }

func (g Gate) String() string {
	return fmt.Sprintf("[%s: p1:%v g1:%v p2:%v g2:%v]",
		g.name, g.p1, g.g1, g.p2, g.g2)
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
		if gate, found := b.gatesByName[name]; found {
			if pPos.LessThan(gate.p1) {
				gate.p1, gate.p2 = pPos, gate.p1
				gate.g1, gate.g2 = gPos, gate.g1
			} else {
				gate.p2 = pPos
				gate.g2 = gPos
			}
			b.gatesByPortalLoc[pPos] = gate
			b.gatesByGateLoc[gPos] = gate
		} else {
			gate = &Gate{name, pPos, pos.P2{-1, -1}, gPos, pos.P2{-1, -1}}
			b.gatesByName[name] = gate
			b.gatesByPortalLoc[pPos] = gate
			b.gatesByGateLoc[gPos] = gate
		}
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
	return *b.gatesByGateLoc[p]
}
