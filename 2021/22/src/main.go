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
	"regexp"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/logger"
	"github.com/simmonmt/aoc/2021/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	axisPattern = regexp.MustCompile(`^([xyz])=([-\d]+)\.\.([-\d]+)$`)
)

type Cube struct {
	xLo, xHi int
	yLo, yHi int
	zLo, zHi int
}

func (c *Cube) String() string {
	return fmt.Sprintf("x=%v..%v,y=%v..%v,z=%v..%v",
		c.xLo, c.xHi, c.yLo, c.yHi, c.zLo, c.zHi)
}

func (c *Cube) Contains(p pos.P3) bool {
	if p.X < c.xLo || p.X > c.xHi {
		return false
	}
	if p.Y < c.yLo || p.Y > c.yHi {
		return false
	}
	if p.Z < c.zLo || p.Z > c.zHi {
		return false
	}
	return true
}

type Command struct {
	On   bool
	Cube Cube
}

func (c *Command) String() string {
	cmd := "on"
	if !c.On {
		cmd = "off"
	}

	return fmt.Sprintf("<%v %v>", cmd, c.Cube.String())
}

func readInput(path string) ([]*Command, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, err
	}

	cmds := []*Command{}
	for i, line := range lines {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("%d: no space", i+1)
		}

		cmd := &Command{
			On: parts[0] == "on",
		}

		for _, r := range strings.Split(parts[1], ",") {
			parts := axisPattern.FindStringSubmatch(r)
			if parts == nil {
				return nil, fmt.Errorf("%d: bad axis: %v", i+1, r)
			}

			axis := parts[1]

			from, err := strconv.Atoi(parts[2])
			if err != nil {
				return nil, fmt.Errorf("%d: bad %v from: %v", i+1, axis, err)
			}

			to, err := strconv.Atoi(parts[3])
			if err != nil {
				return nil, fmt.Errorf("%d: bad %v to: %v", i+1, axis, err)
			}

			switch axis {
			case "x":
				cmd.Cube.xLo, cmd.Cube.xHi = from, to
			case "y":
				cmd.Cube.yLo, cmd.Cube.yHi = from, to
			case "z":
				cmd.Cube.zLo, cmd.Cube.zHi = from, to
			default:
				return nil, fmt.Errorf("%d: unknown axis %v", i+1, axis)
			}
		}

		cmds = append(cmds, cmd)
	}

	return cmds, err
}

func runCommands(p pos.P3, cmds []*Command) bool {
	state := false
	for _, cmd := range cmds {
		if cmd.Cube.xLo < -50 || cmd.Cube.xHi > 50 {
			continue
		}

		if cmd.Cube.Contains(p) {
			state = cmd.On
		}
	}
	return state
}

func solveA(cmds []*Command) {
	numLit := 0

	for z := -50; z <= 50; z++ {
		for y := -50; y <= 50; y++ {
			for x := -50; x <= 50; x++ {
				p := pos.P3{x, y, z}
				state := runCommands(p, cmds)
				if state {
					numLit++
				}
			}
		}
	}

	fmt.Println("A", numLit)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	cmds, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	solveA(cmds)
}
