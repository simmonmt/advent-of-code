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

package elf

import (
	"fmt"

	"logger"
)

type Elf struct {
	Name     uint
	Presents uint
}

func Print(elves []int) {
	for elfNum, numPresents := range elves {
		if numPresents > 0 {
			fmt.Printf("  %d %d\n", elfNum, numPresents)
		}
	}
}

func findToSkip(numWithPresents int, cur int) int {
	toSkip := 0
	if numWithPresents < 3 {
		toSkip = 1
	} else {
		toSkip = numWithPresents / 2
	}

	return toSkip
}

func skip(elves []int, cur int, toSkip int) int {
	for toSkip > 0 {
		cur++
		if cur == len(elves) {
			cur = 0
		}
		if elves[cur] > 0 {
			toSkip--
		}
	}
	return cur
}

func nextWithPresents(elves []int, cur int) int {
	for ; cur < len(elves) && elves[cur] == 0; cur++ {
	}
	return cur
}

func Play(num int) int {
	elves := make([]int, num)
	for i := range elves {
		elves[i] = 1
	}
	numWithPresents := num

	eIdx := 0
	toSkip := findToSkip(numWithPresents, 0)
	nIdx := skip(elves, eIdx, toSkip)

	for round := 1; numWithPresents > 1; round++ {
		if logger.Enabled() || round%1000 == 0 {
			fmt.Printf("round %v, elves: %v\n", round, numWithPresents)
		}

		if round != 1 {
			var nextSkip int
			if numWithPresents%2 == 0 {
				nextSkip = 2
			} else {
				nextSkip = 1
			}
			nIdx = skip(elves, nIdx, nextSkip)
		}

		if elves[nIdx] == 0 {
			panic(fmt.Sprintf("nIdx = %v, zero", nIdx))
		}

		logger.LogF("%v stealing from %v: %v\n", eIdx, nIdx, elves[nIdx])
		elves[eIdx] += elves[nIdx]
		elves[nIdx] = 0
		numWithPresents--
		eIdx = skip(elves, eIdx, 1)

		if logger.Enabled() {
			logger.LogLn()
			Print(elves)
		}
	}

	return nextWithPresents(elves, 0) + 1
}
