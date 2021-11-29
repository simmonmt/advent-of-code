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

import "fmt"

type Char struct {
	Num   int
	IsElf bool
	P     Pos
	HP    int
	AP    int
}

func NewChar(num int, isElf bool, pos Pos, ap int) *Char {
	return &Char{
		Num:   num,
		IsElf: isElf,
		P:     pos,
		HP:    200,
		AP:    ap,
	}
}

func (c Char) String() string {
	t := "Elf"
	if !c.IsElf {
		t = "Gob"
	}

	return fmt.Sprintf("#%d: %s HP:%3d AP:%d %+v", c.Num, t, c.HP, c.AP, c.P)
}

func (c Char) Short() rune {
	if c.IsElf {
		return 'E'
	}
	return 'G'
}

type CharByReadingOrder []Char

func (a CharByReadingOrder) Len() int      { return len(a) }
func (a CharByReadingOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a CharByReadingOrder) Less(i, j int) bool {
	if a[i].P.Y != a[j].P.Y {
		return a[i].P.Y < a[j].P.Y
	}
	return a[i].P.X < a[j].P.X
}
