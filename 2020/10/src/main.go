package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/intmath"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func seekBack(nums []int, start, goal int) int {
	for i := start; i >= 0; i-- {
		num := nums[i]
		switch {
		case num < goal:
			return i + 1
		case num == goal:
			return i
		}
	}
	return 0
}

func solveA(nums []int) {
	used := map[int]int{}
	curIdx := 0
	last := 0

	path := doSolveA(nums, []int{}, used, curIdx, last)
	if path == nil {
		panic("no path found")
	}

	path = append([]int{0}, path...)
	path = append(path, path[len(path)-1]+3)

	numOne, numThree := 0, 0
	for i := 1; i < len(path); i++ {
		switch intmath.Abs(path[i] - path[i-1]) {
		case 1:
			numOne++
		case 3:
			numThree++
		}
	}

	fmt.Printf("A: path is %v\n", path)
	fmt.Printf("A: #1 %d #3 %d => %d\n", numOne, numThree, numOne*numThree)
}

func doSolveA(nums, path []int, used map[int]int, curIdx, last int) []int {
	logger.LogF("doSolveA path %v used %v curIdx %v last %v",
		path, used, curIdx, last)

	if len(used) == len(nums) {
		return path
	}

	curIdx = seekBack(nums, curIdx, last)

	cands := []int{}
	for _, num := range nums[curIdx:] {
		if num < last {
			continue
		}
		if num > last+3 {
			break
		}
		if _, found := used[num]; found {
			continue
		}
		cands = append(cands, num)
	}

	if len(cands) == 0 {
		return nil
	}

	newPath := make([]int, len(path)+1)
	copy(newPath, path)
	for _, cand := range cands {
		used[cand] = len(cands)
		newPath[len(path)] = cand
		retPath := doSolveA(nums, newPath, used, curIdx, cand)
		delete(used, cand)
		if len(retPath) > 0 {
			return retPath
		}
	}

	return nil

}

func solveB(nums []int) {
	knownPaths := map[int]int{
		nums[len(nums)-1]: 1,
	}

	logger.LogF("nums %v %v %v", 0, nums, nums[len(nums)-1]+3)

	for i := len(nums) - 2; i >= 0; i-- {
		findNumPaths(nums[i], nums[i+1:], knownPaths)
	}
	findNumPaths(0, nums, knownPaths)

	fmt.Printf("B: %v\n", knownPaths)
}

func findNumPaths(start int, rest []int, knownPaths map[int]int) {
	logger.LogF("start %d kp %v", start, knownPaths)
	numPaths := 0
	for _, num := range rest {
		if num-3 > start {
			continue
		}

		subNum, found := knownPaths[num]
		if !found {
			logger.LogF("discarding %d to %d; no future", start, num)
			continue
		}

		numPaths += subNum
	}

	if numPaths > 0 {
		logger.LogF("adding %d numPaths %d", start, numPaths)
		knownPaths[start] = numPaths
	}
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	nums, err := filereader.Numbers(*input)
	if err != nil {
		log.Fatal(err)
	}

	sort.Ints(nums)

	//solveA(nums)
	solveB(nums)
}
