package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"intmath"
	"logger"
)

var (
	verbose  = flag.Bool("verbose", false, "verbose")
	distance = flag.Int("distance", 32, "distance")
)

type Point struct {
	X, Y int
}

type Bounds struct {
	minX, maxX int
	minY, maxY int
}

func readPoints() ([]Point, error) {
	points := []Point{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ",", 2)
		x := intmath.AtoiOrDie(parts[0])
		y := intmath.AtoiOrDie(strings.TrimSpace(parts[1]))
		points = append(points, Point{x, y})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return points, nil
}

func manhattanDistance(p1, p2 Point) int {
	return intmath.Abs(p1.X-p2.X) + intmath.Abs(p1.Y-p2.Y)
}

var (
	neighborDeltas = []Point{Point{-1, 0}, Point{1, 0}, Point{0, -1}, Point{0, 1}}
)

func candidates(center Point, points []Point, seen map[Point]bool) []Point {
	found := map[Point]bool{}

	for _, point := range points {
		pointDist := manhattanDistance(point, center)
		for _, delta := range neighborDeltas {
			cand := Point{point.X + delta.X, point.Y + delta.Y}
			if _, wasSeen := seen[cand]; wasSeen {
				continue
			}

			// Is this necessary? does the seen array suffice?
			candDist := manhattanDistance(cand, center)
			if candDist >= pointDist {
				found[cand] = true
			}
		}
	}

	for _, point := range points {
		if _, ok := found[point]; ok {
			panic("found point in found")
		}
	}

	out := make([]Point, len(found))
	i := 0
	for point, _ := range found {
		out[i] = point
		i++
	}

	return out
}

func inBounds(cand Point, points []Point, maxDistance int) bool {
	sum := 0
	for _, point := range points {
		sum += manhattanDistance(cand, point)
	}
	logger.LogF("inbounds cand %+v sum %+v\n", cand, sum)
	return sum < maxDistance
}

func findArea(center Point, points []Point, maxDistance int) int {
	logger.LogF("starting at %+v\n", center)

	counted := map[Point]bool{center: true}
	goodCands := []Point{center}
	area := 1
	for {
		cands := candidates(center, goodCands, counted)
		logger.LogF("points %+v cands %v\n", goodCands, cands)
		goodCands = []Point{}

		for _, cand := range cands {
			if inBounds(cand, points, maxDistance) {
				logger.LogF("-- cand: %+v: yes\n", cand)
				area++
				counted[cand] = true
				goodCands = append(goodCands, cand)
			} else {
				logger.LogF("-- cand: %+v: no\n", cand)
			}
		}

		if len(goodCands) == 0 {
			logger.LogF("-- no good candidates\n")
			break
		}

		logger.LogF("good cands: %+v\n", goodCands)
	}

	return area
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	points, err := readPoints()
	if err != nil {
		log.Fatalf("failed to read points: %v", points)
	}

	var minX int = math.MaxInt32
	var minY int = math.MaxInt32
	maxX := 0
	maxY := 0
	for _, point := range points {
		minX = intmath.IntMin(minX, point.X)
		minY = intmath.IntMin(minY, point.Y)
		maxX = intmath.IntMax(maxX, point.X)
		maxY = intmath.IntMax(maxY, point.Y)
	}

	center := Point{
		X: minX + (maxX-minX)/2,
		Y: minY + (maxY-minY)/2,
	}

	fmt.Println(findArea(center, points, *distance))
}
