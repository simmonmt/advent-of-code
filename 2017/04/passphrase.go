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
	"os"
	"sort"
	"strings"
)

func sortWord(word string) string {
	chars := []int{}
	for _, r := range word {
		chars = append(chars, int(r))
	}
	sort.Ints(chars)

	runes := []rune{}
	for _, c := range chars {
		runes = append(runes, rune(c))
	}

	return string(runes)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	numValid := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		wordList := strings.Split(line, " ")
		valid := true
		words := map[string]bool{}
		for _, word := range wordList {
			sortedWord := sortWord(word)

			if _, found := words[sortedWord]; found {
				valid = false
				break
			}
			words[sortedWord] = true
		}

		if valid {
			fmt.Printf("  valid %s\n", line)
			numValid++
		} else {
			fmt.Printf("invalid %s\n", line)
		}

		//fmt.Printf("%v %s\n", valid, wordList)
	}

	fmt.Printf("num valid %d\n", numValid)
}
