package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose  = flag.Bool("verbose", false, "verbose")
	input    = flag.String("input", "", "input file")
	numTurns = flag.Int("num_turns", 2020, "number of turns")
)

func solve(startingNums []int) int {
	lastByNum := map[int][]int{}
	lastSpoken := 0

	for i := 0; ; i++ {
		turn := i + 1
		num := 0

		if i < len(startingNums) {
			num = startingNums[i]
		} else {
			times := lastByNum[lastSpoken]
			if len(times) == 1 {
				// that was the first time
				num = 0
			} else {
				num = times[0] - times[1]
			}
		}

		logger.LogF("Turn %d: %v", turn, num)

		if turn == *numTurns {
			return num
		}

		lastSpoken = num
		if _, found := lastByNum[num]; !found {
			lastByNum[num] = []int{i}
		} else {
			lastByNum[num] = []int{i, lastByNum[num][0]}
		}

		//logger.LogF("%d: %v last now %v", turn, num, lastByNum[num])
	}
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	nums, err := filereader.OneRowOfNumbers(*input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Answer: %v\n", solve(nums))
}
