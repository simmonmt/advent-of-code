package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime/pprof"

	"github.com/simmonmt/aoc/2024/common/dir"
	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/logger"
	"github.com/simmonmt/aoc/2024/common/pos"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

	inputPattern = regexp.MustCompile(`^[0-9]+A$`)
)

type Keypad interface {
	PressFrom(from, to rune) []string
	AllButtons() []rune
}

var (
	// 789
	// 456
	// 123
	//  0A
	numPadPosns = map[rune]pos.P2{
		'7': pos.P2{X: 0, Y: 0}, '8': pos.P2{X: 1, Y: 0}, '9': pos.P2{X: 2, Y: 0},
		'4': pos.P2{X: 0, Y: 1}, '5': pos.P2{X: 1, Y: 1}, '6': pos.P2{X: 2, Y: 1},
		'1': pos.P2{X: 0, Y: 2}, '2': pos.P2{X: 1, Y: 2}, '3': pos.P2{X: 2, Y: 2},
		'?': pos.P2{X: 0, Y: 3}, '0': pos.P2{X: 1, Y: 3}, 'A': pos.P2{X: 2, Y: 3},
	}
	numPadButtons = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A'}

	numAvoid = pos.P2{X: 0, Y: 3}

	//  ^A
	// <v>
	dirPadPosns = map[rune]pos.P2{
		'?': pos.P2{X: 0, Y: 0}, '^': pos.P2{X: 1, Y: 0}, 'A': pos.P2{X: 2, Y: 0},
		'<': pos.P2{X: 0, Y: 1}, 'v': pos.P2{X: 1, Y: 1}, '>': pos.P2{X: 2, Y: 1},
	}
	dirPadButtons = []rune{'<', '>', '^', 'v', 'A'}

	dirAvoid = pos.P2{X: 0, Y: 0}
)

func press(cur, dest, avoid pos.P2, soFar string) []string {
	dist := cur.ManhattanDistance(dest)
	if dist == 0 {
		return []string{soFar + "A"}
	}

	out := []string{}
	for _, d := range dir.AllDirs {
		n := d.From(cur)
		if n.Equals(avoid) {
			continue
		}
		nDist := n.ManhattanDistance(dest)
		if nDist > dist {
			continue
		}

		out = append(out, press(n, dest, avoid, soFar+string(d.Icon()))...)
	}
	return out
}

type NumPad struct{}

func NewNumPad() *NumPad {
	return &NumPad{}
}

func (np *NumPad) padPosn(label rune) pos.P2 {
	p, found := numPadPosns[label]
	if !found {
		panic("bad label " + string(label))
	}
	return p
}

func (np *NumPad) PressFrom(from, to rune) []string {
	return press(np.padPosn(from), np.padPosn(to), numAvoid, "")
}

func (np *NumPad) AllButtons() []rune {
	return numPadButtons
}

type DirPad struct{}

func NewDirPad() *DirPad {
	return &DirPad{}
}

func (dp *DirPad) padPosn(label rune) pos.P2 {
	p, found := dirPadPosns[label]
	if !found {
		panic("bad label " + string(label))
	}
	return p
}

func (dp *DirPad) PressFrom(from, to rune) []string {
	return press(dp.padPosn(from), dp.padPosn(to), dirAvoid, "")
}

func (dp *DirPad) AllButtons() []rune {
	return dirPadButtons
}

func parseInput(lines []string) ([]string, error) {
	for i, line := range lines {
		if matched := inputPattern.MatchString(line); !matched {
			return nil, fmt.Errorf("%d: bad input", i+1)
		}
	}

	return lines, nil
}

// NP1   DP1   DP2  DP3
// 789    ^A    ^A   ^A
// 456   <v>   <v>  <v>
// 123
// _0A
//
// Last entered on NP1 was 3
// Want to enter 5 on NP1
// R1 (operating NP1) is
//  - controlled by DP1
//  - positioned over 3 on NP1
// R2 (operating DP1) is
//  - controlled by DP2
//  - positioned over v on DP1
// R3 (operating DP2) is
//  - controlled by DP3
//  - posiitoved over < on DP2
// Human (operating DP3)
//
// To enter 5 on NP1
//   - R1 has to move from 3 to 5, as directed by DP1/R2
//   - R2 has to execute commands as directed by DP2/R3
//   - R3 has to execute commands as directed by DP3/Human
//
// DP1: R1 was on 3/NP1, so R2 send <^A or ^<A             to activate 5/NP1
// DP2: R2 was on v/DP1, so R3 send <A>^A>A                to enter <^A/DP1
//                          or send ^Av<A>>^A or ^Av<A>^>A to enter ^<A/DP1
// DP3: R3 was on </DP2, so H  send
//
//
// node
//   from, to rune
//   pad
//   parent
//   next
//   paths []*node
//

type Node struct {
	To           rune
	Level        int
	Parent, Next *Node
	Paths        map[rune][]*Node
}

func buildPathNodes(line string, level int, parent *Node, stack []Keypad) *Node {
	var start, last *Node
	var pad Keypad
	if len(stack) > 0 {
		pad = stack[0]
	}
	cur := '?'

	for _, r := range line {
		node := &Node{To: r, Level: level, Parent: parent}

		if start == nil {
			start = node
			last = node
		} else {
			last.Next = node
			last = node
		}

		if pad != nil {
			starts := []rune{cur}
			if cur == '?' {
				starts = pad.AllButtons()
			}

			node.Paths = map[rune][]*Node{}
			for _, start := range starts {
				paths := pad.PressFrom(start, node.To)
				pathNodes := make([]*Node, len(paths))
				for i := range paths {
					pathNodes[i] = buildPathNodes(paths[i], level+1, node, stack[1:])
				}
				node.Paths[start] = pathNodes
			}
		}

		cur = r
	}

	return start
}

func walkNodesInOrder(start *Node) {
	// begin with start, cur = A

}

func solveALine(line string, stack []Keypad) int {
	// start := buildPathNodes(line, 1, nil, stack)
	// fmt.Println(start)
	return 0
}

func solveA(lines []string) int64 {
	stack := []Keypad{
		NewNumPad(),
		NewDirPad(),
		NewDirPad(),
	}

	sum := 0
	for _, line := range lines {
		sum += solveALine(line, stack)
	}

	return int64(sum)
}

func numericCodePart(line string) int {
	num := 0
	for _, r := range line {
		if r >= '0' && r <= '9' {
			num = num*10 + int(r-'0')
		}
	}
	return num
}

func solveB(input []string) int64 {
	return -1
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
