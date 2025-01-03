package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/simmonmt/aoc/2024/common/dir"
	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/grid"
	"github.com/simmonmt/aoc/2024/common/lineio"
	"github.com/simmonmt/aoc/2024/common/logger"
	"github.com/simmonmt/aoc/2024/common/pos"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

type Input struct {
	G     *grid.Grid[rune]
	Start pos.P2
	Cmds  []dir.Dir
}

func parseGrid(lines []string) (g *grid.Grid[rune], start pos.P2, err error) {
	foundStart := false
	g, err = grid.NewFromLines[rune](lines, func(p pos.P2, r rune) (rune, error) {
		if r == '@' {
			foundStart = true
			start = p
			return '.', nil
		}
		return r, nil
	})

	if !foundStart {
		return nil, pos.P2{}, fmt.Errorf("no start found")
	}
	return
}

func parseInput(lines []string) (*Input, error) {
	groups := lineio.BlankSeparatedGroups(lines)
	if len(groups) != 2 {
		return nil, fmt.Errorf("bad number of groups: %d", len(groups))
	}

	g, start, err := parseGrid(groups[0])
	if err != nil {
		return nil, fmt.Errorf("bad grid: %v", err)
	}

	cmds := []dir.Dir{}
	for _, line := range groups[1] {
		for _, r := range line {
			cmd, ok := dir.ParseIcon(r)
			if !ok {
				return nil, fmt.Errorf("bad dir %v", string(r))
			}
			cmds = append(cmds, cmd)
		}
	}

	return &Input{G: g, Start: start, Cmds: cmds}, nil
}

func isBlocked(g *grid.Grid[rune], cur pos.P2, cmd dir.Dir) (blocked bool, toPush int) {
	for cur = cmd.From(cur); ; cur = cmd.From(cur) {
		r, _ := g.Get(cur)
		if r == '#' {
			blocked = true
			return
		} else if r == 'O' {
			toPush++
		} else if r == '.' {
			blocked = false
			return
		}
	}
}

func doMoveA(g *grid.Grid[rune], cur pos.P2, cmd dir.Dir) pos.P2 {
	blocked, toPush := isBlocked(g, cur, cmd)
	if blocked {
		return cur
	}

	next := cmd.From(cur)
	if toPush != 0 {
		g.Set(next, '.')

		end := next
		for i := 0; i < toPush; i++ {
			end = cmd.From(end)
		}

		if r, _ := g.Get(end); r != '.' {
			panic("not end")
		}
		g.Set(end, 'O')
	}

	return next
}

func dumpGrid(g *grid.Grid[rune], cur pos.P2) {
	g.Dump(true, func(p pos.P2, r rune, _ bool) string {
		if p.Equals(cur) {
			return "@"
		}
		return string(r)
	})
}

func checkGrid(g *grid.Grid[rune]) {
	g.Walk(func(p pos.P2, r rune) {
		if r == '[' {
			nr := get(g, pos.P2{X: p.X + 1, Y: p.Y})
			if nr != ']' {
				panic(fmt.Sprintf("unclosed box at %v", p))
			}
		} else if r == ']' {
			nr := get(g, pos.P2{X: p.X - 1, Y: p.Y})
			if nr != '[' {
				panic(fmt.Sprintf("unopened box at %v", p))
			}
		}
	})
}

func solveA(input *Input) int {
	cur := input.Start
	g := input.G.Clone()
	for _, cmd := range input.Cmds {
		cur = doMoveA(g, cur, cmd)
	}

	sum := 0
	g.Walk(func(p pos.P2, r rune) {
		if r == 'O' {
			sum += p.X + 100*p.Y
		}
	})

	return sum
}

func get(g *grid.Grid[rune], p pos.P2) rune {
	r, _ := g.Get(p)
	return r
}

func getBox(boxes map[pos.P2]*Box, p pos.P2) *Box {
	b := boxes[p]
	if b == nil {
		panic("no box found")
	}
	return b
}

type Box struct {
	Left pos.P2
}

func (b *Box) right() pos.P2 {
	return pos.P2{X: b.Left.X + 1, Y: b.Left.Y}
}

func (b *Box) Blocked(g *grid.Grid[rune], boxes map[pos.P2]*Box, cmd dir.Dir) (blocked bool, toMove []*Box) {
	if cmd == dir.DIR_EAST {
		return b.blockedE(g, boxes)
	} else if cmd == dir.DIR_WEST {
		return b.blockedW(g, boxes)
	} else {
		return b.blockedNS(g, boxes, cmd)
	}
}

func (b *Box) blockedE(g *grid.Grid[rune], boxes map[pos.P2]*Box) (blocked bool, toMove []*Box) {
	next := pos.P2{X: b.Left.X + 2, Y: b.Left.Y}

	r := get(g, next)
	if r == '#' {
		return true, nil
	} else if r == '.' {
		return false, nil
	}

	nb := getBox(boxes, next)
	toMove = []*Box{nb}
	if nextBlocked, nextToMove := nb.blockedE(g, boxes); nextBlocked {
		return true, nil
	} else {
		toMove = append(toMove, nextToMove...)
		return false, toMove
	}
}

func (b *Box) blockedW(g *grid.Grid[rune], boxes map[pos.P2]*Box) (blocked bool, toMove []*Box) {
	next := pos.P2{X: b.Left.X - 1, Y: b.Left.Y}

	r := get(g, next)
	if r == '#' {
		return true, nil
	} else if r == '.' {
		return false, nil
	}

	nb := getBox(boxes, pos.P2{X: next.X - 1, Y: next.Y})
	toMove = []*Box{nb}
	if nextBlocked, nextToMove := nb.blockedW(g, boxes); nextBlocked {
		return true, nil
	} else {
		toMove = append(toMove, nextToMove...)
		return false, toMove
	}
}

func (b *Box) blockedNS(g *grid.Grid[rune], boxes map[pos.P2]*Box, cmd dir.Dir) (blocked bool, toMove []*Box) {
	lNext, rNext := cmd.From(b.Left), cmd.From(b.right())
	lNextRune, rNextRune := get(g, lNext), get(g, rNext)

	if lNextRune == '#' || rNextRune == '#' {
		return true, nil
	}
	if lNextRune == '.' && rNextRune == '.' {
		return false, nil
	}

	// Blocked by something
	toMove = []*Box{}
	if lNextRune == ']' {
		toMove = append(toMove, getBox(boxes, pos.P2{X: lNext.X - 1, Y: lNext.Y}))
	} else if lNextRune == '[' {
		toMove = append(toMove, getBox(boxes, lNext))
	}
	if rNextRune == '[' {
		toMove = append(toMove, getBox(boxes, rNext))
	}

	added := []*Box{}
	for _, nb := range toMove {
		nextBlocked, nextToMove := nb.blockedNS(g, boxes, cmd)
		if nextBlocked {
			return true, nil
		}
		added = append(added, nextToMove...)
	}

	for _, ab := range added {
		found := false
		for _, mb := range toMove {
			if ab.Left.Equals(mb.Left) {
				found = true
				break
			}
		}
		if !found {
			toMove = append(toMove, ab)
		}
	}

	return false, toMove
}

func (b *Box) RemoveFromGrid(g *grid.Grid[rune]) {
	g.Set(b.Left, '.')
	g.Set(b.right(), '.')
}

func (b *Box) AddToGrid(g *grid.Grid[rune]) {
	g.Set(b.Left, '[')
	g.Set(b.right(), ']')
}

func findBoxes(g *grid.Grid[rune]) map[pos.P2]*Box {
	out := map[pos.P2]*Box{}
	g.Walk(func(p pos.P2, r rune) {
		if r == '[' {
			out[p] = &Box{Left: p}
		}
	})
	return out
}

func doMoveB(g *grid.Grid[rune], boxes map[pos.P2]*Box, cur pos.P2, cmd dir.Dir) pos.P2 {
	next := cmd.From(cur)
	var blockingBox *Box
	if r := get(g, next); r == '#' {
		return cur
	} else if r == '[' || r == ']' {
		bp := next
		if r == ']' {
			bp.X--
		}

		var found bool
		blockingBox, found = boxes[bp]
		if !found {
			panic("no box")
		}
	}

	if blockingBox != nil {
		blocked, toMove := blockingBox.Blocked(g, boxes, cmd)
		if blocked {
			return cur
		}

		toMove = append(toMove, blockingBox)

		for _, mb := range toMove {
			delete(boxes, mb.Left)
			mb.RemoveFromGrid(g)
		}

		for _, mb := range toMove {
			mb.Left = cmd.From(mb.Left)
			mb.AddToGrid(g)
			boxes[mb.Left] = mb
		}
	}

	// We should only reach here if nothing was there or if boxes got moved
	// out of the way. Verify that.
	if r := get(g, next); r != '.' {
		panic("not empty")
	}

	return next
}

func solveB(input *Input) int {
	g := grid.New[rune](input.G.Width()*2, input.G.Height())
	input.G.Walk(func(p pos.P2, r rune) {
		var out []rune
		switch r {
		case '#':
			out = []rune{'#', '#'}
		case '.':
			out = []rune{'.', '.'}
		case 'O':
			out = []rune{'[', ']'}
		default:
			panic("unexpected rune " + string(r))
		}

		g.Set(pos.P2{X: p.X * 2, Y: p.Y}, out[0])
		g.Set(pos.P2{X: p.X*2 + 1, Y: p.Y}, out[1])
	})

	boxes := findBoxes(g)
	cur := pos.P2{X: input.Start.X * 2, Y: input.Start.Y}

	for _, cmd := range input.Cmds {
		cur = doMoveB(g, boxes, cur, cmd)

		// fmt.Println(cur, cmd)
		// dumpGrid(g, cur)
		// fmt.Println()
		checkGrid(g)
	}

	//dumpGrid(g, cur)

	sum := 0
	g.Walk(func(p pos.P2, r rune) {
		if r == '[' {
			sum += p.X + 100*p.Y
		}
	})

	return sum
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *input == "" {
		logger.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		logger.Fatalf("failed to read input: %v", err)
	}

	input, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
