package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/logger"
	"github.com/simmonmt/aoc/2020/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type State uint8

const (
	STATE_FLOOR State = iota
	STATE_EMPTY
	STATE_FULL
)

func (s State) String() string {
	switch s {
	case STATE_FLOOR:
		return "."
	case STATE_EMPTY:
		return "L"
	case STATE_FULL:
		return "#"
	default:
		return "?"
	}
}

type Board struct {
	w, h  int
	cells []State
}

func newBoard(lines []string) *Board {
	w := len(lines[0])
	h := len(lines)

	b := &Board{
		w:     w,
		h:     h,
		cells: make([]State, w*h),
	}
	for y, line := range lines {
		for x, r := range line {
			p := pos.P2{X: x, Y: y}

			var state State
			switch r {
			case 'L':
				state = STATE_EMPTY
			case '.':
				state = STATE_FLOOR
			case '#':
				state = STATE_FULL
			default:
				panic("bad char")
			}

			b.Set(p, state)
		}
	}
	return b
}

func (b *Board) Width() int {
	return b.w
}

func (b *Board) Height() int {
	return b.h
}

func (b *Board) Set(p pos.P2, state State) {
	off := p.Y*b.w + p.X
	b.cells[off] = state
}

func (b *Board) Get(p pos.P2) State {
	off := p.Y*b.w + p.X
	return b.cells[off]
}

func (b *Board) NumOfState(state State) int {
	num := 0
	for _, s := range b.cells {
		if s == state {
			num++
		}
	}
	return num
}

func (b *Board) Walk(cb func(p pos.P2, state State) int) int {
	sum := 0
	for y := 0; y < b.Height(); y++ {
		for x := 0; x < b.Width(); x++ {
			p := pos.P2{X: x, Y: y}
			sum += cb(p, b.Get(p))
		}
	}
	return sum
}

func (b *Board) Dump() {
	for y := 0; y < b.h; y++ {
		for x := 0; x < b.w; x++ {
			fmt.Print(b.Get(pos.P2{X: x, Y: y}))
		}
		fmt.Println()
	}
}

func evolve(src, dest *Board, fullToEmptyMin int, neighborCounter func(*Board, pos.P2) int) int {
	return src.Walk(func(p pos.P2, state State) int {
		changed := 0
		neighbors := neighborCounter(src, p)

		if state == STATE_EMPTY && neighbors == 0 {
			state = STATE_FULL
			changed = 1
		} else if state == STATE_FULL && neighbors >= fullToEmptyMin {
			state = STATE_EMPTY
			changed = 1
		}

		dest.Set(p, state)
		return changed
	})
}

func playGame(lines []string, fullToEmptyMin int, neighborCounter func(*Board, pos.P2) int) int {
	boards := []*Board{newBoard(lines), newBoard(lines)}
	cur := 0

	var step int
	for step = 1; ; step++ {
		next := (cur + 1) % 2
		if logger.Enabled() {
			fmt.Printf("\nstep %d:\n", step)
			boards[cur].Dump()
		}

		numChanged := evolve(boards[cur], boards[next], fullToEmptyMin,
			neighborCounter)
		cur = next
		if numChanged == 0 {
			logger.LogF("no change; breaking")
			break
		}
	}

	step-- // subtract the no-change step
	logger.LogF("%d steps\n", step)
	return boards[cur].NumOfState(STATE_FULL)
}

func neighborCounterA(b *Board, p pos.P2) int {
	num := 0
	for _, n := range p.AllNeighbors(true) {
		if n.X < 0 || n.Y < 0 || n.X >= b.Width() || n.Y >= b.Height() {
			continue
		}
		if b.Get(n) == STATE_FULL {
			num++
		}
	}
	return num
}

var (
	deltas = []pos.P2{
		pos.P2{-1, -1},
		pos.P2{0, -1},
		pos.P2{1, -1},
		pos.P2{-1, 0},
		pos.P2{1, 0},
		pos.P2{-1, 1},
		pos.P2{0, 1},
		pos.P2{1, 1},
	}
)

func neighborCounterB(b *Board, p pos.P2) int {
	numOccupied := 0
	for _, delta := range deltas {
		p2 := p
		for {
			p2.Add(delta)
			if p2.X < 0 || p2.Y < 0 || p2.X >= b.Width() || p2.Y >= b.Height() {
				logger.LogF("ran off from %v in %v at %v",
					p, delta, p2)
				break // ran off the board
			}

			p2State := b.Get(p2)
			if p2State == STATE_FLOOR {
				continue
			}

			if p2State == STATE_FULL {
				logger.LogF("found from %v in %v at %v",
					p, delta, p2)
				numOccupied++
			}
			break
		}
	}
	return numOccupied
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		log.Fatal(err)
	}

	numOccupiedA := playGame(lines, 4, neighborCounterA)
	fmt.Printf("A: %d\n", numOccupiedA)

	numOccupiedB := playGame(lines, 5, neighborCounterB)
	fmt.Printf("B: %d\n", numOccupiedB)
}
