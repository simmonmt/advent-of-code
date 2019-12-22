package puzzle

import (
	"fmt"

	"github.com/simmonmt/aoc/2019/common/pos"
)

type Tile int

const (
	TILE_OPEN Tile = iota
	TILE_WALL
	TILE_KEY
	TILE_DOOR
)

type Board struct {
	w, h int
	c    []Tile

	keys       map[string]pos.P2
	keysByLoc  map[pos.P2]string
	doors      map[string]pos.P2
	doorsByLoc map[pos.P2]string
}

func NewBoard(lines []string) (*Board, pos.P2) {
	b := &Board{
		w:          len(lines[0]),
		h:          len(lines),
		keys:       map[string]pos.P2{},
		keysByLoc:  map[pos.P2]string{},
		doors:      map[string]pos.P2{},
		doorsByLoc: map[pos.P2]string{},
	}

	b.c = make([]Tile, b.w*b.h)

	var start pos.P2
	for y, line := range lines {
		for x, r := range line {
			p := pos.P2{x, y}
			if r == '@' {
				start = p
				b.set(p, '.')
			} else {
				b.set(p, r)
			}
		}
	}

	return b, start
}

func (b *Board) set(p pos.P2, r rune) {
	i := p.Y*b.w + p.X

	switch {
	case r == '#':
		b.c[i] = TILE_WALL
		break
	case r == '.':
		b.c[i] = TILE_OPEN
		break
	case r >= 'a' && r <= 'z':
		b.c[i] = TILE_KEY
		b.keys[string(r)] = p
		b.keysByLoc[p] = string(r)
		break
	case r >= 'A' && r <= 'Z':
		b.c[i] = TILE_DOOR
		b.doors[string(r)] = p
		b.doorsByLoc[p] = string(r)
		break
	default:
		panic(fmt.Sprintf("unexpected char %d %c", r, r))
	}
}

func (b *Board) Get(p pos.P2) Tile {
	i := p.Y*b.w + p.X
	if i < 0 || i >= len(b.c) {
		panic("get out of range")
	}
	return b.c[i]
}

func (b *Board) Keys() []string {
	out := []string{}
	for key := range b.keys {
		out = append(out, key)
	}
	return out
}

func (b *Board) KeyLoc(key string) pos.P2 {
	if p, ok := b.keys[key]; ok {
		return p
	}
	panic(fmt.Sprintf("bad key %s", key))
}

func (b *Board) KeyAtLoc(p pos.P2) string {
	if k, ok := b.keysByLoc[p]; ok {
		return k
	}
	panic(fmt.Sprintf("no key at %v", p))
}

func (b *Board) DoorLoc(door string) pos.P2 {
	if p, ok := b.doors[door]; ok {
		return p
	}
	panic(fmt.Sprintf("bad dor %s", door))
}

func (b *Board) DoorAtLoc(p pos.P2) string {
	if d, ok := b.doorsByLoc[p]; ok {
		return d
	}
	panic(fmt.Sprintf("no door at %v", p))
}
