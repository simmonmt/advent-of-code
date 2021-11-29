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
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Matrix struct {
	elems [][]bool
}

func newMatrix(height, width int) *Matrix {
	elems := make([][]bool, height)
	for i := range elems {
		elems[i] = make([]bool, width)
	}

	return &Matrix{elems}
}

func (m *Matrix) Height() int {
	return len(m.elems[0])
}

func (m *Matrix) Width() int {
	return len(m.elems)
}

func (m *Matrix) Encode() uint {
	val := uint(0)
	for _, row := range m.elems {
		for _, cell := range row {
			val = (val << 1)
			if cell {
				val |= 1
			}
		}
	}
	return val
}

func (m *Matrix) FlipHorizontal() *Matrix {
	out := make([][]bool, len(m.elems))

	for rowNum, row := range m.elems {
		newRow := make([]bool, len(row))
		for colNum, cell := range row {
			newRow[len(row)-colNum-1] = cell
		}
		out[rowNum] = newRow
	}

	return &Matrix{out}
}

func (m *Matrix) FlipVertical() *Matrix {
	out := make([][]bool, len(m.elems))

	for rowNum, row := range m.elems {
		newRow := make([]bool, len(row))
		for colNum, cell := range row {
			newRow[colNum] = cell
		}
		out[len(m.elems)-rowNum-1] = newRow
	}

	return &Matrix{out}
}

func (m *Matrix) Rotate() *Matrix {
	newNumRows := len(m.elems[0])
	newNumCols := len(m.elems)

	out := make([][]bool, newNumRows)
	for i := range out {
		out[i] = make([]bool, newNumCols)
	}

	// 12   31
	// 34   42
	for rowNum, row := range m.elems {
		for colNum, cell := range row {
			newRow := colNum
			newCol := newNumCols - rowNum - 1
			out[newRow][newCol] = cell
		}
	}

	return &Matrix{out}
}

func (m *Matrix) Subset(xOrigin, yOrigin, xSize, ySize int) *Matrix {
	out := make([][]bool, ySize)
	for y := 0; y < ySize; y++ {
		out[y] = make([]bool, xSize)
		for x := 0; x < xSize; x++ {
			out[y][x] = m.elems[y+yOrigin][x+xOrigin]
		}
	}
	return &Matrix{out}
}

func (m *Matrix) CopyFrom(from *Matrix, xOrigin, yOrigin int) {
	for rowNum, row := range from.elems {
		for colNum, cell := range row {
			m.elems[rowNum+yOrigin][colNum+xOrigin] = cell
		}
	}
}

func (m *Matrix) ToString() string {
	out := ""

	for rowNum, row := range m.elems {
		if rowNum != 0 {
			out += "/"
		}
		for _, cell := range row {
			if cell {
				out += "#"
			} else {
				out += "."
			}
		}
	}

	out += fmt.Sprintf(":%04x", m.Encode())
	return out
}

func (m *Matrix) Dump(out io.Writer) {
	for _, row := range m.elems {
		for _, cell := range row {
			if cell {
				fmt.Fprintf(out, "#")
			} else {
				fmt.Fprintf(out, ".")
			}
		}
		fmt.Fprintf(out, "\n")
	}
}

func (m *Matrix) Count() int {
	num := 0
	for _, row := range m.elems {
		for _, cell := range row {
			if cell {
				num++
			}
		}
	}
	return num
}

func parseMatrix(str string) (*Matrix, error) {
	out := [][]bool{}

	parts := strings.Split(str, "/")
	for _, part := range parts {
		row := make([]bool, len(part))
		for i, c := range part {
			if c == '#' {
				row[i] = true
			}
		}
		out = append(out, row)
	}

	for i := 1; i < len(out); i++ {
		if len(out[i]) != len(out[0]) {
			return nil, fmt.Errorf("uneven matrix")
		}
	}

	return &Matrix{out}, nil
}

func parseLine(str string) (match, replace *Matrix, err error) {
	parts := strings.Split(str, " => ")

	if match, err = parseMatrix(parts[0]); err != nil {
		err = fmt.Errorf("failed to parse match in %v", str)
	} else if replace, err = parseMatrix(parts[1]); err != nil {
		err = fmt.Errorf("failed to parse replace in %v", str)
	}

	return
}

func addEncodings(matchers map[uint]*Matrix, match, replace *Matrix) error {
	rots := make([]*Matrix, 4)
	rots[0] = match

	for i := 1; i < 4; i++ {
		rots[i] = rots[i-1].Rotate()
	}

	encs := map[uint]bool{}
	for _, rot := range rots {
		encs[rot.Encode()] = true
		encs[rot.FlipHorizontal().Encode()] = true
		encs[rot.FlipVertical().Encode()] = true
	}

	for enc, _ := range encs {
		if existing, found := matchers[enc]; found {
			return fmt.Errorf(
				"enc %v for %v repl %v already in matchers for %v",
				enc, match.ToString(), replace.ToString(),
				existing.ToString())
		}
		matchers[enc] = replace
	}

	return nil
}

func readInput(in io.Reader) (twoMatchers, threeMatchers map[uint]*Matrix, err error) {
	twoMatchers = map[uint]*Matrix{}
	threeMatchers = map[uint]*Matrix{}

	reader := bufio.NewReader(in)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		match, replace, err := parseLine(line)
		if err != nil {
			return nil, nil, err
		}

		var matchers map[uint]*Matrix
		switch match.Width() {
		case 2:
			matchers = twoMatchers
			break
		case 3:
			matchers = threeMatchers
			break
		default:
			panic(fmt.Sprintf("unexpected matcher dim %v: %v",
				match.Width(), match.ToString()))
		}

		if err := addEncodings(matchers, match, replace); err != nil {
			return nil, nil, err
		}
	}

	return twoMatchers, threeMatchers, nil
}

func evolve(pat *Matrix, size int, matchers map[uint]*Matrix) *Matrix {
	repls := [][]*Matrix{}
	for rowNum := 0; rowNum < pat.Height(); rowNum += size {
		row := []*Matrix{}
		for colNum := 0; colNum < pat.Width(); colNum += size {
			sub := pat.Subset(colNum, rowNum, size, size)

			// fmt.Printf("looking at %v,%v pat:\n", colNum, rowNum)
			// sub.Dump(os.Stdout)

			replace, found := matchers[sub.Encode()]
			if !found {
				panic(fmt.Sprintf("failed to find %v in matchers",
					sub.ToString()))
			}

			// fmt.Println("replacing with:")
			// replace.Dump(os.Stdout)

			row = append(row, replace)
		}
		repls = append(repls, row)
	}

	// fmt.Println("pat in:")
	// pat.Dump(os.Stdout)
	// for y, row := range repls {
	// 	for x, cell := range row {
	// 		fmt.Printf("(%v,%v):\n", x, y)
	// 		cell.Dump(os.Stdout)
	// 	}
	// }

	out := newMatrix(repls[0][0].Height()*len(repls),
		repls[0][0].Width()*len(repls[0]))

	for rowNum, row := range repls {
		outY := rowNum * repls[0][0].Height()
		for colNum, cell := range row {
			outX := colNum * repls[0][0].Width()
			out.CopyFrom(cell, outX, outY)
		}
	}

	return out
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %v niters", os.Args[0])
	}
	nIters, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("failed to parse niters %v: %v", os.Args[1], err)
	}

	twoMatchers, threeMatchers, err := readInput(os.Stdin)
	if err != nil {
		log.Fatalf("failed to parse input: %v", err)
	}

	fmt.Printf("read %d twoMatchers\n", len(twoMatchers))
	fmt.Printf("read %d threeMatchers\n", len(threeMatchers))

	pat, _ := parseMatrix(".#./..#/###")
	// fmt.Println("initial:")
	// pat.Dump(os.Stdout)

	for i := 0; i < nIters; i++ {
		if pat.Width()%2 == 0 {
			pat = evolve(pat, 2, twoMatchers)
		} else {
			pat = evolve(pat, 3, threeMatchers)
		}

		// fmt.Println()
		// fmt.Println("after evolve")
		// pat.Dump(os.Stdout)
	}

	fmt.Printf("on cells: %d\n", pat.Count())
}
