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

package main

import (
	"container/list"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose      = flag.Bool("verbose", false, "verbose")
	input        = flag.String("input", "", "input file")
	numMovesFlag = flag.Int("num_moves", -1,
		"if set, override the number of moves")
	solveAFlag = flag.Bool("solvea", true, "solve A")
	solveBFlag = flag.Bool("solveb", true, "solve b")
)

type CupSet struct {
	l        *list.List
	cur      *list.Element
	min, max int
	toElems  map[int]*list.Element
}

func newCupSet(nums []int) *CupSet {
	cs := &CupSet{
		l:       list.New(),
		toElems: map[int]*list.Element{},
	}

	min, max := -1, -1
	for _, num := range nums {
		cs.toElems[num] = cs.l.PushBack(num)

		if min == -1 || num < min {
			min = num
		}
		if num > max {
			max = num
		}
	}

	cs.cur = cs.l.Front()
	cs.min = min
	cs.max = max

	return cs
}

func (cs *CupSet) nextCW(elem *list.Element) *list.Element {
	elem = elem.Next()
	if elem != nil {
		return elem
	}
	return cs.l.Front()
}

func (cs *CupSet) MoveCW() {
	cs.cur = cs.cur.Next()
	if cs.cur == nil {
		cs.cur = cs.l.Front()
	}
}

func (cs *CupSet) PopCW() int {
	cwElem := cs.nextCW(cs.cur)
	val := cs.l.Remove(cwElem).(int)
	delete(cs.toElems, val)
	return val
}

func (cs *CupSet) PushCWAfter(mark, num int) {
	if !cs.InSet(mark) {
		panic("unknown mark")
	}
	if cs.InSet(num) {
		panic("re-add")
	}
	cs.toElems[num] = cs.l.InsertAfter(num, cs.toElems[mark])
}

func (cs *CupSet) InSet(num int) bool {
	_, found := cs.toElems[num]
	return found
}

func (cs *CupSet) Cur() int {
	return cs.cur.Value.(int)
}

func (cs *CupSet) SetCur(num int) {
	elem, found := cs.toElems[num]
	if !found {
		panic("setcur not in set")
	}
	cs.cur = elem
}

func (cs *CupSet) String() string {
	out := []string{}
	for e := cs.l.Front(); e != nil; e = e.Next() {
		num := e.Value.(int)
		var str string
		if e == cs.cur {
			str = fmt.Sprintf("(%d)", num)
		} else {
			str = strconv.Itoa(num)
		}
		out = append(out, str)
	}
	return strings.Join(out, " ")
}

func (cs *CupSet) Summarize() string {
	out := ""

	one := cs.toElems[1]
	start := cs.nextCW(one)

	for e := start; e != one; e = cs.nextCW(e) {
		out += strconv.Itoa(e.Value.(int))
	}

	return out
}

func PlayGame(cs *CupSet, numMoves int) {
	for moveNum := 1; moveNum <= numMoves; moveNum++ {
		logger.LogF("\n-- move %d --", moveNum)
		logger.LogF("cups: %s", cs)

		var removed [3]int
		removed[0], removed[1], removed[2] = cs.PopCW(), cs.PopCW(), cs.PopCW()

		logger.LogF("pick up: %v", removed)

		dest := cs.Cur() - 1
		for {
			if dest < cs.min {
				dest = cs.max - (cs.min - dest - 1)
			}
			if cs.InSet(dest) {
				break
			}
			dest--
		}

		logger.LogF("destination: %v", dest)
		cs.PushCWAfter(dest, removed[2])
		cs.PushCWAfter(dest, removed[1])
		cs.PushCWAfter(dest, removed[0])
		cs.MoveCW()
	}
}

func summarize(cs *CupSet) string {
	logger.LogF("\n-- final --")
	logger.LogF("cups: %v", cs)

	return cs.Summarize()
}

func solveA(nums []int) string {
	numMoves := 100
	if *numMovesFlag != -1 {
		numMoves = *numMovesFlag
	}

	cupSet := newCupSet(nums)
	PlayGame(cupSet, numMoves)
	return summarize(cupSet)
}

func solveB(nums []int) string {
	numMoves := 10 * 1000 * 1000
	if *numMovesFlag != -1 {
		numMoves = *numMovesFlag
	}

	start := 0
	for _, num := range nums {
		if num > start {
			start = num
		}
	}
	start++

	allNums := make([]int, 1000000)
	copy(allNums, nums)
	for i := len(nums); i < 1000000; i++ {
		allNums[i] = start
		start++
	}

	cupSet := newCupSet(allNums)
	PlayGame(cupSet, numMoves)

	cupSet.SetCur(1)
	num1 := cupSet.PopCW()
	num2 := cupSet.PopCW()
	logger.LogF("num1 %v num2 %v\n", num1, num2)
	return strconv.Itoa(num1 * num2)
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
	for _, r := range lines[0] {
		nums = append(nums, int(r-'0'))
	}

	if *solveAFlag {
		fmt.Printf("A: %s\n", solveA(nums))
	}
	if *solveBFlag {
		fmt.Printf("B: %s\n", solveB(nums))
	}
}
