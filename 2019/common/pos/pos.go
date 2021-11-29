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

package pos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2019/common/intmath"
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
	return P2{vs[0], vs[1]}, nil
}

func (p *P2) Equals(o P2) bool {
	return p.X == o.X && p.Y == o.Y
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

func (p *P2) ManhattanDistance(o P2) int {
	return intmath.Abs(o.X-p.X) + intmath.Abs(o.Y-p.Y)
}

func (p *P2) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

type P3 struct {
	X, Y, Z int
}

func (p *P3) Equals(o P3) bool {
	return p.X == o.X && p.Y == o.Y && p.Z == o.Z
}

func P3FromString(str string) (P3, error) {
	vs, err := fromString(str, 3)
	if err != nil {
		return P3{}, err
	}
	return P3{vs[0], vs[1], vs[2]}, nil
}
