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

	"logger"
	"maze"
)

var (
	startX      = flag.Int("start_x", 1, "start x")
	startY      = flag.Int("start_y", 1, "start y")
	goalX       = flag.Int("goal_x", -1, "goal x")
	goalY       = flag.Int("goal_y", -1, "goal y")
	magicNumber = flag.Int("magic_number", -1, "magic number")
	verbose     = flag.Bool("verbose", false, "verbose mode")
)

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *goalX == -1 || *goalY == -1 {
		log.Fatalf("--goal_x and --goal_y are required")
	}
	if *magicNumber == -1 {
		log.Fatalf("--magic_number is required")
	}

	positions := maze.WalkMaze(*magicNumber, *startX, *startY, *goalX, *goalY)
	fmt.Println(len(positions) - 1)
}
