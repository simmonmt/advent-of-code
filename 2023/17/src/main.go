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
	"image"
	"image/color"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2023/common/astar"
	"github.com/simmonmt/aoc/2023/common/dir"
	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/grid"
	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/pos"
)

var (
	verbose        = flag.Bool("verbose", false, "verbose")
	input          = flag.String("input", "", "input file")
	dumpAPath      = flag.String("dump_a_path", "", "If present, dump A solution to path")
	dumpBPath      = flag.String("dump_b_path", "", "If present, dump B solution to path")
	dumpGrid       = flag.Bool("dump_grid", false, "If true dump the grid")
	numAStarRounds = flag.Int("num_rounds", -1, "If >0 stop after this many rounds")
	cpuprofile     = flag.String("cpuprofile", "", "write cpu profile to file")
)

func parseInput(lines []string) (*grid.Grid[int], error) {
	return grid.NewFromLines(lines, func(p pos.P2, r rune) (int, error) {
		if r < '0' || r > '9' {
			return -1, fmt.Errorf("bad digit at %v", p)
		}
		return int(r - '0'), nil
	})
}

type Node struct {
	p           pos.P2
	d           dir.Dir
	numStraight int
}

func dumpSolver(solver *astar.AStar[*Node], height, width int, path string) {
	cb := func(node *Node, score, maxScore uint, img *image.NRGBA) (pos.P2, color.NRGBA) {
		val := uint8((float64(score) / float64(maxScore)) * 255.0)

		col := img.NRGBAAt(node.p.X, node.p.Y)
		if col.R != 0 {
			col.B = max(col.R, val)
		}
		col.R = max(col.R, val)

		return node.p, col
	}

	white := color.NRGBA{0, 0, 0, 255}

	if err := solver.Dump(path, height, width, white, cb); err != nil {
		logger.Errorf("failed write astar dump: %v", err)
	}
}

type BaseClient struct {
	g *grid.Grid[int]
}

func (c *BaseClient) EstimateDistance(start, end *Node) uint {
	return uint(start.p.ManhattanDistance(end.p))
}

func (c *BaseClient) NeighborDistance(n1, n2 *Node) uint {
	cost, _ := c.g.Get(n2.p)
	return uint(cost)
}

func (c *BaseClient) Serialize(n *Node) string {
	return fmt.Sprintf("%s/%s/%d", n.p.String(), n.d.String(), n.numStraight)
}

func (c *BaseClient) Deserialize(s string) (*Node, error) {
	parts := strings.SplitN(s, "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("bad split")
	}

	p, err := pos.P2FromString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("bad pos")
	}

	d := dir.Parse(parts[1])

	num, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("bad straight num")
	}

	return &Node{p, d, num}, nil
}

func solve(g *grid.Grid[int], client astar.ClientInterface[*Node], dumpPath string) int {
	// Starting with a -1 numStraight tells the client this is the beginning
	// of the path. The beginning is special because we haven't committed to
	// a direction yet (the DIR_EAST in the start node is a
	// lie). AllNeighbors therefore knows to try both of them in this
	// special case.
	start := &Node{pos.P2{X: 0, Y: 0}, dir.DIR_EAST, -1}
	goal := &Node{pos.P2{X: g.Width() - 1, Y: g.Height() - 1}, dir.DIR_EAST, 0}

	solver := astar.New(start, goal, client)
	if *numAStarRounds > 0 {
		solver.SetNumRounds(*numAStarRounds)
	}
	path := solver.Solve()

	if dumpPath != "" {
		dumpSolver(solver, g.Height(), g.Width(), dumpPath)
	}

	if *dumpGrid {
		posns := map[pos.P2]bool{}
		for _, node := range path {
			posns[node.p] = true
		}

		g.Dump(true, func(p pos.P2, v int, _ bool) string {
			if _, found := posns[p]; found {
				return "."
			}
			return strconv.Itoa(v)
		})
	}

	if path == nil {
		panic("no path found")
	}

	sum := 0
	for i := 0; i < len(path)-1; i++ {
		node := path[i]
		loss, _ := g.Get(node.p)
		sum += loss
	}
	return sum
}

type AClient struct {
	BaseClient
}

func (c *AClient) AllNeighbors(node *Node) []*Node {
	if node.numStraight == -1 {
		// This is the starting node. It can be facing east or down.
		return []*Node{
			&Node{dir.DIR_EAST.From(node.p), dir.DIR_EAST, 1},
			&Node{dir.DIR_SOUTH.From(node.p), dir.DIR_SOUTH, 1},
		}
	}

	out := []*Node{}
	if node.numStraight < 3 {
		nextPos := node.d.From(node.p)
		if c.g.IsValid(nextPos) {
			out = append(out, &Node{nextPos, node.d, node.numStraight + 1})
		}
	}

	leftDir, rightDir := node.d.Left(), node.d.Right()

	if left := leftDir.From(node.p); c.g.IsValid(left) {
		out = append(out, &Node{left, leftDir, 1})
	}
	if right := rightDir.From(node.p); c.g.IsValid(right) {
		out = append(out, &Node{right, rightDir, 1})
	}

	return out
}

func (c *AClient) GoalReached(cand, goal *Node) bool {
	return cand.p.Equals(goal.p)
}

func solveA(g *grid.Grid[int]) int {
	client := &AClient{BaseClient{g: g}}

	return solve(g, client, *dumpAPath)
}

type BClient struct {
	BaseClient
}

func (c *BClient) AllNeighbors(node *Node) []*Node {
	if node.numStraight == -1 {
		// This is the starting node. It can be facing east or down.
		return []*Node{
			&Node{dir.DIR_EAST.From(node.p), dir.DIR_EAST, 1},
			&Node{dir.DIR_SOUTH.From(node.p), dir.DIR_SOUTH, 1},
		}
	}

	out := []*Node{}
	if node.numStraight < 10 {
		nextPos := node.d.From(node.p)
		if c.g.IsValid(nextPos) {
			out = append(out, &Node{nextPos, node.d, node.numStraight + 1})
		}
	}

	if node.numStraight >= 4 {
		leftDir, rightDir := node.d.Left(), node.d.Right()

		if left := leftDir.From(node.p); c.g.IsValid(left) {
			out = append(out, &Node{left, leftDir, 1})
		}
		if right := rightDir.From(node.p); c.g.IsValid(right) {
			out = append(out, &Node{right, rightDir, 1})
		}
	}

	return out
}

func (c *BClient) GoalReached(cand, goal *Node) bool {
	return cand.p.Equals(goal.p) && cand.numStraight >= 4
}

// 1286 too high
// 1261 too low
func solveB(g *grid.Grid[int]) int {
	client := &BClient{BaseClient{g: g}}

	return solve(g, client, *dumpBPath)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			logger.Fatalf("can't make profile file: %v", err)
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
