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
	"bytes"
	"container/list"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/simmonmt/aoc/2021/common/astar"
	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/grid"
	"github.com/simmonmt/aoc/2021/common/logger"
	"github.com/simmonmt/aoc/2021/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type CellType int

const (
	CT_UNKNOWN CellType = iota
	CT_NORMAL
	CT_INTERSECT
	CT_ROOMA
	CT_ROOMB
	CT_ROOMC
	CT_ROOMD
)

func (c CellType) IsRoom() bool {
	if c < CT_ROOMA || c > CT_ROOMD {
		return false
	}

	return true
}

func (c CellType) String() string {
	switch c {
	case CT_NORMAL:
		return "."
	case CT_INTERSECT:
		return "x"
	case CT_ROOMA:
		return "a"
	case CT_ROOMB:
		return "b"
	case CT_ROOMC:
		return "c"
	case CT_ROOMD:
		return "d"
	default:
		panic("bad cell type")
	}
}

type CharType int

const (
	CHAR_A CharType = iota
	CHAR_B
	CHAR_C
	CHAR_D
	CHAR_SENTINEL
)

func (c CharType) String() string {
	if c < CHAR_SENTINEL {
		return string('A' + c)
	}
	return "?"
}

func (c CharType) RoomType() CellType {
	return CellType(int(CT_ROOMA) + int(c))
}

func (c CharType) MoveCost() int {
	return []int{1, 10, 100, 1000}[c]
}

type GameState struct {
	roomHeight int
	locs       [4][]pos.P2
	locsByPos  map[pos.P2]CharType
}

func NewGameState(roomHeight int, locs [4][]pos.P2) *GameState {
	locsByPos := map[pos.P2]CharType{}
	for i, ps := range locs {
		for _, p := range ps {
			locsByPos[p] = CharType(i)
		}
	}

	gs := &GameState{
		roomHeight: roomHeight,
		locs:       locs,
		locsByPos:  locsByPos,
	}
	gs.sortLocs()
	return gs
}

func DeserializeGameState(ser string) (*GameState, error) {
	strs := strings.Split(ser, "|")
	if len(strs)%4 != 1 {
		return nil, fmt.Errorf("bad splits (%v) in %v",
			len(strs), ser)
	}

	roomHeight, err := strconv.Atoi(strs[0])
	if err != nil {
		return nil, fmt.Errorf("bad roomsize: %v", err)
	}

	strs = strs[1:]

	locs := [4][]pos.P2{}
	for i, s := range strs {
		charNum := i / roomHeight
		if locs[charNum] == nil {
			locs[charNum] = []pos.P2{}
		}

		p, err := pos.P2FromString(s)
		if err != nil {
			return nil, fmt.Errorf("bad pos in %v", ser)
		}
		locs[charNum] = append(locs[charNum], p)
	}

	return NewGameState(roomHeight, locs), nil
}

func (gs *GameState) sortLocs() {
	for i := 0; i < len(gs.locs); i++ {
		sort.Slice(gs.locs[i], func(j, k int) bool {
			return gs.locs[i][j].LessThan(gs.locs[i][k])
		})
	}
}

func (gs *GameState) Serialize() string {
	roomHeight := len(gs.locs[0])
	out := make([]string, 4*roomHeight+1)
	out[0] = strconv.Itoa(roomHeight)
	for i := 0; i < 4; i++ {
		for j, p := range gs.locs[i] {
			out[i*roomHeight+j+1] = p.String()
		}
	}
	return strings.Join(out, "|")
}

func (gs *GameState) IsOccupied(p pos.P2) (bool, CharType) {
	i, found := gs.locsByPos[p]
	if !found {
		return false, CHAR_A
	}

	return true, CharType(i)
}

func (gs *GameState) Move(char CharType, from, to pos.P2) *GameState {
	toLocs := [4][]pos.P2{}
	for i := 0; i < 4; i++ {
		toLocs[i] = make([]pos.P2, len(gs.locs[i]))
		if i == int(char) {
			found := true
			for j, p := range gs.locs[i] {
				if p.Equals(from) {
					found = true
					p = to
				}
				toLocs[i][j] = p
			}
			if !found {
				panic("from not found")
			}
		} else {
			copy(toLocs[i], gs.locs[i])
		}
	}

	return NewGameState(gs.roomHeight, toLocs)
}

func (gs *GameState) AllChars() map[pos.P2]CharType {
	return gs.locsByPos
}

func (gs *GameState) RoomHeight() int {
	return gs.roomHeight
}

type Board struct {
	g         *grid.Grid
	roomCells map[CellType][]pos.P2
}

func NewBoard(roomHeight int) *Board {
	g := grid.New(11, 1+roomHeight)
	roomCells := map[CellType][]pos.P2{}

	for x := 0; x < 11; x++ {
		g.Set(pos.P2{x, 0}, CT_NORMAL)
	}
	for i, x := range []int{2, 4, 6, 8} {
		g.Set(pos.P2{x, 0}, CT_INTERSECT)

		roomType := CellType(int(CT_ROOMA) + i)

		for y := 1; y <= roomHeight; y++ {
			p := pos.P2{x, y}
			g.Set(p, roomType)
			roomCells[roomType] = append(roomCells[roomType], p)
		}
	}

	return &Board{
		g:         g,
		roomCells: roomCells,
	}
}

func (b *Board) RoomCells(roomType CellType) []pos.P2 {
	return b.roomCells[roomType]
}

func (b *Board) Get(p pos.P2) CellType {
	if p.Y >= b.g.Height() || p.X >= b.g.Width() {
		return CT_UNKNOWN
	}

	v := b.g.Get(p)
	if v == nil {
		return CT_UNKNOWN
	}
	return v.(CellType)
}

func (b *Board) cellToString(p pos.P2, gs *GameState) string {
	if found, char := gs.IsOccupied(p); found {
		return char.String()
	}

	v := b.g.Get(p)
	if c, ok := v.(CellType); ok {
		return c.String()
	} else {
		return "?"
	}
}

func (b *Board) AllNeighbors(p pos.P2) []pos.P2 {
	out := []pos.P2{}
	for _, n := range b.g.AllNeighbors(p, false) {
		ct := b.Get(n)
		if ct == CT_UNKNOWN {
			continue
		}
		out = append(out, n)
	}
	return out
}

func (b *Board) DumpTo(o io.Writer, gs *GameState) {
	fmt.Fprintln(o, "#############")
	fmt.Fprint(o, "#")
	for x := 0; x < 11; x++ {
		p := pos.P2{x, 0}
		fmt.Fprintf(o, "%s", b.cellToString(p, gs))
	}
	fmt.Fprintln(o, "#")

	for y := 1; y <= gs.roomHeight; y++ {
		if y == 1 {
			fmt.Fprint(o, "###")
		} else {
			fmt.Fprint(o, "  #")
		}

		for x := 2; x <= 8; x += 2 {
			p := pos.P2{x, y}
			fmt.Fprintf(o, "%s", b.cellToString(p, gs))
			fmt.Fprint(o, "#")
		}
		if y == 1 {
			fmt.Fprint(o, "##")
		}
		fmt.Fprintln(o)
	}

	fmt.Fprintln(o, "  #########")
}

func (b *Board) DumpToString(gs *GameState) string {
	buf := bytes.Buffer{}
	b.DumpTo(&buf, gs)
	return buf.String()
}

func (b *Board) Dump(gs *GameState) {
	b.DumpTo(os.Stdout, gs)
}

func readInput(path string) ([]string, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, err
	}
	return lines, nil
}

func parseExpandedInput(lines []string) (*Board, *GameState) {
	expanded := []string{}
	expanded = append(expanded, lines[0:3]...)
	expanded = append(expanded, "  #D#C#B#A#  ", "  #D#B#A#C#  ")
	expanded = append(expanded, lines[3:]...)

	return parseInput(expanded)
}

func parseInput(lines []string) (*Board, *GameState) {
	locs := [4][]pos.P2{}
	for i := range locs {
		locs[i] = []pos.P2{}
	}

	saveChar := func(p pos.P2, ch byte) {
		c := int(ch - 'A')
		if c >= 0 && c <= 3 {
			locs[c] = append(locs[c], p)
		}
	}

	for lx := 1; lx < 12; lx++ {
		p := pos.P2{lx - 1, 0}
		saveChar(p, lines[1][lx])
	}

	roomHeight := len(lines) - 3
	for lineNo := 2; lineNo < len(lines)-1; lineNo++ {
		y := lineNo - 1

		for _, lx := range []int{3, 5, 7, 9} {
			p := pos.P2{lx - 1, y}
			saveChar(p, lines[lineNo][lx])
		}
	}

	return NewBoard(roomHeight), NewGameState(roomHeight, locs)
}

type astarClient struct {
	b *Board
}

func (c *astarClient) roomOpen(gs *GameState, roomType CellType) bool {
	for _, p := range c.b.RoomCells(roomType) {
		if occupied, char := gs.IsOccupied(p); occupied {
			if char.RoomType() != roomType {
				return false
			}
		}
	}
	return true
}

func (c *astarClient) dfs(gs *GameState, charLoc pos.P2, cb func(*list.List, pos.P2, CellType) bool) {
	visited := map[pos.P2]bool{}
	path := &list.List{}
	c.doDfs(gs, charLoc, visited, path, cb)
}

func (c *astarClient) doDfs(gs *GameState, p pos.P2, visited map[pos.P2]bool, path *list.List, cb func(*list.List, pos.P2, CellType) bool) {
	if _, found := visited[p]; found {
		return
	}

	visited[p] = true
	path.PushBack(p)
	deeper := cb(path, p, c.b.Get(p))

	//logger.LogF("%v deeper? %v", p, deeper)

	if deeper {
		for _, n := range c.b.AllNeighbors(p) {
			if _, found := visited[n]; found {
				continue
			}

			c.doDfs(gs, n, visited, path, cb)
		}
	}

	path.Remove(path.Back())
}

func (c *astarClient) makeMove(gs *GameState, path *list.List, char CharType) *GameState {
	from := path.Front().Value.(pos.P2)
	to := path.Back().Value.(pos.P2)
	return gs.Move(char, from, to)
}

func (c *astarClient) atBottomOfRoom(gs *GameState, loc pos.P2, roomType CellType) bool {
	roomCells := c.b.RoomCells(roomType)
	for i := len(roomCells) - 1; i >= 0; i-- {
		rp := roomCells[i]
		if rp.Equals(loc) {
			return true
		} else if occupied, who := gs.IsOccupied(rp); !occupied {
			return false
		} else if who.RoomType() != roomType {
			return false
		}
	}

	panic("didn't find loc")
}

func (c *astarClient) allMovesForChar(gs *GameState, charLoc pos.P2, char CharType) []string {
	//logger.LogF("all moves for %v at %v", char, charLoc)

	// No further moves if it's as far down in its room as it can get.
	cellType := c.b.Get(charLoc)
	if cellType == char.RoomType() && c.atBottomOfRoom(gs, charLoc, char.RoomType()) {
		return []string{}
	}

	startsFromRoom := c.b.Get(charLoc).IsRoom()

	outs := []string{}
	c.dfs(gs, charLoc, func(path *list.List, cur pos.P2, cellType CellType) bool {
		if cur.Equals(charLoc) {
			return true
		}

		if occupied, _ := gs.IsOccupied(cur); occupied {
			return false
		}

		switch cellType {
		case CT_INTERSECT:
			// Rule 1: can't move here, but keep traversing
			return true

		case CT_NORMAL:
			if !startsFromRoom {
				// Rule 3: Once an amphipod stops moving in the
				// hallway, it will stay in that spot until it
				// can move into a room. Keep traversing.
				return true
			}

		default:
			if !cellType.IsRoom() {
				panic("unexpected cell type")
			}

			wantRoomType := char.RoomType()

			if cellType != wantRoomType {
				// Rule 2a: Amphipods will never move from the
				// hallway into a room unless that room is their
				// destination room ... so keep traversing
				return true
			}

			if !c.roomOpen(gs, wantRoomType) {
				// Rule 2b: ... and that room contains no
				// amphipods which do not also have that room as
				// their own destination.
				//
				// Open rooms don't have amphipods of non-want
				// types.
				return true
			}

			// This is the right room type, and it's empty (or
			// otherwise occupied only by the same type of
			// amphipod). Make sure we're at the bottom, as we'll
			// have to get there eventually. No reason to do it in
			// multiple steps.
			below := pos.P2{X: cur.X, Y: cur.Y + 1}
			if c.b.Get(below) == wantRoomType {
				if occupied, _ := gs.IsOccupied(below); !occupied {
					// Below is empty, so keep traversing.
					return true
				}
			}
		}

		newState := c.makeMove(gs, path, char)
		outs = append(outs, newState.Serialize())
		return true
	})

	//logger.LogF("moves for %v: %v", char, len(outs))
	return outs
}

func (c *astarClient) AllNeighbors(start string) []string {
	gs, err := DeserializeGameState(start)
	if err != nil {
		panic(err.Error())
	}

	outs := []string{}
	for charLoc, char := range gs.AllChars() {
		neighbors := c.allMovesForChar(gs, charLoc, char)
		// if logger.Enabled() {
		// 	if char == CHAR_B2 {
		// 		for _, n := range neighbors {
		// 			ns, _ := DeserializeGameState(n)
		// 			fmt.Println(n)
		// 			c.b.Dump(ns)
		// 		}
		// 	}
		// }
		outs = append(outs, neighbors...)
	}
	return outs
}

func (c *astarClient) EstimateDistance(cur, goal string) uint {
	curState, err := DeserializeGameState(cur)
	if err != nil {
		panic(err.Error())
	}

	cost := 0
	for charLoc, char := range curState.AllChars() {
		roomCells := c.b.RoomCells(char.RoomType())

		inRoom := false
		for _, c := range roomCells {
			if c == charLoc {
				break
			}
		}
		if inRoom {
			continue
		}

		closestRoomCell := roomCells[0]
		distToRoom := charLoc.ManhattanDistance(closestRoomCell)

		cost += distToRoom * char.MoveCost()
	}

	return uint(cost)
}

// 	// NeighborDistance returns the distance between two known direct
// 	// neighbors (i.e. a pair derived using AllNeighbors).
// 	NeighborDistance(n1, n2 string) uint
func (c *astarClient) NeighborDistance(from, to string) uint {
	fromState, err := DeserializeGameState(from)
	if err != nil {
		panic("bad from state")
	}
	toState, err := DeserializeGameState(to)
	if err != nil {
		panic("bad to state")
	}

	// We only move one piece at a time, so we can find the moved piece by
	// looking for the difference between the AllChars maps.
	var fromPos, toPos pos.P2
	var fromChar CharType
	fromCharsByLoc := fromState.AllChars()
	toCharsByLoc := toState.AllChars()

	for loc, char := range fromCharsByLoc {
		if _, found := toCharsByLoc[loc]; !found {
			fromPos = loc
			fromChar = char
			break
		}
	}
	for loc := range toCharsByLoc {
		if _, found := fromCharsByLoc[loc]; !found {
			toPos = loc
			break
		}
	}

	foundPathLen := -1
	c.dfs(fromState, fromPos, func(path *list.List, cur pos.P2, cellType CellType) bool {
		if cur.Equals(toPos) {
			if foundPathLen != -1 {
				panic("repath")
			}

			// if logger.Enabled() {
			// 	fmt.Print("found path: ")
			// 	for e := path.Front(); e != nil; e = e.Next() {
			// 		fmt.Print(e.Value.(pos.P2), " ")
			// 	}
			// 	fmt.Println()
			// }

			foundPathLen = path.Len()
			return false
		}
		return true
	})

	if foundPathLen == -1 {
		panic("didn't find a path")
	}

	return uint(fromChar.MoveCost() * (foundPathLen - 1))
}

func (c *astarClient) GoalReached(cand, goal string) bool {
	return cand == goal
}

func makeGoalState(board *Board, initial *GameState) *GameState {
	locs := [4][]pos.P2{}

	for _, char := range initial.AllChars() {
		if locs[char] != nil {
			continue
		}

		locs[char] = board.RoomCells(char.RoomType())[:]
	}

	return NewGameState(initial.RoomHeight(), locs)
}

func solve(board *Board, gameState *GameState) uint {
	board.Dump(gameState)

	goal := makeGoalState(board, gameState)

	client := &astarClient{b: board}
	path := astar.AStar(gameState.Serialize(), goal.Serialize(), client)
	if path == nil {
		fmt.Println("no path found")
		return 0
	}

	cost := uint(0)
	for i := 0; i < len(path)-1; i++ {
		cost += client.NeighborDistance(path[i], path[i+1])
	}

	logger.LogF("result", strings.Join(path, "\n"))
	return cost
}

func solveA(lines []string) {
	board, gameState := parseInput(lines)
	cost := solve(board, gameState)

	fmt.Println("A", cost)
}

func solveB(lines []string) {
	board, gameState := parseExpandedInput(lines)
	cost := solve(board, gameState)

	fmt.Println("B", cost)
}

func timeSolve(fn func()) {
	start := time.Now()
	fn()
	end := time.Now()

	fmt.Println("elapsed:", end.Sub(start))
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	timeSolve(func() { solveA(lines) })
	timeSolve(func() { solveB(lines) })
}
