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
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var (
	commandPattern = regexp.MustCompile(`(turn on|turn off|toggle) ([0-9]+),([0-9]+) through ([0-9]+),([0-9]+)`)
)

type Pos struct {
	X, Y int
}

type Matrix struct {
	mat  [][]bool
	w, h int
}

func NewMatrix(w, h int) *Matrix {
	mat := make([][]bool, h)
	for y := 0; y < h; y++ {
		mat[y] = make([]bool, w)
	}

	return &Matrix{
		mat: mat,
		w:   w,
		h:   h,
	}
}

func (m *Matrix) Set(a, b Pos, val bool) {
	for y := a.Y; y <= b.Y; y++ {
		for x := a.X; x <= b.X; x++ {
			m.mat[y][x] = val
		}
	}
}

func (m *Matrix) Toggle(a, b Pos) {
	for y := a.Y; y <= b.Y; y++ {
		for x := a.X; x <= b.X; x++ {
			m.mat[y][x] = !m.mat[y][x]
		}
	}
}

func (m *Matrix) Count() int {
	numOn := 0

	for _, row := range m.mat {
		for _, val := range row {
			if val {
				numOn++
			}
		}
	}

	return numOn
}

func (m *Matrix) Dump() {
	for _, row := range m.mat {
		for _, val := range row {
			if val {
				fmt.Printf("O")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
}

func mkPos(xStr, yStr string) (Pos, error) {
	x, err := strconv.Atoi(xStr)
	if err != nil {
		return Pos{}, err
	}
	y, err := strconv.Atoi(yStr)
	if err != nil {
		return Pos{}, err
	}

	return Pos{X: x, Y: y}, nil
}

func main() {
	mat := NewMatrix(1000, 1000)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		match := commandPattern.FindStringSubmatch(line)
		if match == nil {
			log.Fatalf("failed to parse %v", match)
		}
		fmt.Println(match)

		command := match[1]
		a, err := mkPos(match[2], match[3])
		if err != nil {
			log.Fatalf("failed to parse a in %v: %v", line, err)
		}
		b, err := mkPos(match[4], match[5])
		if err != nil {
			log.Fatalf("failed to parse b in %v: %v", line, err)
		}

		switch command {
		case "turn on":
			mat.Set(a, b, true)
			break
		case "turn off":
			mat.Set(a, b, false)
			break
		case "toggle":
			mat.Toggle(a, b)
			break
		default:
			log.Fatalf("unknown command %v in %v", command, line)
		}
		// mat.Dump()
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("reading stdin: %v", err)
	}

	fmt.Println(mat.Count())
}
