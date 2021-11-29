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

type Pos struct {
	X, Y int
}

func PosLess(a, b Pos) bool {
	if a.Y != b.Y {
		return a.Y < b.Y
	}
	return a.X < b.X
}

type PosByReadingOrder []Pos

func (a PosByReadingOrder) Len() int           { return len(a) }
func (a PosByReadingOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PosByReadingOrder) Less(i, j int) bool { return PosLess(a[i], a[j]) }
