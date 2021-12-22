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

	separatorPattern = regexp.MustCompile(`^-+ scanner (\d+) -+$`)
)

type Field struct {
	Beacons    []pos.P3
	ScannerNum int
	ScannerPos pos.P3
}

func NewField(num int, pos pos.P3, beacons []pos.P3) *Field {
	return &Field{
		Beacons:    beacons,
		ScannerNum: num,
		ScannerPos: pos,
	}
}

func (f *Field) Rotate(mat []int, deg int) *Field {
	nb := make([]pos.P3, len(f.Beacons))
	for i, p := range f.Beacons {
		nb[i] = rotatePos(mat, deg, p)
	}

	return &Field{
		Beacons:    nb,
		ScannerNum: f.ScannerNum,
		ScannerPos: rotatePos(mat, deg, f.ScannerPos),
	}
}

func readInput(path string) (map[int]*Field, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, err
	}

	scannerNum := 0
	beacons := map[int][]pos.P3{}

	for i, line := range lines {
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "--") {
			parts := separatorPattern.FindStringSubmatch(line)
			if parts == nil {
				return nil, fmt.Errorf("%d: failed to parse separator: %v", i, line)
			}

			scannerNum, _ = strconv.Atoi(parts[1])
			beacons[scannerNum] = []pos.P3{}
			continue
		}

		p, err := pos.P3FromString(line)
		if err != nil {
			return nil, fmt.Errorf("%d: failed to parse pos: %v", i, err)
		}

		beacons[scannerNum] = append(beacons[scannerNum], p)
	}

	scanners := map[int]*Field{}
	for scannerNum, beaconPosns := range beacons {
		scanners[scannerNum] = NewField(scannerNum, pos.P3{0, 0, 0}, beaconPosns)
	}

	return scanners, err
}

var (
	RX = []int{1, 0, 0, 0, 0, -1, 0, 1, 0} // 90 deg
	RY = []int{0, 0, 1, 0, 1, 0, -1, 0, 0} // 90 deg
	RZ = []int{0, -1, 0, 1, 0, 0, 0, 0, 1} // 90 deg
)

func rotatePos(mat []int, deg int, p pos.P3) pos.P3 {
	for ; deg > 0; deg -= 90 {
		nX := mat[0]*p.X + mat[1]*p.Y + mat[2]*p.Z
		nY := mat[3]*p.X + mat[4]*p.Y + mat[5]*p.Z
		nZ := mat[6]*p.X + mat[7]*p.Y + mat[8]*p.Z
		p.X, p.Y, p.Z = nX, nY, nZ
	}
	return p
}

func allAxisRotations(f *Field, mat []int) []*Field {
	out := []*Field{f}
	for i := 0; i < 3; i++ {
		f = f.Rotate(mat, 90)
		out = append(out, f)
	}
	return out
}

func allOrientations(b *Field) []*Field {
	out := []*Field{}

	// starts facing x [0-3]
	out = append(out, allAxisRotations(b, RX)...)

	// face -x, rotate around x [4-7]
	b = b.Rotate(RZ, 180)
	out = append(out, allAxisRotations(b, RX)...)

	// face a z, rotate around z [8-11]
	b = b.Rotate(RY, 90)
	out = append(out, allAxisRotations(b, RZ)...)

	// face other z, rotate around z [12-15]
	b = b.Rotate(RY, 180)
	out = append(out, allAxisRotations(b, RZ)...)

	// face a y, rotate around y [16-19]
	b = b.Rotate(RX, 90)
	out = append(out, allAxisRotations(b, RY)...)

	// face other y, rotate around y [20-23]
	b = b.Rotate(RX, 180)
	out = append(out, allAxisRotations(b, RY)...)

	return out
}

func isMatchWithDelta(ref, cand *Field, refB, candB pos.P3, delta pos.P3) *Field {
	candsInRefSpaceMap := map[pos.P3]bool{}
	candsInRefSpace := []pos.P3{}
	for _, c := range cand.Beacons {
		o := c
		o.Add(delta)
		candsInRefSpaceMap[o] = true
		candsInRefSpace = append(candsInRefSpace, o)
	}

	matches := []pos.P3{}
	for _, r := range ref.Beacons {
		if _, found := candsInRefSpaceMap[r]; found {
			matches = append(matches, r)

			// we test them all (i.e. don't early exit
			// here) to ease testing
		}
	}

	if len(matches) < 12 {
		return nil
	}

	nsp := cand.ScannerPos
	nsp.Add(delta)

	return NewField(cand.ScannerNum, nsp, candsInRefSpace)
}

func matchOrientedFields(ref, cand *Field) *Field {
	for _, candB := range cand.Beacons {
		for _, refB := range ref.Beacons {
			// If we equate candB and refB, do we have 12 matches?

			delta := pos.P3{
				X: refB.X - candB.X,
				Y: refB.Y - candB.Y,
				Z: refB.Z - candB.Z,
			}

			if adjB := isMatchWithDelta(ref, cand, refB, candB, delta); adjB != nil {
				return adjB
			}
		}
	}

	return nil
}

func matchFields(a, b *Field) (adjB *Field) {
	for _, rotB := range allOrientations(b) {
		if adjB = matchOrientedFields(a, rotB); adjB != nil {
			return
		}
	}

	return nil
}

func mapSpace(scanners map[int]*Field) {
	adjusted := map[int]bool{0: true}

	for len(adjusted) != len(scanners) {
		for candNum, cand := range scanners {
			if _, found := adjusted[candNum]; found {
				continue // already adjusted
			}

			for refNum, _ := range adjusted {
				ref := scanners[refNum]
				adjCand := matchFields(ref, cand)
				if adjCand != nil {
					scanners[candNum] = adjCand
					adjusted[candNum] = true

					logger.LogF("match ref %v to cand %v (%v left)",
						refNum, candNum, len(scanners)-len(adjusted))
					break
				}
			}
		}
	}
}

func solveA(scanners map[int]*Field) {
	uniques := map[pos.P3]bool{}
	for _, s := range scanners {
		for _, b := range s.Beacons {
			uniques[b] = true
		}
	}

	fmt.Println("A", len(uniques))
}

func solveB(scanners map[int]*Field) {
	maxDist := 0

	for _, s1 := range scanners {
		for _, s2 := range scanners {
			if s1.ScannerNum == s2.ScannerNum {
				continue
			}

			dist := s1.ScannerPos.ManhattanDistance(s2.ScannerPos)
			if dist > maxDist {
				maxDist = dist
			}
		}
	}

	fmt.Println("B", maxDist)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	scanners, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	mapSpace(scanners)

	solveA(scanners)
	solveB(scanners)
}
