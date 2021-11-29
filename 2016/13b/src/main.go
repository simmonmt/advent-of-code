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
	"log"

	"intmath"
	"logger"
	"maze"
)

var (
	startX      = flag.Int("start_x", 1, "start x")
	startY      = flag.Int("start_y", 1, "start y")
	distance    = flag.Uint("dist", 50, "distance")
	magicNumber = flag.Int("magic_number", -1, "magic number")
	verbose     = flag.Bool("verbose", false, "verbose mode")
)

func dist(x1, y1, x2, y2 int) uint {
	return uint(intmath.Abs(x1-x2) + intmath.Abs(y1-y2))
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *magicNumber == -1 {
		log.Fatalf("--magic_number is required")
	}

	goals := [][2]int{}

	for x := *startX - int(*distance); x <= *startX+int(*distance); x++ {
		for y := *startY - int(*distance); y <= *startY+int(*distance); y++ {
			if !maze.IsOpenSpace(*magicNumber, x, y) {
				continue
			}

			// if dist(*startX, *startY, x, y) > *distance {
			// 	continue
			// }

			goals = append(goals, [2]int{x, y})
		}
	}

	fmt.Printf("goals: %v\n", len(goals))

	numReachedGoals := 0
	allPositions := map[string]bool{}
	for i, goal := range goals {
		logger.LogF("checking goal %v of %v: %v,%v\n", i, len(goals), goal[0], goal[1])
		positions := maze.WalkMaze(*magicNumber, *startX, *startY, goal[0], goal[1])
		logger.LogF("checking goal result %v\n", len(positions))
		if len(positions) == 0 || len(positions)-1 > int(*distance) {
			continue
		}

		numReachedGoals++
		for _, pos := range positions {
			allPositions[pos] = true
		}
	}

	fmt.Printf("positions: %v\n", len(allPositions))
	fmt.Printf("reached goals: %v\n", numReachedGoals)
}
