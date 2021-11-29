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

func isNice(line string) bool {
	numVowels, numDups, numBad := 0, 0, 0

	var last rune
	for _, c := range line {
		if c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u' {
			numVowels++
		}

		if last == c {
			numDups++
		}

		if (last == 'a' && c == 'b') ||
			(last == 'c' && c == 'd') ||
			(last == 'p' && c == 'q') ||
			(last == 'x' && c == 'y') {
			numBad++
		}

		last = c
	}

	return numVowels >= 3 && numDups > 0 && numBad == 0
}

func main() {
	numNice := 0

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		if isNice(line) {
			fmt.Printf("nice: %v\n", line)
			numNice++
		} else {
			fmt.Printf("naughty: %v\n", line)
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("reading stdin: %v", err)
	}

	fmt.Println(numNice)
}
