// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/simmonmt/aoc/2023/common/dir"
	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/grid"
	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	commandPattern = regexp.MustCompile(`^([UDRL]) (\d+) \(#([0-9a-f]+)\)$`)
)

type Node struct {
	Wall  bool
	Color uint32
}

type Command struct {
	Dir   rune
	Num   int
	Color uint32
}

func parseCommand(in string) (Command, error) {
	parts := commandPattern.FindStringSubmatch(in)
	if parts == nil {
		return Command{}, fmt.Errorf("regexp mismatch")
	}

	dir := rune(parts[1][0])

	num, err := strconv.Atoi(parts[2])
	if err != nil {
		return Command{}, fmt.Errorf("bad num: %v", err)
	}

	color, err := strconv.ParseUint(parts[3], 16, 32)
	if err != nil {
		return Command{}, fmt.Errorf("bad color: %v", err)
	}

	return Command{
		Dir:   dir,
		Num:   num,
		Color: uint32(color),
	}, nil
}

func parseInput(lines []string) ([]Command, error) {
	out := []Command{}
	for i, line := range lines {
		command, err := parseCommand(line)
		if err != nil {
			return nil, fmt.Errorf("%d: bad command: %v", i+1, err)
		}
		out = append(out, command)
	}
	return out, nil
}

func runeToDir(r rune) dir.Dir {
	switch r {
	case 'U':
		return dir.DIR_NORTH
	case 'D':
		return dir.DIR_SOUTH
	case 'L':
		return dir.DIR_WEST
	case 'R':
		return dir.DIR_EAST
	default:
		panic("unknown")
	}
}

func isWall(g *grid.SparseGrid[*Node], p pos.P2) bool {
	node, found := g.Get(p)
	if !found {
		return false
	}

	return node.Wall
}

func fill(g *grid.SparseGrid[*Node]) int {
	start, end := g.Start(), g.End()
	filled := 0

	for y := start.Y; y <= end.Y; y++ {
		out := true
		lastUp := false

		for x := start.X - 1; x <= end.X; x++ {
			p := pos.P2{X: x, Y: y}

			if isWall(g, p) {
				up := isWall(g, dir.DIR_NORTH.From(p))
				down := isWall(g, dir.DIR_SOUTH.From(p))
				next := isWall(g, dir.DIR_EAST.From(p))

				if up && down {
					out = !out
				} else if up || down {
					if !next {
						if up != lastUp {
							out = !out
						}
					}
					lastUp = up
				}
			} else if !out {
				g.Set(p, &Node{Wall: false})
				filled++
			}
		}
	}

	return filled
}

func oldSolveA(commands []Command) int {
	g := grid.NewSparseGrid[*Node]()

	walls := 0
	p := pos.P2{X: 0, Y: 0}
	for _, cmd := range commands {
		for i := 0; i < cmd.Num; i++ {
			p = runeToDir(cmd.Dir).From(p)
			g.Set(p, &Node{Wall: true, Color: cmd.Color})
			walls++
		}
	}

	filled := fill(g)

	if *verbose {
		g.Dump(false, func(_ pos.P2, node *Node, found bool) string {
			if found {
				if node.Wall {
					return "#"
				} else {
					return "."
				}
			} else {
				return " "
			}
		})
	}

	return walls + filled
}

type LocatedCommand struct {
	Dir        rune
	Num        int
	Start, End pos.P2
}

func (lc *LocatedCommand) String() string {
	return fmt.Sprintf("%c/%d/%s/%s", lc.Dir, lc.Num, lc.Start, lc.End)
}

func makeLocatedCommands(commands []Command) []*LocatedCommand {
	locatedCommands := []*LocatedCommand{}
	p := pos.P2{}
	for _, cmd := range commands {
		off := runeToDir(cmd.Dir).From(pos.P2{})
		off.X *= cmd.Num
		off.Y *= cmd.Num
		end := p
		end.Add(off)

		lc := &LocatedCommand{
			Dir:   cmd.Dir,
			Num:   cmd.Num,
			Start: p,
			End:   end,
		}

		locatedCommands = append(locatedCommands, lc)
		p = end
	}
	return locatedCommands
}

func drawAndDump(lcs []*LocatedCommand) {
	g := grid.NewSparseGrid[*Node]()

	p := pos.P2{X: 0, Y: 0}
	for _, cmd := range lcs {
		for i := 0; i < cmd.Num; i++ {
			p = runeToDir(cmd.Dir).From(p)
			g.Set(p, &Node{Wall: true})
		}
	}

	g.Dump(false, func(_ pos.P2, node *Node, found bool) string {
		if found {
			return "#"
		} else {
			return ","
		}
	})
}

type Range struct {
	From, To pos.P2
}

func (r *Range) String() string {
	return fmt.Sprintf("%s-%s", r.From, r.To)
}

func findIntersections(y int, lcs []*LocatedCommand) (current, future []*Range) {
	current, future = []*Range{}, []*Range{}

	for _, lc := range lcs {
		a, b := lc.Start, lc.End
		if b.LessThan(a) {
			a, b = b, a
		}

		if a.Y > y {
			future = append(future, &Range{a, b})
		} else { // start <= y
			if b.Y >= y {
				current = append(current, &Range{a, b})
			}
		}
	}

	return
}

func rowSize(y int, lcs []*LocatedCommand) int {
	// Get list of ranges that intersect with current Y, sorted by From.X, with horizontal pieces first
	ranges, _ := findIntersections(y, lcs)
	sort.Slice(ranges, func(i, j int) bool {
		ri, rj := ranges[i], ranges[j]

		if ri.From.X == rj.From.X {
			return ri.From.X != ri.To.X
		} else {
			return ri.From.X < rj.From.X
		}
	})

	//logger.Infof("y=%d ranges %v", y, ranges)

	isUp := func(r *Range) bool { return r.From.Y < y || r.To.Y < y }
	isDown := func(r *Range) bool { return r.From.Y > y || r.To.Y > y }

	size := 0
	out := true
	ri, x := 0, ranges[0].From.X
	for ri < len(ranges) {
		startUp := ri+1 < len(ranges) && isUp(ranges[ri+1]) && ranges[ri].From.X == ranges[ri+1].From.X
		startDown := ri+1 < len(ranges) && isDown(ranges[ri+1]) && ranges[ri].From.X == ranges[ri+1].From.X

		//logger.Infof("y=%d x=%d startUp %v startDown %v", y, x, startUp, startDown)

		if !startUp && !startDown {
			out = !out
			if ri+1 < len(ranges) {
				nextX := ranges[ri+1].From.X

				//logger.Infof("y=%d | out %v nextX %v", y, out, nextX)

				if !out {
					size += nextX - x
				} else {
					size++
				}

				x = nextX
			} else {
				size++ // the wall we just crossed
			}
			ri++

		} else if startUp != startDown {
			endX := ranges[ri].To.X
			toAdd := endX - x + 1
			//logger.Infof("y=%d x=%d adding %d", y, x, toAdd)
			size += toAdd

			nextUp := ri+2 < len(ranges) && isUp(ranges[ri+2])
			if nextUp != startUp {
				out = !out
			}

			//logger.Infof("y=%d x=%d L out %v startUp %v startDown %v nextup %v",
			//y, x, out, startUp, startDown, nextUp)

			if ri+3 < len(ranges) {
				nextX := ranges[ri+3].From.X
				if !out {
					suffix := nextX - endX - 1
					//logger.Infof("y=%d endX=%d adding suffix %d", y, endX, suffix)
					size += suffix
				}
				x = nextX
			}
			ri += 3
		} else {
			panic("unexpected")
		}
	}

	return size
}

func rowRepeats(y int, lcs []*LocatedCommand) int {
	current, future := findIntersections(y, lcs)
	//logger.Infof("y=%d current %v", y, current)
	//logger.Infof("y=%d future %v", y, future)

	minY := current[0].To.Y
	for _, r := range current {
		minY = min(minY, r.To.Y)
	}
	for _, r := range future {
		minY = min(minY, r.From.Y)
	}

	return max(1, minY-y)
}

func solve(lcs []*LocatedCommand) int64 {
	minY, maxY := lcs[0].Start.Y, lcs[0].Start.Y
	for _, lc := range lcs {
		minY = min(minY, lc.Start.Y, lc.End.Y)
		maxY = max(maxY, lc.Start.Y, lc.End.Y)
	}

	sum := int64(0)
	for y := minY; y <= maxY; {
		size := rowSize(y, lcs)
		repeats := rowRepeats(y, lcs)
		sum += int64(size) * int64(repeats)
		y += repeats
	}

	return sum
}

func solveA(commands []Command) int64 {
	lcs := makeLocatedCommands(commands)
	logger.Infof("=== start %v", lcs)

	return solve(lcs)
}

func solveB(commands []Command) int64 {
	newCommands := []Command{}
	for _, cmd := range commands {
		num := cmd.Color >> 4
		d := func() rune {
			switch cmd.Color & 7 {
			case 0:
				return 'R'
			case 1:
				return 'D'
			case 2:
				return 'L'
			case 3:
				return 'U'
			default:
				panic("bad dir")
			}
		}()

		newCommands = append(newCommands, Command{Dir: d, Num: int(num)})
	}

	lcs := makeLocatedCommands(newCommands)
	logger.Infof("=== start %v", lcs)

	return solve(lcs)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

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
