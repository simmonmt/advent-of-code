// Copyright 2022 Google LLC
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
	"container/list"
	"flag"
	"fmt"
	"log"

	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/logger"
	"github.com/simmonmt/aoc/2022/common/mtsmath"
	"github.com/simmonmt/aoc/2022/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func parseInput(lines []string) (map[pos.P3]bool, error) {
	out := map[pos.P3]bool{}
	for _, line := range lines {
		p, err := pos.P3FromString(line)
		if err != nil {
			return nil, fmt.Errorf("bad pos %v: %v", line, err)
		}
		out[p] = true
	}
	return out, nil
}

func solveA(locs map[pos.P3]bool) int {
	sum := 0

	for loc := range locs {
		closed := 0
		for _, n := range loc.AllNeighbors(false) {
			if _, found := locs[n]; found {
				closed++
			}
		}

		logger.LogF("%v closed %v open %v", loc, closed, 6-closed)

		sum += 6 - closed
	}

	return sum
}

func findOutside(start, end pos.P3, locs map[pos.P3]bool) map[pos.P3]bool {
	outside := map[pos.P3]bool{start: true}

	considered := map[pos.P3]bool{start: true}

	queue := list.New()
	queue.PushBack(start)
	for queue.Front() != nil {
		cand := queue.Front().Value.(pos.P3)
		queue.Remove(queue.Front())

		for _, n := range cand.AllNeighbors(false) {
			if n.X < start.X || n.Y < start.Y || n.Z < start.Z {
				continue
			}
			if n.X > end.X || n.Y > end.Y || n.Z > end.Z {
				continue
			}

			if _, found := locs[n]; found {
				continue
			}
			if _, found := considered[n]; found {
				continue
			}
			considered[n] = true

			outside[n] = true
			queue.PushBack(n)
		}
	}

	return outside
}

func solveB(locs map[pos.P3]bool) int {
	var start, end pos.P3
	first := true

	for loc := range locs {
		if first == true {
			start = loc
			end = loc
			first = false
		} else {
			start = pos.P3{
				mtsmath.Min(start.X, loc.X),
				mtsmath.Min(start.Y, loc.Y),
				mtsmath.Min(start.Z, loc.Z),
			}

			end = pos.P3{
				mtsmath.Max(end.X, loc.X),
				mtsmath.Max(end.Y, loc.Y),
				mtsmath.Max(end.Z, loc.Z),
			}
		}
	}

	start = pos.P3{start.X - 1, start.Y - 1, start.Z - 1}
	end = pos.P3{end.X + 1, end.Y + 1, end.Z + 1}

	outside := findOutside(start, end, locs)

	sum := 0
	for loc := range locs {
		watch := false

		open := 0
		for _, n := range loc.AllNeighbors(false) {
			if _, found := locs[n]; found {
				if watch {
					logger.LogF("loc %v n %v touches another", loc, n)
				}
				continue
			}
			if _, found := outside[n]; found {
				if watch {
					logger.LogF("loc %v n %v open", loc, n)
				}
				open++
				continue
			}
			if watch {
				logger.LogF("loc %v n %v inside", loc, n)
			}

		}

		logger.LogF("%v open %v", loc, open)

		sum += open
	}

	return sum
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		log.Fatal(err)
	}

	input, err := parseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
