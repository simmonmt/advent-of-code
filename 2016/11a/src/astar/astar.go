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
	closedSet := map[string]bool{}
	openSet := map[string]bool{start: true}
	cameFrom := map[string]string{}

	gScore := scoreMap{}
	gScore[start] = 0

	fScore := newFScoreMap()
	fScore.Set(start, distanceEstimator.Estimate(start, goal))

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
