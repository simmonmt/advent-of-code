package lib

import "fmt"

type CellType int

const (
	TYPE_OPEN CellType = iota
	TYPE_TREES
	TYPE_LUMBER
)

func (t CellType) String() string {
	switch t {
	case TYPE_OPEN:
		return "open"
	case TYPE_TREES:
		return "trees"
	case TYPE_LUMBER:
		return "lumberyard"
	default:
		panic("unknown")
	}
}

func (t CellType) Short() string {
	switch t {
	case TYPE_OPEN:
		return "."
	case TYPE_TREES:
		return "|"
	case TYPE_LUMBER:
		return "#"
	default:
		panic("unknown")
	}
}

type Board struct {
	cells [][]CellType
}

func NewBoardFromString(lines []string) *Board {
	h := len(lines)
	w := len(lines[0])

	cells := make([][]CellType, h)
	for y := 0; y < h; y++ {
		row := make([]CellType, w)

		for x, c := range lines[y] {
			switch c {
			case '.':
				row[x] = TYPE_OPEN
			case '|':
				row[x] = TYPE_TREES
			case '#':
				row[x] = TYPE_LUMBER
			default:
				panic("bad char")
			}
		}

		cells[y] = row
	}

	return &Board{cells: cells}
}

func (b *Board) Dump() {
	for y := range b.cells {
		for _, c := range b.cells[y] {
			fmt.Print(c.Short())
		}
		fmt.Println()
	}
}

type offset struct {
	dx, dy int
}

var (
	neighborDirs = []offset{offset{-1, -1}, offset{-1, 0}, offset{-1, 1},
		offset{0, -1}, offset{0, 1},
		offset{1, -1}, offset{1, 0}, offset{1, 1},
	}
)

func countNeighbors(cells [][]CellType, x, y int, typ CellType) int {
	get := func(x, y int) CellType {
		if x < 0 || y < 0 || y >= len(cells) || x >= len(cells[0]) {
			return TYPE_OPEN
		}
		return cells[y][x]
	}

	num := 0
	for _, dir := range neighborDirs {
		nx := x + dir.dx
		ny := y + dir.dy
		if get(nx, ny) == typ {
			num++
		}
	}

	return num
}

func (b *Board) Step() {
	old := make([][]CellType, len(b.cells))
	for y := range b.cells {
		old[y] = make([]CellType, len(b.cells[y]))
		copy(old[y], b.cells[y])
	}

	for y := range old {
		for x, c := range old[y] {
			numTrees := countNeighbors(old, x, y, TYPE_TREES)
			numLumber := countNeighbors(old, x, y, TYPE_LUMBER)

			switch c {
			case TYPE_OPEN:
				if numTrees >= 3 {
					b.cells[y][x] = TYPE_TREES
				}
			case TYPE_TREES:
				if numLumber >= 3 {
					b.cells[y][x] = TYPE_LUMBER
				}
			case TYPE_LUMBER:
				if numLumber >= 1 && numTrees >= 1 {
					// stays lumber
				} else {
					b.cells[y][x] = TYPE_OPEN
				}
			}
		}
	}
}

func (b *Board) Score() (numWoods, numLumber int) {
	numWoods, numLumber = 0, 0
	for y := range b.cells {
		for _, c := range b.cells[y] {
			if c == TYPE_TREES {
				numWoods++
			} else if c == TYPE_LUMBER {
				numLumber++
			}
		}
	}
	return
}
