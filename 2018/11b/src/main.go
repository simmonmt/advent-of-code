package main

import (
	"flag"
	"fmt"

	"logger"
)

var (
	verbose      = flag.Bool("verbose", false, "verbose")
	serialNumber = flag.Int("serno", 0, "serial number")
)

type Grid [301][301]int

type Point struct{ X, Y int }

func sumSquare(grid *Grid, startX, startY, w, h int) int {
	sum := 0
	for y := startY; y < startY+h; y++ {
		for x := startX; x < startX+w; x++ {
			sum += grid[y][x]
		}
	}
	return sum
}

func dumpSquare(grid *Grid, startX, startY, w, h int) {
	for y := startY; y < startY+h; y++ {
		for x := startX; x < startX+w; x++ {
			fmt.Printf("%3d ", grid[y][x])
		}
		fmt.Println()
	}
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	grid := &Grid{}
	for y := 1; y < 301; y++ {
		for x := 1; x < 301; x++ {
			rackID := x + 10
			powerLevel := (((rackID*y+*serialNumber)*rackID)/100)%10 - 5
			grid[y][x] = powerLevel
		}
	}

	// fmt.Printf("122,79 %d\n", grid[79][122])
	// fmt.Printf("217,196 %d\n", grid[196][217])
	// fmt.Printf("101,153 %d\n", grid[153][101])

	maxPoint := Point{-1, -1}
	maxSz := -1
	maxSum := -1

	for sz := 1; sz <= 300; sz++ {
		fmt.Println(sz)
		for y := 1; y < 301-(sz-1); y++ {
			for x := 1; x < 301-(sz-1); x++ {
				sum := sumSquare(grid, x, y, sz, sz)
				if maxSum == -1 || sum > maxSum {
					maxPoint = Point{x, y}
					maxSum = sum
					maxSz = sz
					fmt.Printf("sz %v maxSz %v maxPoint %+v\n", sz, maxSz, maxPoint)
				}
			}
		}
	}

	// dumpSquare(grid, 33, 45, 3, 3)
	// fmt.Println(sumSquare(grid, 33, 45, 3, 3))

	fmt.Printf("%v,%v,%v\n", maxPoint.X, maxPoint.Y, maxSz)
}
