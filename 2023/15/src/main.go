// Copyright 2023 Google LLC
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
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func parseInput(lines []string) ([]string, error) {
	if len(lines) > 1 {
		return nil, fmt.Errorf("too many lines")
	}

	out := strings.Split(lines[0], ",")
	return out, nil
}

func hashString(line string) int {
	cur := 0
	for _, r := range line {
		cur += int(r)
		cur *= 17
		cur %= 256
	}
	return cur
}

func solveA(input []string) int {
	sum := 0
	for _, line := range input {
		sum += hashString(line)
	}
	return sum
}

type Lens struct {
	Label string
	Focal int
}

type Box struct {
	lenses   []*Lens
	presence map[string]bool
}

func NewBox() *Box {
	return &Box{
		lenses:   []*Lens{},
		presence: map[string]bool{},
	}
}

func (b *Box) Empty() bool {
	return len(b.lenses) == 0
}

func (b *Box) Remove(label string) {
	if _, found := b.presence[label]; !found {
		return
	}

	delete(b.presence, label)
	ol := []*Lens{}
	for _, l := range b.lenses {
		if l.Label == label {
			continue
		}
		ol = append(ol, l)
	}
	b.lenses = ol
}

func (b *Box) Insert(label string, focal int) {
	if _, found := b.presence[label]; found {
		for _, l := range b.lenses {
			if l.Label == label {
				l.Focal = focal
				return
			}
		}
		panic("not found")
	}

	b.presence[label] = true
	b.lenses = append(b.lenses, &Lens{Label: label, Focal: focal})
}

type Boxes [256]*Box

func NewBoxes() *Boxes {
	a := [256]*Box{}
	for i := 0; i < 256; i++ {
		a[i] = NewBox()
	}
	return (*Boxes)(&a)
}

func (b *Boxes) Dump() {
	for i := 0; i < 256; i++ {
		if b[i].Empty() {
			continue
		}

		fmt.Printf("Box %d:", i)
		for _, lens := range b[i].lenses {
			fmt.Printf(" [%s %d]", lens.Label, lens.Focal)
		}
		fmt.Println()
	}
}

func (b *Boxes) FocusingPower() int {
	sum := 0
	for i := 0; i < 256; i++ {
		box := b[i]
		for j, lens := range box.lenses {
			slot := j + 1
			sum += (i + 1) * slot * lens.Focal
		}
	}
	return sum
}

func solveB(commands []string) int {
	boxes := NewBoxes()

	for _, command := range commands {
		var label string
		var focal int
		remove := command[len(command)-1] == '-'

		if remove {
			label = command[0 : len(command)-1]
		} else {
			var num string
			var ok bool
			label, num, ok = strings.Cut(command, "=")
			if !ok {
				panic("bad command")
			}

			var err error
			focal, err = strconv.Atoi(num)
			if err != nil {
				panic("bad number")
			}
		}

		boxNum := hashString(label)
		box := boxes[boxNum]
		if remove {
			box.Remove(label)
		} else {
			box.Insert(label, focal)
		}

		// fmt.Printf("After \"%s\":\n", command)
		// boxes.Dump()
		// fmt.Println()
	}

	return boxes.FocusingPower()
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		logger.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		logger.Fatalf("failed to read input: %v", err)
	}

	input, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
