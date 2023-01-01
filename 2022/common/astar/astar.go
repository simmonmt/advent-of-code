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

// AStar search algorithm
//
// Implemented from pseudocode on
// https://en.wikipedia.org/wiki/A*_search_algorithm

package astar

import (
	"math"

	"github.com/simmonmt/aoc/2022/common/collections"
	"github.com/simmonmt/aoc/2022/common/logger"
)

type ClientInterface interface {
	AllNeighbors(node string) []string
	EstimateDistance(start, end string) uint

	// NeighborDistance returns the distance between two known direct
	// neighbors (i.e. a pair derived using AllNeighbors).
	NeighborDistance(n1, n2 string) uint
	GoalReached(cand, goal string) bool
}

type scoreMap map[string]uint

func (m *scoreMap) Get(key string) uint {
	if v, found := (*m)[key]; found {
		return v
	} else {
		return math.MaxUint32
	}
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

func AStar(start, goal string, client ClientInterface) []string {
	logger.LogF("astar start %v goal %v", start, goal)

	openSet := collections.NewPriorityQueue[string](collections.LessThan)
	cameFrom := map[string]string{}

	gScore := scoreMap{}
	gScore[start] = 0

	fScore := scoreMap{}

	startEstimate := client.EstimateDistance(start, goal)
	openSet.Insert(start, int(startEstimate))
	fScore[start] = startEstimate

	for round := 0; !openSet.IsEmpty(); round++ {
		logger.LogF("===round %v", round)
		logger.LogF("open set %v", openSet)
		logger.LogF("gScore %+v", gScore)
		logger.LogF("fScore %+v", fScore)

		current, _ := openSet.Next()
		if client.GoalReached(current, goal) {
			return reconstructPath(cameFrom, current)
		}

		currentGScore := gScore.Get(current)

		neighbors := client.AllNeighbors(current)
		logger.LogF("neighbors of %v: %v", current, neighbors)
		for _, neighbor := range neighbors {
			neighborGScore := currentGScore +
				client.NeighborDistance(current, neighbor)

			if neighborGScore >= gScore.Get(neighbor) {
				logger.LogF("%v to %v isn't better", current, neighbor)
				continue // not a better path
			}

			// this path is the best until now. record it!
			cameFrom[neighbor] = current
			gScore[neighbor] = neighborGScore

			neighborFScore := neighborGScore + client.EstimateDistance(neighbor, goal)
			fScore[neighbor] = neighborFScore
			openSet.Insert(neighbor, int(neighborFScore))
		}
	}

	return nil // no path found
}
