// AStar search algorithm
//
// Implemented from pseudocode on
// https://en.wikipedia.org/wiki/A*_search_algorithm

package astar

import (
	"math"

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
	closedSet := map[string]bool{}
	openSet := map[string]bool{start: true}
	cameFrom := map[string]string{}

	gScore := scoreMap{}
	gScore[start] = 0

	fScore := scoreMap{}
	fScore[start] = distanceEstimator.Estimate(start, goal)

	for round := 0; len(openSet) > 0; round++ {
		logger.LogF("===round %v\n", round)
		logger.LogF("closed set %v\n", closedSet)
		logger.LogF("open set %v\n", openSet)
		logger.LogF("gScore %+v\n", gScore)
		logger.LogF("fScore %+v\n", fScore)

		// This is pretty inefficient. We can't simply use Go's
		// container/heap priority queue implementation as-is because
		// that implementation is missing methods that let us:
		//  - find the lowest-priority item that's in a given set. this
		//    probably requires a Walk method that lets us traverse the
		//    heap in order.
		//  - update an arbitrary node in the heap (or remove an
		//    arbitrary node)
		var currentFScore uint = math.MaxUint32
		current := ""
		for open := range openSet {
			if score := fScore.GetWithDefault(open, math.MaxUint32); score < currentFScore {
				current = open
				currentFScore = score
			}
		}

		logger.LogF("current %v\n", current)

		if current == goal {
			return reconstructPath(cameFrom, current)
		}

		delete(openSet, current)
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
			fScore[neighbor] = neighborGScore + distanceEstimator.Estimate(neighbor, goal)
		}
	}

	return nil // no path found
}
