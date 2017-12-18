package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type pair struct {
	x, y int
}

func gridSum(grid *map[pair]int, x, y int) int {
	return (*grid)[pair{x - 1, y - 1}] + (*grid)[pair{x, y - 1}] + (*grid)[pair{x + 1, y - 1}] + //
		(*grid)[pair{x - 1, y}] + (*grid)[pair{x, y}] + (*grid)[pair{x + 1, y}] + //
		(*grid)[pair{x - 1, y + 1}] + (*grid)[pair{x, y + 1}] + (*grid)[pair{x + 1, y + 1}]
}

func dumpGrid(grid *map[pair]int) {
	for y := 5; y >= -5; y-- {
		for x := -5; x <= 5; x++ {
			num, found := (*grid)[pair{x, y}]
			if found {
				fmt.Printf("%3d ", num)
			} else {
				fmt.Printf("    ")
			}
		}
		fmt.Printf("\n")
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s num", os.Args[0])
	}

	num, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to parse '%s' as num", os.Args[1])
	}
	fmt.Printf("num %d\n", num)

	grid := map[pair]int{}

	grid[pair{0, 0}] = 1

	base := 1
	spiralLow := 1
	spiralHigh := 1
	x, y := 0, 0

	for i := 2; true; i++ {
		if i > spiralHigh {
			spiralLow = spiralHigh + 1
			base += 2
			spiralHigh = base * base
			x += 1
			fmt.Printf("i %d x %d y %d\n", i, x, y)
		} else {
			sideLen := base - 1
			sideNum := (i - spiralLow) / sideLen

			switch sideNum {
			case 0:
				y += 1
				break // right side
			case 1:
				x -= 1
				break // top side
			case 2:
				y -= 1
				break // left side
			case 3:
				x += 1
				break // bottom side
			}
		}

		sum := gridSum(&grid, x, y)
		if sum > num {
			fmt.Printf("value %d\n", sum)
			break
		}

		grid[pair{x, y}] = sum
	}

	//dumpGrid(&grid)
}
