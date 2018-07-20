package maze

import (
	"crypto/md5"
	"fmt"

	"astar"
	"intmath"
	"node"
)

type aStarHelper struct {
	width, height int
}

func NewHelper(width, height int) *aStarHelper {
	return &aStarHelper{
		width:  width,
		height: height,
	}
}

func (h *aStarHelper) tryNeighbor(start *node.Node, dir string, xDelta, yDelta int, val byte, neighbors *[]string) {
	if val < 0xb {
		return
	}

	n := &node.Node{
		X:        start.X + xDelta,
		Y:        start.Y + yDelta,
		Passcode: start.Passcode,
		Path:     start.Path + dir,
	}

	if n.X < 0 || n.Y < 0 {
		return
	}
	if n.X >= h.width || n.Y >= h.height {
		return
	}

	*neighbors = append(*neighbors, n.Serialize())
}

func (h *aStarHelper) AllNeighbors(start string) []string {
	neighbors := []string{}

	n, err := node.Deserialize(start)
	if err != nil {
		panic(fmt.Sprintf(`failed to deserialize "%v": %v`, start, err))
	}

	hash := md5.Sum([]byte(n.Passcode + n.Path))
	h.tryNeighbor(n, "U", 0, -1, hash[0]>>4, &neighbors)
	h.tryNeighbor(n, "D", 0, 1, hash[0]&0xf, &neighbors)
	h.tryNeighbor(n, "L", -1, 0, hash[1]>>4, &neighbors)
	h.tryNeighbor(n, "R", 1, 0, hash[1]&0xf, &neighbors)

	return neighbors
}

func (h *aStarHelper) EstimateDistance(start, end string) uint {
	sNode, err := node.Deserialize(start)
	if err != nil {
		panic(fmt.Sprintf("failed to deserialize start %v: %v", start, err))
	}

	eNode, err := node.Deserialize(end)
	if err != nil {
		panic(fmt.Sprintf("failed to deserialize end %v: %v", end, err))
	}

	return uint(intmath.Abs(sNode.X-eNode.X) + intmath.Abs(sNode.Y-eNode.Y))
}

func (h *aStarHelper) NeighborDistance(n1, n2 string) uint {
	if n1 == n2 {
		return 0
	} else {
		return 1
	}
}

func (h *aStarHelper) GoalReached(cand, goal string) bool {
	cNode, err := node.Deserialize(cand)
	if err != nil {
		panic(fmt.Sprintf("failed to deserialize cand %v: %v", cand, err))
	}

	return cNode.X == h.width-1 && cNode.Y == h.height-1
}

func longestPath(start string, client astar.ClientInterface) string {
	openSet := map[string]bool{start: true}
	scores := map[string]uint{start: 0}

	for round := 0; len(openSet) > 0; round++ {
		current := ""
		for k := range openSet {
			current = k
			break
		}

		currentScore := scores[current]
		delete(openSet, current)

		if client.GoalReached(current, current) {
			continue
		}

		neighbors := client.AllNeighbors(current)
		for _, neighbor := range neighbors {
			var scoreToNeighbor uint = currentScore + client.NeighborDistance(current, neighbor)
			scores[neighbor] = scoreToNeighbor
			openSet[neighbor] = true
		}
	}

	// The path is in the node name, so there are many possible nodes in the
	// scores map that terminate at the goal location. Find the one with the
	// highest score.
	var highestScore uint
	highestNode := ""
	for n, s := range scores {
		if !client.GoalReached(n, n) {
			continue
		}
		if s > highestScore {
			highestScore = s
			highestNode = n
		}
	}

	return highestNode
}

func RunMaze(width, height int, passcode string) (found bool, lastStep string) {
	helper := NewHelper(width, height)

	start := node.New(0, 0, passcode)

	// StackExchange suggests there's a way to update A* to make it find the
	// longest path, but I was unable to bend it to my wishes. The graph is
	// directed and acyclic thanks to the paths changing the possible doors
	// when a room is revisited. Also the graph is relatively small (the
	// example paths are in the hundreds of steps), so I went with a
	// brute-force find-all-longest-paths approach.
	longest := longestPath(start.Serialize(), helper)
	if longest == "" {
		return false, ""
	}

	longestNode, err := node.Deserialize(longest)
	if err != nil {
		panic(fmt.Sprintf("failed to deserialize last step %v", longestNode))
	}

	return true, longestNode.Path
}
