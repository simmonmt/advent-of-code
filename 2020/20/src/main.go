// Possible optimization: Some combinations of (dir,fliph,flipv) are
// redundant. It finishes in less than a minute, though, and I never
// want to see another sea monster again.

package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"regexp"

	"github.com/simmonmt/aoc/2020/20/src/tiles"
	"github.com/simmonmt/aoc/2020/common/dir"
	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/intmath"
	"github.com/simmonmt/aoc/2020/common/logger"
	"github.com/simmonmt/aoc/2020/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	titlePattern = regexp.MustCompile(`^Tile ([0-9]+):$`)
)

type Board struct {
	dim  int // height and width
	used map[int]*tiles.OrientedTile
	arr  []*tiles.OrientedTile
}

func newBoard(dim int) *Board {
	return &Board{
		dim:  dim,
		used: map[int]*tiles.OrientedTile{},
		arr:  make([]*tiles.OrientedTile, dim*dim),
	}
}

func (b *Board) Dim() int {
	return b.dim
}

func (b *Board) off(pos pos.P2) int {
	return b.dim*pos.Y + pos.X
}

func (b *Board) Get(pos pos.P2) *tiles.OrientedTile {
	return b.arr[b.off(pos)]
}

func (b *Board) Set(pos pos.P2, t *tiles.OrientedTile) {
	off := b.off(pos)

	if _, found := b.used[t.Num()]; found {
		b.Dump()
		panic(fmt.Sprintf("reuse tile %v", t.Num()))
	}
	b.used[t.Num()] = t
	b.arr[off] = t
}

func (b *Board) Clear(pos pos.P2) {
	off := b.off(pos)
	if b.arr[off] == nil {
		b.Dump()
		panic(fmt.Sprintf("clear unused %v", pos))
	}

	t := b.arr[off]
	if _, found := b.used[t.Num()]; !found {
		panic("clear not used")
	}

	b.arr[off] = nil
	delete(b.used, t.Num())
}

func (b *Board) Used(num int) bool {
	_, found := b.used[num]
	return found
}

func (b *Board) Dump() {
	for y := 0; y < b.dim; y++ {
		for x := 0; x < b.dim; x++ {
			p := pos.P2{X: x, Y: y}

			otStr := ""
			ot := b.arr[b.off(p)]
			if ot != nil {
				otStr = ot.String()
			}

			fmt.Printf("%-12s | ", otStr)
		}
		fmt.Println()
	}
}

func (b *Board) Neighbors(p pos.P2) (neighbors map[dir.Dir]*tiles.OrientedTile) {
	neighbors = map[dir.Dir]*tiles.OrientedTile{}

	for _, dir := range dir.AllDirs {
		np := dir.From(p)
		if np.X < 0 || np.Y < 0 || np.X >= b.dim || np.Y >= b.dim {
			continue
		}

		if ot := b.arr[b.off(np)]; ot != nil {
			neighbors[dir] = ot
		}
	}

	return
}

func makeAllSidesMap(allTiles []*tiles.Tile) map[tiles.Side][]*tiles.OrientedTile {
	allSides := map[tiles.Side][]*tiles.OrientedTile{}

	for _, tile := range allTiles {
		for _, d := range dir.AllDirs {
			for _, flip := range []int{0, 1, 2, 3} {
				flipH := (flip & 1) != 0
				flipV := (flip & 2) != 0

				ot := tiles.NewOrientedTile(tile, d, flipH, flipV)

				for _, sideDir := range dir.AllDirs {
					side := ot.Side(sideDir)

					if ots, found := allSides[side]; !found {
						allSides[side] = []*tiles.OrientedTile{ot}
					} else {
						matched := false
						for _, elem := range ots {
							if elem.Num() == ot.Num() {
								matched = true
								break
							}
						}
						if !matched {
							allSides[side] = append(
								allSides[side], ot)
						}
					}
				}
			}
		}
	}

	return allSides
}

func findBorderTiles(allTiles []*tiles.Tile, allSides map[tiles.Side][]*tiles.OrientedTile) []*tiles.Tile {
	borders := []*tiles.Tile{}

	for _, tile := range allTiles {
		numNeighbors := 0
		matchDirs := []dir.Dir{}
		for _, d := range dir.AllDirs {
			side := tile.Side(d)
			if len(allSides[side]) > 1 {
				numNeighbors++
				matchDirs = append(matchDirs, d)
			}
		}

		if numNeighbors == 2 || numNeighbors == 3 {
			borders = append(borders, tile)
		} else if numNeighbors != 4 {
			panic("bad matches")
		}
	}

	return borders
}

func makePerimiterPath(dim int) (path []pos.P2) {
	path = []pos.P2{}
	for x := 0; x < dim; x++ {
		path = append(path, pos.P2{X: x, Y: 0})
	}
	for y := 1; y < dim; y++ {
		path = append(path, pos.P2{X: dim - 1, Y: y})
	}
	for x := dim - 2; x >= 0; x-- {
		path = append(path, pos.P2{X: x, Y: dim - 1})
	}
	for y := dim - 2; y > 0; y-- {
		path = append(path, pos.P2{X: 0, Y: y})
	}

	return
}

func makeInnerPath(dim int) (path []pos.P2) {
	path = []pos.P2{}

	for y := 1; y < dim-1; y++ {
		for x := 1; x < dim-1; x++ {
			path = append(path, pos.P2{X: x, Y: y})
		}
	}

	return
}

func pieceWillFit(b *Board, p pos.P2, ot *tiles.OrientedTile) bool {
	for dir, neighbor := range b.Neighbors(p) {
		if neighbor.Side(dir.Reverse()) != ot.Side(dir) {
			return false
		}
	}

	return true
}

func placeTiles(b *Board, tilesToPlace []*tiles.Tile, path []pos.P2) bool {
	if len(path) == 0 {
		return true
	}

	// for each candidate tile
	//   for each fitting orientation
	//     recurse

	curPos := path[0]

	for _, tile := range tilesToPlace {
		if b.Used(tile.Num()) {
			continue
		}

		cands := []*tiles.OrientedTile{}
		for _, dir := range dir.AllDirs {
			for _, flip := range []int{0, 1, 2, 3} {
				flipH := (flip & 1) != 0
				flipV := (flip & 2) != 0
				ot := tiles.NewOrientedTile(tile, dir, flipH, flipV)

				if pieceWillFit(b, curPos, ot) {
					cands = append(cands, ot)
				}
			}
		}

		for _, cand := range cands {
			b.Set(curPos, cand)
			if placeTiles(b, tilesToPlace, path[1:]) {
				return true
			}
			b.Clear(curPos)
		}
	}

	return false
}

func SolveBoard(b *Board, allTiles []*tiles.Tile) {
	allSides := makeAllSidesMap(allTiles)
	borders := findBorderTiles(allTiles, allSides)
	perimiter := makePerimiterPath(b.Dim())

	logger.LogF("solving for borders")
	placeTiles(b, borders, perimiter)

	inner := makeInnerPath(b.Dim())

	left := []*tiles.Tile{}
	for _, tile := range allTiles {
		if !b.Used(tile.Num()) {
			left = append(left, tile)
		}
	}

	logger.LogF("solving for inner")
	placeTiles(b, left, inner)
	logger.LogF("done solving")
}

const (
	monsterWidth  = 19
	monsterHeight = 3
)

var (
	monsterPattern = []pos.P2{
		pos.P2{18, 0},
		pos.P2{0, 1},
		pos.P2{5, 1},
		pos.P2{6, 1},
		pos.P2{11, 1},
		pos.P2{12, 1},
		pos.P2{17, 1},
		pos.P2{18, 1},
		pos.P2{19, 1},
		pos.P2{1, 2},
		pos.P2{4, 2},
		pos.P2{7, 2},
		pos.P2{10, 2},
		pos.P2{13, 2},
		pos.P2{16, 2},
	}
)

//   00000000001111111111
//   01234567890123456789
//   +                 O
//   O    OO    OO    OOO
//    O  O  O  O  O  O

func findSeaMonsters(ot *tiles.OrientedTile) (numFound, nonMonster int) {
	matched := map[pos.P2]bool{}
	numMonsters := 0

	for y := 0; y < ot.Dim()-monsterHeight; y++ {
		for x := 0; x < ot.Dim()-monsterWidth; x++ {
			found := true
			thisMatched := map[pos.P2]bool{}
			for _, off := range monsterPattern {
				pos := pos.P2{X: x + off.X, Y: y + off.Y}
				if !ot.Get(pos) {
					found = false
					break
				}
				thisMatched[pos] = true
			}
			if !found {
				continue
			}

			for pos := range thisMatched {
				if _, found := matched[pos]; found {
					panic("overlap")
				}
				matched[pos] = true
			}

			numMonsters++
		}
	}

	if numMonsters == 0 {
		return 0, 0
	}

	numNonMonster := 0
	for y := 0; y < ot.Dim(); y++ {
		for x := 0; x < ot.Dim(); x++ {
			pos := pos.P2{X: x, Y: y}
			if ot.Get(pos) {
				if _, found := matched[pos]; !found {
					numNonMonster++
				}
			}
		}
	}

	return numMonsters, numNonMonster
}

func SolveB(b *Board, tileDim int) int {
	combinedDim := b.Dim() * (tileDim - 2)
	combinedArr := make([][]rune, combinedDim)
	for i := range combinedArr {
		combinedArr[i] = make([]rune, combinedDim)
	}

	for tileY := 0; tileY < b.Dim(); tileY++ {
		for tileX := 0; tileX < b.Dim(); tileX++ {
			tilePos := pos.P2{X: tileX, Y: tileY}
			t := b.Get(tilePos)

			for cellY := 1; cellY < tileDim-1; cellY++ {
				for cellX := 1; cellX < tileDim-1; cellX++ {
					cellPos := pos.P2{X: cellX, Y: cellY}

					combinedPos := pos.P2{
						X: tileX*(tileDim-2) + (cellPos.X - 1),
						Y: tileY*(tileDim-2) + (cellPos.Y - 1),
					}

					val := '.'
					if t.Get(cellPos) {
						val = '#'
					}

					combinedArr[combinedPos.Y][combinedPos.X] = val
				}
			}
		}
	}

	combinedStrs := make([]string, len(combinedArr))
	for y := range combinedArr {
		combinedStrs[y] = string(combinedArr[y])
	}

	combinedTile, err := tiles.NewTile(1, combinedStrs, combinedDim)
	if err != nil {
		log.Fatalf("bad combined (shouldn't happen)")
	}

	answer := -1
	for _, d := range dir.AllDirs {
		for _, flip := range []int{0, 1, 2, 3} {
			flipH := (flip & 1) != 0
			flipV := (flip & 2) != 0

			ot := tiles.NewOrientedTile(combinedTile, d, flipH, flipV)
			if num, nonMonster := findSeaMonsters(ot); num > 0 {
				logger.LogF("found %d monsters, %d non in %s",
					num, nonMonster, ot)
				if answer == -1 {
					answer = nonMonster
				} else if answer != nonMonster {
					panic("disagreement")
				}
			}
		}
	}

	return answer
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

	allTiles := []*tiles.Tile{}
	for {
		title := lines[0]
		parts := titlePattern.FindStringSubmatch(title)
		if parts == nil {
			log.Fatalf("bad title")
		}
		tileNum := intmath.AtoiOrDie(parts[1])

		end := 1
		for ; end < len(lines); end++ {
			if lines[end] == "" {
				break
			}
		}

		body := lines[1:end]

		tile, err := tiles.NewTile(tileNum, body, len(body[0]))
		if err != nil {
			log.Fatalf("failed to parse tile: %v", err)
		}
		allTiles = append(allTiles, tile)

		if end == len(lines) {
			break
		}
		lines = lines[end+1:]
	}
	tileDim := allTiles[0].Dim()

	dim := int(math.Sqrt(float64(len(allTiles))))
	if dim*dim != len(allTiles) {
		panic("non-square")
	}

	b := newBoard(dim)
	SolveBoard(b, allTiles)

	corners := []pos.P2{
		pos.P2{0, 0},
		pos.P2{dim - 1, 0},
		pos.P2{0, dim - 1},
		pos.P2{dim - 1, dim - 1},
	}

	mult := 1
	for _, corner := range corners {
		mult *= b.Get(corner).Num()
	}
	fmt.Printf("A: %v\n", mult)

	fmt.Printf("B: %v\n", SolveB(b, tileDim))
}
