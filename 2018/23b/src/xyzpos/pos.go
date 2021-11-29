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

package xyzpos

import (
	"fmt"
	"strconv"
	"strings"

	"intmath"
)

type Pos struct {
	X, Y, Z int
}

func (p Pos) Eq(o Pos) bool {
	return p.X == o.X && p.Y == o.Y && p.Z == o.Z
}

func (p Pos) Dist(o Pos) int {
	return intmath.Abs(o.X-p.X) + intmath.Abs(o.Y-p.Y) +
		intmath.Abs(o.Z-p.Z)
}

func Parse(str string) (Pos, error) {
	parts := strings.SplitN(str, ",", 3)
	if len(parts) != 3 {
		return Pos{}, fmt.Errorf("bad pos %v", str)
	}

	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return Pos{}, fmt.Errorf("bad pos x: %v", err)
	}

	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return Pos{}, fmt.Errorf("bad pos y: %v", err)
	}

	z, err := strconv.Atoi(parts[2])
	if err != nil {
		return Pos{}, fmt.Errorf("bad pos z: %v", err)
	}

	return Pos{x, y, z}, nil
}
