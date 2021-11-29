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

package puzzle

import (
	"github.com/simmonmt/aoc/2019/common/intmath"
)

type Pos struct {
	X, Y int
}

type ByManhattanOriginDistance []Pos

func (a ByManhattanOriginDistance) Len() int      { return len(a) }
func (a ByManhattanOriginDistance) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByManhattanOriginDistance) Less(i, j int) bool {
	iDist := intmath.Abs(a[i].X) + intmath.Abs(a[i].Y)
	jDist := intmath.Abs(a[j].X) + intmath.Abs(a[j].Y)
	return iDist < jDist
}
