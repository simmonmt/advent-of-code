package main

import (
	"container/list"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/intmath"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	prefixSize = flag.Int("prefix_size", 25, "number of leading numbers")
)

type Roller struct {
	elems *list.List
	sums  map[int]int
}

func NewRoller() *Roller {
	return &Roller{
		elems: list.New(),
		sums:  map[int]int{},
	}
}

func (r *Roller) PushBack(val int) {
	r.elems.PushBack(val)
	for elem := r.elems.Front(); elem != nil; elem = elem.Next() {
		sum := elem.Value.(int) + val
		r.sums[sum]++
	}
}

func (r *Roller) RemoveFirst() {
	if r.elems.Len() == 0 {
		panic("empty list")
	}

	removeVal := r.elems.Front().Value.(int)
	r.elems.Remove(r.elems.Front())

	for elem := r.elems.Front(); elem != nil; elem = elem.Next() {
		sum := elem.Value.(int) + removeVal
		if r.sums[sum] == 1 {
			delete(r.sums, sum)
		} else {
			r.sums[sum]--
		}
	}
}

func (r *Roller) HasSum(sum int) bool {
	_, found := r.sums[sum]
	return found
}

func (r *Roller) String() string {
	out := []string{}
	for elem := r.elems.Front(); elem != nil; elem = elem.Next() {
		out = append(out, strconv.Itoa(elem.Value.(int)))
	}
	return strings.Join(out, ",")
}

func findRange(nums []int, invalidNum int) {
	for i := range nums {
		sum := nums[i]
		smallest, largest := nums[i], nums[i]

		for j := i + 1; j < len(nums); j++ {
			sum += nums[j]
			smallest = intmath.IntMin(smallest, nums[j])
			largest = intmath.IntMax(largest, nums[j])

			if sum > invalidNum {
				logger.LogF("no range from %d", i)
				break
			}
			if sum == invalidNum {
				fmt.Printf("B: range [%d,%d] result %d\n",
					i, j, smallest+largest)
				return
			}
		}
	}

	panic("no range found")
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

	nums := []int{}
	for _, line := range lines {
		nums = append(nums, intmath.AtoiOrDie(line))
	}

	roller := NewRoller()

	prefixNums := nums[0:*prefixSize]
	restNums := nums[*prefixSize:]

	for _, num := range prefixNums {
		roller.PushBack(num)
	}

	invalidNum := 0
	for _, num := range restNums {
		if !roller.HasSum(num) {
			invalidNum = num
			fmt.Printf("A: invalid %d\n", num)
			break
		}

		roller.PushBack(num)
		roller.RemoveFirst()
	}

	if invalidNum == 0 {
		panic("no invalid found")
	}

	findRange(nums, invalidNum)
}
