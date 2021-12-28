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

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/simmonmt/aoc/2021/common/logger"
	"github.com/simmonmt/aoc/2021/common/pos"
)

func TestGameStateSerialize(t *testing.T) {
	gs := NewGameState([8]pos.P2{
		pos.P2{1, 2}, pos.P2{3, 4}, pos.P2{5, 6}, pos.P2{7, 8},
		pos.P2{10, 20}, pos.P2{30, 40}, pos.P2{50, 60}, pos.P2{70, 80},
	})

	ser := gs.Serialize()
	want := "1,2|3,4|5,6|7,8|10,20|30,40|50,60|70,80"
	if ser != want {
		t.Errorf("Serialize = %v, want %v", ser, want)
	}

	gs2, err := DeserializeGameState(ser)
	if err != nil || !reflect.DeepEqual(gs, gs2) {
		t.Errorf("Deserialize = %+v, %v, want %+v, nil",
			gs2, err, gs)
	}
}

func TestGoalReached(t *testing.T) {
	type TestCase struct {
		in      []string
		reached bool
	}

	testCases := []TestCase{
		TestCase{
			in: []string{
				"#############",
				"#...........#",
				"###B#C#B#D###",
				"  #A#D#C#A#  ",
				"  #########  ",
			},
			reached: false,
		},
		TestCase{
			in: []string{
				"#############",
				"#.A.......C.#",
				"###.#B#C#D###",
				"  #A#B#.#D#  ",
				"  #########  ",
			},
			reached: false,
		},
		TestCase{
			in: []string{
				"#############",
				"#...........#",
				"###A#B#C#D###",
				"  #A#B#C#D#  ",
				"  #########  ",
			},
			reached: true,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b, gs := parseInput(tc.in)
			ser := gs.Serialize()

			c := &astarClient{b: b}
			if got := c.GoalReached(ser, ser); got != tc.reached {
				t.Errorf("GoalReached(_,_) = %v, want %v", got, tc.reached)
			}
		})
	}
}

func TestAllNeighbors(t *testing.T) {
	type TestCase struct {
		seq []string
	}

	testCases := []TestCase{
		TestCase{
			seq: []string{
				//A1  A2  B1  B2  C1  C2  D1  D2
				"2,2|8,2|2,1|6,1|4,1|6,2|8,1|4,2",
				"2,2|8,2|2,1|3,0|4,1|6,2|8,1|4,2",
				"2,2|8,2|2,1|3,0|6,1|6,2|8,1|4,2",
				"2,2|8,2|2,1|3,0|6,1|6,2|8,1|5,0",
				"2,2|8,2|2,1|4,2|6,1|6,2|8,1|5,0",
				"2,2|8,2|4,1|4,2|6,1|6,2|8,1|5,0",
				"2,2|8,2|4,1|4,2|6,1|6,2|7,0|5,0",
				"2,2|9,0|4,1|4,2|6,1|6,2|7,0|5,0",
				"2,2|9,0|4,1|4,2|6,1|6,2|8,2|5,0",
				"2,2|9,0|4,1|4,2|6,1|6,2|8,2|8,1",
				"2,2|2,1|4,1|4,2|6,1|6,2|8,2|8,1",
			},
		},
	}

	boardStr := []string{
		"#############",
		"#...........#",
		"###B#C#B#D###",
		"  #A#D#C#A#  ",
		"  #########  ",
	}

	for tcNum, tc := range testCases {
		b, _ := parseInput(boardStr)
		c := &astarClient{b: b}

		for i := 0; i < len(tc.seq)-1; i++ {
			t.Run(fmt.Sprintf("seq %v from %v to %v", tcNum, i, i+1),
				func(t *testing.T) {
					in := tc.seq[i]
					want := tc.seq[i+1]

					inGS, err := DeserializeGameState(in)
					if err != nil {
						t.Fatalf("failed to deserialize in")
					}
					logger.LogF("in %v", in)
					b.Dump(inGS)

					wantGS, err := DeserializeGameState(want)
					if err != nil {
						t.Fatalf("failed to deserialize want")
					}
					logger.LogF("want %v", want)
					b.Dump(wantGS)

					found := false
					for _, n := range c.AllNeighbors(in) {
						if n == want {
							found = true
							break
						}
					}

					if !found {
						t.Errorf("want not found\n\nstart:\n%v\nwant:\n%v",
							b.DumpToString(inGS),
							b.DumpToString(wantGS))
					}
				})
		}
	}
}

func TestNeighborDistance(t *testing.T) {
	boardStr := []string{
		"#############",
		"#...........#",
		"###B#C#B#D###",
		"  #A#D#C#A#  ",
		"  #########  ",
	}
	board, _ := parseInput(boardStr)

	type TestCase struct {
		n1, n2 string
		want   uint
	}

	testCases := []TestCase{
		TestCase{
			n1:   "2,1|2,2|4,1|4,2|6,1|6,2|9,0|8,2",
			n2:   "2,1|2,2|0,0|4,2|6,1|6,2|9,0|8,2",
			want: 50,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			n1, err := DeserializeGameState(tc.n1)
			if err != nil {
				t.Fatalf("bad n1: %v", err)
			}
			board.Dump(n1)

			n2, _ := DeserializeGameState(tc.n2)
			if err != nil {
				t.Fatalf("bad n2: %v", err)
			}
			board.Dump(n2)

			c := &astarClient{b: board}
			if got := c.NeighborDistance(tc.n1, tc.n2); got != tc.want {
				t.Errorf("dist = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
