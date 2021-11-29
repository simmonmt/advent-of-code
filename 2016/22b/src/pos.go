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

import "intmath"

type Pos struct {
	X, Y int
}

func (p Pos) Eq(o Pos) bool {
	return p.X == o.X && p.Y == o.Y
}

func (p Pos) Dist(o Pos) int {
	return intmath.Abs(p.X-o.X) + intmath.Abs(p.Y-o.Y)
}
