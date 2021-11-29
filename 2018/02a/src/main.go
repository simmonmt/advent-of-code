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
)

func count(line string) (int, int) {
	found := map[rune]int{}
	for _, c := range line {
		found[c]++
	}

	numTwo := 0
	numThree := 0
	for _, num := range found {
		if num == 2 {
			numTwo++
		}
		if num == 3 {
			numThree++
		}
	}

	return numTwo, numThree
}

func main() {
	hasTwice := 0
	hasThrice := 0

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		numTwo, numThree := count(line)
		fmt.Printf("line: %v #2: %v #3: %v\n", line, numTwo, numThree)

		if numTwo > 0 {
			hasTwice++
		}
		if numThree > 0 {
			hasThrice++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("read failed: %v", err)
	}

	fmt.Println(hasTwice * hasThrice)
}
