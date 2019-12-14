package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/simmonmt/aoc/2019/common/intmath"
	"github.com/simmonmt/aoc/2019/common/logger"
	"github.com/simmonmt/aoc/2019/common/pos"
	"github.com/simmonmt/aoc/2019/common/vm"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	ramPath    = flag.String("ram", "", "path to file containing ram values")
	startWhite = flag.Bool("start_white", false, "if true, make the starting panel white")
)

type Dir int

const (
	DIR_UP Dir = iota
	DIR_LEFT
	DIR_DOWN
	DIR_RIGHT
	DIR_MAX
)

func (d Dir) String() string {
	switch d {
	case DIR_UP:
		return "up"
	case DIR_LEFT:
		return "left"
	case DIR_DOWN:
		return "down"
	case DIR_RIGHT:
		return "right"
	default:
		return fmt.Sprintf("bad dir %v", d)
	}
}

func (d Dir) Short() string {
	switch d {
	case DIR_UP:
		return "^"
	case DIR_LEFT:
		return "<"
	case DIR_DOWN:
		return "v"
	case DIR_RIGHT:
		return ">"
	default:
		return "*"
	}
}

func executeTurn(curDir Dir, p pos.P2, turnRight bool) (newDir Dir, newP pos.P2) {
	if turnRight {
		newDir = (curDir + DIR_MAX - 1) % DIR_MAX
	} else {
		newDir = (curDir + 1) % DIR_MAX
	}

	switch newDir {
	case DIR_UP:
		return newDir, pos.P2{p.X, p.Y - 1}
	case DIR_DOWN:
		return newDir, pos.P2{p.X, p.Y + 1}
	case DIR_LEFT:
		return newDir, pos.P2{p.X - 1, p.Y}
	case DIR_RIGHT:
		return newDir, pos.P2{p.X + 1, p.Y}
	default:
		panic("bad dir")
	}
}

func nextCommand(async *vm.Async, isWhite bool) (writeWhite, turnRight, keepGoing bool) {
	sendVal := int64(0)
	if isWhite {
		sendVal = 1
	}
	logger.LogF("sending %v", sendVal)
	async.In <- &vm.ChanIOMessage{Val: sendVal}

	readMsg := func() (v, ok bool) {
		msg, ok := <-async.Out
		if !ok {
			return false, false
		}
		if msg.Err != nil {
			panic(fmt.Sprintf("error from vm: %v", msg.Err))
		}
		return msg.Val == 1, true
	}

	writeWhite, ok := readMsg()
	if !ok {
		return false, false, false
	}

	turnRight, ok = readMsg()
	if !ok {
		return false, false, false
	}

	keepGoing = true
	return
}

type Board map[pos.P2]bool

func dumpBoard(b Board, curPos pos.P2, curDir Dir) {
	minX, maxX, minY, maxY := curPos.X, curPos.X, curPos.Y, curPos.Y
	for p := range b {
		minX = intmath.IntMin(minX, p.X)
		minY = intmath.IntMin(minY, p.Y)
		maxX = intmath.IntMax(maxX, p.X)
		maxY = intmath.IntMax(maxY, p.Y)
	}

	// Pad by 2 for visual appeal
	minX, minY = minX-2, minY-2
	maxX, maxY = maxX+2, maxY+2

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			p := pos.P2{x, y}
			if p.Equals(curPos) {
				fmt.Print(curDir.Short())
			} else if b[pos.P2{x, y}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *ramPath == "" {
		log.Fatalf("--ram is required")
	}

	ram, err := vm.NewRamFromFile(*ramPath)
	if err != nil {
		log.Fatal(err)
	}

	board := Board{}
	curPos := pos.P2{0, 0}
	dir := DIR_UP

	if *startWhite {
		board[curPos] = true
	}

	async := vm.RunAsync("vm", ram)
	dumpBoard(board, curPos, dir)
	var turnNo int
	for turnNo = 1; ; turnNo++ {
		//fmt.Printf("turn %d start: p %+v, d %s\n", turnNo, curPos, dir)

		writeWhite, turnRight, keepGoing := nextCommand(async, board[curPos])
		if !keepGoing {
			break
		}

		//fmt.Printf("got write=%v turnright=%v\n", writeWhite, turnRight)

		board[curPos] = writeWhite

		dir, curPos = executeTurn(dir, curPos, turnRight)

		//dumpBoard(board, curPos, dir)
	}

	fmt.Printf("program terminated after %d turns\n", turnNo)
	fmt.Printf("number of cells painted at least once: %v\n", len(board))

	dumpBoard(board, curPos, dir)
}
