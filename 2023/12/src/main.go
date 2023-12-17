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
	"strings"

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type CellState int

const (
	CS_UNKNOWN CellState = iota
	CS_YES
	CS_NO
)

func (cs CellState) String() string {
	switch cs {
	case CS_UNKNOWN:
		return "?"
	case CS_YES:
		return "#"
	case CS_NO:
		return "."
	default:
		return "X"
	}
}

func CellStateFromRune(r rune) CellState {
	switch r {
	case '.':
		return CS_NO
	case '#':
		return CS_YES
	case '?':
		return CS_UNKNOWN
	default:
		panic("bad rune")
	}
}

func CellStatesFromString(str string) []CellState {
	states := make([]CellState, len(str))
	for i, r := range str {
		states[i] = CellStateFromRune(r)
	}
	return states
}

type Spring struct {
	States []CellState
	Sizes  []int
}

func parseInput(lines []string) ([]*Spring, error) {
	out := []*Spring{}
	for i, line := range lines {
		str, rest, ok := strings.Cut(line, " ")
		if !ok {
			return nil, fmt.Errorf("%d: no space", i+1)
		}
		states := CellStatesFromString(str)

		nums, err := filereader.ParseNumbersFromLine(rest)
		if err != nil {
			return nil, fmt.Errorf("%d: bad numbers: %v", i+1, err)
		}

		out = append(out, &Spring{States: states, Sizes: nums})
	}
	return out, nil
}

type Range struct {
	Left, Right int
}

func (r Range) Equals(o Range) bool {
	return r.Left == o.Left && r.Right == o.Right
}

func (r Range) Add(v int) Range {
	return Range{Left: r.Left + v, Right: r.Right + v}
}

func canWork(states []CellState, left, size int, within Range) bool {
	right := left + (size - 1)
	if right > within.Right {
		return false
	}

	for i := left; i <= right; i++ {
		if states[i] == CS_NO {
			return false
		}
	}

	if right < within.Right && states[right+1] == CS_YES {
		return false
	}

	return true
}

func findConstraint(states []CellState, size int, first bool, within Range) Range {
	//logger.Infof("findConstraint states %s, size %d, within %v", states, size, within)

	start := within.Left
	constraint := Range{Left: -1, Right: -1}
	for ; start+(size-1) < len(states); start++ {
		if canWork(states, start, size, within) {
			constraint = Range{Left: start, Right: size - 1}
			break
		}
	}

	if constraint.Left == -1 {
		panic("can't find a starting point")
	}

	maxLeft := start
	left := start
	// We know the leftmost place where this size can work. Try pushing it right.
	for {
		// Will this constraint let us push right?
		if first && states[left] == CS_YES {
			// ?#???, sz=3, within=0,4, constraint=0,3
			// We can't push right because the # needs to be included.
			break
		}

		left++
		right := left + (size - 1)
		if right > within.Right {
			break
		}

		// Will size work if we push right?
		if canWork(states, left, size, within) {
			maxLeft = left
		}
	}

	constraint.Right = maxLeft + (size - 1)
	return constraint
}

func updateStates(states []CellState, size int, r, within Range) bool {
	//logger.Infof("updateStates %v sz %v r %v within %v", states, size, r, within)
	changed := false

	// We know the exact range, so everything inside must be set and we know
	// where the surrounding NOs are.
	if r.Left+(size-1) == r.Right {
		for i := 0; i < size; i++ {
			if idx := r.Left + i; states[idx] != CS_YES {
				changed = true
				states[idx] = CS_YES
			}
		}

		if r.Left > within.Left {
			want := r.Left - 1
			if got := states[want]; got != CS_UNKNOWN && got != CS_NO {
				panic("bad no left insert")
			}
			if states[want] != CS_NO {
				changed = true
				states[want] = CS_NO
			}
		}

		if r.Right < within.Right {
			want := r.Right + 1
			if got := states[want]; got != CS_UNKNOWN && got != CS_NO {
				panic("bad no right insert")
			}
			if states[want] != CS_NO {
				changed = true
				states[want] = CS_NO
			}
		}
		return changed
	}

	// sweep := func(cur, size, inc int) {
	// 	logger.Infof("sweeping %d sz %d inc %d", cur, size, inc)
	// 	found := false
	// 	for i := 0; i < size; i, cur = i+1, cur+inc {
	// 		switch states[cur] {
	// 		case CS_YES:
	// 			found = true
	// 		case CS_UNKNOWN:
	// 			if found {
	// 				states[cur] = CS_YES
	// 				changed = true
	// 			}
	// 		default:
	// 			panic("unexpected cell")
	// 		}
	// 	}
	// }

	// sweep(r.Left, size, 1)

	return changed
}

func isExact(r Range, size int) bool {
	return r.Left+(size-1) == r.Right
}

func constrainSpring(spring *Spring) ([]CellState, []Range) {
	//logger.Infof("constrainSpring %v", spring)

	states := make([]CellState, len(spring.States))
	copy(states, spring.States)

	constraints := make([]Range, len(spring.Sizes))
	for i := range spring.Sizes {
		constraints[i] = Range{Left: -1, Right: -1}
	}

	changed := true
	within := Range{0, len(states) - 1}
	for changed {
		changed = false

		start := 0

		for i, size := range spring.Sizes {
			if i > 0 {
				lastSize := spring.Sizes[i-1]
				lastConstraint := constraints[i-1]
				if isExact(lastConstraint, lastSize) {
					start = lastConstraint.Right + 2
					//logger.Infof("fixing start to %d because exact last %v sz %v",
					//	start, lastConstraint, lastSize)
				} else {
					// There has to be one blank position
					// between the end of the
					// leftmost-possible position of the
					// previous constraint and the beginning
					// of this one.
					start = lastConstraint.Left + (lastSize - 1) + 2
					//logger.Infof("fixing start to %d because of last %v sz %v",
					//	start, lastConstraint, lastSize)
				}
			}

			end := within.Right
			if i < len(spring.Sizes)-1 {
				if next := constraints[i+1]; next.Left != -1 {
					// We've calculated the next
					// constraint. This one must end before
					// the beginning of the rightmost legal
					// position within the next constraint,
					// so adjust it if findConstraint found
					// something longer.
					nextSize := spring.Sizes[i+1]
					nextRightStart := next.Right - (nextSize - 1)

					end = min(end, nextRightStart-2)
					//logger.Infof("fixing end to %v because of right %v", end, next)
				}
			}

			constraint := findConstraint(states, size, i == 0, Range{Left: start, Right: end})
			//logger.Infof("findConstraint returned %v for #%d sz %d", constraint, i, size)

			if !constraint.Equals(constraints[i]) {
				changed = true
				constraints[i] = constraint
			}

			if updateStates(states, size, constraints[i], within) {
				//logger.Infof("updated states: %v", states)
				changed = true
			}
		}
	}

	return states, constraints
}

type RangeIterator struct {
	ranges  []Range
	sizes   []int
	nums    []int
	offsets []int
}

func newRangeIterator(ranges []Range, sizes []int) *RangeIterator {
	//logger.Infof("new ri %v %v", ranges, sizes)

	rc := make([]Range, len(ranges))
	copy(rc, ranges)

	sc := make([]int, len(sizes))
	copy(sc, sizes)

	nums := make([]int, len(ranges))
	for i, size := range sizes {
		nums[i] = ranges[i].Right - (size - 1) - ranges[i].Left + 1
	}

	return &RangeIterator{
		ranges:  rc,
		sizes:   sc,
		nums:    nums,
		offsets: make([]int, len(ranges)),
	}
}

func (ri *RangeIterator) makeRange(i, offset int) Range {
	return Range{
		Left:  ri.ranges[i].Left + offset,
		Right: ri.ranges[i].Left + (ri.sizes[i] - 1) + offset,
	}
}

func (ri *RangeIterator) Next() ([]Range, bool) {
	//logger.Infof("next")

	out := make([]Range, len(ri.ranges))
	for i, offset := range ri.offsets {
		out[i] = ri.makeRange(i, offset)
	}

	//logger.Infof("ret %v", out)

	cur := len(ri.ranges) - 1
	for cur >= 0 {
		if !ri.increment(cur) {
			//logger.Infof("failed increment %d", cur)
			cur--
			continue
		}
		//logger.Infof("ok increment %d now %v", cur, ri.makeRange(cur, ri.offsets[cur]))

		if !ri.resetFrom(cur) {
			//logger.Infof("failed resetfrom %v", cur)
			continue
		}

		//logger.Infof("breaking")
		break
	}

	return out, cur < 0
}

func (ri *RangeIterator) increment(i int) bool {
	ri.offsets[i]++
	if ri.offsets[i] >= ri.nums[i] {
		ri.offsets[i] = 0
		return false
	}
	return true
}

func (ri *RangeIterator) resetFrom(start int) bool {
	after := func(i, j int) bool {
		a := ri.makeRange(i, ri.offsets[i])
		b := ri.makeRange(j, ri.offsets[j])

		return a.Right < b.Left-1
	}

	if start == len(ri.ranges)-1 {
		//logger.Infof("resetfrom returning end")
		return true
	}

	// 1 find a next that's greater than this
	// 2 tell next to reset its children
	// 3 if it can
	// 4   return true
	// 5 else
	// 6   try to increment next
	// 5     if we can't
	// 6       return failure
	// 7   go to 2

	next := start + 1
	ri.offsets[next] = 0

	for !after(start, next) {
		//logger.Infof("%v is not after %v", ri.makeRange(next, ri.offsets[next]), ri.makeRange(start, ri.offsets[start]))
		if !ri.increment(next) {
			return false
		}
	}

	//logger.Infof("found new next %v", ri.makeRange(next, ri.offsets[next]))

	// next is now after start

	for {
		if ri.resetFrom(next) {
			return true
		}

		if !ri.increment(next) {
			return false
		}
	}
}

func isValid(states []CellState, ranges []Range) bool {
	ref := make([]CellState, len(states))
	copy(ref, states)

	for _, r := range ranges {
		for i := r.Left; i <= r.Right; i++ {
			if states[i] == CS_NO {
				//logger.Infof("found no")
				return false
			}
			ref[i] = CS_UNKNOWN
		}
		if r.Left > 0 && states[r.Left-1] == CS_YES {
			//logger.Infof("left yes")
			return false
		}
		if r.Right < len(states)-1 && states[r.Right+1] == CS_YES {
			//logger.Infof("right yes")
			return false
		}
	}

	for _, s := range ref {
		if s == CS_YES {
			return false
		}
	}
	return true
}

func solveSpring(spring *Spring) int {
	states, constraints := constrainSpring(spring)

	//logger.Infof("constrain results: states %s constraints %v sizes %v",
	//	states, constraints, spring.Sizes)

	mine := solveSpringWithConstraints(states, constraints, spring.Sizes)

	// states = spring.States
	// for i := range constraints {
	// 	constraints[i] = Range{0, len(states) - 1}
	// }

	// slow := solveSpringWithConstraints(spring.States, constraints, spring.Sizes)

	// if mine != slow {
	// 	panic(fmt.Sprintf("mismatch, mine %d slow %d", mine, slow))
	// }

	//logger.Infof("spring %v combinations %v", spring, mine)
	return mine
}

func solveSpringWithConstraints(states []CellState, constraints []Range, sizes []int) int {
	ri := newRangeIterator(constraints, sizes)
	num := 0
	for {
		ranges, done := ri.Next()
		//logger.Infof("checking %v", ranges)

		good := true
		for i := 0; i < len(ranges)-1; i++ {
			if ranges[i].Right >= ranges[i+1].Left-1 {
				//logger.Infof("collision %d %d", i, i+1)
				good = false
				break
			}
		}
		if good && !isValid(states, ranges) {
			//logger.Infof("invalid")
			good = false
		}

		if good {
			num++
			//logger.Infof("ok: %v", ranges)
		}

		if done {
			break
		}
	}

	return num
}

func solveA(springs []*Spring) int {
	out := 0
	for i, spring := range springs {
		fmt.Printf("-- solving %d of %d\n", i+1, len(springs))
		out += solveSpring(spring)
	}
	return out
}

func allPossibleRev(r Range, size int) []Range {
	out := []Range{}

	for i := r.Right - (size - 1); i >= r.Left; i-- {
		out = append(out, Range{Left: i, Right: i + (size - 1)})
	}

	return out
}

func solveBSpring(spring *Spring) int {
	// find out how many last solutions there are
	// in a vector indexed by starting position.
	//
	// v[34] = number of solutions with last at left>=34

	// now we need prev
	//
	// new vector vv
	//
	//  vv[23] = number of solutions with prev at left >= 23
	//
	//    which means number of last solutions >= prev.right+2

	// prev = sum(for each possible prev position
	//                  number of possible

	ns := Spring{States: []CellState{}, Sizes: []int{}}
	for j := 0; j < 5; j++ {
		if j != 0 {
			ns.States = append(ns.States, CS_UNKNOWN)
		}
		ns.States = append(ns.States, spring.States...)
		ns.Sizes = append(ns.Sizes, spring.Sizes...)
	}

	states, constraints := constrainSpring(&ns)

	//logger.Infof("states %v", states)
	//logger.Infof("constraints %v", constraints)

	v := make([]int, len(ns.States))

	last := len(constraints) - 1
	for _, r := range allPossibleRev(constraints[last], ns.Sizes[last]) {
		//logger.Infof("last rev %v", r)
		if isValid(states[r.Left:], []Range{Range{0, r.Right - r.Left}}) {
			v[r.Left] = 1
		}
	}

	for i := last - 1; i >= 0; i-- {
		v2 := make([]int, len(ns.States))
		for _, r := range allPossibleRev(constraints[i], ns.Sizes[i]) {
			//logger.Infof("prev rev %v", r)
			for j := r.Right + 2; j < len(states); j++ {
				if isValid(states[r.Left:j], []Range{Range{0, r.Right - r.Left}}) {
					v2[r.Left] += v[j]
				}
			}
		}

		//logger.Infof("i %v ct %v size %v v2 %v", i, constraints[i], ns.Sizes[i], v2)

		v = v2
	}

	sum := 0
	for _, val := range v {
		sum += val
	}

	return sum
}

func solveB(springs []*Spring) int {
	out := 0
	for i, spring := range springs {
		fmt.Printf("-- solving %d of %d\n", i+1, len(springs))
		out += solveBSpring(spring)
	}
	return out
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
