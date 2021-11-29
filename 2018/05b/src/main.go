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
	"strings"
	"unicode"
)

func printArr(arr []byte, start, end, this, that int) {
	for i := start; i <= end; i++ {
		if i < 0 || i >= len(arr) {
			fmt.Printf("=")
			continue
		}

		c := arr[i]
		if i == this {
			fmt.Printf("*%c*", c)
		} else if i == that {
			fmt.Printf("_%c_", c)
		} else {
			fmt.Printf("%c", c)
		}
	}
	fmt.Println()
}

func next(arr []byte) int {
	for i := 0; i < len(arr); i++ {
		if arr[i] != 0 {
			return i
		}
	}
	return -1
}

func prev(arr []byte, start int) int {
	for i := start - 1; i > 0; i-- {
		if arr[i] != 0 {
			return i
		}
	}
	return 0
}

func react(arr []byte) {
	i := 0
	nextI := 0
	for {
		i = nextI
		nextI = i + 1

		if i == len(arr) {
			break
		}

		//fmt.Println(arr)

		if arr[i] == 0 {
			continue
		}

		this := arr[i]
		thatIdx := next(arr[i+1:]) + i + 1
		if thatIdx < 0 {
			return
		}
		that := arr[thatIdx]
		//printArr(arr, i-5, thatIdx+5, i, thatIdx)

		if this != that && strings.ToUpper(string(this)) == strings.ToUpper(string(that)) {
			// they cancel
			//fmt.Printf("%c %c cancel\n", this, that)
			arr[i] = 0
			arr[thatIdx] = 0

			nextI = prev(arr, i)
		}
	}
}

func main() {
	var line string

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("read failed: %v", err)
	}

	arr := []byte(line)

	units := map[rune]bool{}
	for _, c := range arr {
		u := unicode.ToUpper(rune(c))
		units[u] = true
	}

	fmt.Printf("%d units\n", len(units))

	counts := map[rune]int{}

	i := -1
	for unit, _ := range units {
		i++
		fmt.Printf("trying %c (#%d)\n", unit, i)

		sub := []byte(line)
		for i := range sub {
			if unicode.ToUpper(rune(sub[i])) == unit {
				sub[i] = 0
			}
		}

		react(sub)

		num := 0
		for _, c := range sub {
			if c != 0 {
				num++
			}
		}

		fmt.Printf("unit %c num %d\n", unit, num)
		counts[unit] = num
	}

	minCount := -1
	for _, count := range counts {
		if minCount == -1 || count < minCount {
			minCount = count
		}
	}

	fmt.Println(minCount)
}
