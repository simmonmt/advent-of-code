// Copyright 2022 Google LLC
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
	_ "embed"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/simmonmt/aoc/2022/common/dir"
	"github.com/simmonmt/aoc/2022/common/logger"
	"github.com/simmonmt/aoc/2022/common/pos"
)

var (
	//go:embed sample.txt
	rawSample string
	//go:embed sample2.txt
	rawSample2 string

	sampleLines, sample2Lines []string
)

func encodingOnly(lines []string) map[dir.Dir][]*big.Int {
	g, _, _ := parseInput(lines)
	return encodeGrid(g)
}

func makeEncoding(in map[dir.Dir][]string) map[dir.Dir][]*big.Int {
	enc := map[dir.Dir][]*big.Int{}
	for d, a := range in {
		enc[d] = make([]*big.Int, len(a))
		for i, s := range a {
			enc[d][i] = &big.Int{}
			_, err := fmt.Sscanf(s, "%x", enc[d][i])
			if err != nil {
				panic(fmt.Sprintf("bad read from '%s': %v", s, err))
			}
		}
	}
	return enc
}

func assertEncodingsAreEqual(t *testing.T, got, want map[dir.Dir][]*big.Int) {
	t.Helper()
	for _, d := range dir.AllDirs {
		if !reflect.DeepEqual(got[d], want[d]) {
			t.Errorf("%v: got %v, want %v", d, got[d], want[d])
		}
	}
}

func assertEncodingsAreDifferent(t *testing.T, got, want map[dir.Dir][]*big.Int) {
	t.Helper()

	for _, d := range dir.AllDirs {
		if !reflect.DeepEqual(got[d], want[d]) {
			return
		}
	}
	t.Errorf("got %v wanted different", got)
}

func TestParseAndEncodeGrid(t *testing.T) {
	type TestCase struct {
		in   []string
		want map[dir.Dir][]*big.Int
	}

	testCases := []TestCase{
		TestCase{
			in: sampleLines,
			want: makeEncoding(map[dir.Dir][]string{
				dir.DIR_NORTH: []string{"0", "0", "0", "0", "0"},
				dir.DIR_SOUTH: []string{"0", "0", "0", "8", "0"},
				dir.DIR_WEST:  []string{"0", "0", "0", "0", "0"},
				dir.DIR_EAST:  []string{"0", "1", "0", "0", "0"},
			}),
		},
		TestCase{
			in: sample2Lines,
			want: makeEncoding(map[dir.Dir][]string{
				dir.DIR_NORTH: []string{"0", "8", "0", "8", "9", "0"},
				dir.DIR_SOUTH: []string{"0", "4", "8", "0", "0", "0"},
				dir.DIR_WEST:  []string{"28", "32", "10", "1"},
				dir.DIR_EAST:  []string{"3", "0", "29", "20"},
			}),
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			g, _, _ := parseInput(tc.in)
			got := encodeGrid(g)

			assertEncodingsAreEqual(t, got, tc.want)
		})
	}
}

func TestEncodingToGrid(t *testing.T) {
	type TestCase struct {
		in   map[dir.Dir][]*big.Int
		want []string
	}

	testCases := []TestCase{
		TestCase{
			in: makeEncoding(map[dir.Dir][]string{
				dir.DIR_NORTH: []string{"0", "0", "0", "0", "0", "0"},
				dir.DIR_SOUTH: []string{"0", "0", "0", "0", "0", "0"},
				dir.DIR_WEST:  []string{"0", "0", "0", "0"},
				dir.DIR_EAST:  []string{"0", "0", "0", "0"},
			}),
			want: []string{
				"......",
				"......",
				"......",
				"......",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var sb strings.Builder
			dumpEncodedTo(tc.in, &sb)
			got := strings.Split(sb.String(), "\n")
			got = got[0 : len(got)-1] // there was a trailing \n

			if !reflect.DeepEqual(got, tc.want) {
				fmt.Printf("test %d: got:\n  %v\n", i, strings.Join(got, "\n  "))
				fmt.Printf("test %d: want:\n  %v\n", i, strings.Join(tc.want, "\n  "))
				t.Errorf("output mismatch")
			}
		})
	}
}

func TestAdvanceEncoding(t *testing.T) {
	type TestCase struct {
		in   map[dir.Dir][]*big.Int
		want map[dir.Dir][]string
	}

	testCases := []TestCase{
		TestCase{
			in: encodingOnly(sampleLines),
			want: map[dir.Dir][]string{
				dir.DIR_NORTH: []string{"0", "0", "0", "0", "0"},
				dir.DIR_SOUTH: []string{"0", "0", "0", "10", "0"},
				dir.DIR_WEST:  []string{"0", "0", "0", "0", "0"},
				dir.DIR_EAST:  []string{"0", "2", "0", "0", "0"},
			},
		},
		TestCase{
			in: encodingOnly(sample2Lines),
			want: map[dir.Dir][]string{
				// ....         / ...^ => ..^. / .... /
				// ...^ => ..^. / ^..^ => ..^^ / ....
				dir.DIR_NORTH: []string{"0", "4", "0", "4", "c", "0"},
				// ....         / ..v. => ...v / ...v => v... /
				// ....         / ....         / ....
				dir.DIR_SOUTH: []string{"0", "8", "1", "0", "0", "0"},
				// ...<.< => ..<.<. / .<..<< => <..<<. /
				// ....<. => ...<.. / <..... => .....<
				dir.DIR_WEST: []string{"14", "19", "8", "20"},
				// >>.... => .>>... / ......           /
				// >..>.> => >>..>. / .....> => >.....
				dir.DIR_EAST: []string{"6", "0", "13", "1"},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := advanceEncoding(tc.in)
			for _, d := range dir.AllDirs {
				gotArr := []string{}
				for _, n := range got[d] {
					gotArr = append(gotArr, fmt.Sprintf("%x", n))
				}

				if !reflect.DeepEqual(gotArr, tc.want[d]) {
					t.Errorf("%v: got %v, want %v", d, gotArr, tc.want[d])
				}
			}
		})
	}
}

func TestMakeAllEncodings(t *testing.T) {
	type TestCase struct {
		in      map[dir.Dir][]*big.Int
		wantNum int
	}

	testCases := []TestCase{
		TestCase{
			in:      encodingOnly(sampleLines),
			wantNum: 5,
		},
		TestCase{
			in:      encodingOnly(sample2Lines),
			wantNum: 12,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			encs := makeAllEncodings(tc.in)
			if gotNum := len(encs); gotNum != tc.wantNum {
				t.Errorf("#encoding = %v, want %v", gotNum, tc.wantNum)
			}

			// for _, enc := range encs {
			// 	dumpEncoded(enc)
			// 	fmt.Println()
			// }

			assertEncodingsAreEqual(t, encs[0], tc.in)
			assertEncodingsAreDifferent(t, encs[0], encs[len(encs)-1])

			nextEnc := advanceEncoding(encs[len(encs)-1])
			assertEncodingsAreEqual(t, encs[0], nextEnc)
		})
	}
}

// func TestNewBoardStep(t *testing.T) {
// 	enc := encodingOnly(sample2Lines)
// 	startX, endX := 0, 5
// 	id := 1

// 	want := &BoardStep{
// 		id:  id,
// 		enc: enc,
// 		neighborsOfOpens: map[pos.P2][]pos.P2{
// 			pos.P2{0, -1}: []pos.P2{},
// 			pos.P2{2, 0}:  []pos.P2{pos.P2{2, 1}},
// 			pos.P2{0, 1}:  []pos.P2{},
// 			pos.P2{2, 1}:  []pos.P2{pos.P2{3, 1}, pos.P2{2, 0}, pos.P2{2, 2}},
// 			pos.P2{3, 1}:  []pos.P2{pos.P2{2, 1}},
// 			pos.P2{2, 2}:  []pos.P2{pos.P2{2, 1}},
// 		},
// 	}

// 	got := NewBoardStep(id, startX, endX, enc)

// 	if !reflect.DeepEqual(got, want) {
// 		for o, wantNeighbors := range want.neighborsOfOpens {
// 			if !reflect.DeepEqual(got.neighborsOfOpens[o], wantNeighbors) {
// 				t.Errorf("%v: got %v, want %v",
// 					o, got.neighborsOfOpens[o], wantNeighbors)
// 			}
// 		}

// 		t.Errorf("got %+v, want %+v", got, want)
// 	}
// }

// func TestNewBoardStepEntryAndExit(t *testing.T) {
// 	// Move to minute 2
// 	enc := advanceEncoding(advanceEncoding(encodingOnly(sample2Lines)))
// 	startX, endX := 0, 5
// 	id := 1

// 	// We'll only check for the ones in these map, ignoring all others.
// 	wantNeighborsOfOpens := map[pos.P2][]pos.P2{
// 		pos.P2{0, -1}: []pos.P2{pos.P2{0, 0}},
// 		pos.P2{5, 3}:  []pos.P2{pos.P2{5, 4}},
// 	}

// 	got := NewBoardStep(id, startX, endX, enc)

// 	for o, wantNeighbors := range wantNeighborsOfOpens {
// 		if !reflect.DeepEqual(got.neighborsOfOpens[o], wantNeighbors) {
// 			t.Errorf("%v: got %v, want %v",
// 				o, got.neighborsOfOpens[o], wantNeighbors)
// 		}
// 	}
// }

func TestAStarNeighbors(t *testing.T) {
	allEncodings := makeAllEncodings(encodingOnly(sample2Lines))

	steps := make([]*BoardStep, len(allEncodings))
	for i, enc := range allEncodings {
		steps[i] = NewBoardStep(i, enc)
	}

	client := &astarClient{
		startPos: pos.P2{0, -1},
		endPos:   pos.P2{5, 4},
		steps:    steps,
	}

	type TestCase struct {
		cur  string
		want []string
	}

	testCases := []TestCase{
		TestCase{"0/0,-1", []string{"1/0,-1", "1/0,0"}},
		TestCase{"1/0,0", []string{"2/0,0", "2/0,1"}},
		TestCase{"2/0,1", []string{"3/0,1"}},
		TestCase{"3/0,1", []string{"4/0,0"}},
		TestCase{"4/5,3", []string{"5/5,3", "5/5,4"}},
	}

	for _, tc := range testCases {
		t.Run(tc.cur, func(t *testing.T) {
			got := client.AllNeighbors(tc.cur)
			sort.Strings(got)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("AllNeighbors('%v') = %v, want %v",
					tc.cur, got, tc.want)
			}
		})
	}
}

func TestSolveA(t *testing.T) {
	g, startX, endX := parseInput(sample2Lines)

	if got, want := solveA(g, startX, endX), 18; got != want {
		t.Errorf("solveA(sample2) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	g, startX, endX := parseInput(sample2Lines)

	if got, want := solveB(g, startX, endX), 54; got != want {
		t.Errorf("solveB(sample2) = %v, want %v", got, want)
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)

	sampleLines = strings.Split(rawSample, "\n")
	if len(sampleLines) > 0 && sampleLines[len(sampleLines)-1] == "" {
		sampleLines = sampleLines[0 : len(sampleLines)-1]
	}
	sample2Lines = strings.Split(rawSample2, "\n")
	if len(sample2Lines) > 0 && sample2Lines[len(sample2Lines)-1] == "" {
		sample2Lines = sample2Lines[0 : len(sample2Lines)-1]
	}

	os.Exit(m.Run())
}
