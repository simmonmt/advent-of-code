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
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type SpinArr struct {
	elems *list.List
	pos   *list.Element
	step  int
}

func NewSpinArr(step int) *SpinArr {
	elems := list.New()
	elems.PushBack(0)
	pos := elems.Front()

	return &SpinArr{
		elems: elems,
		pos:   pos,
		step:  step,
	}
}

func (s *SpinArr) Insert(val int) {
	for i := 0; i < s.step; i++ {
		s.pos = s.pos.Next()
		if s.pos == nil {
			s.pos = s.elems.Front()
		}
	}

	s.pos = s.elems.InsertAfter(val, s.pos)
}

func (s *SpinArr) Dump(out io.Writer) {
	for elem := s.elems.Front(); elem != nil; elem = elem.Next() {
		if elem == s.elems.Front() {
			fmt.Fprintf(out, " ")
		}
		if elem == s.pos {
			fmt.Fprintf(out, "(")
		} else {
			fmt.Fprintf(out, " ")
		}
		fmt.Fprintf(out, "%v", elem.Value)
		if elem == s.pos {
			fmt.Fprintf(out, ")")
		} else {
			fmt.Fprintf(out, " ")
		}
	}
	fmt.Fprintf(out, "\n")
}

func (s *SpinArr) Next() int {
	next := s.pos.Next()
	if next == nil {
		next = s.elems.Front()
	}
	return next.Value.(int)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %v setp", os.Args[0])
	}
	step, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("failed to parse setp %v: %v\n", os.Args[1], err)
	}

	arr := NewSpinArr(step)

	//arr.Dump(os.Stdout)
	for i := 1; i < 2018; i++ {
		arr.Insert(i)
		//arr.Dump(os.Stdout)
	}

	//arr.Dump(os.Stdout)
	fmt.Printf("next value: %v\n", arr.Next())
}
