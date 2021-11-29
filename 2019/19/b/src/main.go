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

	"github.com/simmonmt/aoc/2019/common/logger"
	"github.com/simmonmt/aoc/2019/common/pos"
	"github.com/simmonmt/aoc/2019/common/vm"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	ramPath = flag.String("ram", "", "path to file containing ram values")
)

func inBeam(ram vm.Ram, p pos.P2) bool {
	io := vm.NewSaverIO(int64(p.X), int64(p.Y))
	if err := vm.Run(ram.Clone(), io); err != nil {
		panic(fmt.Sprintf("program failed: %v", err))
	}
	if w := io.Written(); len(w) != 1 {
		panic(fmt.Sprintf("unexpected output %v", w))
	} else {
		return w[0] == 1
	}
}

func doesFit(ram vm.Ram, p pos.P2) (fits bool, closest pos.P2) {
	// find the bottom
	bottomY := -1
	for y := p.Y; bottomY < 0; y++ {
		if !inBeam(ram, pos.P2{p.X, y}) {
			bottomY = y - 1
		}
	}
	if !inBeam(ram, pos.P2{p.X, bottomY}) {
		panic("bad bottomy")
	}

	topY := bottomY - 99
	if !inBeam(ram, pos.P2{p.X, topY}) {
		return false, pos.P2{}
	}

	rightX := p.X + 99
	if !inBeam(ram, pos.P2{rightX, topY}) {
		return false, pos.P2{}
	}

	for inBeam(ram, pos.P2{p.X, topY - 1}) && inBeam(ram, pos.P2{p.X + 99, topY - 1}) {
		topY--
		bottomY--
		fmt.Printf("refined at %d to y=%d\n", p.X, topY)
	}

	fmt.Printf("box [%d,%d] [%d,%d]\n", p.X, topY, rightX, bottomY)

	return true, pos.P2{p.X, topY}
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

	// Make sure it's at 0,0
	if !inBeam(ram, pos.P2{0, 0}) {
		panic("not in beam at 0,0")
	}

	// Find the center at 50,0
	startY, endY := -1, -1
	for y := 0; startY < 0 || endY < 0; y++ {
		p := pos.P2{50, y}

		if inBeam(ram, p) {
			if startY == -1 {
				startY = y
			}
		} else {
			if startY > 0 && endY == -1 {
				endY = y - 1
			}
		}
	}

	midpointY := startY + (endY+1-startY)/2
	fmt.Printf("at 50 start %d end %d midpoint %d\n", startY, endY, midpointY)

	fitsX := -1
	lastY := -1
	for i := 2; i < 200; i++ {
		p := pos.P2{50 * i, midpointY * i}
		if fits, _ := doesFit(ram, p); fits {
			fmt.Printf("fits at %v\n", p)
			fitsX = p.X
			lastY = p.Y
			break
		}
	}

	for x := fitsX - 1; x > 0; x-- {
		for !inBeam(ram, pos.P2{x, lastY}) {
			lastY--
		}

		if fits, closest := doesFit(ram, pos.P2{x, lastY}); fits {
			fmt.Printf("fits at %v closest %v answer %d\n",
				x, closest, closest.X*10000+closest.Y)
		}
	}
}
