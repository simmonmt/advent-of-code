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

package lib

import (
	"fmt"
	"strconv"
	"strings"
)

type Pos struct {
	X, Y int
}

func (p Pos) Eq(s Pos) bool {
	return p.X == s.X && p.Y == s.Y
}

func PosFromString(str string) (Pos, error) {
	parts := strings.SplitN(str, ",", 2)
	if len(parts) != 2 {
		return Pos{}, fmt.Errorf("invalid pos %v", str)
	}

	x, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return Pos{}, fmt.Errorf("bad x in %v", str)
	}

	y, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return Pos{}, fmt.Errorf("bad y in %v", str)
	}

	return Pos{int(x), int(y)}, nil
}
