// Copyright 2021 Google LLC
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

package puzzle

import (
	"math"
	"sort"

	"github.com/simmonmt/aoc/2019/common/logger"
)

func findVisible(ctr Pos, all map[Pos]bool) map[Pos][]Pos {
	visible := map[Pos][]Pos{}

	for p := range all {
		if p.X == ctr.X && p.Y == ctr.Y {
			continue
		}

		rel := Pos{
			X: p.X - ctr.X,
			Y: p.Y - ctr.Y,
		}

		var slope Pos
		slope.Y, slope.X = Factor(rel.Y, rel.X)

		if _, found := visible[slope]; !found {
			visible[slope] = []Pos{rel}
		} else {
			visible[slope] = append(visible[slope], rel)
		}
	}

	return visible
}

const (
	negOneHalfPi = -math.Pi / 2.0
)

func FindAll(ctr Pos, all map[Pos]bool) []Pos {
	visible := findVisible(ctr, all)
	angleToSlope := map[float64]Pos{}
	angles := []float64{}

	logger.LogF("visible: %v", visible)

	for s := range visible {
		// Sort the groups with the same slope so the first is
		// closest. This will make it easier to decide which to vaporize
		// first.
		sort.Sort(ByManhattanOriginDistance(visible[s]))

		// atan2 returns values (-pi,pi] with 0 on the right (+x,0), pi
		// on the left (-x,0), and -pi/2 on the top (0,-y). It's a nicer
		// version of atan that deals well with zeroes in the numerator
		// or denominator. See https://en.wikipedia.org/wiki/Atan2
		angle := math.Atan2(float64(s.Y), float64(s.X))
		logger.LogF("got angle %f for %+v", angle, s)
		angleToSlope[angle] = s
		angles = append(angles, angle)
	}

	sort.Float64s(angles)
	for _, a := range angles {
		logger.LogF("%f: %+v", a, angleToSlope[a])
	}

	// We want to walk angles CW -pi/2, so find the starting point. Angles
	// is sorted, with values (-pi,pi), so we can stop at the first value >=
	// -pi/2.
	var startIdx int
	for i, theta := range angles {
		if theta >= negOneHalfPi {
			startIdx = i
			break
		}
	}

	vapes := []Pos{}

	vaporized := true
	for vaporized {
		vaporized = false
		for i := 0; i < len(angles); i++ {
			angle := angles[(i+startIdx)%len(angles)]
			slope := angleToSlope[angle]

			if len(visible[slope]) == 0 {
				continue
			}

			relVape := visible[slope][0]
			vapes = append(vapes, Pos{ctr.X + relVape.X, ctr.Y + relVape.Y})
			visible[slope] = visible[slope][1:]
			vaporized = true
		}
	}

	return vapes
}

func FindNumVisible(ctr Pos, all map[Pos]bool) int {
	visible := findVisible(ctr, all)
	return len(visible)
}

func FindBest(all map[Pos]bool) (Pos, int) {
	var bestPos Pos
	bestVisible := -1
	for ctr := range all {
		numVisible := FindNumVisible(ctr, all)
		//logger.LogF("%+v visible %v", ctr, numVisible)
		if bestVisible == -1 || bestVisible < numVisible {
			bestPos = ctr
			bestVisible = numVisible
		}
	}
	return bestPos, bestVisible
}
