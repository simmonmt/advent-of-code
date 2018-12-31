package main

import (
	"fmt"
	"intmath"
	"strings"
)

type PlayState struct {
	Empty, Goal Pos
}

func (ps *PlayState) Encode() string {
	return fmt.Sprintf("%d,%d,%d,%d",
		ps.Empty.X, ps.Empty.Y, ps.Goal.X, ps.Goal.Y)
}

func Decode(str string) *PlayState {
	parts := strings.Split(str, ",")
	ex := intmath.AtoiOrDie(parts[0])
	ey := intmath.AtoiOrDie(parts[1])
	gx := intmath.AtoiOrDie(parts[2])
	gy := intmath.AtoiOrDie(parts[3])

	return &PlayState{
		Empty: Pos{ex, ey},
		Goal:  Pos{gx, gy},
	}
}

type Board struct {
	cells [][]bool
}

func NewBoard(w, h int) *Board {
	cells := make([][]bool, h)
	for y := range cells {
		cells[y] = make([]bool, w)
	}

	return &Board{
		cells: cells,
	}
}

func (b *Board) Dump(state *PlayState) {
	for y := range b.cells {
		for x, c := range b.cells[y] {
			p := Pos{x, y}

			switch {
			case p.Eq(state.Goal):
				fmt.Print("G")
			case p.Eq(state.Empty):
				fmt.Print("_")
			case c == true:
				fmt.Print(".")
			default:
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
}

func (b *Board) Set(pos Pos, moveable bool) {
	b.cells[pos.Y][pos.X] = moveable
}

func (b *Board) IsMoveable(pos Pos) bool {
	return b.cells[pos.Y][pos.X]
}

func (b *Board) Size() (w, h int) {
	return len(b.cells[0]), len(b.cells)
}
