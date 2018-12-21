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
	visited map[Pos]bool
	doors   map[Door]bool
}

func NewBoard() *Board {
	return &Board{
		visited: map[Pos]bool{},
		doors:   map[Door]bool{},
	}
}

var (
	dirs = map[string]Pos{
		"N": Pos{0, -1},
		"S": Pos{0, 1},
		"E": Pos{1, 0},
		"W": Pos{-1, 0},
	}
)

func (b *Board) Move(cur Pos, dir string) Pos {
	b.visited[cur] = true

	off, found := dirs[dir]
	if !found {
		panic("bad dir " + dir)
	}

	new := Pos{X: cur.X + off.X, Y: cur.Y + off.Y}
	b.visited[new] = true

	var door Door
	if cur.Before(new) {
		door = Door{cur, new}
	} else {
		door = Door{new, cur}
	}
	b.doors[door] = true

	return new
}

func (b *Board) Dump(cur Pos) {
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
			if _, found := b.doors[Door{above, p}]; found {
				fmt.Print("-")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println("#")

		for x := xmin; x <= xmax; x++ {
			p := Pos{x, y}
			left := Pos{x - 1, y}

			if _, found := b.doors[Door{left, p}]; found {
				fmt.Print("|")
			} else {
				fmt.Print("#")
			}

			if p.Eq(cur) {
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
