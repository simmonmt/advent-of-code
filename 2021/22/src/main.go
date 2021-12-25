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
	"github.com/simmonmt/aoc/2021/common/intmath"
	"github.com/simmonmt/aoc/2021/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	axisPattern = regexp.MustCompile(`^([xyz])=([-\d]+)\.\.([-\d]+)$`)
)

type Field struct {
	Cubes  []*Cube
	Bounds *Cube
}

func NewField(bounds *Cube) *Field {
	return &Field{
		Cubes:  []*Cube{},
		Bounds: bounds,
	}
}

func (f *Field) Dump() {
	fmt.Println("bounds:", f.Bounds)
	for _, c := range f.Cubes {
		fmt.Println("-", c)
	}
	fmt.Println("#cubes", len(f.Cubes), "size", f.Size())
}

// Add the `add` cube to the field, removing any overlaps.
func (f *Field) Add(add *Cube) {
	res := []*Cube{}

	for _, c := range f.Cubes {
		if add.Contains(c) {
			continue // add will replace c, so drop c
		} else if add.Overlaps(c) {
			// add will partially replace c. copy over the part of c
			// that isn't overlapped by add.
			left := c.Sub(*add)
			res = append(res, left...)
		} else {
			res = append(res, c)
		}
	}

	res = append(res, add)
	f.Cubes = res
}

func (f *Field) Sub(sub *Cube) {
	res := []*Cube{}

	for _, c := range f.Cubes {
		if sub.Overlaps(c) {
			left := c.Sub(*sub)
			res = append(res, left...)
		} else {
			res = append(res, c)
		}
	}

	f.Cubes = res
}

func (f *Field) Size() int64 {
	sum := int64(0)
	for _, c := range f.Cubes {
		sum += c.Size()
	}
	return sum
}

type Cube struct {
	xLo, xHi int
	yLo, yHi int
	zLo, zHi int
}

func (c *Cube) String() string {
	return fmt.Sprintf("x=%v..%v,y=%v..%v,z=%v..%v",
		c.xLo, c.xHi, c.yLo, c.yHi, c.zLo, c.zHi)
}

func (c *Cube) Contains(c2 *Cube) bool {
	return c.xLo <= c2.xLo && c.xHi >= c2.xHi &&
		c.yLo <= c2.yLo && c.yHi >= c2.yHi &&
		c.zLo <= c2.zLo && c.zHi >= c2.zHi
}

func (c *Cube) Overlaps(c2 *Cube) bool {
	axisOverlaps := func(cLo, cHi, c2Lo, c2Hi int) bool {
		return c2Lo <= cHi && c2Hi >= cLo
	}

	if !axisOverlaps(c.xLo, c.xHi, c2.xLo, c2.xHi) {
		return false
	}
	if !axisOverlaps(c.yLo, c.yHi, c2.yLo, c2.yHi) {
		return false
	}
	if !axisOverlaps(c.zLo, c.zHi, c2.zLo, c2.zHi) {
		return false
	}
	return true
}

func (c *Cube) Sub(sub Cube) []*Cube {
	sub.zLo = intmath.IntMax(sub.zLo, c.zLo)
	sub.zHi = intmath.IntMin(sub.zHi, c.zHi)
	sub.yLo = intmath.IntMax(sub.yLo, c.yLo)
	sub.yHi = intmath.IntMin(sub.yHi, c.yHi)
	sub.xLo = intmath.IntMax(sub.xLo, c.xLo)
	sub.xHi = intmath.IntMin(sub.xHi, c.xHi)

	out := []*Cube{}

	if sub.zHi < c.zHi { // top
		out = append(out, &Cube{c.xLo, c.xHi, c.yLo, c.yHi, sub.zHi + 1, c.zHi})
	}
	if c.zLo < sub.zLo { // bottom
		out = append(out, &Cube{c.xLo, c.xHi, c.yLo, c.yHi, c.zLo, sub.zLo - 1})
	}

	if sub.xHi < c.xHi { // right
		out = append(out, &Cube{sub.xHi + 1, c.xHi, c.yLo, c.yHi, sub.zLo, sub.zHi})
	}
	if c.xLo < sub.xLo { // left
		out = append(out, &Cube{c.xLo, sub.xLo - 1, c.yLo, c.yHi, sub.zLo, sub.zHi})
	}

	if sub.yHi < c.yHi { // back
		out = append(out, &Cube{sub.xLo, sub.xHi, sub.yHi + 1, c.yHi, sub.zLo, sub.zHi})
	}
	if c.yLo < sub.yLo { // front
		out = append(out, &Cube{sub.xLo, sub.xHi, c.yLo, sub.yLo - 1, sub.zLo, sub.zHi})
	}

	return out
}

func (c *Cube) Size() int64 {
	return int64(c.xHi-c.xLo+1) * int64(c.yHi-c.yLo+1) * int64(c.zHi-c.zLo+1)
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

func solve(cmds []*Command, field *Field) int64 {
	for _, cmd := range cmds {
		logger.LogF("cmd %v", cmd)

		if cmd.On {
			field.Add(&cmd.Cube)
		} else {
			field.Sub(&cmd.Cube)
		}

		if logger.Enabled() {
			field.Dump()
		}
	}

	return field.Size()
}

func solveA(cmds []*Command) {
	filtered := []*Command{}
	for _, cmd := range cmds {
		if cmd.Cube.xLo < -50 || cmd.Cube.xHi > 50 {
			continue
		}
		filtered = append(filtered, cmd)
	}
	logger.LogF("input #cmds %v filtered %v", len(cmds), len(filtered))
	cmds = filtered

	field := NewField(&Cube{-50, 50, -50, 50, -50, 50})
	numLit := solve(cmds, field)

	fmt.Println("A", numLit)
}

func solveB(cmds []*Command) {
	bounds := cmds[0].Cube
	for i := 1; i < len(cmds); i++ {
		cmd := cmds[i]
		bounds.xLo = intmath.IntMin(bounds.xLo, cmd.Cube.xLo)
		bounds.xHi = intmath.IntMax(bounds.xHi, cmd.Cube.xHi)
		bounds.yLo = intmath.IntMin(bounds.yLo, cmd.Cube.yLo)
		bounds.yHi = intmath.IntMax(bounds.yHi, cmd.Cube.yHi)
		bounds.zLo = intmath.IntMin(bounds.zLo, cmd.Cube.zLo)
		bounds.zHi = intmath.IntMax(bounds.zHi, cmd.Cube.zHi)
	}

	field := NewField(&bounds)
	numLit := solve(cmds, field)

	fmt.Println("B", numLit)
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
	solveB(cmds)
}
