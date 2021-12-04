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

package board

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

var (
	start = [5][5]int{
		[5]int{22, 13, 17, 11, 0},
		[5]int{8, 2, 23, 4, 24},
		[5]int{21, 9, 14, 16, 7},
		[5]int{6, 10, 3, 18, 5},
		[5]int{1, 12, 20, 15, 19},
	}
)

func TestBoard(t *testing.T) {
	b := New(start)

	expected := strings.Join([]string{
		" 22   13   17   11    0 ",
		"  8    2   23    4   24 ",
		" 21    9   14   16    7 ",
		"  6   10    3   18    5 ",
		"  1   12   20   15   19 ",
	}, "\n") + "\n"

	buf := bytes.Buffer{}
	b.DumpTo(&buf)
	if expected != buf.String() {
		t.Errorf("dump want\n%v\n, got\n%v\n", expected, buf.String())
	}

	if won := b.Mark(24); won {
		t.Errorf("unexpected win")
		return
	}
	if won := b.Mark(3); won {
		t.Errorf("unexpected win")
		return
	}

	expected = strings.Join([]string{
		" 22   13   17   11    0 ",
		"  8    2   23    4  *24*",
		" 21    9   14   16    7 ",
		"  6   10  * 3*  18    5 ",
		"  1   12   20   15   19 ",
	}, "\n") + "\n"

	buf.Reset()
	b.DumpTo(&buf)
	if expected != buf.String() {
		t.Errorf("2mark want\n%v\n, got\n%v\n", expected, buf.String())
	}
}

func TestWon(t *testing.T) {
	type TestCase struct {
		moves []int
	}

	testCases := []TestCase{
		TestCase{moves: []int{22, 13, 17, 11, 0}},
		TestCase{moves: []int{7, 16, 14, 9, 21}},
		TestCase{moves: []int{13, 2, 9, 10, 12}},
		TestCase{moves: []int{15, 18, 16, 4, 11}},
		TestCase{moves: []int{6, 11, 10, 4, 3, 16, 18, 15}},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b := New(start)
			for i, move := range tc.moves {
				expectedWin := i == len(tc.moves)-1
				if won := b.Mark(move); won != expectedWin {
					b.Dump()
					t.Errorf("move %d=%d, won=%v, want %v",
						i, move, won, expectedWin)
					return
				}
			}
		})
	}
}

func TestScore(t *testing.T) {
	b := New([5][5]int{
		[5]int{14, 21, 17, 24, 4},
		[5]int{10, 16, 15, 9, 19},
		[5]int{18, 8, 23, 26, 20},
		[5]int{22, 11, 13, 6, 5},
		[5]int{2, 0, 12, 3, 7},
	})

	moves := []int{7, 4, 9, 5, 11, 17, 23, 2, 0, 14, 21, 24}
	for i, move := range moves[0 : len(moves)-1] {
		if won := b.Mark(move); won {
			b.Dump()
			t.Errorf("move %d=%d, unexpected win", i, move)
			return
		}
	}

	lastMove := moves[len(moves)-1]
	if won := b.Mark(lastMove); !won {
		b.Dump()
		t.Errorf("last %d, unexpected nonwin", lastMove)
		return
	}

	expectedScore := 4512
	if got := b.Score(lastMove); got != expectedScore {
		t.Errorf("score = %v, want %v", got, expectedScore)
		return
	}
}
