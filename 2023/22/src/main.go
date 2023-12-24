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
	"sort"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/grid"
	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/mtsmath"
	"github.com/simmonmt/aoc/2023/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Range struct {
	From, To pos.P3
}

type Block struct {
	Num int
	R   Range
}

func parseInput(lines []string) ([]*Block, error) {
	out := []*Block{}
	for i, line := range lines {
		a, b, ok := strings.Cut(line, "~")
		if !ok {
			return nil, fmt.Errorf("%d: bad split", i+1)
		}

		from, err := pos.P3FromString(a)
		if err != nil {
			return nil, fmt.Errorf("%d: bad from: %v", i+1, err)
		}

		to, err := pos.P3FromString(b)
		if err != nil {
			return nil, fmt.Errorf("%d: bad to: %v", i+1, err)
		}

		if to.LessThan(from) {
			from, to = to, from
		}

		out = append(out, &Block{Num: i + 1, R: Range{from, to}})
	}

	return out, nil
}

func calcIncr(from, to pos.P2) pos.P2 {
	incr := pos.P2{X: to.X - from.X, Y: to.Y - from.Y}
	incr.X /= max(1, mtsmath.Abs(incr.X))
	incr.Y /= max(1, mtsmath.Abs(incr.Y))
	return incr
}

func walkRange(from, to pos.P2, cb func(pos.P2) bool) {
	incr := calcIncr(from, to)
	for p := from; ; p.Add(incr) {
		if !cb(p) {
			return
		}
		if p.Equals(to) {
			return
		}
	}
}

func walkXYRange(r Range, cb func(pos.P2) bool) {
	from, to := xyRange(r)
	walkRange(from, to, cb)
}

func xyRange(r Range) (pos.P2, pos.P2) {
	return pos.P2{X: r.From.X, Y: r.From.Y},
		pos.P2{X: r.To.X, Y: r.To.Y}
}

type Highests struct {
	a [][]int
}

func NewHighests(x, y int) *Highests {
	a := make([][]int, y)
	for i := range a {
		a[i] = make([]int, x)
	}

	return &Highests{a: a}
}

func (h *Highests) Get(p pos.P2) int {
	return h.a[p.Y][p.X]
}

func (h *Highests) GetRange(from, to pos.P2) int {
	highest := 0
	walkRange(from, to, func(p pos.P2) bool {
		highest = max(highest, h.Get(p))
		return true
	})
	return highest
}

func (h *Highests) Set(p pos.P2, z int) {
	h.a[p.Y][p.X] = z
}

func (h *Highests) SetRange(from, to pos.P2, val int) {
	walkRange(from, to, func(p pos.P2) bool {
		if cur := h.Get(p); cur >= val {
			panic(fmt.Sprintf("%v already at %v, trying to set %v", p, cur, val))
		}
		h.Set(p, val)
		return true
	})
}

func drop(blocks []*Block) int {
	maxX, maxY := blocks[0].R.To.X, blocks[0].R.To.Y
	for _, block := range blocks {
		maxX = max(maxX, block.R.From.X, block.R.To.X)
		maxY = max(maxY, block.R.From.Y, block.R.To.Y)
	}

	highests := NewHighests(maxX+1, maxY+1)
	numMoved := 0

	for _, block := range blocks {
		//logger.Infof("%d: %v", i, block)

		if block.R.From.X == block.R.To.X && block.R.From.Y == block.R.To.Y {
			//logger.Infof("vertical")

			p := pos.P2{X: block.R.From.X, Y: block.R.From.Y}
			height := block.R.From.Z - block.R.To.Z + 1
			highest := highests.Get(p)

			z := highest + 1

			if block.R.To.Z != z {
				numMoved++
				block.R.To.Z = z
				block.R.From.Z = z + (height - 1)
			}

			highests.Set(p, block.R.From.Z)

		} else {
			//logger.Infof("horizontal")

			xyFrom, xyTo := xyRange(block.R)

			highest := highests.GetRange(xyFrom, xyTo)
			z := highest + 1

			if block.R.From.Z != z {
				numMoved++
				block.R.From.Z, block.R.To.Z = z, z
			}
			highests.SetRange(xyFrom, xyTo, z)
		}

		//logger.Infof("%d: now %v", i, block)
	}

	return numMoved
}

func restsOn(blocks []*Block, rester, restee int) bool {
	er := map[pos.P2]bool{}

	{
		xyFrom, xyTo := xyRange(blocks[rester].R)
		walkRange(xyFrom, xyTo, func(p pos.P2) bool {
			er[p] = true
			return true
		})
	}

	matched := false
	{
		xyFrom, xyTo := xyRange(blocks[restee].R)
		walkRange(xyFrom, xyTo, func(p pos.P2) bool {
			if _, found := er[p]; found {
				matched = true
				return false
			}
			return true
		})
	}

	return matched
}

func numRestsOn(rester *Block, cands []*Block) int {
	er := map[pos.P2]bool{}

	{
		xyFrom, xyTo := xyRange(rester.R)
		walkRange(xyFrom, xyTo, func(p pos.P2) bool {
			er[p] = true
			return true
		})
	}

	num := 0
	for _, cand := range cands {
		xyFrom, xyTo := xyRange(cand.R)
		walkRange(xyFrom, xyTo, func(p pos.P2) bool {
			if _, found := er[p]; found {
				num++
				return false
			}
			return true
		})
	}
	return num
}

func dump(blocks []*Block, removeable map[Block]bool, focus int) {
	var focusBlock *Block

	maxX, maxY, maxZ := blocks[0].R.To.X, blocks[0].R.To.Y, 0
	zs := map[int][]*Block{}
	for _, block := range blocks {
		maxX = max(maxX, block.R.From.X, block.R.To.X)
		maxY = max(maxY, block.R.From.Y, block.R.To.Y)
		maxZ = max(maxZ, block.R.From.Z)

		if block.Num == focus {
			focusBlock = block
		}

		for z := block.R.To.Z; ; z++ {
			zs[z] = append(zs[z], block)
			if z == block.R.From.Z {
				break
			}
		}
	}

	startZ, endZ := 1, maxZ
	if focusBlock != nil {
		startZ, endZ = focusBlock.R.To.Z-1, focusBlock.R.From.Z+1
	}

	for z := startZ; z <= endZ; z++ {
		fmt.Printf("z=%d\n", z)

		g := grid.New[int](maxX+1, maxY+1)
		for _, block := range zs[z] {
			name := block.Num
			if _, found := removeable[*block]; found {
				name = -name
			}

			xyFrom, xyTo := xyRange(block.R)
			walkRange(xyFrom, xyTo, func(p pos.P2) bool {
				g.Set(p, name)
				return true
			})
		}

		g.Dump(false, func(p pos.P2, v int, _ bool) string {
			return strconv.Itoa(v) + " "
		})
		fmt.Println()
	}
}

func findRemoveableHard(blocks []*Block, focus int) map[Block]bool {
	fromZ := map[int][]*Block{}
	for _, block := range blocks {
		fromZ[block.R.From.Z] = append(fromZ[block.R.From.Z], block)
	}

	removeableBlocks := map[Block]bool{}
	for i, block := range blocks {
		focused := block.Num == focus

		if focused {
			logger.Infof("block %v", block)
		}
		above := block.R.From.Z + 1
		canRemove := true
		for j := i + 1; j < len(blocks); j++ {
			// Look at all blocks one level above 'block's top.
			if z := blocks[j].R.To.Z; z == above {
				if focused {
					logger.Infof("considering block %v", blocks[j])
				}
				if restsOn(blocks, j, i) {
					// block[j] rests on block[i]. Is block[i] its only support?
					if focused {
						logger.Infof("rests on %v", block)
					}

					// Count the number of blocks blocks[j] rests on.
					nro := numRestsOn(blocks[j], fromZ[block.R.From.Z])
					if focused {
						logger.Infof("%v rests on %v blocks", blocks[j], nro)
					}

					if nro == 1 {
						canRemove = false
						break
					}
				}
			}

		}

		if canRemove {
			if focused {
				logger.Infof("block %v can be removed", block)
			}
			removeableBlocks[*block] = true
		} else if focused {
			logger.Infof("block %v cannot be removed", block)
		}
	}

	return removeableBlocks
}

func findRemoveableEasy(blocks []*Block, focus int) map[Block]bool {
	maxX, maxY, maxZ := blocks[0].R.To.X, blocks[0].R.To.Y, 0
	for _, block := range blocks {
		maxX = max(maxX, block.R.From.X, block.R.To.X)
		maxY = max(maxY, block.R.From.Y, block.R.To.Y)
		maxZ = max(maxZ, block.R.From.Z)
	}

	g := make([][][]int, maxZ+1)
	for z := range g {
		g[z] = make([][]int, maxY+1)
		for y := range g[z] {
			g[z][y] = make([]int, maxX+1)
		}
	}

	get := func(p pos.P3) int {
		if p.Z < 0 || p.Y < 0 || p.X < 0 {
			panic("<0")
		}
		if p.Y >= len(g[0]) || p.X >= len(g[0][0]) {
			panic("xy")
		}
		if p.Z >= len(g) {
			return 0
		}
		return g[p.Z][p.Y][p.X]
	}

	set := func(p pos.P3, v int) {
		g[p.Z][p.Y][p.X] = v
	}

	byNum := map[int]*Block{}
	for _, block := range blocks {
		byNum[block.Num] = block
		for z := block.R.To.Z; z <= block.R.From.Z; z++ {
			walkXYRange(block.R, func(p pos.P2) bool {
				set(pos.P3{X: p.X, Y: p.Y, Z: z}, block.Num)
				return true
			})
		}
	}

	removeable := map[Block]bool{}
	for _, block := range blocks {
		focused := block.Num == focus

		if block.R.From.X == block.R.To.X && block.R.From.Y == block.R.To.Y {
			p := block.R.From // the top
			p.Z++             // where any supported block would be
			num := get(p)
			if num == 0 {
				removeable[*block] = true
				continue
			}

			supported := byNum[num]

			// `block` has another block atop it. Is `block` that
			// block's only support?
			supporters := map[int]bool{}
			walkXYRange(supported.R, func(p pos.P2) bool {
				below := pos.P3{X: p.X, Y: p.Y, Z: supported.R.To.Z - 1}
				if v := get(below); v != 0 {
					supporters[v] = true
				}
				return true
			})

			if len(supporters) > 1 {
				// There are other supporters
				removeable[*block] = true
			}
			continue
		}

		// Which blocks are on the next higher level at x,y that equal
		// the x.y in block?
		supports := map[int]bool{}
		walkXYRange(block.R, func(p pos.P2) bool {
			above := pos.P3{X: p.X, Y: p.Y, Z: block.R.From.Z + 1}
			if v := get(above); v != 0 {
				supports[v] = true
			}
			return true
		})

		if focused {
			logger.Infof("block %v supports %v", block, supports)
		}

		if len(supports) == 0 {
			// Nothing is above this block, so it can be removed
			removeable[*block] = true
			continue
		}

		// For each block supported by `block`
		allMultiSupported := true
		for num := range supports {
			supported := byNum[num]

			// Find the blocks that support it
			supporters := map[int]bool{}
			walkXYRange(supported.R, func(p pos.P2) bool {
				below := pos.P3{X: p.X, Y: p.Y, Z: supported.R.To.Z - 1}
				if v := get(below); v != 0 {
					supporters[v] = true
				}
				return true
			})

			if focused {
				logger.Infof("block %v supported by %v", supported, supporters)
			}

			// If multiple blocks support 'supported', including
			// 'block', then 'block' can be removed.
			if _, found := supporters[block.Num]; !found {
				panic("missing block")
			}
			if l := len(supporters); l == 1 {
				allMultiSupported = false
				break
			} else if l == 0 {
				panic("unexpected no support")
			}
		}

		if allMultiSupported {
			removeable[*block] = true
		}
	}

	return removeable
}

func copyBlocks(in []*Block) []*Block {
	out := make([]*Block, len(in))
	for i, b := range in {
		b2 := *b
		out[i] = &b2
	}
	return out
}

// 1230 too high
// 513 too high
// 422 too low
// 428 just right
func solveA(input []*Block) int {
	blocks := copyBlocks(input)

	// Sort by bottom Z, with lowest first. This means iterating through the
	// resulting slice goes from closest to the ground to furthest from it.
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].R.To.Z < blocks[j].R.To.Z
	})

	drop(blocks)

	// Re-sort since ordering may have changed during dropping.
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].R.To.Z < blocks[j].R.To.Z
	})

	focus := -1

	easyRemoveable := findRemoveableEasy(blocks, focus)
	hardRemoveable := findRemoveableHard(blocks, focus)

	for easy := range easyRemoveable {
		if _, found := hardRemoveable[easy]; !found {
			logger.Infof("%v in easy but not in hard", easy)
		}
	}
	for hard := range hardRemoveable {
		if _, found := easyRemoveable[hard]; !found {
			if focus == -1 || hard.Num == focus {
				logger.Infof("hard thinks %v is removeable; easy does not", hard)
			}
		}
	}

	logger.Infof("easy: %v", len(easyRemoveable))
	logger.Infof("hard: %v", len(hardRemoveable))

	if focus != -1 {
		dump(blocks, easyRemoveable, focus)
	}

	return len(easyRemoveable)
}

func tryRemove(in []*Block, toRemove int) int {
	blocks := make([]*Block, len(in)-1)
	i, j := 0, 0
	for i < len(in) {
		if in[i].Num != toRemove {
			b := *in[i]
			blocks[j] = &b
			j++
		}
		i++
	}

	n := drop(blocks)
	logger.Infof("removing %d drops %v blocks", toRemove, n)
	return n
}

func solveB(input []*Block) int {
	blocks := copyBlocks(input)

	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].R.To.Z < blocks[j].R.To.Z
	})

	drop(blocks)

	// Re-sort since ordering may have changed during dropping.
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].R.To.Z < blocks[j].R.To.Z
	})

	sum := 0
	for _, toRemove := range blocks {
		sum += tryRemove(blocks, toRemove.Num)
	}
	return sum
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
