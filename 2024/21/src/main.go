package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime/pprof"
	"strings"

	"github.com/simmonmt/aoc/2024/common/collections"
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

type Node struct {
	To           rune
	Level        int
	Parent, Next *Node
	Paths        map[rune][]*Node
}

func (n *Node) String() string {
	yn := func(arg bool) rune {
		if arg {
			return 'y'
		}
		return 'n'
	}

	paths := []string{}
	for r, pl := range n.Paths {
		paths = append(paths, fmt.Sprintf("%s:%p", string(r), pl))
	}

	return fmt.Sprintf("{To:%c;L:%d;Par?%c;Nxt?%c;Paths:%s}",
		n.To, n.Level, yn(n.Parent != nil), yn(n.Next != nil),
		strings.Join(paths, ","))
}

func (n *Node) Nexts() string {
	out := ""
	for ; n != nil; n = n.Next {
		out += string(n.To)
	}
	return out
}

func makeTopNodes(line string) *Node {
	nodes := make([]*Node, len(line))
	for i, r := range line {
		n := &Node{To: r, Level: 1}
		nodes[i] = n
		if i > 0 {
			nodes[i-1].Next = n
		}
	}
	return nodes[0]
}

func makeNodePaths(from, to rune, parent *Node, pad Keypad) []*Node {
	paths := pad.PressFrom(from, to)
	pathList := make([]*Node, len(paths))
	for i, path := range paths {
		pathNodes := make([]*Node, len(path))
		for j := len(path) - 1; j >= 0; j-- {
			pathNodes[j] = &Node{To: rune(path[j]), Level: parent.Level + 1, Parent: parent}
			if j != len(path)-1 {
				pathNodes[j].Next = pathNodes[j+1]
			}
		}
		pathList[i] = pathNodes[0]
	}
	return pathList
}

type Cost struct {
	To   string
	Cost int
}

func (c *Cost) String() string {
	return fmt.Sprintf("/%v=%d", c.To, c.Cost)
}

type CacheKey struct {
	Level int
	From  string
	Nexts string
}

var (
	cacheHit  = 0
	cacheMiss = 0
)

func findMinCosts(start *Node, from string, stack []Keypad, cache map[CacheKey][]*Cost) []*Cost {
	log := false

	key := CacheKey{Level: start.Level, From: from, Nexts: start.Nexts()}
	if costs, found := cache[key]; found {
		cacheHit++
		return costs
	} else {
		cacheMiss++
	}

	if from == "" {
		panic("no from")
	}

	cookie := fmt.Sprintf("%sfindMinCosts L%d To %c from %s",
		strings.Repeat(" ", (start.Level-1)*2), start.Level, start.To, from)
	if log {
		fmt.Println(cookie)
	}

	curCosts := []*Cost{&Cost{To: from, Cost: 0}}

	for n := start; n != nil; n = n.Next {
		// The costs from start through n. Keyed by the stack of 'to'
		// values.
		nCosts := map[string]*Cost{}

		for _, preCost := range curCosts {
			// The beginning positions used when evaluating n. These
			// are the final positions from the previous n.
			nFrom := rune(preCost.To[0])
			subFrom := string(preCost.To[1:])

			if len(stack) == 0 {
				to := string(n.To)
				nCosts[to] = &Cost{To: to, Cost: preCost.Cost + 1}
				continue
			}

			pathList := makeNodePaths(nFrom, n.To, n, stack[0])

			for i, subN := range pathList {
				if log {
					fmt.Printf("%s: sub %d evaluating %s\n", cookie, i, subN.Nexts())
				}
				for _, subCost := range findMinCosts(subN, subFrom, stack[1:], cache) {
					to := string(n.To) + subFrom
					nCost := &Cost{To: to, Cost: preCost.Cost + subCost.Cost}
					if existing := nCosts[to]; existing == nil || existing.Cost > nCost.Cost {
						nCosts[to] = nCost
					}
				}
			}
		}

		curCosts = collections.MapValues(nCosts)
		if log {
			fmt.Printf("%s: costs %v\n", cookie, curCosts)
		}
	}

	if log {
		fmt.Println(cookie, ": done")
	}

	cache[key] = curCosts

	return curCosts
}

func solveALine(line string, stack []Keypad) int {
	start := makeTopNodes(line)
	cache := map[CacheKey][]*Cost{}
	costs := findMinCosts(start, strings.Repeat("A", len(stack)+1), stack, cache)

	minCost := -1
	for _, cost := range costs {
		if minCost == -1 || cost.Cost < minCost {
			minCost = cost.Cost
		}
	}

	num := numericCodePart(line)
	result := minCost * num
	logger.Infof("line %v min %v numeric %v result %v", line, minCost, num, result)
	return result
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

func solveA(lines []string) int {
	stack := []Keypad{
		NewNumPad(),
		NewDirPad(),
		NewDirPad(),
	}

	sum := 0
	for _, line := range lines {
		sum += solveALine(line, stack)
	}

	return sum
}

func solveB(lines []string) int {
	stack := []Keypad{NewNumPad()}
	for range 25 {
		stack = append(stack, NewDirPad())
	}

	sum := 0
	for _, line := range lines {
		sum += solveALine(line, stack)
	}

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
