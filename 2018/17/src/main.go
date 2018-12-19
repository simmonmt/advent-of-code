// 37580 too low
// 37511 too low
// 38453 too high
// 38452 no

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"image/png"
	"intmath"
	"lib"
	"logger"
)

var (
	verbose  = flag.Bool("verbose", false, "verbose")
	spring   = flag.String("spring", "500,0", "spring location")
	dump     = flag.Bool("dump", false, "dump the parsed board and exit")
	imageOut = flag.String("image_out", "", "where to write a png of the result")
)

type Direction int

const (
	DIR_LEFT Direction = iota
	DIR_RIGHT
	DIR_UP
	DIR_DOWN
)

func parseRange(str string) (min, max int) {
	parts := strings.SplitN(str, "..", 2)
	if len(parts) == 1 {
		val := intmath.AtoiOrDie(parts[0])
		return val, val
	} else {
		return intmath.AtoiOrDie(parts[0]), intmath.AtoiOrDie(parts[1])
	}
}

func readInput() ([]lib.InputLine, error) {
	lines := []string{}
	inputLines := []lib.InputLine{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")

		xpart := strings.TrimSpace(parts[0])
		ypart := strings.TrimSpace(parts[1])
		if xpart > ypart {
			xpart, ypart = ypart, xpart
		}

		if xpart[0] != 'x' || ypart[0] != 'y' {
			return nil, fmt.Errorf("invalid line %v", line)
		}

		xmin, xmax := parseRange(xpart[2:])
		ymin, ymax := parseRange(ypart[2:])

		inputLines = append(inputLines, lib.InputLine{xmin, xmax, ymin, ymax})

		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return inputLines, nil
}

func move(pos lib.Pos, dir Direction) lib.Pos {
	switch dir {
	case DIR_UP:
		return lib.Pos{pos.X, pos.Y - 1}
	case DIR_DOWN:
		return lib.Pos{pos.X, pos.Y + 1}
	case DIR_LEFT:
		return lib.Pos{pos.X - 1, pos.Y}
	case DIR_RIGHT:
		return lib.Pos{pos.X + 1, pos.Y}
	default:
		panic("bad dir")
	}
}

func progress(board *lib.Board, pos lib.Pos) (newPos lib.Pos, newStarts []lib.Pos) {
	// By default we don't move, and say we're done.

	cell := board.Get(pos)
	if cell == lib.TYPE_SPRING {
		return move(pos, DIR_DOWN), nil
	}

	if cell != lib.TYPE_FLOW {
		panic(fmt.Sprintf("unknown cell type %s at %v", cell, pos))
	}

	candPos := move(pos, DIR_DOWN)
	if !board.InBounds(candPos) {
		// It ran off the edge of the board, so no further analysis.
		// needed.
		return candPos, nil
	}

	candCell := board.Get(candPos)

	// If it's open below, pour one step down.
	if candCell == lib.TYPE_OPEN {
		return candPos, nil
	}

	// It's not open below. If we're blocked to one side, flow to the other
	// side.
	if candCell == lib.TYPE_FLOW {
		// There's already flow below, so no progress can be
		// made. Return the original pos.
		return pos, nil
	} else if candCell != lib.TYPE_WALL && candCell != lib.TYPE_FILLED {
		panic("implement non-wall/filled")
	}

	left := move(pos, DIR_LEFT)
	right := move(pos, DIR_RIGHT)

	leftCell := board.GetWithDefault(left, lib.TYPE_OPEN)
	rightCell := board.GetWithDefault(right, lib.TYPE_OPEN)

	if leftCell == lib.TYPE_OPEN && rightCell == lib.TYPE_OPEN {
		return right, []lib.Pos{pos}
	} else if leftCell != lib.TYPE_OPEN && rightCell == lib.TYPE_OPEN {
		return right, nil
	} else if leftCell == lib.TYPE_OPEN && rightCell != lib.TYPE_OPEN {
		return left, nil
	} else {
		// closed on both sides; no progress can be made, so return the
		// initial pos
		return pos, nil
	}
}

func pour(board *lib.Board, pos lib.Pos) []lib.Pos {
	allNewStarts := []lib.Pos{}

	if *verbose {
		logger.LogF("started at %v", pos)
	}

	for {
		if *verbose {
			logger.LogF("pos now %v", pos)
			board.DumpWithFocus(pos)
		}

		newPos, newStarts := progress(board, pos)
		allNewStarts = append(allNewStarts, newStarts...)

		if !board.InBounds(newPos) {
			logger.LogF("new pos out of bounds %v", newPos)
			break
		}

		if newPos.Eq(pos) {
			logger.LogF("pos stuck %v", pos)
			break
		}

		board.Set(newPos, lib.TYPE_FLOW)
		pos = newPos
	}

	if *verbose {
		logger.LogF("done at %v", pos)
		board.Dump()
	}

	return allNewStarts
}

func pourUntilDone(board *lib.Board) bool {
	madeProgress := false
	for {
		curCursor, ok := board.GetACursor()
		if !ok {
			break
		}

		madeProgress = true

		newCursors := pour(board, curCursor)
		board.DeleteCursor(curCursor)
		for _, cursor := range newCursors {
			board.AddCursor(cursor)
		}
	}

	return madeProgress
}

func findContainment(board *lib.Board, pos lib.Pos) (contained, seen []lib.Pos) {
	contained = []lib.Pos{pos}
	seen = []lib.Pos{pos}

	// logger.LogF("looking for containment at %v", pos)

	hasWall := func(dir Direction) bool {
		for next := pos; ; next = move(next, dir) {
			seen = append(seen, next)

			downType := board.GetWithDefault(move(next, DIR_DOWN), lib.TYPE_OPEN)
			if downType != lib.TYPE_WALL && downType != lib.TYPE_FILLED {
				return false
			}

			cell := board.GetWithDefault(next, lib.TYPE_OPEN)
			// logger.LogF("next is %v %s", next, cell)
			if cell == lib.TYPE_FLOW {
				// logger.LogF("is seen")
				contained = append(contained, next)
				continue
			} else {
				// logger.LogF("non-flow wall %v", cell == lib.TYPE_WALL)
				return cell == lib.TYPE_WALL
			}
		}
	}

	if !hasWall(DIR_LEFT) {
		// logger.LogF("no left containment")
		return nil, seen
	}
	if !hasWall(DIR_RIGHT) {
		// logger.LogF("no right containment")
		return nil, seen
	}

	logger.LogF("contained")
	return contained, seen
}

func fill(board *lib.Board) (bool, []lib.Pos) {
	allSeen := map[lib.Pos]bool{}
	madeProgress := false

	newStarts := []lib.Pos{}

	logger.LogF("filling")

	board.Visit(lib.TYPE_FLOW, func(pos lib.Pos) {
		if _, found := allSeen[pos]; found {
			return
		}

		// See if this flow is contained by walls
		contained, seen := findContainment(board, pos)
		for _, p := range seen {
			allSeen[p] = true
		}

		if contained == nil {
			return
		}

		madeProgress = true

		// It is contained, so turn these nodes to filled.
		for _, p := range contained {
			board.Set(p, lib.TYPE_FILLED)
		}

		// make new cursors
		for _, p := range contained {
			up := move(p, DIR_UP)
			if board.GetWithDefault(up, lib.TYPE_OPEN) == lib.TYPE_FLOW {
				newStarts = append(newStarts, up)
			}
		}
	})

	if *verbose {
		logger.LogF("after filling")
		board.Dump()
	}

	return madeProgress, newStarts
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	springPos, err := lib.PosFromString(*spring)
	if err != nil {
		log.Fatalf("failed to parse spring pos: %v", err)
	}

	lines, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	board := lib.NewBoard(springPos, lines)
	board.AddCursor(springPos)

	if *dump {
		xmin, xmax, ymin, ymax := board.Bounds()
		fmt.Printf("bounds: %v,%v %v,%v\n", xmin, ymin, xmax, ymax)
		board.Dump()
		os.Exit(0)
	}

	for {
		pourProgress := pourUntilDone(board)
		fillProgress, fillNewStarts := fill(board)

		for _, p := range fillNewStarts {
			board.AddCursor(p)
		}

		if !pourProgress && !fillProgress {
			break
		}
	}

	numFlow, numFilled := board.Score()
	fmt.Printf("%v+%v=%d\n", numFlow, numFilled, numFlow+numFilled)

	if *imageOut != "" {
		img := board.ToImage()

		f, err := os.Create(*imageOut)
		if err != nil {
			log.Fatal(err)
		}

		if err := png.Encode(f, img); err != nil {
			f.Close()
			log.Fatal(err)
		}

		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
