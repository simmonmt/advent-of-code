// Copyright 2024 Google LLC
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
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/simmonmt/aoc/2025/common/collections"
	"github.com/simmonmt/aoc/2025/common/logger"
	"github.com/simmonmt/aoc/2025/common/pos"
)

type ClientInterface[T any] interface {
	AllNeighbors(node T) []T
	EstimateDistance(start, end T) uint

	// NeighborDistance returns the distance between two known direct
	// neighbors (i.e. a pair derived using AllNeighbors).
	NeighborDistance(n1, n2 T) uint
	GoalReached(cand, goal T) bool

	Serialize(val T) string
	Deserialize(val string) (T, error)
}

type scoreMap map[string]uint

func (m *scoreMap) Get(key string) uint {
	if v, found := (*m)[key]; found {
		return v
	} else {
		return math.MaxUint32
	}
}

func reconstructPath[T any](cameFrom map[string]string, current T, client ClientInterface[T]) []T {
	totalPath := []T{current}
	for {
		nextStr, found := cameFrom[client.Serialize(current)]
		if !found {
			break
		}

		next, err := client.Deserialize(nextStr)
		if err != nil {
			panic("bad nextStr")
		}

		totalPath = append(totalPath, next)
		current = next

	}
	return totalPath
}

type AStar[T any] struct {
	client      ClientInterface[T]
	start, goal T
	numRounds   int

	gScore   scoreMap
	openSet  collections.PriorityQueue[string]
	cameFrom map[string]string
}

func New[T any](start, goal T, client ClientInterface[T]) *AStar[T] {
	return &AStar[T]{
		client:    client,
		start:     start,
		goal:      goal,
		numRounds: -1,
	}
}

func (a *AStar[T]) SetNumRounds(numRounds int) {
	if a.openSet != nil {
		panic("reuse")
	}
	a.numRounds = numRounds
}

func (a *AStar[T]) Solve() []T {
	if a.openSet != nil {
		panic("reuse")
	}

	startStr, goalStr := a.client.Serialize(a.start), a.client.Serialize(a.goal)
	logger.Infof("astar start %v goal %v", startStr, goalStr)

	a.openSet = collections.NewPriorityQueue[string](collections.LessThan)
	a.cameFrom = map[string]string{}

	a.gScore = scoreMap{}
	a.gScore[startStr] = 0

	startEstimate := a.client.EstimateDistance(a.start, a.goal)
	a.openSet.Insert(startStr, int(startEstimate))

	for round := 0; !a.openSet.IsEmpty() && (a.numRounds < 0 || round < a.numRounds); round++ {
		logger.Infof("===round %v", round)
		logger.Infof("open set %v", a.openSet)
		logger.Infof("gScore %+v", a.gScore)

		currentStr, _ := a.openSet.Next()
		current, err := a.client.Deserialize(currentStr)
		if err != nil {
			panic("bad current")
		}

		if a.client.GoalReached(current, a.goal) {
			return reconstructPath(a.cameFrom, current, a.client)
		}

		currentGScore := a.gScore.Get(currentStr)

		neighbors := a.client.AllNeighbors(current)
		neighborStrs := make([]string, len(neighbors))
		for i, neighbor := range neighbors {
			neighborStrs[i] = a.client.Serialize(neighbor)
		}

		logger.Infof("neighbors of %v: %v", currentStr, neighborStrs)
		for i, neighbor := range neighbors {
			neighborStr := neighborStrs[i]

			neighborGScore := currentGScore +
				a.client.NeighborDistance(current, neighbor)

			if neighborGScore >= a.gScore.Get(neighborStr) {
				logger.Infof("%v to %v isn't better", currentStr, neighborStr)
				continue // not a better path
			}

			// this path is the best until now. record it!
			a.cameFrom[neighborStr] = currentStr
			a.gScore[neighborStr] = neighborGScore

			neighborFScore := neighborGScore + a.client.EstimateDistance(neighbor, a.goal)
			a.openSet.Insert(neighborStr, int(neighborFScore))
		}
	}

	return nil // no path found
}

func (a *AStar[T]) Dump(path string, height, width int, background color.NRGBA, cb func(val T, score, maxScore uint, img *image.NRGBA) (pos.P2, color.NRGBA)) error {
	if a.openSet == nil {
		panic("unsolved")
	}

	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, background)
		}
	}

	maxScore := uint(0)
	for _, score := range a.gScore {
		maxScore = max(maxScore, score)
	}

	for key, score := range a.gScore {
		node, err := a.client.Deserialize(key)
		if err != nil {
			return fmt.Errorf("bad key %v", key)
		}

		p, newColor := cb(node, score, maxScore, img)
		img.Set(p.X, p.Y, newColor)
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}

	return f.Close()
}
