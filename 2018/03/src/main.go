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
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var (
	pattern = regexp.MustCompile(`^#([0-9]+) @ ([0-9]+),([0-9]+): ([0-9]+)x([0-9]+)$`)
)

type Fabric struct {
	W, H  int
	Elems [][]int
}

func NewFabric(w, h int) *Fabric {
	f := &Fabric{
		W:     w,
		H:     h,
		Elems: make([][]int, w*h),
	}

	for i := 0; i < w*h; i++ {
		f.Elems[i] = []int{}
	}

	return f
}

func (f *Fabric) AddClaim(loff, toff, w, h int, id int) {
	for y := toff; y < toff+h; y++ {
		for x := loff; x < loff+w; x++ {
			off := y*f.W + x
			f.Elems[off] = append(f.Elems[off], id)
		}
	}
}

func (f *Fabric) Dump() {
	for y := 0; y < f.H; y++ {
		for x := 0; x < f.W; x++ {
			off := y*f.W + x
			claims := f.Elems[off]

			switch len(claims) {
			case 0:
				fmt.Print(".")
				break
			case 1:
				fmt.Print(claims[0])
				break
			default:
				fmt.Print("X")
			}
		}
		fmt.Println()
	}
}

func AtoiOrDie(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse %v: %v", s, err))
	}
	return val
}

func main() {
	f := NewFabric(1000, 1000)
	fmt.Println("fabric built")

	overlapsByID := map[int]bool{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		parts := pattern.FindStringSubmatch(line)
		if parts == nil {
			log.Fatalf("failed to parse %v", line)
		}

		id := AtoiOrDie(parts[1])
		loff := AtoiOrDie(parts[2])
		toff := AtoiOrDie(parts[3])
		w := AtoiOrDie(parts[4])
		h := AtoiOrDie(parts[5])

		f.AddClaim(loff, toff, w, h, id)
		overlapsByID[id] = false
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("read failed: %v", err)
	}

	numDouble := 0
	for _, claims := range f.Elems {
		if len(claims) > 1 {
			for _, claim := range claims {
				overlapsByID[claim] = true
			}

			numDouble++
		}
	}
	fmt.Printf("num overlap: %v\n", numDouble)

	for id, overlapped := range overlapsByID {
		if !overlapped {
			fmt.Printf("no overlap: %v\n", id)
		}
	}
}
