package main

import (
	"fmt"
	"math"

	"intmath"
)

type Door struct {
	From, To Pos
}

type Board struct {
	origin  Pos
	visited map[Pos]bool
	doors   map[Pos][]Pos
}

func NewBoard(origin Pos) *Board {
	return &Board{
		origin:  origin,
		visited: map[Pos]bool{},
		doors:   map[Pos][]Pos{},
	}
}

var (
	dirs = map[rune]Pos{
		'N': Pos{0, -1},
		'S': Pos{0, 1},
		'E': Pos{1, 0},
		'W': Pos{-1, 0},
	}
)

func containsPos(l []Pos, p Pos) bool {
	if l == nil {
		return false
	}
	for _, e := range l {
		if e.Eq(p) {
			return true
		}
	}
	return false
}

func (b *Board) addDoor(p1, p2 Pos) {

	add := func(from, to Pos) {
		if b.doors[from] == nil {
			b.doors[from] = []Pos{to}
		} else if !containsPos(b.doors[from], to) {
			b.doors[from] = append(b.doors[from], to)
		}
	}

	add(p1, p2)
	add(p2, p1)
}

func (b *Board) hasDoor(p1, p2 Pos) bool {
	return containsPos(b.doors[p1], p2)
}

func (b *Board) Move(cur Pos, dir rune) Pos {
	b.visited[cur] = true

	off, found := dirs[dir]
	if !found {
		panic("bad dir " + string(dir))
	}

	new := Pos{X: cur.X + off.X, Y: cur.Y + off.Y}
	b.visited[new] = true
	b.addDoor(cur, new)
	return new
}

func (b *Board) BFS(start Pos, visitor func(p Pos, neighbors []Pos)) {
	visited := map[Pos]bool{}
	todo := []Pos{start}

	for len(todo) > 0 {
		cur := todo[0]
		todo = todo[1:]
		visitor(cur, b.doors[cur])
		visited[cur] = true

		for _, n := range b.doors[cur] {
			if _, found := visited[n]; found {
				continue
			}
			todo = append(todo, n)
		}
	}
}

func (b *Board) Dump() {
	xmin, xmax := math.MaxInt32, math.MinInt32
	ymin, ymax := math.MaxInt32, math.MinInt32

	for p := range b.visited {
		xmin = intmath.IntMin(xmin, p.X)
		xmax = intmath.IntMax(xmax, p.X)
		ymin = intmath.IntMin(ymin, p.Y)
		ymax = intmath.IntMax(ymax, p.Y)
	}

	for y := ymin; y <= ymax; y++ {
		for x := xmin; x <= xmax; x++ {
			p := Pos{x, y}
			above := Pos{x, y - 1}

			fmt.Print("#")
			if b.hasDoor(above, p) {
				fmt.Print("-")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println("#")

		for x := xmin; x <= xmax; x++ {
			p := Pos{x, y}
			left := Pos{x - 1, y}

			if b.hasDoor(left, p) {
				fmt.Print("|")
			} else {
				fmt.Print("#")
			}

			if p.Eq(b.origin) {
				fmt.Print("X")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("#")
	}

	for x := xmin; x <= xmax; x++ {
		fmt.Printf("##")
	}
	fmt.Println("#")
}
