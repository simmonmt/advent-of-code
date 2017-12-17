package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type SpinArr struct {
	// TODO(simmonmt): make a node struct with {val,next}, then
	// make an array of that. better cache locality. also rewrite
	// defrag to be aware of the new structure.
	vals    []int
	nexts   []int
	nextVal int
	pos     int
	step    int
	size    int
}

func NewSpinArr(size int, step int) *SpinArr {
	vals := make([]int, size)
	vals[0] = 0

	nexts := make([]int, size)
	nexts[0] = -1

	return &SpinArr{
		vals:    vals,
		nexts:   nexts,
		nextVal: 1,
		pos:     0,
		step:    step,
		size:    size,
	}
}

func (s *SpinArr) Insert(val int) {
	//fmt.Printf("inserting %v\n", val)

	for i := 0; i < s.step; i++ {
		s.pos = s.nexts[s.pos]
		if s.pos == -1 {
			s.pos = 0
		}
		//fmt.Printf("s.pos now %v\n", s.pos)
	}

	// insert after s.pos
	s.vals[s.nextVal] = val
	s.nexts[s.nextVal] = s.nexts[s.pos]
	s.nexts[s.pos] = s.nextVal
	s.pos = s.nextVal
	s.nextVal++
}

func (s *SpinArr) Dump(out io.Writer) {
	for i := 0; i != -1; i = s.nexts[i] {
		if i == 0 {
			fmt.Fprintf(out, " ")
		}
		if i == s.pos {
			fmt.Fprintf(out, "(")
		} else {
			fmt.Fprintf(out, " ")
		}
		fmt.Fprintf(out, "%v", s.vals[i])
		if i == s.pos {
			fmt.Fprintf(out, ")")
		} else {
			fmt.Fprintf(out, " ")
		}
	}
	fmt.Fprintf(out, "\n")
}

func (s *SpinArr) Next() int {
	nextPos := s.nexts[s.pos]
	if nextPos == -1 {
		nextPos = 0
	}
	return s.vals[nextPos]
}

func (s *SpinArr) After(after int) int {
	for i := 0; i != -1; i = s.nexts[i] {
		if s.vals[i] == after {
			nextPos := s.nexts[i]
			if nextPos == -1 {
				nextPos = 0
			}
			return s.vals[nextPos]
		}
	}
	return -1
}

// func (s *SpinArr) Defrag() {
// 	newVals := make([]int, s.nextVal)

// 	newValsIdx := 0
// 	for i := 0; i != -1; i = s.nexts[i] {
// 		newVals[newValsIdx] = s.vals[i]
// 		newValsIdx++
// 	}

// 	copy(s.vals[0:newValsIdx], newVals)
// 	for i := 0; i < newValsIdx-1; i++ {
// 		s.nexts[i] = i+1
// 	}
// 	s.nexts[newValsIdx-1] = -1
// }

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %v setp", os.Args[0])
	}
	step, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("failed to parse step %v: %v\n", os.Args[1], err)
	}

	iters := 50 * 1000 * 1000
	//iters := 100

	arr := NewSpinArr(iters+1, step)

	//arr.Dump(os.Stdout)
	for i := 1; i <= iters; i++ {
		if i%100000 == 0 {
			fmt.Println(i)
		}
		arr.Insert(i)
		//arr.Dump(os.Stdout)
	}

	// arr.Dump(os.Stdout)
	fmt.Printf("next value: %v\n", arr.Next())
	fmt.Printf("after zero: %v\n", arr.After(0))
}
