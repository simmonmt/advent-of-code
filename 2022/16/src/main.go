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

// NOTES:
//  - This part 2 works, but it's incredibly slow. Hours slow. The only saving
//    grace is that the sorting it uses (prioritizing short trips to high-rate
//    nodes), helped by the short-circuiting here and there, leads it to the
//    highest-value path early on with my input.
//  - It's silly to store node names and use a map for visited when we could
//    just use a bitmask.
//  - It's not necessary (per reddit) to model two simultaneous actors. That
//    also adds time. Instead of trying to simulate passing time, generate all
//    possible paths then look for the best one? (we're going to go AA BB DD so
//    figure out however long that takes).
//  - This solution goes to great lengths to figure out the possible next
//    neighbor for each node. Except we also went to great lengths to make a
//    simplified graph where every node has a link to every other graph. So the
//    set of possible neighbors for a given node is just all-visited.

package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/graph"
	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	inputPattern = regexp.MustCompile(
		`^Valve ([^ ]+) has flow rate=(\d+); tunnels? leads? to valves? ([^ ,]+(?:, [^ ,]+)*)`)
)

type InputNode struct {
	Name  string
	Rate  int
	Dests []string
}

func parseInput(lines []string) ([]*InputNode, error) {
	out := []*InputNode{}
	for i, line := range lines {
		parts := inputPattern.FindStringSubmatch(line)
		if len(parts) == 0 {
			return nil, fmt.Errorf("%d: bad match", i+1)
		}

		name := parts[1]
		rate, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("%d: bad rate: %v", i+1, err)
		}

		dests := strings.Split(parts[3], ", ")

		out = append(out, &InputNode{name, rate, dests})
	}
	return out, nil
}

type InputGraph struct {
	nodes map[graph.NodeID]*InputNode
}

func (g *InputGraph) NeighborDistance(from, to graph.NodeID) int {
	return 1
}

func (g *InputGraph) Neighbors(id graph.NodeID) []graph.NodeID {
	inputNode, found := g.nodes[id]
	if !found {
		panic("unknown")
	}

	//logger.LogF("%v: input node %+v", id, inputNode)

	out := []graph.NodeID{}
	for _, dest := range inputNode.Dests {
		out = append(out, graph.NodeID(dest))
	}

	//logger.LogF("neighbors of %v: %v", id, out)
	return out
}

type SimpleGraph struct {
	edges map[graph.NodeID]map[graph.NodeID]int
	rates map[graph.NodeID]int
}

func NewSimpleGraph(nodes []*InputNode) *SimpleGraph {
	g := &SimpleGraph{
		edges: map[graph.NodeID]map[graph.NodeID]int{},
		rates: map[graph.NodeID]int{},
	}

	for _, node := range nodes {
		if node.Rate > 0 {
			g.rates[graph.NodeID(node.Name)] = node.Rate
		}
	}

	return g
}

func (g *SimpleGraph) AddEdge(from, to graph.NodeID, cost int) {
	if _, found := g.edges[from]; !found {
		g.edges[from] = map[graph.NodeID]int{}
	}
	g.edges[from][to] = cost
}

func (g *SimpleGraph) HasEdge(from, to graph.NodeID) bool {
	sub, found := g.edges[from]
	if !found {
		return false
	}
	_, found = sub[to]
	return found
}

func (g *SimpleGraph) EdgeCost(from, to graph.NodeID) int {
	sub, found := g.edges[from]
	if !found {
		panic("bad from")
	}
	cost, found := sub[to]
	if !found {
		panic("bad to")
	}
	return cost
}

func (g *SimpleGraph) AllEdges(from graph.NodeID) map[graph.NodeID]int {
	sub, found := g.edges[from]
	if !found {
		panic("bad from")
	}
	return sub
}

func (g *SimpleGraph) Rate(id graph.NodeID) int {
	return g.rates[id]
}

func simplifyInputGraph(nodes []*InputNode) *SimpleGraph {
	allNodes := map[graph.NodeID]*InputNode{}
	usableNodes := []graph.NodeID{"AA"}
	for _, node := range nodes {
		allNodes[graph.NodeID(node.Name)] = node
		if node.Rate > 0 {
			usableNodes = append(usableNodes, graph.NodeID(node.Name))
		}
	}

	inputGraph := &InputGraph{allNodes}
	simpleGraph := NewSimpleGraph(nodes)
	for _, from := range usableNodes {
		for _, to := range usableNodes {
			if from == to || to == "AA" {
				continue
			}
			if simpleGraph.HasEdge(from, to) {
				continue
			}

			//logger.LogF("%v to %v", from, to)

			path := graph.ShortestPath(graph.NodeID(from), graph.NodeID(to), inputGraph)
			if path == nil {
				panic(fmt.Sprintf("disconnect; no path from %v to %v", from, to))
			}

			cost := len(path) + 1 // + 1 to turn on dest valve

			simpleGraph.AddEdge(from, to, cost)
			if from != "AA" {
				simpleGraph.AddEdge(to, from, cost)
			}
		}
	}

	return simpleGraph
}

type PlayerState struct {
	dest         graph.NodeID
	minsToArrive int
}

func (s *PlayerState) Clone() *PlayerState {
	out := &PlayerState{
		dest:         s.dest,
		minsToArrive: s.minsToArrive,
	}

	return out
}

type Action struct {
	Player int
	Dest   graph.NodeID
}

func findAvailableNeighbors(g *SimpleGraph, id graph.NodeID, left int, claimed map[graph.NodeID]bool) []graph.NodeID {
	out := []graph.NodeID{}
	for dest, cost := range g.AllEdges(id) {
		if _, found := claimed[dest]; found {
			continue
		}
		if cost > left {
			continue
		}
		out = append(out, dest)
	}

	sort.Slice(out, func(i, j int) bool {
		if iCost, jCost := g.EdgeCost(id, out[i]), g.EdgeCost(id, out[j]); iCost < jCost {
			return true
		} else if iCost > jCost {
			return false
		}

		return g.Rate(out[i]) > g.Rate(out[j])
	})

	return out
}

func allPossibleActions(destsPerPlayer [][]graph.NodeID) [][]Action {
	// TODO: write general version of this
	groups := [][]Action{}

	if len(destsPerPlayer) == 1 {
		for _, neighbor := range destsPerPlayer[0] {
			groups = append(groups, []Action{Action{0, neighbor}})
		}
		return append(groups, []Action{}) // no action
	} else if len(destsPerPlayer) == 2 {
		iDests := destsPerPlayer[0]
		jDests := destsPerPlayer[1]

		singles := [][]Action{}

		for i := 0; i <= len(iDests); i++ {
			iEmpty := i == len(iDests)

			for j := 0; j <= len(jDests); j++ {
				jEmpty := j == len(jDests)

				if !iEmpty && !jEmpty && iDests[i] == jDests[j] {
					continue
				}

				if iEmpty && jEmpty {
					continue
				}
				group := []Action{}
				if !iEmpty {
					group = append(group, Action{0, iDests[i]})
				}
				if !jEmpty {
					group = append(group, Action{1, jDests[j]})
				}

				if len(group) == 1 {
					singles = append(singles, group)
				} else {
					groups = append(groups, group)
				}
			}
		}

		groups = append(groups, singles...)
		groups = append(groups, []Action{})

		return groups
	} else {
		panic("bad number of players")
	}

}

type Release struct {
	total int
	rate  int
}

func executeMinute(g *SimpleGraph, curMin, leftMin int, release *Release, players []PlayerState, claimed map[graph.NodeID]bool) (futures [][]Action) {
	// t1 start to BB             (end mtg=1)
	// t2 going to BB             (start mtg=1, end mtg=0 => rate++)
	// t3 BB on, start to CC      (start mtg=0 => new dest, end mtg=>1)
	// t4 BB on, going to CC      (start mtg=1, end mtg=0 => rate++)
	// t5 BB,CC on, ...           (start mtg=0 => new dest...)

	// total += rate

	// for player := range players {
	// 	if player.mtg == 0 {
	// 		player.dest = newDest
	// 		player.mtg = newMtg
	// 	} else if player.mtg > 0 {
	// 		player.mtg--
	// 		if player.mtg == 0 {
	// 			rate++
	// 		}
	// 	}
	// }

	release.total += release.rate

	perPlayerDests := [][]graph.NodeID{}

	hasDests := false
	for i := 0; i < len(players); i++ {
		state := &players[i]

		dests := []graph.NodeID{}
		if state.minsToArrive == 0 {
			cur := state.dest // we're already there
			dests = findAvailableNeighbors(g, cur, leftMin, claimed)
			hasDests = true
		} else {
			state.minsToArrive--
			if state.minsToArrive == 0 {
				// was >0, now 0, arrived at dest
				release.rate += g.Rate(state.dest)
			}
		}

		perPlayerDests = append(perPlayerDests, dests)
	}

	if !hasDests {
		return nil // all still in motion
	}

	futures = allPossibleActions(perPlayerDests)

	if curMin == 1 && len(players) == 1 {
		fmt.Printf("in %v\n", len(futures))
		out := [][]Action{}
		for _, future := range futures {
			if len(future) == 0 {
				continue
			}
			out = append(out, future)
		}
		fmt.Printf("out %v\n", len(out))
		return out
	}

	if curMin == 1 && len(players) == 2 {
		fmt.Printf("in %v\n", len(futures))
		singles := map[graph.NodeID]bool{}
		doubles := map[string]bool{}

		out := [][]Action{}
		for _, future := range futures {
			if len(future) == 0 {
				continue
			}

			if len(future) == 1 {
				if _, found := singles[future[0].Dest]; !found {
					singles[future[0].Dest] = true
					out = append(out, future)
				}
				continue
			}

			if len(future) != 2 {
				panic("bad future")
			}
			a, b := future[0].Dest, future[1].Dest
			if a > b {
				a, b = b, a
			}

			key := fmt.Sprintf("%v/%v", a, b)
			if _, found := doubles[key]; !found {
				doubles[key] = true
				out = append(out, future)
			}
		}

		fmt.Printf("out %v\n", len(out))

		return out
	}

	//fmt.Println(futures)
	return futures
}

type ProgressTracker struct {
	maxRelease  int
	allOpenRate int
	abandoned   int
	minAllOpen  int
}

func runWorld(g *SimpleGraph, curMin, maxMin int, numPlayers int, players *[2]PlayerState, release Release, progress *ProgressTracker, claimed map[graph.NodeID]bool) {
	if curMin > maxMin {
		if release.total > progress.maxRelease {
			fmt.Println("new max release", release.total)
			progress.maxRelease = release.total
		}
		return
	}

	if progress.maxRelease > 0 {
		minToGo := -1
		for i := 0; i < numPlayers; i++ {
			if minToGo == -1 || players[i].minsToArrive < minToGo {
				minToGo = players[i].minsToArrive
			}
		}

		total := release.total
		if minToGo > 0 {
			total += release.rate * minToGo
		}
		needed := progress.maxRelease - total
		if (maxMin-curMin+minToGo+1)*progress.allOpenRate < needed {
			// There's no way this subtree can set a new
			// maxRelease so give up on it.
			//logger.LogF("drop")
			return
		}
	}

	if release.rate == progress.allOpenRate {
		//fmt.Println("claimed", len(claimed))
		for i := 0; i < numPlayers; i++ {
			if players[i].minsToArrive != 0 {
				panic("bad")
			}
		}

		total := release.total + release.rate*(maxMin-curMin+1)
		if total > progress.maxRelease {
			fmt.Println("new max release", total, "all open", curMin)
			progress.maxRelease = total
		}
		return
	}

	futures := executeMinute(g, curMin, maxMin-curMin, &release,
		(*players)[0:numPlayers], claimed)

	if futures == nil {
		runWorld(g, curMin+1, maxMin, numPlayers, players, release, progress, claimed)
		return
	}

	// Set up the future (we're not in it until the call to runWorld). There
	// will always be a future in which nobody does anything.
	for i, future := range futures {
		if curMin == 1 {
			fmt.Printf("future %v of %v\n", i+1, len(futures))
		}

		start := len(claimed)

		futurePlayers := *players
		for _, action := range future {
			fp := &futurePlayers[action.Player]
			fp.minsToArrive = g.EdgeCost(fp.dest, action.Dest) - 1
			fp.dest = action.Dest

			if _, found := claimed[action.Dest]; found {
				panic("reclaim")
			}
			claimed[action.Dest] = true
		}

		if false && logger.Enabled() {
			logger.LogF("== Minute %d == ", curMin)
			logger.LogF("Total release %d (cur rate %d, max rate %d, max release %v)",
				release.total, release.rate, progress.allOpenRate,
				progress.maxRelease)
			for i := 0; i < numPlayers; i++ {
				if players[i].dest != futurePlayers[i].dest {
					logger.LogF("player %d was %v now %v",
						i, players[i], futurePlayers[i])
				} else {
					logger.LogF("player %d %v", i, futurePlayers[i])
				}
			}
		}

		// runs that world until the end
		runWorld(g, curMin+1, maxMin, numPlayers, &futurePlayers,
			release, progress, claimed)

		for _, action := range future {
			if _, found := claimed[action.Dest]; !found {
				panic("bad unclaim")
			}
			delete(claimed, action.Dest)
		}

		if len(claimed) != start {
			panic("mismatch")
		}
	}
}

func solve(g *SimpleGraph, maxMin int, numPlayers int, start graph.NodeID) int {
	players := [2]PlayerState{}
	for i := 0; i < numPlayers; i++ {
		players[i].dest = start
	}

	allOpenRate := 0
	for _, rate := range g.rates {
		allOpenRate += rate
	}

	progress := &ProgressTracker{
		maxRelease:  0,
		allOpenRate: allOpenRate,
	}

	runWorld(g, 1, maxMin, numPlayers, &players, Release{}, progress,
		map[graph.NodeID]bool{start: true})
	return progress.maxRelease
}

func solveA(nodes []*InputNode) int {
	g := simplifyInputGraph(nodes)
	return solve(g, 30, 1, "AA")
}

func solveB(nodes []*InputNode) int {
	g := simplifyInputGraph(nodes)
	return solve(g, 26, 2, "AA")
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		log.Fatal(err)
	}

	nodes, err := parseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("A", solveA(nodes))
	fmt.Println("B", solveB(nodes))
}
