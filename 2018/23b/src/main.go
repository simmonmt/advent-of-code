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
	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"

	"intmath"
	"logger"
	"xyzpos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")

	inputPattern = regexp.MustCompile(`pos=<([^>]+)>, r=(\d+)`)
)

type Bot struct {
	Pos    xyzpos.Pos
	Radius int
}

func (b *Bot) Intersects(o *Bot) bool {
	dist := b.Pos.Dist(o.Pos)
	return dist <= b.Radius+o.Radius
}

func readInput() ([]*Bot, error) {
	bots := []*Bot{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		parts := inputPattern.FindStringSubmatch(line)
		if parts == nil {
			return nil, fmt.Errorf("failed to parse %v", line)
		}

		pos, err := xyzpos.Parse(parts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse pos: %v", err)
		}
		radius, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("failed to parse radiusz: %v", err)
		}

		bots = append(bots, &Bot{pos, radius})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return bots, nil
}

func findBounds(bots []*Bot) (xmin, xmax, ymin, ymax, zmin, zmax int) {
	xmin, ymin, zmin = math.MaxInt32, math.MaxInt32, math.MaxInt32
	xmax, ymax, zmax = math.MinInt32, math.MinInt32, math.MinInt32
	for _, bot := range bots {
		xmin = intmath.IntMin(xmin, bot.Pos.X)
		xmax = intmath.IntMax(xmax, bot.Pos.X)
		ymin = intmath.IntMin(ymin, bot.Pos.Y)
		ymax = intmath.IntMax(ymax, bot.Pos.Y)
		zmin = intmath.IntMin(zmin, bot.Pos.Z)
		zmax = intmath.IntMax(zmax, bot.Pos.Z)
	}

	return
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	bots, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	xmin, xmax, ymin, ymax, zmin, zmax := findBounds(bots)

	sz := 1
	for sz < xmax+1-xmin || sz < ymax+1-ymin || sz < zmax+1-zmin {
		sz *= 2
	}

	fmt.Printf("[%d,%d] [%d,%d] [%d,%d] %d\n", xmin, xmax, ymin, ymax, zmin, zmax, sz)

	ac := NewAllocatedCube(bots, &SearchCube{Min: xyzpos.Pos{xmin, ymin, zmin}, Size: sz})
	if len(ac.Bots) != len(bots) {
		for _, b := range bots {
			found := false
			for _, acb := range ac.Bots {
				if acb.Pos == b.Pos {
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("missing %v\n", b.Pos)
			}
		}

		panic(fmt.Sprintf("missing %d vs %d", len(ac.Bots), len(bots)))
	}

	cq := CubeQueue{}
	heap.Push(&cq, &cubeQueueItem{cube: ac})

	maxPos := []xyzpos.Pos{}
	maxNum := 0
	for cq.Len() > 0 {
		ac := heap.Pop(&cq).(*cubeQueueItem).cube

		if ac.Cube.Size == 1 {
			if len(ac.Bots) > maxNum {
				maxNum = len(ac.Bots)
				maxPos = []xyzpos.Pos{ac.Cube.Min}
			} else if len(ac.Bots) == maxNum {
				maxPos = append(maxPos, ac.Cube.Min)
			} else {
				break
			}
		} else {
			for _, sub := range ac.Divide() {
				heap.Push(&cq, &cubeQueueItem{cube: sub})
			}
		}
	}

	fmt.Println(maxPos, maxNum)
	for _, p := range maxPos {
		fmt.Printf("pos %v dist %v\n", p, p.Dist(xyzpos.Pos{0, 0, 0}))
	}
}
