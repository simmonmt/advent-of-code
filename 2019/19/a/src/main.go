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

func query(ram vm.Ram, p pos.P2) int {
	io := vm.NewSaverIO(int64(p.X), int64(p.Y))
	if err := vm.Run(ram, io); err != nil {
		panic(fmt.Sprintf("program failed: %v", err))
	}
	if w := io.Written(); len(w) != 1 {
		panic(fmt.Sprintf("unexpected output %v", w))
	} else {
		return int(w[0])
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

	numAffected := 0
	//board := map[pos.P2]bool{}
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			p := pos.P2{x, y}
			if query(ram.Clone(), p) == 1 {
				//board[p] = true
				numAffected++
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}

	fmt.Printf("num affected: %d\n", numAffected)

}
