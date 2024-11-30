// Copyright 2024 Google LLC
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

package pos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2024/common/mtsmath"
)

func fromString(str string, wantParts int) ([]int, error) {
	parts := strings.Split(str, ",")
	if len(parts) != wantParts {
		return nil, fmt.Errorf("invalid input")
	}
	vs := make([]int, wantParts)
	for i := range vs {
		v, err := strconv.Atoi(parts[i])
		if err != nil {
			return nil, fmt.Errorf("invalid coord %v", parts[i])
		}
		vs[i] = v
	}
	return vs, nil
}

type P2 struct {
	X, Y int
}

func P2FromString(str string) (P2, error) {
	vs, err := fromString(str, 2)
	if err != nil {
		return P2{}, err
	}
	return P2{X: vs[0], Y: vs[1]}, nil
}

func (p *P2) Equals(o P2) bool {
	return p.X == o.X && p.Y == o.Y
}

func (p *P2) Add(o P2) {
	p.X += o.X
	p.Y += o.Y
}

func (p *P2) LessThan(o P2) bool {
	if p.X < o.X {
		return true
	} else if p.X > o.X {
		return false
	} else {
		return p.Y < o.Y
	}
}

func (p P2) ManhattanDistance(o P2) int {
	return mtsmath.Abs(o.X-p.X) + mtsmath.Abs(o.Y-p.Y)
}

func (p P2) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

func (p P2) AllNeighbors(includeDiag bool) []P2 {
	num := 4
	if includeDiag {
		num = 8
	}

	out := make([]P2, num)
	out[0] = P2{p.X - 1, p.Y}
	out[1] = P2{p.X + 1, p.Y}
	out[2] = P2{p.X, p.Y - 1}
	out[3] = P2{p.X, p.Y + 1}

	if includeDiag {
		out[4] = P2{p.X - 1, p.Y - 1}
		out[5] = P2{p.X + 1, p.Y - 1}
		out[6] = P2{p.X - 1, p.Y + 1}
		out[7] = P2{p.X + 1, p.Y + 1}
	}

	return out
}

func WalkP2(numY, numX int, cb func(p P2)) {
	for y := 0; y < numY; y++ {
		for x := 0; x < numX; x++ {
			cb(P2{Y: y, X: x})
		}
	}
}

type P3 struct {
	X, Y, Z int
}

func (p *P3) Equals(o P3) bool {
	return p.X == o.X && p.Y == o.Y && p.Z == o.Z
}

func (p P3) LessThan(o P3) bool {
	if o.X != p.X {
		return o.X < p.X
	}
	if o.Y != p.Y {
		return o.Y < p.Y
	}
	return o.Z < p.Z
}

func P3FromString(str string) (P3, error) {
	vs, err := fromString(str, 3)
	if err != nil {
		return P3{}, err
	}
	return P3{vs[0], vs[1], vs[2]}, nil
}

func (p P3) String() string {
	return fmt.Sprintf("%d,%d,%d", p.X, p.Y, p.Z)
}

func (p P3) AllNeighbors(includeDiag bool) []P3 {
	num := 6
	if includeDiag {
		num = 26
	}

	out := make([]P3, num)
	out[0] = P3{p.X - 1, p.Y, p.Z}
	out[1] = P3{p.X + 1, p.Y, p.Z}
	out[2] = P3{p.X, p.Y + 1, p.Z}
	out[3] = P3{p.X, p.Y - 1, p.Z}
	out[4] = P3{p.X, p.Y, p.Z + 1}
	out[5] = P3{p.X, p.Y, p.Z - 1}

	if includeDiag {
		out[6] = P3{p.X + 1, p.Y - 1, p.Z}
		out[7] = P3{p.X + 1, p.Y + 1, p.Z}
		out[8] = P3{p.X + 1, p.Y - 1, p.Z - 1}
		out[9] = P3{p.X + 1, p.Y, p.Z - 1}
		out[10] = P3{p.X + 1, p.Y + 1, p.Z - 1}
		out[11] = P3{p.X + 1, p.Y - 1, p.Z + 1}
		out[12] = P3{p.X + 1, p.Y, p.Z + 1}
		out[13] = P3{p.X + 1, p.Y + 1, p.Z + 1}

		out[14] = P3{p.X - 1, p.Y - 1, p.Z}
		out[15] = P3{p.X - 1, p.Y + 1, p.Z}
		out[16] = P3{p.X - 1, p.Y - 1, p.Z - 1}
		out[17] = P3{p.X - 1, p.Y, p.Z - 1}
		out[18] = P3{p.X - 1, p.Y + 1, p.Z - 1}
		out[19] = P3{p.X - 1, p.Y - 1, p.Z + 1}
		out[20] = P3{p.X - 1, p.Y, p.Z + 1}
		out[21] = P3{p.X - 1, p.Y + 1, p.Z + 1}

		out[22] = P3{p.X, p.Y - 1, p.Z - 1}
		out[23] = P3{p.X, p.Y - 1, p.Z + 1}
		out[24] = P3{p.X, p.Y + 1, p.Z - 1}
		out[25] = P3{p.X, p.Y + 1, p.Z + 1}
	}

	return out
}

type P4 struct {
	X, Y, Z, W int
}

func (p P4) AllNeighbors() []P4 {
	out := make([]P4, 80)

	span := []int{-1, 0, 1}
	i := 0

	for _, w := range span {
		for _, z := range span {
			for _, y := range span {
				for _, x := range span {
					if w == 0 && z == 0 && y == 0 && x == 0 {
						continue
					}
					out[i] = P4{
						X: p.X + x,
						Y: p.Y + y,
						Z: p.Z + z,
						W: p.W + w,
					}
					i++
				}
			}
		}
	}

	return out
}
