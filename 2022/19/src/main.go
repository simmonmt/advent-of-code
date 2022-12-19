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

	"github.com/simmonmt/aoc/2022/common/collections"
	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/logger"
	"github.com/simmonmt/aoc/2022/common/mtsmath"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	inputPattern = regexp.MustCompile(
		`Blueprint (\d+): Each ore robot costs (\d+) ore. Each clay robot costs (\d+) ore. Each obsidian robot costs (\d+) ore and (\d+) clay. Each geode robot costs (\d+) ore and (\d+) obsidian.`)
)

type Blueprint struct {
	ID           int
	OreCostOre   int
	ClayCostOre  int
	ObsCostOre   int
	ObsCostClay  int
	GeodeCostOre int
	GeodeCostObs int
}

func parseInput(lines []string) ([]Blueprint, error) {
	out := []Blueprint{}

	for i, line := range lines {
		parts := inputPattern.FindStringSubmatch(line)
		if parts == nil {
			return nil, fmt.Errorf("bad blueprint: %v", line)
		}

		id, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("bad blueprint on line %d: %v",
				i+1, err)
		}
		oreCostOre, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("bad blueprint %d: ore robot ore cost: %v",
				id, err)
		}
		clayCostOre, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, fmt.Errorf("bad blueprint %d: clay robot ore cost: %v",
				id, err)
		}
		obsCostOre, err := strconv.Atoi(parts[4])
		if err != nil {
			return nil, fmt.Errorf("bad blueprint %d: obs robot ore cost: %v",
				id, err)
		}
		obsCostClay, err := strconv.Atoi(parts[5])
		if err != nil {
			return nil, fmt.Errorf("bad blueprint %d: obs robot clay cost: %v",
				id, err)
		}
		geodeCostOre, err := strconv.Atoi(parts[6])
		if err != nil {
			return nil, fmt.Errorf("bad blueprint %d: geode robot ore cost: %v",
				id, err)
		}
		geodeCostObs, err := strconv.Atoi(parts[7])
		if err != nil {
			return nil, fmt.Errorf("bad blueprint %d: geode robot obs ore cost: %v",
				id, err)
		}

		out = append(out, Blueprint{
			ID:           id,
			OreCostOre:   oreCostOre,
			ClayCostOre:  clayCostOre,
			ObsCostOre:   obsCostOre,
			ObsCostClay:  obsCostClay,
			GeodeCostOre: geodeCostOre,
			GeodeCostObs: geodeCostObs,
		})
	}

	return out, nil
}

type State [12]byte

func NewState(minute, numOre, numOreRobots, numClay, numClayRobots, numObs, numObsRobots, numGeode, numGeodeRobots int) State {
	if numOre > 0xffff || numClay > 0xffff || numObs > 0xffff || numGeode > 0xff {
		panic("num too big")
	}

	o := [12]byte{}
	o[0] = byte((numOre >> 8) & 0xff)
	o[1] = byte(numOre & 0xff)
	o[2] = byte((numClay >> 8) & 0xff)
	o[3] = byte(numClay & 0xff)
	o[4] = byte((numObs >> 8) & 0xff)
	o[5] = byte(numObs & 0xff)
	o[6] = byte(numGeode & 0xff)
	o[7] = byte(numOreRobots)
	o[8] = byte(numClayRobots)
	o[9] = byte(numObsRobots)
	o[10] = byte(numGeodeRobots)
	o[11] = byte(minute)

	return o
}

func DecodeState(s string) State {
	if len(s) != 12 {
		panic("bad string")
	}
	o := [12]byte{}
	copy(o[:], []byte(s))
	return o
}

func (s *State) String() string {
	return fmt.Sprintf("[min:%d ore:%d[%d], clay:%d[%d], obs:%d[%d], geo:%d[%d]]",
		s.Min(),
		s.NumOre(), s.NumOreRobots(),
		s.NumClay(), s.NumClayRobots(),
		s.NumObs(), s.NumObsRobots(),
		s.NumGeode(), s.NumGeodeRobots())
}

func (s *State) Encode() string {
	return string(s[:])
}

func (s *State) Min() int     { return int(s[11]) }
func (s *State) SetMin(m int) { s[11] = byte(m) }

func (s *State) NumOre() int         { return int(s[0]<<8) | int(s[1]) }
func (s *State) NumClay() int        { return int(s[2]<<8) | int(s[3]) }
func (s *State) NumObs() int         { return int(s[4]<<8) | int(s[5]) }
func (s *State) NumGeode() int       { return int(s[6]) }
func (s *State) NumOreRobots() int   { return int(s[7]) }
func (s *State) NumClayRobots() int  { return int(s[8]) }
func (s *State) NumObsRobots() int   { return int(s[9]) }
func (s *State) NumGeodeRobots() int { return int(s[10]) }

func (s *State) SetOre(n int)   { s[0] = byte((n >> 8) & 0xff); s[1] = byte(n & 0xff) }
func (s *State) SetClay(n int)  { s[2] = byte((n >> 8) & 0xff); s[3] = byte(n & 0xff) }
func (s *State) SetObs(n int)   { s[4] = byte((n >> 8) & 0xff); s[5] = byte(n & 0xff) }
func (s *State) SetGeode(n int) { s[6] = byte(n & 0xff) }

func (s *State) SetOreRobots(n int)   { s[7] = byte(n) }
func (s *State) SetClayRobots(n int)  { s[8] = byte(n) }
func (s *State) SetObsRobots(n int)   { s[9] = byte(n) }
func (s *State) SetGeodeRobots(n int) { s[10] = byte(n) }

func viable(s *State, max int, timeLimit int) bool {
	numGeodes := s.NumGeode()
	numrobots := s.NumGeodeRobots()
	sum := numGeodes

	// I know there's a formula. Just too tired to implement+test it.
	for i := s.Min() + 1; i <= timeLimit; i++ {
		sum += numrobots
		numrobots++
	}
	return sum > max
}

func (s *State) Neighbors(bp *Blueprint) []State {
	out := []State{}
	if numOre, numObs := s.NumOre(), s.NumObs(); numOre >= bp.GeodeCostOre && numObs >= bp.GeodeCostObs {
		new := *s
		new.SetOre(numOre - bp.GeodeCostOre)
		new.SetObs(numObs - bp.GeodeCostObs)
		new.SetGeodeRobots(new.NumGeodeRobots() + 1)
		out = append(out, new)
	}
	if numOre, numClay := s.NumOre(), s.NumClay(); numOre >= bp.ObsCostOre && numClay >= bp.ObsCostClay {
		new := *s
		new.SetOre(numOre - bp.ObsCostOre)
		new.SetClay(numClay - bp.ObsCostClay)
		new.SetObsRobots(new.NumObsRobots() + 1)
		out = append(out, new)
	}
	if numOre := s.NumOre(); numOre >= bp.ClayCostOre {
		new := *s
		new.SetOre(numOre - bp.ClayCostOre)
		new.SetClayRobots(new.NumClayRobots() + 1)
		out = append(out, new)
	}
	if numOre := s.NumOre(); numOre >= bp.OreCostOre {
		new := *s
		new.SetOre(numOre - bp.OreCostOre)
		new.SetOreRobots(new.NumOreRobots() + 1)
		out = append(out, new)
	}

	out = append(out, *s) // just produce

	for i := 0; i < len(out); i++ {
		new := &out[i]
		new.SetMin(s.Min() + 1)
		new.SetOre(new.NumOre() + s.NumOreRobots())
		new.SetClay(new.NumClay() + s.NumClayRobots())
		new.SetObs(new.NumObs() + s.NumObsRobots())
		new.SetGeode(new.NumGeode() + s.NumGeodeRobots())
	}

	return out
}

func solveBlueprint(blueprint Blueprint, timeLimit int) int {
	queue := collections.NewPriorityQueue[string](collections.GreaterThan)
	start := NewState(0, 0, 1, 0, 0, 0, 0, 0, 0)
	queue.Insert(start.Encode(), 0)

	visited := map[string]bool{}

	maxGeodes := 0
	for !queue.IsEmpty() {
		encoded, _ := queue.Next()
		state := DecodeState(encoded)

		if logger.Enabled() {
			logger.LogF("cur %v", state.String())
		}

		neighbors := []State{}
		if state.Min() != timeLimit {
			neighbors = state.Neighbors(&blueprint)
		}

		for i := 0; i < len(neighbors); i++ {
			neighbor := &neighbors[i]
			if logger.Enabled() {
				logger.LogF("  neighbor %v", neighbor.String())
			}

			if !viable(neighbor, maxGeodes, timeLimit) {
				continue
			}

			encodedNeighbor := neighbor.Encode()
			if _, found := visited[encodedNeighbor]; found {
				continue
			}

			neighborDist := neighbor.NumGeode()*1000 + neighbor.NumObsRobots()
			queue.Insert(encodedNeighbor, neighborDist)
		}

		if state.NumGeode() > maxGeodes {
			//fmt.Println("new max", maxGeodes)
			maxGeodes = state.NumGeode()
		}

		visited[encoded] = true
	}

	return maxGeodes
}

func solveA(blueprints []Blueprint) int {
	sum := 0
	for _, bp := range blueprints {
		fmt.Println(bp.ID)
		sum += solveBlueprint(bp, 24) * bp.ID
	}
	return sum
}

func solveB(blueprints []Blueprint) int {
	prod := 1
	for _, bp := range blueprints[0:mtsmath.Min(len(blueprints), 3)] {
		fmt.Println(bp.ID)
		max := solveBlueprint(bp, 32)
		fmt.Println("=>", max)
		prod *= max
	}
	return prod
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

	input, err := parseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
