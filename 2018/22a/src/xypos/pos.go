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

package xypos

import (
	"fmt"
	"strconv"
	"strings"
)

type Pos struct {
	X, Y int
}

func (p Pos) Eq(o Pos) bool {
	return p.X == o.X && p.Y == o.Y
}

func (p Pos) Before(o Pos) bool {
	if p.Y < o.Y {
		return true
	} else if p.Y > o.Y {
		return false
	} else {
		return p.X < o.X
	}
}

func Parse(str string) (Pos, error) {
	parts := strings.SplitN(str, ",", 2)
	if len(parts) != 2 {
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

	return Pos{x, y}, nil
}
