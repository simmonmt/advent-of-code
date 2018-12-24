package main

import (
	"flag"
	"fmt"
	"log"

	"logger"
	"xypos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	depth   = flag.Int("depth", -1, "depth")
	target  = flag.String("target", "", "target")
)

type Fill int

const (
	FILL_ROCKY Fill = iota
	FILL_WET
	FILL_NARROW
)

func erosion(geo int) int {
	return (geo + *depth) % 20183
}

func erosionToFill(geo int) Fill {
	switch geo % 3 {
	case 0:
		return FILL_ROCKY
	case 1:
		return FILL_WET
	case 2:
		return FILL_NARROW
	default:
		panic("unknown")
	}
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *target == "" {
		log.Fatal("--target is required")
	}
	if *depth == -1 {
		log.Fatal("--depth is required")
	}

	target, err := xypos.Parse(*target)
	if err != nil {
		log.Fatal(err)
	}

	w := target.X + 1
	h := target.Y + 1

	geo := make([][]int, h)
	for y := range geo {
		geo[y] = make([]int, w)
		for x := range geo[y] {
			pos := xypos.Pos{x, y}
			val := -1

			if target.Eq(pos) {
				val = 0
			} else if x == 0 && y == 0 {
				val = 0
			} else if y == 0 {
				val = 16807 * x
			} else if x == 0 {
				val = 48271 * y
			} else {
				val = erosion(geo[y][x-1]) * erosion(geo[y-1][x])
			}

			geo[y][x] = val
		}
	}

	posns := []xypos.Pos{
		xypos.Pos{0, 0},
		xypos.Pos{1, 0},
		xypos.Pos{0, 1},
		xypos.Pos{1, 1},
		xypos.Pos{10, 10},
	}

	for _, p := range posns {
		g := geo[p.Y][p.X]
		fmt.Printf("%v geo %v erosion %v\n", p, g, erosion(g))
	}

	sum := 0
	for y := range geo {
		for x := range geo[y] {
			switch erosionToFill(erosion(geo[y][x])) {
			case FILL_ROCKY:
				sum += 0
			case FILL_WET:
				sum += 1
			case FILL_NARROW:
				sum += 2
			default:
				panic("unknown")
			}
		}
	}

	fmt.Println(sum)
}
