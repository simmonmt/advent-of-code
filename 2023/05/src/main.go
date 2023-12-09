// Copyright 2023 Google LLC
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
	"sort"
	"strings"

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/strutil"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Range struct {
	Lo, Hi int64
}

func (r Range) Before(or Range) bool {
	return r.Hi < or.Lo
}

func (r Range) ContainedBy(or Range) bool {
	return or.Lo <= r.Lo && or.Hi >= r.Hi
}

func (r Range) ContainsVal(v int64) bool {
	return v >= r.Lo && v <= r.Hi
}

func (r Range) Overlaps(or Range) bool {
	if or.Hi < r.Lo || or.Lo > r.Hi {
		return false
	}

	return true
}

func (r Range) Merge(or Range) Range {
	return Range{
		Lo: min(r.Lo, or.Lo),
		Hi: max(r.Hi, or.Hi),
	}
}

type RangeMap struct {
	Src, Dest Range
}

type GardenMap struct {
	From, To string
	Ranges   []*RangeMap
}

func (gm *GardenMap) String() string {
	ranges := []string{}
	for _, r := range gm.Ranges {
		ranges = append(ranges, fmt.Sprintf("%+v", r))
	}

	return fmt.Sprintf("%s=>%s: %s", gm.From, gm.To, strings.Join(ranges, ","))
}

func parseGardenMap(lines []string) (*GardenMap, error) {
	if len(lines) < 2 {
		return nil, fmt.Errorf("map too short")
	}

	names, _, found := strings.Cut(lines[0], " ")
	parts := strings.Split(names, "-")
	if !found || len(parts) != 3 {
		return nil, fmt.Errorf("bad description line")
	}

	gm := &GardenMap{
		From:   parts[0],
		To:     parts[2],
		Ranges: []*RangeMap{},
	}

	for i, line := range lines[1:] {
		nums, err := strutil.ListOfNumbers(line)
		if err != nil || len(nums) != 3 {
			return nil, fmt.Errorf("bad map line %d for %s", i+1, lines[0])
		}

		destStart, srcStart, rangeLen := nums[0], nums[1], nums[2]

		gr := &RangeMap{
			Src:  Range{Lo: int64(srcStart), Hi: int64(srcStart + rangeLen - 1)},
			Dest: Range{Lo: int64(destStart), Hi: int64(destStart + rangeLen - 1)},
		}
		gm.Ranges = append(gm.Ranges, gr)
	}

	return gm, nil
}

func findGardenMaps(lines []string) [][]string {
	out := [][]string{}

	start := 0
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			out = append(out, lines[start:i])
			start = i + 1
		}
	}

	out = append(out, lines[start:])
	return out
}

func parseInput(lines []string) ([]int64, []*GardenMap, error) {
	if len(lines) < 3 {
		return nil, nil, fmt.Errorf("too short input")
	}

	_, numsStr, found := strings.Cut(lines[0], ": ")
	nums, err := strutil.ListOfNumbers(numsStr)
	if !found || err != nil {
		return nil, nil, fmt.Errorf("bad seeds line")
	}

	gardenMaps := []*GardenMap{}
	for i, mapLines := range findGardenMaps(lines[2:]) {
		gardenMap, err := parseGardenMap(mapLines)
		if err != nil {
			return nil, nil, fmt.Errorf("bad garden map %d: %v", i+1, err)
		}

		gardenMaps = append(gardenMaps, gardenMap)
	}

	outNums := []int64{}
	for _, num := range nums {
		outNums = append(outNums, int64(num))
	}

	return outNums, gardenMaps, nil
}

func mapRange(r Range, rm *RangeMap) Range {
	return Range{
		Lo: r.Lo - rm.Src.Lo + rm.Dest.Lo,
		Hi: r.Hi - rm.Src.Hi + rm.Dest.Hi,
	}
}

func solveLevelForVal(val int64, gardenMap *GardenMap) int64 {
	for _, r := range gardenMap.Ranges {
		if r.Src.ContainsVal(val) {
			return val - r.Src.Lo + r.Dest.Lo
		}
	}
	return val
}

func solveLevel(ranges []Range, gardenMap *GardenMap) []Range {
	logger.Infof("solving %v with %v", ranges, gardenMap)

	out := []Range{}
	gmri := 0

	consumeRange0 := func() {
		if len(ranges) == 1 {
			ranges = []Range{}
		} else {
			ranges = ranges[1:]
		}
	}

	for gmri < len(gardenMap.Ranges) && len(ranges) > 0 {
		gmr := gardenMap.Ranges[gmri]

		if ranges[0].Before(gmr.Src) {
			logger.Infof("range %+v is before gmr %+v", ranges[0], gmr)
			out = append(out, ranges[0])
			consumeRange0()
			continue
		}

		// Handle parts of ranges[0] that precede gmr
		if ranges[0].Lo < gmr.Src.Lo {
			logger.Infof("range %+v starts before but overlaps gmr %+v", ranges[0], gmr)
			out = append(out, Range{ranges[0].Lo, gmr.Src.Lo - 1})
			ranges[0].Lo = gmr.Src.Lo
			continue
		}

		if ranges[0].ContainedBy(gmr.Src) {
			logger.Infof("range %+v contained by gmr %+v", ranges[0], gmr)
			// ends at or before gmr end; map and consume all
			out = append(out, mapRange(ranges[0], gmr))
			consumeRange0()
			continue
		}

		if ranges[0].Overlaps(gmr.Src) {
			// extends beyond gmr so consume what we can
			logger.Infof("range %+v extends beyond gmr %+v", ranges[0], gmr)
			out = append(out, mapRange(Range{Lo: ranges[0].Lo, Hi: gmr.Src.Hi}, gmr))
			ranges[0].Lo = gmr.Src.Hi + 1
			gmri++
			continue
		}

		// ranges[0] comes after gmr
		logger.Infof("range %+v comes after gmr %+v", ranges[0], gmr)
		gmri++
	}

	for len(ranges) > 0 {
		logger.Infof("range %+v after all", ranges[0])
		out = append(out, ranges[0])
		consumeRange0()
	}

	logger.Infof("%v-to-%v returning %+v", gardenMap.From, gardenMap.To, out)
	return out
}

func combineRanges(in []Range) []Range {
	out := []Range{in[0]}

	ii, oi := 1, 0
	for ii < len(in) {
		if in[ii].Overlaps(out[oi]) {
			logger.Infof("merging %v into %v", in[ii], out[oi])
			out[oi] = out[oi].Merge(in[ii])
		} else {
			out = append(out, in[ii])
			oi += 1
		}
		ii += 1
	}

	return out
}

func checkResult(ranges, gotRanges []Range, gardenMap *GardenMap) {
	for i, r := range ranges {
		if r.Lo != r.Hi {
			continue
		}

		gr := gotRanges[i].Lo

		if got := solveLevelForVal(r.Lo, gardenMap); got != gr {
			logger.Fatalf("%d = full %d single %d from %v", r.Lo, gr, got, gardenMap)
		}
	}
}

func solveForLocation(ranges []Range, gardenMaps map[string]*GardenMap) []Range {
	fromName := "seed"
	for {
		gardenMap := gardenMaps[fromName]
		if gardenMap == nil {
			panic("bad from name")
		}

		gotRanges := solveLevel(ranges, gardenMap)
		if gardenMap.To == "location" {
			return gotRanges
		}
		checkResult(ranges, gotRanges, gardenMap)

		sortRanges(gotRanges)
		ranges = combineRanges(gotRanges)

		fromName = gardenMap.To
	}
}

type BySourceRange []*RangeMap

func (r BySourceRange) Len() int { return len(r) }
func (r BySourceRange) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r BySourceRange) Less(i, j int) bool {
	return r[i].Src.Lo < r[j].Src.Lo
}

func sortGardenMap(gm *GardenMap) {
	sort.Sort(BySourceRange(gm.Ranges))
}

func sortRanges(r []Range) {
	sort.Slice(r, func(i, j int) bool { return r[i].Lo < r[j].Lo })
}

func solve(fromRanges []Range, gardenMaps []*GardenMap) int64 {
	froms := map[string]*GardenMap{}
	for _, gm := range gardenMaps {
		sortGardenMap(gm)
		froms[gm.From] = gm
	}

	sortRanges(fromRanges)
	locationRanges := solveForLocation(fromRanges, froms)
	sortRanges(locationRanges)

	return locationRanges[0].Lo
}

func solveA(seedNums []int64, gardenMaps []*GardenMap) int64 {
	seedRanges := []Range{}
	for i := 0; i < len(seedNums); i++ {
		seedRanges = append(seedRanges, Range{Lo: seedNums[i], Hi: seedNums[i]})
	}

	return solve(seedRanges, gardenMaps)
}

func solveB(seedNums []int64, gardenMaps []*GardenMap) int64 {
	seedRanges := []Range{}
	for i := 0; i < len(seedNums); i += 2 {
		seedRanges = append(seedRanges, Range{
			Lo: seedNums[i],
			Hi: seedNums[i] + seedNums[i+1] - 1,
		})
	}

	return solve(seedRanges, gardenMaps)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		logger.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		logger.Fatalf("failed to read input: %v", err)
	}

	seedNums, gardenMaps, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(seedNums, gardenMaps))
	fmt.Println("B", solveB(seedNums, gardenMaps))
}
