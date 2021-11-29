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

// AStar search algorithm
//
// Implemented from pseudocode on
// https://en.wikipedia.org/wiki/A*_search_algorithm

package astar

import (
	"math"

	"github.com/google/btree"

	"logger"
)

type NeighborDiscoverer interface {
	AllNeighbors(start string) []string
}

type DistanceEstimator interface {
	Estimate(start, end string) uint
	NeighborDistance(n1, n2 string) uint
}

type scoreMap map[string]uint

func (m *scoreMap) GetWithDefault(key string, def uint) uint {
	if v, found := (*m)[key]; found {
		return v
	} else {
		return math.MaxUint32
	}
}

type fScoreItem struct {
	Name  string
	Value uint
}

func (i *fScoreItem) Less(x btree.Item) bool {
	than := x.(*fScoreItem)

	// We implement greater-than because we want the tree to store small-to-large
	if i.Value > than.Value {
		return true
	} else if i.Value < than.Value {
		return false
	} else {
		return i.Name > than.Name
	}
}

type fScoreMap struct {
	btree *btree.BTree
}

func newFScoreMap() *fScoreMap {
	return &fScoreMap{
		btree: btree.New(2000),
	}
}

func (m *fScoreMap) Walk(visitor func(item *fScoreItem) bool) {
	m.btree.Descend(func(item btree.Item) bool {
		return visitor(item.(*fScoreItem))
	})
}

func (m *fScoreMap) Set(name string, value uint) {
	m.btree.ReplaceOrInsert(&fScoreItem{Name: name, Value: value})
}

func (m *fScoreMap) Delete(name string, value uint) {
	m.btree.Delete(&fScoreItem{Name: name, Value: value})
}

func reconstructPath(cameFrom map[string]string, current string) []string {
	totalPath := []string{current}
	for {
		next, found := cameFrom[current]
		if !found {
			break
		}

		totalPath = append(totalPath, next)
		current = next

	}
	return totalPath
}

func AStar(start, goal string, neighborDiscoverer NeighborDiscoverer, distanceEstimator DistanceEstimator) []string {
	logger.LogF("astar start %v goal %v\n", start, goal)

	closedSet := map[string]bool{}
	openSet := map[string]bool{start: true}
	cameFrom := map[string]string{}

	gScore := scoreMap{}
	gScore[start] = 0

	fScore := newFScoreMap()
	fScore.Set(start, distanceEstimator.Estimate(start, goal))

	// A subset of fScore that contains only the nodes currently in openSet.
	//
	// One of the most expensive parts of this algorithm is searching fScore
	// for the item from openSet that has the lowest value. The pseudocode
	// from Wikipedia has a single fScore map that contains values from all
	// nodes -- open and not open. If there are a lot of not open nodes
	// whose scores are lower than those of the open nodes, the search for
	// the open nodes with the lowest score will take longer and longer each
	// time. This map stores only the scores that are part of openSet. It is
	// a subset of fScore.
	openFScore := newFScoreMap()
	openFScore.Set(start, distanceEstimator.Estimate(start, goal))

	for round := 0; len(openSet) > 0; round++ {
		logger.LogF("===round %v\n", round)
		logger.LogF("closed set %v\n", closedSet)
		logger.LogF("open set %v\n", openSet)
		logger.LogF("gScore %+v\n", gScore)
		logger.LogF("fScore %+v\n", fScore)

		current := ""
		var currentFScore uint
		openFScore.Walk(func(item *fScoreItem) bool {
			if _, found := openSet[item.Name]; found {
				current = item.Name
				currentFScore = item.Value
				return false
			}
			return true
		})
		if current == "" {
			panic("nothing found in fscore")
		}

		logger.LogF("current %v\n", current)

		if current == goal {
			return reconstructPath(cameFrom, current)
		}

		delete(openSet, current)
		openFScore.Delete(current, currentFScore)
		closedSet[current] = true

		currentGScore := gScore.GetWithDefault(current, math.MaxUint32)

		neighbors := neighborDiscoverer.AllNeighbors(current)
		for _, neighbor := range neighbors {
			if _, found := closedSet[neighbor]; found {
				continue
			}

			neighborGScore := currentGScore +
				distanceEstimator.NeighborDistance(current, neighbor)

			if _, found := openSet[neighbor]; !found {
				openSet[neighbor] = true
			} else if neighborGScore >= gScore.GetWithDefault(neighbor, math.MaxUint32) {
				logger.LogF("%v to %v isn't better\n", current, neighbor)
				continue // not a better path
			}

			// this path is the best until now. record it!
			cameFrom[neighbor] = current
			gScore[neighbor] = neighborGScore

			neighborFScore := neighborGScore + distanceEstimator.Estimate(neighbor, goal)
			fScore.Set(neighbor, neighborFScore)
			openFScore.Set(neighbor, neighborFScore)
		}
	}

	return nil // no path found
}
