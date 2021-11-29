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

package puzzle

import (
	"sort"

	"github.com/simmonmt/aoc/2019/common/collections"
	"github.com/simmonmt/aoc/2019/common/dir"
	"github.com/simmonmt/aoc/2019/common/pos"
)

type Path struct {
	Dest  string
	Dist  int
	Doors []string
}

type PathsByDest []Path

func (a PathsByDest) Len() int           { return len(a) }
func (a PathsByDest) Less(i, j int) bool { return a[i].Dest < a[j].Dest }
func (a PathsByDest) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type searchElem struct {
	Pos   pos.P2
	Dist  int
	Doors []string
}

func FindAllPathsToKey(b *Board, start pos.P2, key string) []Path {
	dest := b.KeyLoc(key)

	curPath := collections.NewStack()
	curPath.Push(&searchElem{start, 0, []string{}})

	seen := map[pos.P2]bool{start: true}

	paths := []Path{}
	foundPath := func(elem *searchElem) {
		paths = append(paths, Path{
			Dest:  key,
			Dist:  elem.Dist,
			Doors: elem.Doors[:],
		})
	}
	doFindAllPathsToKey(b, curPath, seen, dest, foundPath)
	return paths
}

func doFindAllPathsToKey(b *Board, curPath *collections.Stack, seen map[pos.P2]bool, dest pos.P2, found func(*searchElem)) {
	curElem := curPath.Peek().(*searchElem)

	if curElem.Pos.Equals(dest) {
		found(curElem)
		return
	}

	doors := curElem.Doors
	if t := b.Get(curElem.Pos); t == TILE_DOOR {
		curDoor := b.DoorAtLoc(curElem.Pos)

		if len(doors) == 0 {
			doors = []string{curDoor}
		} else {
			newDoors := make([]string, len(doors)+1)
			copy(newDoors, doors)
			newDoors[len(doors)] = b.DoorAtLoc(curElem.Pos)
			doors = newDoors
		}
	}

	for _, n := range allNeighbors(b, curElem.Pos) {
		if _, found := seen[n]; found {
			continue
		}

		t := b.Get(n)
		if t == TILE_KEY && !n.Equals(dest) {
			continue // can't go through other keys
		}

		nElem := &searchElem{
			Pos:   n,
			Dist:  curElem.Dist + 1,
			Doors: doors,
		}

		seen[n] = true
		curPath.Push(nElem)
		doFindAllPathsToKey(b, curPath, seen, dest, found)
		curPath.Pop()
		delete(seen, n)
	}
}

func allNeighbors(b *Board, p pos.P2) []pos.P2 {
	out := []pos.P2{}
	for _, d := range dir.AllDirs {
		np := d.From(p)
		if t := b.Get(np); t != TILE_WALL {
			out = append(out, np)
		}
	}
	return out
}

func FindAllPaths(b *Board, start pos.P2) map[string][]Path {
	allPathsFromAllKeys := map[string][]Path{}

	sources := b.Keys()
	sources = append(sources, "@")

	for _, source := range sources {
		var sourceLoc pos.P2
		if source == "@" {
			sourceLoc = start
		} else {
			sourceLoc = b.KeyLoc(source)
		}

		allPathsFromAllKeys[source] = []Path{}
		for _, dest := range b.Keys() {
			if source == dest {
				continue
			}

			paths := FindAllPathsToKey(b, sourceLoc, dest)
			if len(paths) == 0 {
				continue
			}

			best := paths[0]
			for _, path := range paths {
				if path.Dist < best.Dist {
					best = path
				}
			}

			allPathsFromAllKeys[source] =
				append(allPathsFromAllKeys[source], best)

			sort.Sort(PathsByDest(allPathsFromAllKeys[source]))
		}
	}

	return allPathsFromAllKeys
}

func FindAllPathsMulti(board *Board, starts []pos.P2) map[pos.P2]map[string][]Path {
	graphs := map[pos.P2]map[string][]Path{}
	for _, start := range starts {
		graphs[start] = FindAllPaths(board, start)
	}
	return graphs
}
