package puzzle

import (
	"github.com/simmonmt/aoc/2019/common/dir"
	"github.com/simmonmt/aoc/2019/common/pos"
)

type Board struct {
	w, h int
	c    []bool
}

func makeBoard(w, h int) *Board {
	b := &Board{
		w: w,
		h: h,
	}

	b.c = make([]bool, b.w*b.h)
	return b
}

func NewBoard(lines []string) *Board {
	b := makeBoard(len(lines[0]), len(lines))

	for y := range lines {
		for x, r := range lines[y] {
			b.set(pos.P2{x, y}, r == '#')
		}
	}

	return b
}

func (b *Board) set(p pos.P2, val bool) {
	if p.X < 0 || p.Y < 0 {
		panic("bad pos")
	}

	i := p.Y*b.w + p.X
	if i >= len(b.c) {
		panic("bat pos")
	}
	b.c[i] = val
}

func (b *Board) Get(p pos.P2) bool {
	if p.X < 0 || p.Y < 0 || p.X >= b.w || p.Y >= b.h {
		return false
	}

	i := p.Y*b.w + p.X
	if i >= len(b.c) {
		return false
	}
	return b.c[i]
}

func (b *Board) Evolve() *Board {
	nb := makeBoard(b.w, b.h)

	for y := 0; y < b.h; y++ {
		for x := 0; x < b.w; x++ {
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

func (b *Board) Hash() string {
	h := make([]rune, len(b.c))
	for i, v := range b.c {
		if v {
			h[i] = '#'
		} else {
			h[i] = '.'
		}
	}
	return string(h)
}

func (b *Board) Strings() []string {
	out := []string{}
	for y := 0; y < b.h; y++ {
		line := ""
		for x := 0; x < b.w; x++ {
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

func (b *Board) Biodiversity() int {
	pow := 1
	out := 0
	for _, v := range b.c {
		if v {
			out += pow
		}
		pow *= 2
	}
	return out
}
