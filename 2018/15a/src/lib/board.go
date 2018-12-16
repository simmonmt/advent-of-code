package lib

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"intmath"
	"logger"

	"github.com/soniakeys/graph"
)

type Board struct {
	w     int
	cells [][]rune
	chars map[Pos]*Char
	graph *Graph
}

func NewBoard(rows []string) *Board {
	w := len(rows[0])
	h := len(rows)

	cells := make([][]rune, h)
	for y, row := range rows {
		if len(row) != w {
			panic("uneven rows")
		}
		cells[y] = []rune(row)
	}

	chars := map[Pos]*Char{}
	charNum := 0
	graph := NewGraph(h * w)

	for y := range cells {
		for x, r := range cells[y] {
			pos := Pos{x, y}
			isElf := false
			switch r {
			case '.':
				for _, neighbor := range openSurroundingPositions(cells, pos) {
					if !hasEdge(graph, w, pos, neighbor) {
						addEdge(graph, w, pos, neighbor)
					}
				}
			case 'E':
				isElf = true
				fallthrough
			case 'G':
				pos := Pos{x, y}
				chars[pos] = NewChar(charNum, isElf, pos)
				charNum++
			}
		}
	}

	return &Board{
		w:     w,
		cells: cells,
		chars: chars,
		graph: graph,
	}
}

func hasEdge(graph *Graph, w int, a, b Pos) bool {
	return graph.HasEdge(posToNI(a, w), posToNI(b, w))
}

func addEdge(graph *Graph, w int, a, b Pos) {
	graph.AddEdge(posToNI(a, w), posToNI(b, w))
}

func removeEdge(graph *Graph, w int, a, b Pos) {
	graph.RemoveEdge(posToNI(a, w), posToNI(b, w))
}

func (b *Board) getCell(pos Pos) rune {
	return b.cells[pos.Y][pos.X]
}

func (b *Board) setCell(pos Pos, val rune) {
	b.cells[pos.Y][pos.X] = val
}

func (b *Board) Dump() {
	b.DumpWithDecorations(nil, ' ')
}

func (b *Board) DumpWithDecoration(decoration Pos, decorationChar rune) {
	b.DumpWithDecorations([]Pos{decoration}, decorationChar)
}

func (b *Board) DumpWithDecorations(decorations []Pos, decorationChar rune) {
	decsByPos := map[Pos]bool{}
	for _, decPos := range decorations {
		decsByPos[decPos] = true
	}

	fmt.Printf(" ")
	for x := range b.cells[0] {
		fmt.Print(strconv.Itoa(x % 10))
	}
	fmt.Println()

	for y := range b.cells {
		fmt.Print(strconv.Itoa(y % 10))

		charsOnLine := []*Char{}
		for x, r := range b.cells[y] {
			pos := Pos{x, y}

			if len(decsByPos) > 0 {
				if _, found := decsByPos[pos]; found {
					fmt.Print(string(decorationChar))
					continue
				}
			}

			char, found := b.chars[pos]
			if !found {
				fmt.Print(string(r))
				continue
			}

			charsOnLine = append(charsOnLine, char)
			fmt.Print(string(char.Short()))
		}

		strs := []string{}
		for _, char := range charsOnLine {
			strs = append(strs, fmt.Sprintf("%s(%d)", string(char.Short()), char.HP))
		}
		fmt.Printf("   %s\n", strings.Join(strs, ", "))
	}
}

func (b *Board) DumpEdges() {
	graphEdges := b.graph.AllEdges()
	knownEdges := []posEdge{}
	for _, e := range graphEdges {
		knownEdges = append(knownEdges, makePosEdge(b.niToPos(e.From), b.niToPos(e.To)))
	}

	sort.Sort(byPosEdgeReadingOrder(knownEdges))
	for _, e := range knownEdges {
		fmt.Printf("%v <-> %v\n", e.a, e.b)
	}
}

func niToPos(ni graph.NI, w int) Pos {
	y := int(ni) / w
	x := int(ni) % w
	return Pos{x, y}
}

func posToNI(pos Pos, w int) graph.NI {
	return graph.NI(pos.Y*w + pos.X)
}

func (b *Board) niToPos(ni graph.NI) Pos {
	return niToPos(ni, b.w)
}

func (b *Board) posToNI(pos Pos) graph.NI {
	return posToNI(pos, b.w)
}

var validDirs = []Pos{Pos{0, -1}, Pos{-1, 0}, Pos{1, 0}, Pos{0, 1}}

func visitNeighbors(cells [][]rune, pos Pos, visitor func(pos Pos, contents rune)) {
	for _, validDir := range validDirs {
		cand := Pos{pos.X + validDir.X, pos.Y + validDir.Y}
		visitor(cand, cells[pos.Y][pos.X])
	}
}

func (b *Board) SurroundingCharacters(pos Pos) []Char {
	out := []Char{}
	visitNeighbors(b.cells, pos, func(cand Pos, contents rune) {
		if char, found := b.chars[cand]; found {
			out = append(out, *char)
		}
	})
	return out
}

func openSurroundingPositions(cells [][]rune, pos Pos) []Pos {
	out := []Pos{}
	visitNeighbors(cells, pos, func(cand Pos, contents rune) {
		if cells[cand.Y][cand.X] == '.' {
			out = append(out, cand)
		}
	})
	return out
}

// surrounding positions that are '.' (I.e. not walls and not characters)
func (b *Board) OpenSurroundingPositions(pos Pos) []Pos {
	return openSurroundingPositions(b.cells, pos)
}

type posEdge struct {
	a, b Pos
}

type byPosEdgeReadingOrder []posEdge

func (a byPosEdgeReadingOrder) Len() int      { return len(a) }
func (a byPosEdgeReadingOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byPosEdgeReadingOrder) Less(i, j int) bool {
	if PosLess(a[i].a, a[j].a) {
		return true
	}
	if PosLess(a[j].a, a[i].a) {
		return false
	}
	return PosLess(a[i].b, a[j].b)
}

func makePosEdge(a, b Pos) posEdge {
	if PosLess(a, b) {
		return posEdge{a, b}
	}
	return posEdge{b, a}
}

func (b *Board) Validate() bool {
	graphEdges := b.graph.AllEdges()

	knownEdges := map[posEdge]bool{}
	for _, e := range graphEdges {
		knownEdges[makePosEdge(b.niToPos(e.From), b.niToPos(e.To))] = false
	}

	valid := true
	for y := range b.cells {
		for x, r := range b.cells[y] {
			pos := Pos{x, y}

			if r != '.' {
				// If pos isn't open it we don't
				// expect any inbound/outbound edges.
				continue
			}

			// We expect edges between pos and each open surrounding position
			for _, neighbor := range openSurroundingPositions(b.cells, pos) {
				expectedEdge := makePosEdge(pos, neighbor)
				if _, found := knownEdges[expectedEdge]; !found {
					fmt.Printf("expected edge %v <-> %v\n", pos, neighbor)
				} else {
					knownEdges[expectedEdge] = true
				}
			}
		}
	}

	for e, found := range knownEdges {
		if !found {
			fmt.Printf("unattributed edge %v <-> %v\n", e.a, e.b)
			valid = false
		}
	}

	if valid {
		fmt.Println("valid board")
	}
	return valid
}

func (b *Board) Chars() []Char {
	chars := make([]Char, len(b.chars))
	i := 0
	for _, c := range b.chars {
		chars[i] = *c
		i++
	}
	return chars
}

func (b *Board) niDistance(from, to graph.NI) int {
	fromPos := b.niToPos(from)
	toPos := b.niToPos(to)
	return intmath.Abs(fromPos.X-toPos.X) + intmath.Abs(fromPos.Y-toPos.Y)
}

func (b *Board) ShortestPath(from, to Pos) []Pos {
	edgesAdded := []posEdge{}

	if r := b.getCell(from); r != '.' {
		if r != 'E' && r != 'G' {
			panic("bad from")
		}

		// We're being asked to route from a cell with a
		// character. This cell won't have any routability, so
		// we need to add some.
		for _, neighbor := range openSurroundingPositions(b.cells, from) {
			edge := posEdge{from, neighbor}
			b.graph.AddEdge(b.posToNI(from), b.posToNI(neighbor))
			edgesAdded = append(edgesAdded, edge)
		}
	}

	niPath := b.graph.ShortestPath(b.posToNI(from), b.posToNI(to),
		func(from, to graph.NI) int { return b.niDistance(from, to) })

	// Remove any edges we added to allow routing from 'from'
	for _, edge := range edgesAdded {
		b.graph.RemoveEdge(b.posToNI(edge.a), b.posToNI(edge.b))
	}

	b.Validate()

	if len(niPath) == 0 {
		return nil
	}

	path := []Pos{}
	for _, ni := range niPath {
		path = append(path, b.niToPos(ni))
	}
	return path
}

func (b *Board) MoveChar(src Char, destPos Pos) Char {
	dest := src
	dest.P = destPos

	b.RemoveChar(src)
	b.AddChar(dest)
	b.Validate()

	return dest
}

func (b *Board) RemoveChar(char Char) {
	logger.LogF("removing %v", char)
	if r := b.getCell(char.P); r != 'G' && r != 'E' {
		panic("no char there")
	}
	b.setCell(char.P, '.')

	delete(b.chars, char.P)

	for _, neighbor := range b.OpenSurroundingPositions(char.P) {
		addEdge(b.graph, b.w, char.P, neighbor)
	}
}

func (b *Board) AddChar(char Char) {
	logger.LogF("adding %v", char)
	if b.getCell(char.P) != '.' {
		panic("not open")
	}
	b.setCell(char.P, char.Short())

	newChar := &Char{}
	*newChar = char

	b.chars[char.P] = newChar

	for _, neighbor := range b.OpenSurroundingPositions(char.P) {
		removeEdge(b.graph, b.w, char.P, neighbor)
	}
}

func (b *Board) Attack(er, ee Char) (Char, bool) {
	if _, found := b.chars[ee.P]; !found {
		panic("ee isn't in chars")
	}
	if _, found := b.chars[er.P]; !found {
		panic("er isn't in chars")
	}

	b.chars[ee.P].HP -= b.chars[er.P].AP
	ee = *b.chars[ee.P]
	return ee, ee.HP < 0
}
