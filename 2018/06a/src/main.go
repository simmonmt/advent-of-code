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
	verbose = flag.Bool("verbose", false, "verbose")
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

func outOfBounds(point Point, bounds *Bounds) bool {
	return point.X < bounds.minX || point.X > bounds.maxX ||
		point.Y < bounds.minY || point.Y > bounds.maxY
}

func findClosest(cand Point, points []Point) (bool, Point) {
	var minPoint Point
	minDouble := false
	minDist := -1

	for _, point := range points {
		dist := manhattanDistance(cand, point)
		//logger.LogF("-- dist %+v to %+v is %+v\n", cand, point, dist)
		if minDist == -1 || dist < minDist {
			minDouble = false
			minPoint = point
			minDist = dist
		} else if dist == minDist {
			minDouble = true
		}
	}

	return minDouble, minPoint
}

func findArea(point Point, points []Point, bounds *Bounds) (bool, int) {
	logger.LogF("\npoint %+v\n", point)

	seen := map[Point]bool{point: true}
	goodCands := []Point{point}
	area := 1
	for {
		cands := candidates(point, goodCands, seen)
		logger.LogF("points %+v cands %v\n", goodCands, cands)
		goodCands = []Point{}

		for _, cand := range cands {
			seen[cand] = true

			equidistant, closest := findClosest(cand, points)

			if equidistant {
				logger.LogF("-- cand %+v: equidistant to closest\n", cand)
				continue
			}

			if closest != point {
				logger.LogF("-- cand %+v: %+v closest\n", cand, closest)
				continue
			}

			// closest == point
			if outOfBounds(cand, bounds) {
				logger.LogF("-- cand %+v: oob (bounds %+v)\n", cand, bounds)
				return false, 0
			}

			logger.LogF("-- cand %+v: point closest\n", cand)
			area++
			goodCands = append(goodCands, cand)
		}

		if len(goodCands) == 0 {
			logger.LogF("-- no good candidates\n")
			break
		}

		logger.LogF("good cands: %+v\n", goodCands)
	}

	logger.LogF("-- finite, area %v\n", area)
	return true, area
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

	bounds := Bounds{minX, maxX, minY, maxY}

	maxArea := 0
	for _, point := range points {
		finite, area := findArea(point, points, &bounds)
		if !finite {
			fmt.Printf("point %+v infinite\n", point)
			continue
		}

		fmt.Printf("point %+v area %v\n", point, area)
		maxArea = intmath.IntMax(maxArea, area)
	}

	fmt.Println(maxArea)
}
