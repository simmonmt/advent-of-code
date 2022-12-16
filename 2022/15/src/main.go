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

package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/simmonmt/aoc/2022/common/area"
	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/logger"
	"github.com/simmonmt/aoc/2022/common/mtsmath"
	"github.com/simmonmt/aoc/2022/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	inputPattern = regexp.MustCompile(
		`^Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)`)
)

type Sensor struct {
	Loc    pos.P2
	Beacon pos.P2
}

func parseInput(lines []string) ([]Sensor, error) {
	sensors := []Sensor{}
	for i, line := range lines {
		parts := inputPattern.FindStringSubmatch(line)
		if len(parts) != 5 {
			return nil, fmt.Errorf("%d: bad match", i+1)
		}

		locX, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("%d: bad loc x: %v", i+1, err)
		}
		locY, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("%d: bad loc y: %v", i+1, err)
		}
		beaconX, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, fmt.Errorf("%d: bad beacon x: %v", i+1, err)
		}
		beaconY, err := strconv.Atoi(parts[4])
		if err != nil {
			return nil, fmt.Errorf("%d: bad beacon y: %v", i+1, err)
		}

		sensors = append(sensors, Sensor{
			Loc:    pos.P2{locX, locY},
			Beacon: pos.P2{beaconX, beaconY},
		})
	}
	return sensors, nil
}

func furthestIntersections(sensor Sensor, row int) (left, right int) {
	beaconDist := sensor.Loc.ManhattanDistance(sensor.Beacon)
	rowDist := mtsmath.Abs(row - sensor.Loc.Y)
	delta := beaconDist - rowDist
	if delta == 0 {
		return sensor.Loc.X, sensor.Loc.X
	}

	left = sensor.Loc.X - delta
	right = sensor.Loc.X + delta
	return
}

func findOccupiedInRow(sensors []Sensor, row int) []area.Area1D {
	ranges := []area.Area1D{}
	for _, sensor := range sensors {
		beaconDist := sensor.Loc.ManhattanDistance(sensor.Beacon)

		dist := sensor.Loc.ManhattanDistance(pos.P2{X: sensor.Loc.X, Y: row})
		if dist > beaconDist {
			//logger.LogF("sensor %v too far", sensor)
			continue
		}
		//logger.LogF("sensor %v close enough", sensor)

		left, right := furthestIntersections(sensor, row)
		//logger.LogF("intersections %v, %v", left, right)

		ranges = append(ranges, area.Area1D{left, right})
	}

	//logger.LogF("%d ranges: %v", len(ranges), ranges)
	ranges = area.Merge1DRanges(ranges)
	//logger.LogF("post-merge %d ranges: %v", len(ranges), ranges)
	return ranges
}

func solveA(sensors []Sensor, activeRow int) int {
	beacons := map[pos.P2]bool{}
	for _, sensor := range sensors {
		beacons[sensor.Beacon] = true
	}

	ranges := findOccupiedInRow(sensors, activeRow)

	numInRow := 0
	for beacon := range beacons {
		if beacon.Y != activeRow {
			continue
		}

		for _, r := range ranges {
			if r.ContainsVal(beacon.X) {
				numInRow++
				break
			}
		}
	}
	logger.LogF("found %d beacons in row", numInRow)

	out := -numInRow
	for _, r := range ranges {
		out += r.To - r.From + 1
	}

	return out
}

func filterRanges(in []area.Area1D, min, max int) []area.Area1D {
	out := []area.Area1D{}
	for _, r := range in {
		if r.To < min || r.From > max {
			continue
		}
		r.From = mtsmath.Max(min, r.From)
		r.To = mtsmath.Min(max, r.To)
		out = append(out, r)
	}
	return out
}

func solveB(sensors []Sensor, min, max int) int {
	for y := min; y <= max; y++ {
		// The row with a blank spot will have one (and only one) value
		// not covered by a range. So we look for rows with two ranges
		// that are within [min,max].
		ranges := findOccupiedInRow(sensors, y)
		ranges = filterRanges(ranges, min, max)

		// Rows with no beacon location will have one range that's
		// full-width.
		if len(ranges) == 1 && ranges[0].Size() == max-min+1 {
			continue
		}

		// From this point on we'll just focus on our assumption that
		// we've got a two-range row with a gap of 1 between them, and
		// we'll panic if we see anything different.

		logger.LogF("y=%d not full-width %v", y, ranges)
		if len(ranges) != 2 {
			panic("multi range")
		}

		gap := 0
		if ranges[0].To+2 == ranges[1].From {
			gap = ranges[0].To + 1
		} else if ranges[1].To+2 == ranges[0].From {
			gap = ranges[1].To + 1
		} else {
			panic("big gap")
		}

		p := pos.P2{gap, y}
		logger.LogF("pos %v", p)
		return p.X*4000000 + y
	}

	return -1
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

	sensors, err := parseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("A", solveA(sensors, 2000000))
	fmt.Println("B", solveB(sensors, 0, 4000000))
}
