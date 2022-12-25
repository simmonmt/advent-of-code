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
	"container/list"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2022/common/dir"
	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/grid"
	"github.com/simmonmt/aoc/2022/common/logger"
	"github.com/simmonmt/aoc/2022/common/mtsmath"
	"github.com/simmonmt/aoc/2022/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Command int

const (
	TURN_LEFT  Command = -1
	TURN_RIGHT Command = -2
)

func (c Command) String() string {
	switch {
	case c >= 0:
		return fmt.Sprintf("forward %d", c)
	case c == TURN_LEFT:
		return fmt.Sprintf("turn left")
	case c == TURN_RIGHT:
		return fmt.Sprintf("turn right")
	default:
		return "unknown"
	}
}

func parseCommands(s string) ([]Command, error) {
	cmds := []Command{}

	numStart := -1
	for i := 0; i < len(s); i++ {
		r := s[i]
		if r >= '0' && r <= '9' {
			if numStart == -1 {
				numStart = i
			}
			continue
		}

		if numStart != -1 {
			var numStr string = s[numStart:i]
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, fmt.Errorf(
					"bad number '%v': %v", numStr, err)
			}
			cmds = append(cmds, Command(num))
			numStart = -1
		}

		if r == 'L' {
			cmds = append(cmds, TURN_LEFT)
		} else if r == 'R' {
			cmds = append(cmds, TURN_RIGHT)
		} else {
			return nil, fmt.Errorf("bad dir %v", string(r))
		}
	}

	if numStart != -1 { // trailing number
		var numStr string = s[numStart:len(s)]
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, fmt.Errorf(
				"bad number '%v': %v", numStr, err)
		}
		cmds = append(cmds, Command(num))
	}

	return cmds, nil
}

type CellType int

const (
	BLANK CellType = iota
	OPEN
	WALL
)

func (t CellType) String() string {
	switch t {
	case BLANK:
		return " "
	case OPEN:
		return "."
	case WALL:
		return "#"
	default:
		return "?"
	}
}

func ParseCellType(r rune) (CellType, bool) {
	switch r {
	case ' ':
		return BLANK, true
	case '.':
		return OPEN, true
	case '#':
		return WALL, true
	default:
		return BLANK, false
	}
}

type Board struct {
	g *grid.Grid[CellType]
}

func (b *Board) Get(p pos.P2) (CellType, bool) {
	return b.g.Get(p)
}

type FlatBoard struct {
	Board
}

func NewFlatBoard(g *grid.Grid[CellType]) *FlatBoard {
	return &FlatBoard{Board{g: g}}
}

func (b *FlatBoard) FindStart() pos.P2 {
	for p := (pos.P2{0, 0}); ; p.X++ {
		v, _ := b.Get(p)
		if v == OPEN {
			return p
		}
	}
}

func (b *FlatBoard) Move(p pos.P2, d dir.Dir) (pos.P2, CellType) {
	if _, ok := b.Get(p); !ok {
		panic("Move called with bad pos")
	}

	p = d.From(p)
	v, ok := b.Get(p)
	if ok && v != BLANK {
		return p, v
	}

	firstNotBlank := func(p pos.P2, inc pos.P2) pos.P2 {
		for {
			if v, ok := b.Get(p); !ok {
				panic("not ok")
			} else if v != BLANK {
				return p
			}
			p.Add(inc)
		}
	}

	// not ok or blank. either way we have to wrap
	if d == dir.DIR_EAST { // wraps off right to x=0
		p = firstNotBlank(pos.P2{0, p.Y}, pos.P2{1, 0})
	} else if d == dir.DIR_WEST { // wraps off left to +x
		p = firstNotBlank(pos.P2{b.g.Width() - 1, p.Y}, pos.P2{-1, 0})
	} else if d == dir.DIR_NORTH { // wraps off top to +y
		p = firstNotBlank(pos.P2{p.X, b.g.Height() - 1}, pos.P2{0, -1})
	} else if d == dir.DIR_SOUTH { // wraps off bottom to y=0
		p = firstNotBlank(pos.P2{p.X, 0}, pos.P2{0, 1})
	} else {
		panic("bad dir")
	}

	v, _ = b.Get(p)
	return p, v
}

type CubeBoard struct {
	Board
	faceSize   int
	facesByNum map[int]*FaceInfo
	facesByPos map[pos.P2]*FaceInfo
}

func NewCubeBoard(g *grid.Grid[CellType]) *CubeBoard {
	faceSize := detectCubeFaceSize(g)
	cubeSpec := makeCubeSpec(g, faceSize)
	facesByNum := findFaceRotations(findCubeFaces(cubeSpec))

	facesByPos := map[pos.P2]*FaceInfo{}
	for _, f := range facesByNum {
		facesByPos[f.pos] = f
	}

	return &CubeBoard{
		Board:      Board{g: g},
		faceSize:   faceSize,
		facesByNum: facesByNum,
		facesByPos: facesByPos,
	}
}

func (b *CubeBoard) FindStart() (rp pos.P2, f *FaceInfo) {
	p := pos.P2{0, 0}
	for {
		v, _ := b.Get(p)
		if v == OPEN {
			break
		}
		p.X++
	}

	// We now have an absolute position. Turn it into a relative position on
	// a particular face.
	rp = pos.P2{
		p.X % b.faceSize,
		p.Y % b.faceSize,
	}

	f = b.facesByPos[pos.P2{p.X / b.faceSize, p.Y / b.faceSize}]
	return
}

func (b *CubeBoard) Face(n int) *FaceInfo {
	return b.facesByNum[n]
}

func inFace(p pos.P2, sz int) bool {
	return p.X >= 0 && p.Y >= 0 && p.X < sz && p.Y < sz
}

func (b *CubeBoard) AbsPos(rp pos.P2, f *FaceInfo) pos.P2 {
	return pos.P2{
		X: f.pos.X*b.faceSize + rp.X,
		Y: f.pos.Y*b.faceSize + rp.Y,
	}
}

func (b *CubeBoard) Move(curRel pos.P2, curFace *FaceInfo, moveDir dir.Dir) (pos.P2, *FaceInfo, dir.Dir, CellType) {
	newRel := moveDir.From(curRel)
	if inFace(newRel, b.faceSize) {
		v, _ := b.Get(b.AbsPos(newRel, curFace))
		return newRel, curFace, moveDir, v
	}

	newFace := b.facesByNum[curFace.neighbors[moveDir]]

	// We know the new face but now need to figure out a) the new direction
	// of travel on the new face and b) the new relative coordinates on the
	// new face.
	//
	// Assume we have sides AB which B rotated 90deg CCW relative to A. An
	// east exit from A means several things from B's point of view:
	//   1) We've entered from B's north (stored in entrySide)
	//   2) We've changed our direction of travel, as A's east is B's south
	//      (stored in newDir)
	//   3) Our relative position changes as A's NE/SW side is B's NW/NE
	//      side (stored in newRel).

	// First, calculate entrySide and newDir (1 and 2 above).
	var entrySide, newDir dir.Dir
	for d, num := range newFace.neighbors {
		if num == curFace.num {
			entrySide = d
			newDir = d.Reverse()
		}
	}

	// Updating newRel (point 3 above) is harder. The current value of
	// newRel has moved beyond the bounds of curFace, and can be thought of
	// as the position on a newFace that's the same orientation as curFace
	// but measured relative to A's 0,0. We start by adjusting newRel so
	// it's relative to 0,0 on a newFace that has the same orientation as
	// curFace.
	switch moveDir {
	case dir.DIR_NORTH:
		newRel.Y = b.faceSize - 1
	case dir.DIR_SOUTH:
		newRel.Y = 0
	case dir.DIR_WEST:
		newRel.X = b.faceSize - 1
	case dir.DIR_EAST:
		newRel.X = 0
	}

	// Rotate newRel to account for newFace's orientation relative to
	// curFace. We start with the entry *side* (not direction) as it would
	// be seen by a newFace that has the same orientation as curFace then
	// rotate until that side matches the actual entry side relative to
	// newFace.
	for d := moveDir.Reverse(); d != entrySide; d = d.Right() {
		switch d {
		case dir.DIR_WEST:
			newRel = pos.P2{b.faceSize - 1 - newRel.Y, 0}
		case dir.DIR_NORTH:
			newRel = pos.P2{b.faceSize - 1, newRel.X}
		case dir.DIR_EAST:
			newRel = pos.P2{b.faceSize - 1 - newRel.Y, b.faceSize - 1}
		case dir.DIR_SOUTH:
			newRel = pos.P2{0, newRel.X}
		}
		fmt.Printf("rot from %v, newRel now %v\n", d, newRel)
	}

	logger.LogF("in %v face %v dir %v exited to %v face %v going %v",
		curRel, curFace.num, moveDir, newRel, newFace.num, newDir)

	v, _ := b.Get(b.AbsPos(newRel, newFace))
	return newRel, newFace, newDir, v
}

func parseInput(lines []string) (*grid.Grid[CellType], []Command, error) {
	gridLines := []string{}
	maxLen := -1
	cmdLine := ""
	for i, line := range lines {
		if line == "" {
			if i+1 >= len(lines) {
				return nil, nil, fmt.Errorf("no command line")
			}
			cmdLine = lines[i+1]
			break
		}
		gridLines = lines[0 : i+1]
		maxLen = mtsmath.Max(maxLen, len(line))
	}

	if cmdLine == "" {
		return nil, nil, fmt.Errorf("no break for command line")
	}

	for i, line := range gridLines {
		if l := len(line); l < maxLen {
			gridLines[i] = gridLines[i] + strings.Repeat(" ", maxLen-l)
		}
	}

	g, err := grid.NewFromLines[CellType](gridLines, func(p pos.P2, r rune) (CellType, error) {
		t, ok := ParseCellType(r)
		if !ok {
			return BLANK, fmt.Errorf("bad cell type %v", string(r))
		}
		return t, nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing grid: %v", err)
	}

	cmds, err := parseCommands(cmdLine)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing commands: %v", err)
	}

	return g, cmds, nil
}

var (
	facingMap = map[dir.Dir]int{
		dir.DIR_EAST:  0,
		dir.DIR_SOUTH: 1,
		dir.DIR_WEST:  2,
		dir.DIR_NORTH: 3,
	}
)

func solveA(g *grid.Grid[CellType], cmds []Command) int {
	board := NewFlatBoard(g)
	cur := board.FindStart()
	curDir := dir.DIR_EAST

	for i, cmd := range cmds {
		logger.LogF("%d: start %v %v; command %v", i+1, cur, curDir, cmd)

		if cmd == TURN_LEFT {
			curDir = curDir.Left()
			continue
		} else if cmd == TURN_RIGHT {
			curDir = curDir.Right()
			continue
		}

		// advance
		for n := int(cmd); n > 0; n-- {
			next, v := board.Move(cur, curDir)
			if v == WALL {
				break
			}
			cur = next
		}

		logger.LogF("  moved to %v", cur)
	}

	// The puzzle starts the grid at 1,1, while we use 0,0 because
	// 1-indexing is the tool of the devil. Adjust the position
	// into AoC-space.
	aoc := pos.P2{cur.X + 1, cur.Y + 1}
	return aoc.Y*1000 + aoc.X*4 + facingMap[curDir]
}

func detectCubeFaceSize(g *grid.Grid[CellType]) int {
	lens := map[int]bool{}

	isBlank := func(p pos.P2) bool {
		v, _ := g.Get(p)
		return v == BLANK
	}

	for y := 0; y <= g.Height(); y++ {
		runStart := 0
		runIsBlank := isBlank(pos.P2{0, y})
		for x := 1; x <= g.Width(); x++ {
			p := pos.P2{x, y}
			if blank := isBlank(p); blank == runIsBlank {
				continue
			} else {
				lens[p.X-runStart] = true
				runStart = p.X
				runIsBlank = blank
			}
		}
		if runStart != g.Width() {
			lens[g.Width()-runStart] = true
		}
	}

	// A horizontal scan seems to be enough for the test data I have...

	size := -1
	for n := range lens {
		if size == -1 || n < size {
			size = n
		}
	}

	// Double check that it evenly divides everything else
	for n := range lens {
		if n%size != 0 {
			panic("picked bad size")
		}
	}

	return size
}

func makeCubeSpec(g *grid.Grid[CellType], faceSize int) [][]bool {
	out := [][]bool{}
	for y := 0; y < g.Height(); y += faceSize {
		row := []bool{}
		for x := 0; x < g.Width(); x += faceSize {
			p := pos.P2{x, y}
			v, _ := g.Get(p)
			row = append(row, v != BLANK)
		}
		out = append(out, row)
	}
	return out
}

func dumpSpec(spec [][]bool, faces map[pos.P2]*PosFaceInfo) {
	for y := 0; y < len(spec); y++ {
		for x := 0; x < len(spec[0]); x++ {
			p := pos.P2{x, y}
			if spec[y][x] {
				if f, found := faces[p]; found && f.num != -1 {
					fmt.Print(f.num)
				} else {
					fmt.Print("X")
				}
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
}

func inSpec(spec [][]bool, p pos.P2) bool {
	if p.X < 0 || p.X >= len(spec[0]) || p.Y < 0 || p.Y >= len(spec) {
		return false
	}
	return spec[p.Y][p.X]
}

type PosFaceInfo struct {
	num       int
	pos       pos.P2
	neighbors map[dir.Dir]pos.P2
}

func (f *PosFaceInfo) String() string {
	out := fmt.Sprintf("[face %d %v", f.num, f.pos)
	for _, d := range dir.AllDirs {
		out += fmt.Sprintf(" %v:", d)
		p, found := f.neighbors[d]
		if found {
			out += fmt.Sprintf("%v", p)
		} else {
			out += "   "
		}
	}
	return out + "]"
}

func NewPosFaceInfo(p pos.P2) *PosFaceInfo {
	return &PosFaceInfo{
		num:       -1,
		pos:       p,
		neighbors: map[dir.Dir]pos.P2{},
	}
}

type FaceInfo struct {
	num              int
	pos              pos.P2
	neighbors        map[dir.Dir]int
	numCWRotsToIdeal int
}

func (f *FaceInfo) String() string {
	out := fmt.Sprintf("[face %d %v", f.num, f.pos)
	for _, d := range dir.AllDirs {
		out += fmt.Sprintf(" %v:", d)
		n, found := f.neighbors[d]
		if found {
			out += fmt.Sprintf("%d", n)
		} else {
			out += " "
		}
	}
	out += fmt.Sprintf(" rots:%d", f.numCWRotsToIdeal)
	return out + "]"
}

func topLeft(spec [][]bool) pos.P2 {
	for y := 0; y < len(spec); y++ {
		for x := 0; x < len(spec[0]); x++ {
			if spec[y][x] {
				return pos.P2{x, y}
			}
		}
	}
	panic("none found")
}

func findOpposite(front pos.P2, faces map[pos.P2]*PosFaceInfo) pos.P2 {
	for _, f := range faces {
		isNeighbor := false
		if f.pos.Equals(front) {
			continue
		}
		for _, n := range faces[front].neighbors {
			if f.pos.Equals(n) {
				isNeighbor = true
				break
			}
		}

		if !isNeighbor {
			return f.pos
		}
	}

	panic("opposite not found")
}

type FindJob struct {
	p pos.P2
	d dir.Dir
}

func copySides(searchDirs, copyDirs [2]dir.Dir, to *PosFaceInfo, faces map[pos.P2]*PosFaceInfo, spec [][]bool) {
	refs := [2]pos.P2{pos.P2{-1, -1}, pos.P2{-1, -1}}

	top := to.pos
	for p := searchDirs[0].From(to.pos); inSpec(spec, p); p = searchDirs[0].From(p) {
		top = p
	}
	for p := top; inSpec(spec, p); p = searchDirs[1].From(p) {
		for i, refDir := range copyDirs {
			if refPos := refDir.From(p); inSpec(spec, refPos) {
				refs[i] = refPos
			}
		}
	}

	for i, copyDir := range copyDirs {
		ref := refs[i]
		if ref.X == -1 {
			continue
		}
		for p := top; inSpec(spec, p); p = searchDirs[1].From(p) {
			if f, found := faces[p]; found {
				f.neighbors[copyDir] = ref
			}
		}
	}
}

func walkSpecFaces(start pos.P2, spec [][]bool, callback func(pos.P2, dir.Dir, pos.P2)) {
	queue := list.New()

	for _, d := range dir.AllDirs {
		if inSpec(spec, d.From(start)) {
			queue.PushBack(&FindJob{start, d})
		}
	}

	for queue.Front() != nil {
		job := queue.Front().Value.(*FindJob)
		queue.Remove(queue.Front())

		dest := job.d.From(job.p)
		callback(job.p, job.d, dest)

		for _, d := range dir.AllDirs {
			if d != job.d.Reverse() && inSpec(spec, d.From(dest)) {
				queue.PushBack(&FindJob{dest, d})
			}
		}
	}
}

type IdealFace struct {
	num   int
	cwRot int
}

var (
	idealNeighborsByFaceNum = map[int]map[dir.Dir]IdealFace{
		1: map[dir.Dir]IdealFace{
			dir.DIR_NORTH: IdealFace{4, 2}, dir.DIR_SOUTH: IdealFace{2, 0},
			dir.DIR_WEST: IdealFace{3, 3}, dir.DIR_EAST: IdealFace{5, 1},
		},
		2: map[dir.Dir]IdealFace{
			dir.DIR_NORTH: IdealFace{1, 0}, dir.DIR_SOUTH: IdealFace{6, 0},
			dir.DIR_WEST: IdealFace{3, 0}, dir.DIR_EAST: IdealFace{5, 0},
		},
		3: map[dir.Dir]IdealFace{
			dir.DIR_NORTH: IdealFace{1, 1}, dir.DIR_SOUTH: IdealFace{6, 3},
			dir.DIR_WEST: IdealFace{4, 0}, dir.DIR_EAST: IdealFace{2, 0},
		},
		4: map[dir.Dir]IdealFace{
			dir.DIR_NORTH: IdealFace{1, 2}, dir.DIR_SOUTH: IdealFace{6, 2},
			dir.DIR_WEST: IdealFace{5, 0}, dir.DIR_EAST: IdealFace{3, 0},
		},
		5: map[dir.Dir]IdealFace{
			dir.DIR_NORTH: IdealFace{1, 3}, dir.DIR_SOUTH: IdealFace{6, 1},
			dir.DIR_WEST: IdealFace{2, 0}, dir.DIR_EAST: IdealFace{4, 0},
		},
		6: map[dir.Dir]IdealFace{
			dir.DIR_NORTH: IdealFace{2, 0}, dir.DIR_SOUTH: IdealFace{4, 2},
			dir.DIR_WEST: IdealFace{3, 1}, dir.DIR_EAST: IdealFace{5, 3},
		},
	}
)

func allCWRotations(in map[dir.Dir]int) []map[dir.Dir]int {
	outs := []map[dir.Dir]int{}

	for i := 0; i < 4; i++ {
		ref := in
		if i > 0 {
			ref = outs[i-1]
		}

		out := map[dir.Dir]int{}
		for d, n := range ref {
			if i != 0 {
				d = d.Right()
			}
			out[d] = n
		}
		outs = append(outs, out)
	}

	return outs
}

func findCubeFaces(spec [][]bool) map[pos.P2]*PosFaceInfo {
	start := topLeft(spec)
	faces := map[pos.P2]*PosFaceInfo{
		start: NewPosFaceInfo(start),
	}

	// Do our best to find face neighbors by walking the spec and
	// propagating simple relationships.
	walkSpecFaces(start, spec, func(fromPos pos.P2, d dir.Dir, toPos pos.P2) {
		from := faces[fromPos]
		if _, found := faces[toPos]; !found {
			faces[toPos] = NewPosFaceInfo(toPos)
		}
		to := faces[toPos]

		from.neighbors[d] = toPos
		to.neighbors[d.Reverse()] = from.pos

		// Consider a spec that looks like this:
		//
		//   A
		//   BCD
		//
		// A is a neighbor in some direction from B. It's also a
		// neighbor in the same direction from C and D. The following
		// calls make that association in either orientation and on
		// either side.
		switch {
		case d == dir.DIR_SOUTH || d == dir.DIR_NORTH:
			copySides([2]dir.Dir{dir.DIR_NORTH, dir.DIR_SOUTH},
				[2]dir.Dir{dir.DIR_EAST, dir.DIR_WEST}, to, faces, spec)
		case d == dir.DIR_EAST || d == dir.DIR_WEST:
			copySides([2]dir.Dir{dir.DIR_EAST, dir.DIR_WEST},
				[2]dir.Dir{dir.DIR_NORTH, dir.DIR_NORTH}, to, faces, spec)
		}
	})

	// Find a face that has four neighbors. Hopefully the propagation was
	// enough to make at least one. If not, well, we'll have a bug to
	// fix. It works on the test data though.
	//
	// Assuming we find such a face it'll be face 1. Having four neighbors
	// we can automatically identify face 6 (its opposite).
	face1 := pos.P2{-1, -1}
	for _, f := range faces {
		if len(f.neighbors) == 4 {
			face1 = f.pos
			break
		}
	}

	if face1.X == -1 {
		panic("no cell with all neighbors")
	}

	face6 := findOpposite(face1, faces)

	// Now that we know faces 1 and 6 we can label all six faces. The
	// remaining four are the four neighbors of face 1.
	faces[face1].num = 1
	faces[face6].num = 6

	facesByNum := map[int]*PosFaceInfo{
		1: faces[face1],
		6: faces[face6],
	}
	for i, d := range []dir.Dir{dir.DIR_SOUTH, dir.DIR_WEST, dir.DIR_NORTH, dir.DIR_EAST} {
		num := i + 2
		n := faces[faces[face1].neighbors[d]]
		n.num = num
		facesByNum[num] = n
	}

	// Fill in unknown neighbors. We know what the neighbors should be and
	// in which order -- they're just rotated (not flipped). One known
	// neighbor on a given face tells us how to orient the list of
	// neighbors.
	for _, f := range faces {
		if len(f.neighbors) == 4 {
			continue
		}

		var anchorDir dir.Dir
		var anchorNum int
		for _, d := range dir.AllDirs {
			if p, found := f.neighbors[d]; found {
				anchorDir = d
				anchorNum = faces[p].num
			}
		}

		var match map[dir.Dir]int
		idealNeighbors := map[dir.Dir]int{}
		for d, nf := range idealNeighborsByFaceNum[f.num] {
			idealNeighbors[d] = nf.num
		}

		for _, rotation := range allCWRotations(idealNeighbors) {
			if rotation[anchorDir] == anchorNum {
				match = rotation
				break
			}
		}

		if match == nil {
			panic("no match")
		}

		for _, d := range dir.AllDirs {
			if np, found := f.neighbors[d]; found {
				if faces[np].num != match[d] {
					panic("mismatch")
				}
			} else {
				f.neighbors[d] = facesByNum[match[d]].pos
			}
		}
	}

	return faces
}

func findFaceRotations(posFaces map[pos.P2]*PosFaceInfo) map[int]*FaceInfo {
	out := map[int]*FaceInfo{}
	for _, pf := range posFaces {
		f := &FaceInfo{
			num:       pf.num,
			pos:       pf.pos,
			neighbors: map[dir.Dir]int{},
		}

		for d, np := range pf.neighbors {
			f.neighbors[d] = posFaces[np].num
		}

		ideal := idealNeighborsByFaceNum[f.num]
		for i, rotation := range allCWRotations(f.neighbors) {
			if rotation[dir.DIR_NORTH] == ideal[dir.DIR_NORTH].num {
				f.numCWRotsToIdeal = i
				break
			}
		}

		out[f.num] = f
	}

	return out
}

func solveB(g *grid.Grid[CellType], cmds []Command) int {
	b := NewCubeBoard(g)
	curRelPos, curFace := b.FindStart()
	curDir := dir.DIR_EAST

	for i, cmd := range cmds {
		logger.LogF("%d: start %v on %v (%v) %v; command %v",
			i+1, curRelPos, curFace.num, b.AbsPos(curRelPos, curFace),
			curDir, cmd)

		if cmd == TURN_LEFT {
			curDir = curDir.Left()
			continue
		} else if cmd == TURN_RIGHT {
			curDir = curDir.Right()
			continue
		}

		// advance
		for n := int(cmd); n > 0; n-- {
			nextRelPos, nextFace, nextDir, v := b.Move(
				curRelPos, curFace, curDir)
			if v == WALL {
				break
			}
			curRelPos, curFace, curDir = nextRelPos, nextFace, nextDir
		}

		logger.LogF("  moved to %v on %v (abs %v)",
			curRelPos, curFace.num, b.AbsPos(curRelPos, curFace))
	}

	// The puzzle starts the grid at 1,1, while we use 0,0 because
	// 1-indexing is the tool of the devil. Adjust the position
	// into AoC-space.
	abs := b.AbsPos(curRelPos, curFace)
	aoc := pos.P2{abs.X + 1, abs.Y + 1}
	return aoc.Y*1000 + aoc.X*4 + facingMap[curDir]
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		log.Fatal(err)
	}

	g, commands, err := parseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("A", solveA(g, commands))
	fmt.Println("B", solveB(g, commands))
}
