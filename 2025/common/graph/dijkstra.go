package graph

import "github.com/simmonmt/aoc/2025/common/collections"

func reverseSlice[T any](in []T) []T {
	out := make([]T, len(in))
	for i, j := 0, len(in)-1; j >= 0; i, j = i+1, j-1 {
		out[j] = in[i]
	}
	return out
}

func ShortestPath(start, end NodeID, graph Graph) []NodeID {
	visited := map[NodeID]bool{}

	queue := collections.NewPriorityQueue[NodeID](collections.LessThan)
	queue.Insert(start, 0)

	distances := map[NodeID]int{}
	distances[start] = 0

	froms := map[NodeID]NodeID{}

	for !queue.IsEmpty() {
		cur, _ := queue.Next()
		curDist := distances[cur]

		for _, neighbor := range graph.Neighbors(cur) {
			if _, found := visited[neighbor]; found {
				continue
			}

			throughCurDist := curDist + graph.NeighborDistance(cur, neighbor)
			neighborDist, found := distances[neighbor]
			if !found || throughCurDist < neighborDist {
				distances[neighbor] = throughCurDist
				queue.Insert(neighbor, throughCurDist)
				froms[neighbor] = cur
			}
		}

		if cur == end {
			revPath := []NodeID{}
			for id := end; id != start; id = froms[id] {
				revPath = append(revPath, id)
			}
			return reverseSlice(revPath)
		}

		visited[cur] = true
	}

	return nil // no path found
}
