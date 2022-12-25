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
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/simmonmt/aoc/2022/common/dir"
	"github.com/simmonmt/aoc/2022/common/logger"
	"github.com/simmonmt/aoc/2022/common/pos"
)

var (
	//go:embed sample.txt
	rawSample   string
	sampleLines []string

	//go:embed input.txt
	rawInput   string
	inputLines []string

	sampleFoldingSpec = "..1./234./..56"
)

func TestCommands(t *testing.T) {
	type TestCase struct {
		in   string
		want []Command
	}

	testCases := []TestCase{
		TestCase{"10R5L5R", []Command{10, TURN_RIGHT, 5, TURN_LEFT, 5, TURN_RIGHT}},
		TestCase{"10R5", []Command{10, TURN_RIGHT, 5}},
		TestCase{"", nil},
		TestCase{"foo", nil},
		TestCase{"10R5Z", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			wantErr, wantErrStr := false, "nil"
			if tc.want == nil {
				wantErr, wantErrStr = true, "err"
			}

			if got, err := parseCommands(tc.in); (err != nil) == wantErr && !reflect.DeepEqual(got, tc.want) {
				t.Errorf("parseCommands('%v') = %v, %v, want %v, %v",
					tc.in, got, err, tc.want, wantErrStr)
			}
		})
	}
}

func TestFlatBoardMove(t *testing.T) {
	type TestCase struct {
		p        pos.P2
		d        dir.Dir
		wantPos  pos.P2
		wantCell CellType
	}

	testCases := []TestCase{
		TestCase{pos.P2{5, 4}, dir.DIR_NORTH, pos.P2{5, 7}, OPEN}, // D wrap
		TestCase{pos.P2{5, 7}, dir.DIR_SOUTH, pos.P2{5, 4}, OPEN}, // C wrap
		TestCase{pos.P2{11, 6}, dir.DIR_EAST, pos.P2{0, 6}, OPEN}, // A wrap
		TestCase{pos.P2{0, 6}, dir.DIR_WEST, pos.P2{11, 6}, OPEN}, // B wrap
		TestCase{pos.P2{0, 4}, dir.DIR_WEST, pos.P2{11, 4}, WALL},
		TestCase{pos.P2{0, 6}, dir.DIR_EAST, pos.P2{1, 6}, OPEN},  // east nonwrap
		TestCase{pos.P2{0, 6}, dir.DIR_NORTH, pos.P2{0, 5}, OPEN}, // north nonwrap
		TestCase{pos.P2{0, 6}, dir.DIR_SOUTH, pos.P2{0, 7}, OPEN}, // south nonwrap
		TestCase{pos.P2{5, 7}, dir.DIR_WEST, pos.P2{4, 7}, OPEN},  // west nonwrap
	}

	g, _, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}
	b := NewFlatBoard(g)

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if gotPos, gotCell := b.Move(tc.p, tc.d); !gotPos.Equals(tc.wantPos) || gotCell != tc.wantCell {
				t.Errorf("b.Move(%v, %v) = %v, %v, want %v, %v",
					tc.p, tc.d, gotPos, gotCell, tc.wantPos, tc.wantCell)
			}
		})
	}
}

func TestCubeBoardMove(t *testing.T) {
	type TestCase struct {
		relPos     pos.P2
		inFace     int
		inDir      dir.Dir
		wantRelPos pos.P2
		wantFace   int
		wantDir    dir.Dir
		wantCell   CellType
	}

	g, _, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	b := NewCubeBoard(g)

	testCases := []TestCase{
		TestCase{ // within 1 going up
			pos.P2{1, 2}, 1, dir.DIR_NORTH,
			pos.P2{1, 1}, 1, dir.DIR_NORTH, OPEN},
		TestCase{ // within 1 going to the side
			pos.P2{1, 2}, 1, dir.DIR_EAST,
			pos.P2{2, 2}, 1, dir.DIR_EAST, OPEN},
		TestCase{ // north off 1 => north on 4
			pos.P2{0, 0}, 1, dir.DIR_NORTH,
			pos.P2{0, 3}, 4, dir.DIR_NORTH, OPEN},
		TestCase{ // south off 1 => south on 2
			pos.P2{0, 3}, 1, dir.DIR_SOUTH,
			pos.P2{0, 0}, 2, dir.DIR_SOUTH, OPEN},
		TestCase{ // east off 1 => south on 5
			pos.P2{3, 0}, 1, dir.DIR_EAST,
			pos.P2{3, 0}, 5, dir.DIR_SOUTH, OPEN},
		TestCase{ // west off 6 => north on 5
			pos.P2{0, 0}, 6, dir.DIR_WEST,
			pos.P2{3, 3}, 5, dir.DIR_NORTH, OPEN},
		TestCase{ // west off 1 => west on 3
			pos.P2{0, 0}, 1, dir.DIR_WEST,
			pos.P2{3, 0}, 3, dir.DIR_WEST, OPEN},

		// The examples from part 2 with sample data
		TestCase{ // A east to B (south)
			pos.P2{3, 1}, 1, dir.DIR_EAST,
			pos.P2{2, 0}, 5, dir.DIR_SOUTH, OPEN},
		TestCase{ // C south to D (north)
			pos.P2{2, 3}, 2, dir.DIR_SOUTH,
			pos.P2{1, 3}, 6, dir.DIR_NORTH, OPEN},
		TestCase{ // E north to wall (east)
			pos.P2{2, 0}, 3, dir.DIR_NORTH,
			pos.P2{0, 2}, 4, dir.DIR_EAST, WALL},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			fmt.Printf("%d: %+v\n", i, tc)

			gotRelPos, gotFace, gotDir, gotCell := b.Move(
				tc.relPos, b.Face(tc.inFace), tc.inDir)

			gotFaceNum := -1
			if gotFace != nil {
				gotFaceNum = gotFace.num
			}
			if !gotRelPos.Equals(tc.wantRelPos) || gotFace.num != tc.wantFace || gotDir != tc.wantDir || gotCell != tc.wantCell {
				t.Errorf("b.Move(%v, %v, %v) = %v, %v, %v, %v, want %v, %v, %v, %v",
					tc.relPos, tc.inFace, tc.inDir,
					gotRelPos, gotFaceNum, gotDir, gotCell,
					tc.wantRelPos, tc.wantFace, tc.wantDir, tc.wantCell)
			}

		})
	}
}

func TestCubeBoardAbsPos(t *testing.T) {
	type TestCase struct {
		rel     pos.P2
		relFace int
		want    pos.P2
	}

	testCases := []TestCase{
		TestCase{pos.P2{0, 0}, 1, pos.P2{8, 4}},
		TestCase{pos.P2{0, 0}, 2, pos.P2{8, 8}},
		TestCase{pos.P2{0, 0}, 3, pos.P2{4, 4}},
		TestCase{pos.P2{0, 0}, 4, pos.P2{8, 0}},
		TestCase{pos.P2{0, 0}, 5, pos.P2{12, 8}},
		TestCase{pos.P2{0, 0}, 6, pos.P2{0, 4}},
	}

	g, _, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	b := NewCubeBoard(g)

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v:%d", tc.rel, tc.relFace), func(t *testing.T) {
			if got := b.AbsPos(tc.rel, b.Face(tc.relFace)); !got.Equals(tc.want) {
				t.Errorf("AbsPos(%v, %v) = %v, _, want %v, _",
					tc.rel, tc.relFace, got, tc.want)
			}
		})
	}
}

func TestDetectCubeFaceSize(t *testing.T) {
	g, _, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := detectCubeFaceSize(g), 4; got != want {
		t.Errorf("detectCubeSize() = %v, want %v", got, want)
	}
}

func TestMakeCubeSpec(t *testing.T) {
	g, _, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	want := [][]bool{
		[]bool{false, false, true, false},
		[]bool{true, true, true, false},
		[]bool{false, false, true, true},
	}

	if got := makeCubeSpec(g, 4); !reflect.DeepEqual(got, want) {
		t.Errorf("makeCubeSpec = %v, want %v", got, want)
	}
}

func neighborMap(s string) map[dir.Dir]int {
	out := map[dir.Dir]int{}
	var d dir.Dir
	for i := 0; i < 8; i++ {
		if i%2 == 0 {
			d = dir.Parse(string(s[i]))
		} else {
			num, err := strconv.Atoi(string(s[i]))
			if err != nil {
				panic("bad num")
			}
			out[d] = num
		}
	}
	return out
}

func TestFindFaceRotations(t *testing.T) {
	type TestCase struct {
		in   []string
		want map[int]*FaceInfo
	}

	testCases := []TestCase{
		TestCase{
			sampleLines,
			map[int]*FaceInfo{
				1: &FaceInfo{1, pos.P2{2, 1}, neighborMap("N4S2W3E5"), 0},
				2: &FaceInfo{2, pos.P2{2, 2}, neighborMap("N1S6W3E5"), 0},
				3: &FaceInfo{3, pos.P2{1, 1}, neighborMap("N4S2W6E1"), 3},
				4: &FaceInfo{4, pos.P2{2, 0}, neighborMap("N6S1W3E5"), 2},
				5: &FaceInfo{5, pos.P2{3, 2}, neighborMap("N1S6W2E4"), 0},
				6: &FaceInfo{6, pos.P2{0, 1}, neighborMap("N4S2W5E3"), 2},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			g, _, err := parseInput(tc.in)
			if err != nil {
				t.Fatal(err)
			}

			spec := makeCubeSpec(g, detectCubeFaceSize(g))
			posFaces := findCubeFaces(spec)
			faces := findFaceRotations(posFaces)

			if len(faces) != 6 {
				t.Fatalf("want 6 faces found %v: %v",
					len(faces), faces)
			}

			for num, wantFace := range tc.want {
				gotFace := faces[num]
				if !reflect.DeepEqual(gotFace, wantFace) {
					t.Errorf("face %v: got %v, want %v",
						num, gotFace, wantFace)
				}
			}
		})
	}
}

func TestSolveA(t *testing.T) {
	board, cmds, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(board, cmds), 6032; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	g, cmds, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(g, cmds), 5031; got != want {
		t.Errorf("solveB(sample) = %v, want %v", got, want)
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)

	sampleLines = strings.Split(rawSample, "\n")
	if len(sampleLines) > 0 && sampleLines[len(sampleLines)-1] == "" {
		sampleLines = sampleLines[0 : len(sampleLines)-1]
	}

	inputLines = strings.Split(rawInput, "\n")
	if len(inputLines) > 0 && inputLines[len(inputLines)-1] == "" {
		inputLines = inputLines[0 : len(inputLines)-1]
	}

	os.Exit(m.Run())
}
