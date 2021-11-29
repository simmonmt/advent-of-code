// Copyright 2021 Google LLC
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
	maxSum := -1

	for y := 1; y < 301-2; y++ {
		for x := 1; x < 301-2; x++ {
			sum := sumSquare(grid, x, y, 3, 3)
			if maxSum == -1 || sum > maxSum {
				maxPoint = Point{x, y}
				maxSum = sum
			}
		}
	}

	// dumpSquare(grid, 33, 45, 3, 3)
	// fmt.Println(sumSquare(grid, 33, 45, 3, 3))

	fmt.Println(maxPoint)
}
