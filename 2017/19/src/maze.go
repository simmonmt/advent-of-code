package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	CellNone = iota
	CellVert
	CellHoriz
	CellLetter
	CellIntersection

	DirNorth = iota
	DirSouth
	DirEast
	DirWest
)

type Point struct {
	X, Y int
}

type Grid struct {
	cells [][]int
}

func NewGrid(cells [][]int) *Grid {
	return &Grid{cells}
}

func (g *Grid) Cell(p Point) int {
	if !g.IsIn(p) {
		return CellNone
	}

	return g.cells[p.Y][p.X]
}

func (g *Grid) IsIn(p Point) bool {
	if p.X < 0 || p.Y < 0 || p.Y >= len(g.cells) {
		return false
	}

	if p.X >= len(g.cells[p.Y]) {
		return false
	}

	return true
}

func readGrid(in io.Reader) (*Grid, map[Point]rune) {
	reader := bufio.NewReader(in)
	grid := [][]int{}
	letters := map[Point]rune{}

	for y := 0; ; y++ {
		grid = append(grid, []int{})

		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimRight(line, "\n")

		for x, c := range line {
			cell := CellNone
			switch c {
			case '|':
				cell = CellVert
				break
			case '-':
				cell = CellHoriz
				break
			case '+':
				cell = CellIntersection
				break
			default:
				if c != ' ' {
					cell = CellLetter
					letters[Point{x, y}] = c
				}
			}
			grid[y] = append(grid[y], cell)
		}
	}

	return NewGrid(grid), letters
}

func findStart(grid *Grid) Point {
	for x := 0; ; x++ {
		pos := Point{x, 0}
		if grid.Cell(pos) != CellNone {
			return pos
		}
	}
	panic("no start")
}

func nextCell(pos Point, dir int) Point {
	switch dir {
	case DirNorth:
		return Point{pos.X, pos.Y - 1}
	case DirSouth:
		return Point{pos.X, pos.Y + 1}
	case DirEast:
		return Point{pos.X + 1, pos.Y}
	case DirWest:
		return Point{pos.X - 1, pos.Y}
	default:
		panic(fmt.Sprintf("unexpected dir %v", dir))
	}
}

func nextDir(grid *Grid, pos Point, lastDir int) (int, Point) {
	dirs := map[int]int{
		DirNorth: DirSouth,
		DirSouth: DirNorth,
		DirEast:  DirWest,
		DirWest:  DirEast,
	}

	for lookDir, incompatDir := range dirs {
		if lastDir != incompatDir {
			nextPos := nextCell(pos, lookDir)
			if grid.Cell(nextPos) != CellNone {
				return lookDir, nextPos
			}
		}
	}

	panic("no next found")
}

func dirToString(dir int) string {
	switch dir {
	case DirNorth:
		return "north"
	case DirSouth:
		return "south"
	case DirEast:
		return "east"
	case DirWest:
		return "west"
	default:
		panic(fmt.Sprintf("unexpected dir %v", dir))
	}
}

func advance(grid *Grid, letters map[Point]rune, pos Point, dir int) (Point, int, rune) {
	// Can we move in the direction we want?
	nextPos := nextCell(pos, dir)
	if cell := grid.Cell(nextPos); cell == CellLetter {
		fmt.Printf("found %c\n", letters[nextPos])
		return nextPos, dir, letters[nextPos]
	} else if cell != CellNone {
		return nextPos, dir, rune(0)
	}

	if cell := grid.Cell(pos); cell != CellIntersection {
		// They've advanced off the end
		return nextPos, dir, rune(0)
	}

	// Find next direction
	nextDir, nextPos := nextDir(grid, pos, dir)
	if cell := grid.Cell(nextPos); cell == CellLetter {
		fmt.Printf("found %c\n", letters[nextPos])
		return nextPos, nextDir, letters[nextPos]
	} else {
		return nextPos, nextDir, rune(0)
	}
}

func main() {
	grid, letters := readGrid(os.Stdin)
	pos := findStart(grid)
	dir := DirSouth

	foundLetters := []rune{}
	steps := 0
	for grid.IsIn(pos) {
		steps++
		//fmt.Printf("pos %v dir %v\n", pos, dirToString(dir))

		var letter rune
		pos, dir, letter = advance(grid, letters, pos, dir)
		if letter != rune(0) {
			foundLetters = append(foundLetters, letter)
		}
	}

	for _, foundLetter := range foundLetters {
		fmt.Printf("%c", foundLetter)
	}
	fmt.Println()
	fmt.Println(steps - 1)
}
