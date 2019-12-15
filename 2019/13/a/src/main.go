package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/simmonmt/aoc/2019/common/intmath"
	"github.com/simmonmt/aoc/2019/common/logger"
	"github.com/simmonmt/aoc/2019/common/vm"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	ramPath = flag.String("ram", "", "path to file containing ram values")
)

type Tile int

const (
	TILE_EMPTY Tile = iota
	TILE_WALL
	TILE_BLOCK
	TILE_HPADDLE
	TILE_BALL
)

func (t Tile) String() string {
	switch t {
	case TILE_EMPTY:
		return " "
	case TILE_WALL:
		return "#"
	case TILE_BLOCK:
		return "X"
	case TILE_HPADDLE:
		return "-"
	case TILE_BALL:
		return "o"
	default:
		return "?"
	}
}

type Board struct {
	a    []Tile
	H, W int
}

func (b *Board) Get(x, y int) Tile {
	return b.a[y*b.W+x]
}

func (b *Board) Set(x, y int, t Tile) {
	b.a[y*b.W+x] = t
	//fmt.Printf("%d %d %d\n", x, y, t)
}

func makeBoard(vs []int64) *Board {
	maxX, maxY := 0, 0
	for i := 0; i < len(vs); i += 3 {
		x, y := int(vs[i]), int(vs[i+1])
		maxX = intmath.IntMax(maxX, x)
		maxY = intmath.IntMax(maxY, y)
	}

	h := maxY + 1
	w := maxX + 1

	b := Board{
		a: make([]Tile, h*w),
		H: h,
		W: w,
	}

	for i := 0; i < len(vs); i += 3 {
		x, y, t := int(vs[i]), int(vs[i+1]), Tile(vs[i+2])
		b.Set(x, y, t)
	}

	return &b
}

func printBoard(b *Board) {
	fmt.Printf("     ")
	for x := 0; x < b.W; x++ {
		fmt.Printf("%d", x%10)
	}
	fmt.Println()

	for y := 0; y < b.H; y++ {
		fmt.Printf("%3d: ", y)
		for x := 0; x < b.W; x++ {
			fmt.Print(b.Get(x, y))
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

	io := vm.NewSaverIO()
	if err := vm.Run(ram, io); err != nil {
		log.Fatal(err)
	}

	board := makeBoard(io.Written())
	printBoard(board)

	numBlocks := 0
	for y := 0; y < board.H; y++ {
		for x := 0; x < board.W; x++ {
			if board.Get(x, y) == TILE_BLOCK {
				numBlocks++
			}
		}
	}
	fmt.Printf("numBlocks=%d\n", numBlocks)

}
