package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"

	"intmath"
	"logger"

	"github.com/soniakeys/graph"
)

var (
	verbose  = flag.Bool("verbose", false, "verbose")
	numTurns = flag.Int("num_turns", 1, "num turns")
)

type Pos struct {
	X, Y int
}

type Char struct {
	IsElf bool
	P     Pos
}

type CharByReadingOrder []Char

func (a CharByReadingOrder) Len() int      { return len(a) }
func (a CharByReadingOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a CharByReadingOrder) Less(i, j int) bool {
	if a[i].P.Y != a[j].P.Y {
		return a[i].P.Y < a[j].P.Y
	}
	return a[i].P.X < a[j].P.X
}

type PosByReadingOrder []Pos

func (a PosByReadingOrder) Len() int      { return len(a) }
func (a PosByReadingOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a PosByReadingOrder) Less(i, j int) bool {
	if a[i].Y != a[j].Y {
		return a[i].Y < a[j].Y
	}
	return a[i].X < a[j].X
}

type Board [][]rune

func (b *Board) Dump(chars []Char) {
	b.DumpWithDecoration(chars, nil, ' ')
}

func (b *Board) DumpWithDecoration(chars []Char, decorations []Pos, decorationChar rune) {
	charsByPos := map[Pos]Char{}
	for _, char := range chars {
		charsByPos[char.P] = char
	}

	decsByPos := map[Pos]bool{}
	for _, decPos := range decorations {
		decsByPos[decPos] = true
	}

	for y, row := range *b {
		for x, r := range row {
			pos := Pos{x, y}

			if len(decsByPos) > 0 {
				if _, found := decsByPos[pos]; found {
					fmt.Print(string(decorationChar))
					continue
				}
			}

			char, found := charsByPos[Pos{x, y}]
			if !found {
				fmt.Print(string(r))
				continue
			}

			if char.IsElf {
				fmt.Print("E")
			} else {
				fmt.Print("G")
			}
		}
		fmt.Println()
	}
}

func readInput() (Board, []Char, error) {
	board := Board{}
	chars := []Char{}

	scanner := bufio.NewScanner(os.Stdin)
	for y := 0; scanner.Scan(); y++ {
		line := scanner.Text()

		rs := []rune(line)
		for x, r := range rs {
			isElf := false
			switch r {
			case 'E':
				isElf = true
				fallthrough
			case 'G':
				pos := Pos{x, y}
				chars = append(chars, Char{isElf, pos})
				rs[x] = '.'
			}
		}

		board = append(board, rs)
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("read failed: %v", err)
	}

	return board, chars, nil
}

func posToIdx(board Board, pos Pos) graph.NI {
	width := len(board[0])
	return graph.NI(width*pos.Y + pos.X)
}

func idxToPos(board Board, idx graph.NI) Pos {
	width := len(board[0])
	x := int(idx) % width
	y := int(idx) / width
	return Pos{x, y}
}

func pathToPosns(board Board, path graph.LabeledPath) []Pos {
	nodes := []Pos{}
	for _, n := range path.Path {
		nodes = append(nodes, idxToPos(board, n.To))
	}
	return nodes
}

func pathToStr(board Board, path graph.LabeledPath) string {
	return fmt.Sprintf("%+v", pathToPosns(board, path))
}

func boardToAdjacencyList(board Board) graph.LabeledAdjacencyList {
	width := len(board[0])
	height := len(board)
	al := make([][]graph.Half, width*height)

	for y, row := range board {
		for x, c := range row {
			idx := width*y + x
			pos := Pos{x, y}

			if c != '.' {
				continue
			}

			al[idx] = []graph.Half{}
			for _, open := range openPosns(board, pos) {
				h := graph.Half{
					To:    posToIdx(board, open),
					Label: graph.LI(1),
				}
				al[idx] = append(al[idx], h)
			}
		}
	}

	return graph.LabeledAdjacencyList(al)
}

// Adding a character to an adjacency list means making that location unroutable.
// We need to remove all outgoing links as well as all incoming ones.
func addCharToAdjacencyList(board Board, char *Char, al graph.LabeledAdjacencyList) {
	charIdx := posToIdx(board, char.P)

	// outgoing
	al[charIdx] = nil

	// incoming
	for _, open := range openPosns(board, char.P) {
		idx := posToIdx(board, open)
		filtered := []graph.Half{}
		for _, h := range al[idx] {
			if h.To != charIdx {
				filtered = append(filtered, h)
			}
		}
		al[idx] = filtered
	}
}

// Removing a character from an adjacency list means making that location routable again.
func removeCharFromAdjacencyList(board Board, char *Char, al graph.LabeledAdjacencyList) {
	charIdx := posToIdx(board, char.P)
	outgoing := []graph.Half{}

	for _, open := range openPosns(board, char.P) {
		openIdx := posToIdx(board, open)
		al[openIdx] = append(al[openIdx], graph.Half{charIdx, graph.LI(1)})
		outgoing = append(outgoing, graph.Half{openIdx, graph.LI(1)})
	}

	al[charIdx] = outgoing
}

func astarHeuristic(board Board, to Pos, fromIdx graph.NI) float64 {
	from := idxToPos(board, fromIdx)
	return float64(intmath.Abs(from.X-to.X) + intmath.Abs(from.Y-to.Y))
}

func findShortestPaths(board Board, char *Char, inRange []Pos, al graph.LabeledAdjacencyList) (map[Pos][]Pos, int) {
	paths := map[Pos][]Pos{}

	// We have to remove char or we can't find paths from it. We
	// need to add it back before returning.
	removeCharFromAdjacencyList(board, char, al)
	defer addCharToAdjacencyList(board, char, al)

	shortestLen := math.MaxInt32
	for _, cand := range inRange {
		heuristic := func(fromIdx graph.NI) float64 {
			return astarHeuristic(board, cand, fromIdx)
		}

		idxPath, _ := al.AStarAPath(
			posToIdx(board, char.P),
			posToIdx(board, cand), heuristic, func(l graph.LI) float64 { return 1 })

		path := pathToPosns(board, idxPath)
		if len(path) == 0 || len(path) > shortestLen {
			continue
		}

		if len(path) < shortestLen {
			paths = map[Pos][]Pos{}
			shortestLen = len(path)
		}
		paths[cand] = path
	}

	return paths, shortestLen
}

var validDirs = []Pos{Pos{0, -1}, Pos{-1, 0}, Pos{1, 0}, Pos{0, 1}}

func validPosns(pos Pos) []Pos {
	cands := make([]Pos, 4)
	for i, validDir := range validDirs {
		cands[i] = Pos{pos.X + validDir.X, pos.Y + validDir.Y}
	}
	return cands
}

func openPosns(board Board, pos Pos) []Pos {
	open := []Pos{}
	for _, cand := range validPosns(pos) {
		if board[cand.Y][cand.X] == '.' {
			open = append(open, cand)
		}
	}
	return open
}

func neighborToAttack(board Board, char *Char, others map[Pos]Char) *Char {
	for _, candPos := range validPosns(char.P) {
		if other, found := others[candPos]; found {
			if char.IsElf != other.IsElf {
				return &other
			}
		}
	}

	return nil
}

func findInRange(board Board, others map[Pos]Char) []Pos {
	inRangeMap := map[Pos]bool{}

	for _, other := range others {
		for _, candPos := range openPosns(board, other.P) {
			inRangeMap[candPos] = true
		}
	}

	inRange := make([]Pos, len(inRangeMap))
	i := 0
	for p := range inRangeMap {
		inRange[i] = p
		i++
	}
	return inRange
}

func findAllPaths(board Board, start, end Pos, maxLen int) [][]Pos {
	//fmt.Printf("findallpaths %v -> %v maxlen %v\n", start, end, maxLen)

	visited := map[Pos]bool{}
	path := make([]Pos, maxLen+1)
	pathIdx := 0

	out := [][]Pos{}

	visitor := func(path []Pos) {
		//fmt.Printf("visitor called with %v\n", path)
		cp := make([]Pos, len(path)-1)
		copy(cp, path[1:])
		out = append(out, cp)
	}

	findAllPathsHelper(board, start, end, visited, path, pathIdx, maxLen+1, visitor)

	//	fmt.Println("findallpaths done")

	return out
}

func findAllPathsHelper(board Board, cur, end Pos, visited map[Pos]bool, path []Pos, pathIdx, maxLen int, visitor func(path []Pos)) {
	//fmt.Printf("findallpathshelper %v -> %v maxlen %v path %v pathIdx %v\n",
	//cur, end, maxLen, path, pathIdx)

	visited[cur] = true
	defer delete(visited, cur)

	path[pathIdx] = cur
	pathIdx++

	if cur.X == end.X && cur.Y == end.Y {
		visitor(path)
	} else if pathIdx < maxLen {
		for _, open := range openPosns(board, cur) {
			if _, found := visited[open]; !found {
				findAllPathsHelper(board, open, end, visited, path, pathIdx,
					maxLen, visitor)
			}
		}
	} else {
		//fmt.Println("giving up too long")
	}
}

func findNextMove(board Board, char *Char, chars []Char, others map[Pos]Char, al graph.LabeledAdjacencyList) (Pos, bool) {
	// No neighbor; have to move then attack
	inRange := findInRange(board, others)

	if *verbose {
		fmt.Println("In range:")
		fmt.Println(inRange)
		board.DumpWithDecoration(chars, inRange, '?')
	}

	paths, shortestLen := findShortestPaths(board, char, inRange, al)
	if *verbose {
		fmt.Println("Shortest paths:")
		for dest, path := range paths {
			fmt.Printf("%v: %v\n", dest, path)
		}
	}

	if len(paths) > 0 {
		closestInRange := []Pos{}
		for dest := range paths {
			closestInRange = append(closestInRange, dest)
		}

		if *verbose {
			fmt.Println("Nearest:")
			board.DumpWithDecoration(chars, closestInRange, '!')
		}

		sort.Sort(PosByReadingOrder(closestInRange))
		chosen := closestInRange[0]

		if *verbose {
			fmt.Println("Chosen:")
			board.DumpWithDecoration(chars, []Pos{chosen}, '+')
		}

		allPaths := findAllPaths(board, char.P, chosen, shortestLen)

		firstStepsMap := map[Pos]bool{}
		for _, path := range allPaths {
			firstStepsMap[path[0]] = true
		}

		firstSteps := []Pos{}
		for step := range firstStepsMap {
			firstSteps = append(firstSteps, step)
		}
		sort.Sort(PosByReadingOrder(firstSteps))

		nextStep := firstSteps[0]
		logger.LogF("chosen step: %v", nextStep)

		return nextStep, true
	}

	return Pos{}, false
}

func turnForChar(board Board, charNum int, chars []Char, al graph.LabeledAdjacencyList) {
	char := &chars[charNum]

	others := map[Pos]Char{}
	for otherNum, other := range chars {
		if otherNum != charNum && other.IsElf != char.IsElf {
			others[other.P] = other
		}
	}

	toAttack := neighborToAttack(board, char, others)
	if toAttack == nil {
		// No neighbor to attack, so let's try to move
		newPos, shouldMove := findNextMove(board, char, chars, others, al)

		if !shouldMove {
			return // nothing we can do
		}

		// Move the character, both in the chars
		// list and in the adjacency list
		removeCharFromAdjacencyList(board, char, al)
		char.P = newPos
		addCharToAdjacencyList(board, char, al)

		toAttack = neighborToAttack(board, char, others)
		if toAttack == nil {
			logger.LogF("nobody to attack; turn done")
			return
		}
	}

	logger.LogF("-- chosen to attack: %v\n", toAttack)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	board, chars, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	//board.Dump(chars)

	al := boardToAdjacencyList(board)

	// Add all characters to the adjacency list
	for _, char := range chars {
		addCharToAdjacencyList(board, &char, al)
	}

	// heuristic := func(fromIdx graph.NI) float64 {
	// 	from := idxToPos(board, fromIdx)
	// 	return float64(intmath.Abs(from.X-end.X) + intmath.Abs(from.Y-end.Y))
	// }

	// weight := func(l graph.LI) float64 { return 1 }

	// start := Pos{1, 1}
	// end := Pos{5, 3}
	// path, dist := al.AStarAPath(
	// 	posToIdx(board, start),
	// 	posToIdx(board, end), heuristic, weight)
	// fmt.Printf("path %v, dist %v\n", pathToStr(board, path), dist)

	for turn := 0; turn < *numTurns; turn++ {
		if *verbose {
			fmt.Printf("turn %d: start\n", turn)
			board.Dump(chars)
		}

		sort.Sort(CharByReadingOrder(chars))
		logger.LogF("character order: %+v\n", chars)

		for charNum, _ := range chars {
			logger.LogF("-- character %d: %+v\n", charNum, chars[charNum])

			turnForChar(board, charNum, chars, al)

		}
	}

	fmt.Println("End")
	board.Dump(chars)
}
