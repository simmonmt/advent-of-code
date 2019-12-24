package puzzle

import (
	"github.com/simmonmt/aoc/2019/common/dir"
	"github.com/simmonmt/aoc/2019/common/pos"
)

type Board struct {
	level    int
	c        [5][5]bool
	up, down *Board
}

func makeBoard(level int, up, down *Board) *Board {
	return &Board{
		level: level,
		up:    up,
		down:  down,
	}
}

func NewBoard(lines []string) *Board {
	if len(lines[0]) != 5 || len(lines) != 5 {
		panic("bad size")
	}

	b := makeBoard(0, nil, nil)
	for y := range lines {
		for x, r := range lines[y] {
			b.set(pos.P2{x, y}, r == '#')
		}
	}

	return b
}

func (b *Board) set(p pos.P2, val bool) {
	if p.X < 0 || p.Y < 0 || p.X >= 5 || p.Y >= 5 {
		panic("bad pos")
	}

	b.c[p.Y][p.X] = val
}

func (b *Board) Get(p pos.P2) bool {
	if p.X < 0 || p.Y < 0 || p.X >= 5 || p.Y >= 5 {
		return false
	}

	return b.c[p.Y][p.X]
}

func (b *Board) Evolve() *Board {
	nb := makeBoard(b.level, b.up, b.down)

	for y := range b.c {
		for x := range b.c[0] {
			p := pos.P2{x, y}
			neighbors := 0
			for _, dir := range dir.AllDirs {
				if np := dir.From(p); b.Get(np) {
					neighbors++
				}
			}

			if b.Get(p) {
				if neighbors == 1 {
					nb.set(p, true) // bug lives
				} else {
					// bug dies. no action because nb was
					// initalized to no bugs.
				}
			} else {
				if neighbors == 1 || neighbors == 2 {
					nb.set(p, true) // bug created
				}
			}
		}
	}

	return nb
}

func (b *Board) Strings() []string {
	out := []string{}
	for y := range b.c {
		line := ""
		for x := range b.c[0] {
			p := pos.P2{x, y}
			if b.Get(p) {
				line += "#"
			} else {
				line += "."
			}
		}
		out = append(out, line)
	}
	return out
}
