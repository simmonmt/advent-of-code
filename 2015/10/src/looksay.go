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
	"fmt"
	"log"
	"os"
	"strconv"
)

func addStreak(arr []rune, streakLen int, r rune) int {
	added := 0
	streakStr := strconv.Itoa(streakLen)
	for _, streakDigit := range streakStr {
		arr[added] = streakDigit
		added++
	}
	arr[added] = r
	added++
	return added
}

func encode(input string) string {
	out := make([]rune, len(input)*2)
	outLen := 0

	streak := 1
	for i, c := range input {
		if i+1 == len(input) {
			//fmt.Printf("last %c %d\n", c, streak)
			outLen += addStreak(out[outLen:], streak, c)
			continue
		}

		next := rune(input[i+1])
		if c == next {
			streak++
		} else {
			//fmt.Printf("in   %c %d\n", c, streak)
			outLen += addStreak(out[outLen:], streak, c)
			streak = 1
		}
	}

	return string(out[0:outLen])
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("usage: %v input len", os.Args[0])
	}
	input := os.Args[1]
	nIters, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("failed to parse len: %v", err)
	}

	for i := 1; i <= nIters; i++ {
		input = encode(input)
		fmt.Printf("iter %d: %d\n", i, len(input))
		//fmt.Println(input)
	}
}
