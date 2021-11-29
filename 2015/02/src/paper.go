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
	"sort"
	"strconv"
	"strings"
)

type Dim struct {
	L, W, H int
}

func parseDim(str string) (*Dim, error) {
	pieces := strings.Split(str, "x")
	if len(pieces) != 3 {
		return nil, fmt.Errorf("wanted 3 pieces, got %d", len(pieces))
	}

	dim := &Dim{}
	var err error

	dim.L, err = strconv.Atoi(pieces[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse length %v: %v",
			pieces[0], err)
	}

	dim.W, err = strconv.Atoi(pieces[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse width %v: %v",
			pieces[1], err)
	}

	dim.H, err = strconv.Atoi(pieces[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse height %v: %v",
			pieces[2], err)
	}

	return dim, nil
}

func readDims(in io.Reader) ([]*Dim, error) {
	reader := bufio.NewReader(in)

	dims := []*Dim{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		dim, err := parseDim(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse dim %v: %v",
				line, err)
		}

		dims = append(dims, dim)
	}

	return dims, nil
}

func calcSide(a, b int) (area, perim int) {
	area = a * b
	perim = 2*a + 2*b
	return
}

func main() {
	dims, err := readDims(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	area := 0
	ribbon := 0
	for _, dim := range dims {
		areas := make([]int, 3)
		perims := make([]int, 3)

		areas[0], perims[0] = calcSide(dim.L, dim.W)
		areas[1], perims[1] = calcSide(dim.W, dim.H)
		areas[2], perims[2] = calcSide(dim.H, dim.L)

		sort.Ints(areas)
		sort.Ints(perims)

		area += 2*areas[0] + 2*areas[1] + 2*areas[2] + areas[0]
		ribbon += perims[0] + dim.L*dim.W*dim.H
	}

	fmt.Printf("total area: %d\n", area)
	fmt.Printf("total ribbon: %d\n", ribbon)
}
