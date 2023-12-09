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
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Trial struct {
	Time int
	Dist int
}

func parseInput(lines []string) ([]string, error) {
	return lines, nil
}

func runTrial(time, press int) int {
	return press * (time - press)
}

func findCompletingTime(trial Trial) int {
	t := trial.Time / 2
	if dist := runTrial(trial.Time, t); dist > trial.Dist {
		return t
	}

	inc := max(t/10, 1)
	for t = 1; t <= trial.Time; t += inc {
		if dist := runTrial(trial.Time, t); dist > trial.Dist {
			return t
		}
	}

	panic("no completion found")
}

func search(start, end int, cb func(v int) int) {
	for start <= end {
		mid := start + (end-start)/2
		logger.Infof("start %v mid %v end %v", start, mid, end)

		if res := cb(mid); res < 0 {
			end = mid - 1
		} else if res > 0 {
			start = mid + 1
		} else {
			return
		}
	}
}

func solveTrial(trial Trial) int {
	// The v time range looks like this
	//
	//  --------XXXXXXXX----------
	//  1                        t
	//
	// We want to find the bounds of the X region.

	completingTime := findCompletingTime(trial)
	logger.Infof("completing time %v", completingTime)

	num := 0

	minTime := -1
	search(1, completingTime, func(v int) int {
		num++
		if num > 100 {
			panic("too much")
		}
		dist := runTrial(trial.Time, v)
		logger.Infof("want min, test %v = %v", v, dist)

		if dist <= trial.Dist {
			return 1 // go higher
		}

		if minTime == -1 || v < minTime {
			minTime = v
		}
		return -1
	})
	logger.Infof("found min %v", minTime)

	num = 0
	maxTime := -1
	search(completingTime, trial.Time, func(v int) int {
		num++
		if num > 100 {
			panic("too much")
		}

		dist := runTrial(trial.Time, v)
		logger.Infof("want max, test %v = %v", v, dist)

		if dist <= trial.Dist {
			return -1
		}

		if v > maxTime {
			maxTime = v
		}

		return 1
	})
	logger.Infof("found max %v", maxTime)

	return maxTime - minTime + 1
}

func parseATrials(lines []string) []Trial {
	times := strings.Fields(lines[0])
	dists := strings.Fields(lines[1])

	trials := make([]Trial, len(times)-1)
	for i := 1; i < len(times); i++ {
		num, err := strconv.Atoi(times[i])
		if err != nil {
			panic("bad time")
		}
		trials[i-1].Time = num

		num, err = strconv.Atoi(dists[i])
		if err != nil {
			panic("bad dist")
		}
		trials[i-1].Dist = num
	}

	return trials
}

func solveA(lines []string) int {
	trials := parseATrials(lines)

	out := 1
	for _, trial := range trials {
		out *= solveTrial(trial)
	}
	return out
}

func solveB(lines []string) int {
	findNum := func(line string) int {
		_, b, _ := strings.Cut(line, ":")
		num := 0
		for _, c := range b {
			if c >= '0' && c <= '9' {
				num = num*10 + int(c-'0')
			}
		}
		return num
	}

	trial := Trial{Time: findNum(lines[0]), Dist: findNum(lines[1])}
	return solveTrial(trial)
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

	input, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
