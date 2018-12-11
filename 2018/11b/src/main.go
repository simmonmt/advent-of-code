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
	fmt.Printf("grid %v,%v %v,%v\n", startX, startY, w, h)
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

	maxPoint := Point{-1, -1}
	maxSz := -1
	maxSum := -1

	grids := [301]Grid{}

	grids[1] = *grid
	for y := 1; y < 301; y++ {
		for x := 1; x < 301; x++ {
			if maxSz == -1 || grid[y][x] > maxSum {
				maxPoint = Point{x, y}
				maxSum = grid[y][x]
				maxSz = 1
			}
		}
	}

	// fmt.Println("orig")
	// dumpSquare(grid, 1, 1, 5, 5)
	// fmt.Println("sz=1")
	// dumpSquare(&grids[1], 1, 1, 5, 5)

	// Each element x,y in grids[5] contains the 5x5 sums rooted at x,y To
	// compute 0,0 in grids[6], we need 0,0 from grids[5] plus an additional
	// strip around the right and bottom edges -- i.e. the extra bit that
	// turns a 5x5 grid centered on 0,0 into a 6x6 grid centered on 0,0.

	for sz := 2; sz <= 300; sz++ {
		// for sz := 2; sz <= 2; sz++ {
		// fmt.Println(sz)
		prev := &grids[sz-1]

		for y := 1; y < 301-(sz-1); y++ {
			for x := 1; x < 301-(sz-1); x++ {
				// prev[y][x] gives us all but the surrounding
				// strip.
				sum := prev[y][x]

				// stripx = x+sz-1, vary y : y to y+sz-2
				for offY := 0; offY <= sz-2; offY++ {
					stripX := x + sz - 1
					stripY := y + offY
					// if y == 2 && x == 2 {
					// 	fmt.Printf("strip %v,%v %v\n", stripX, stripY,
					// 		grid[stripY][stripX])
					// }
					sum += grid[stripY][stripX]
				}

				// stripy = y+sz-1, vary x : x to x+sz-1
				for offX := 0; offX <= sz-1; offX++ {
					stripX := x + offX
					stripY := y + sz - 1
					// if y == 2 && x == 2 {
					// 	fmt.Printf("strip %v,%v %v\n", stripX, stripY,
					// 		grid[stripY][stripX])
					// }
					sum += grid[stripY][stripX]
				}

				grids[sz][y][x] = sum

				if sum > maxSum {
					maxPoint = Point{x, y}
					maxSum = sum
					maxSz = sz
					fmt.Printf("sz %v maxSz %v maxPoint %+v\n", sz, maxSz, maxPoint)
				}
			}
		}

		// fmt.Printf("sz=%v\n", sz)
		// dumpSquare(&grids[sz], 1, 1, 5, 5)
	}

	// dumpSquare(grid, 33, 45, 3, 3)
	// fmt.Println(sumSquare(grid, 33, 45, 3, 3))

	fmt.Printf("%v,%v,%v\n", maxPoint.X, maxPoint.Y, maxSz)
}
