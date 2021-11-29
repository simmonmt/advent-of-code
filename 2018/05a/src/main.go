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

func shrink(arr []byte) {
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
		printArr(arr, i-5, thatIdx+5, i, thatIdx)

		if this != that && strings.ToUpper(string(this)) == strings.ToUpper(string(that)) {
			// they cancel
			fmt.Printf("%c %c cancel\n", this, that)
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
	shrink(arr)

	num := 0
	for _, c := range arr {
		if c != 0 {
			num++
		}
	}

	fmt.Println(num)
}
