package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"

	"intmath"
	"logger"
	"maze"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
)

func readInput() ([]string, error) {
	lines := []string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lines, nil
}

func addPair(dists map[int]map[int]int, from, to, dist int) {
	if _, found := dists[from]; !found {
		dists[from] = map[int]int{}
	}
	dists[from][to] = dist
}

func doShortestRoute(start, end int, s map[int]bool, dists map[int]map[int]int) int {
	if len(s) == 2 {
		// If size of S is 2, then S must be {1, i},
		//  C(S, i) = dist(1, i)
		start, end := -1, -1
		for n := range s {
			if start == -1 {
				start = n
			} else {
				end = n
			}
		}

		dist, found := dists[start][end]
		if !found {
			panic("unknown")
		}
		return dist
	}

	minDist := math.MaxInt32

	// // Else if size of S is greater than 2.
	// //  C(S, i) = min { C(S-{i}, j) + dis(j, i)} where j belongs to S, j != i and j != 1.
	for num := range s {
		if num == start || num == end {
			continue
		}

		sub := map[int]bool{}
		for n := range s {
			if n != end {
				sub[n] = true
			}
		}

		dist, found := dists[num][end]
		if !found {
			continue
		}

		minDist = intmath.IntMin(minDist, doShortestRoute(start, num, sub, dists)+dist)
	}

	return minDist
}

func shortestRoute(start int, dists map[int]map[int]int) int {
	dist := math.MaxInt32

	s := map[int]bool{}
	for n := range dists {
		s[n] = true
	}

	for end := range dists {
		if end == start {
			continue
		}

		dist = intmath.IntMin(dist, doShortestRoute(start, end, s, dists))
	}

	return dist
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	lines, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	board := maze.NewBoard(lines)
	board.Dump()

	dists := map[int]map[int]int{}

	nums := board.Nums()
	for i := range nums {
		for j := i + 1; j < len(nums); j++ {
			from := nums[i]
			to := nums[j]

			steps, found := board.ShortestPath(from, to)
			if !found {
				continue
			}

			addPair(dists, from, to, steps)
			addPair(dists, to, from, steps)
		}
	}

	fmt.Println(shortestRoute(0, dists))
}
