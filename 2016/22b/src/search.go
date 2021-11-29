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
	"image"
	"image/color"
	"logger"
)

// heuristic
//
// moving left: manhattan distance from goal to 0,0
// moving to goal: (ml + manhattan distance empty to goal)

type aStarHelper struct {
	board    *Board
	visited  map[Pos]bool
	num, mod int
	images   []*image.Paletted
}

func NewAStarHelper(board *Board, saveImages bool) *aStarHelper {
	var visited map[Pos]bool
	if saveImages {
		visited = map[Pos]bool{}
	}

	return &aStarHelper{
		board:   board,
		visited: visited,
		num:     0,
		mod:     100,
		images:  []*image.Paletted{},
	}
}

var (
	dirs = []Pos{
		Pos{-1, 0},
		Pos{1, 0},
		Pos{0, -1},
		Pos{0, 1},
	}
)

func (h *aStarHelper) addImage() {
	width, height := h.board.Size()
	arr := make([]uint8, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if _, found := h.visited[Pos{x, y}]; found {
				arr[y*width+x] = 1
			}
		}
	}

	image := &image.Paletted{
		Pix:     arr,
		Stride:  width,
		Rect:    image.Rect(0, 0, width, height),
		Palette: []color.Color{color.White, color.Black},
	}

	h.images = append(h.images, image)
}

func (h *aStarHelper) AllNeighbors(start string) []string {
	ps := Decode(start)

	if h.visited != nil {
		h.num++
		h.visited[ps.Empty] = true
		if h.num%10 == 0 {
			h.addImage()
		}
	}

	width, height := h.board.Size()

	outs := []*PlayState{}

	for _, dir := range dirs {
		cand := Pos{
			X: ps.Empty.X + dir.X,
			Y: ps.Empty.Y + dir.Y,
		}

		if cand.X < 0 || cand.Y < 0 || cand.X >= width || cand.Y >= height {
			continue
		}

		if !h.board.IsMoveable(cand) {
			continue
		}

		var out *PlayState
		if cand.Eq(ps.Goal) {
			out = &PlayState{Empty: ps.Goal, Goal: ps.Empty}
		} else {
			out = &PlayState{Empty: cand, Goal: ps.Goal}
		}
		outs = append(outs, out)
	}

	if *verbose {
		logger.LogLn("AllNeighbors in:")
		h.board.Dump(ps)
		logger.LogLn("AllNeighbors out:")
		for i, out := range outs {
			logger.LogF("out %d:", i)
			h.board.Dump(out)
		}
	}

	neighbors := make([]string, len(outs))
	for i, out := range outs {
		neighbors[i] = out.Encode()
	}
	return neighbors
}

func (h *aStarHelper) EstimateDistance(start, goal string) uint {
	startPs := Decode(start)
	goalPs := Decode(goal)

	return uint(startPs.Empty.Dist(startPs.Goal) +
		startPs.Goal.Dist(goalPs.Goal))
}

func (h *aStarHelper) NeighborDistance(n1, n2 string) uint {
	return 1
}

func (h *aStarHelper) GoalReached(cand, goal string) bool {
	ps := Decode(cand)
	return ps.Goal.X == 0 && ps.Goal.Y == 0
}

func (h *aStarHelper) PrintableKey(key string) string {
	return key
}

func (h *aStarHelper) MarkClosed(key string) {
}
