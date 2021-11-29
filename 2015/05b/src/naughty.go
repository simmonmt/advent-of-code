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
	numFoundRepeatingPairs, numFoundSeparatedPairs := 0, 0
	pairs := map[string]int{}

	var l2, l3 rune
	for i, r := range line {
		if l2 != 0 {
			pair := string([]rune{l2, r})
			if foundLoc, found := pairs[pair]; found {
				//fmt.Printf("found %v now %d was %d\n", pair, i, foundLoc)
				if i != foundLoc+1 {
					numFoundRepeatingPairs++
				}
			} else {
				//fmt.Printf("adding pair %v\n", pair)
				pairs[pair] = i
			}

			if l3 != 0 {
				if r == l3 {
					numFoundSeparatedPairs++
				}
			}
		}

		l3 = l2
		l2 = r
	}

	return numFoundRepeatingPairs > 0 && numFoundSeparatedPairs > 0
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
