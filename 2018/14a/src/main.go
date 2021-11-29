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
	"flag"
	"fmt"

	"logger"
)

var (
	verbose       = flag.Bool("verbose", false, "verbose")
	maxNumRecipes = flag.Int("max_recipes", -1, "max recipes")
)

func dumpRecipes(recipes []int, elves []int) {
	elfMap := map[int]int{}
	for i, elf := range elves {
		elfMap[elf] = i
	}

	for i, recipe := range recipes {
		if num, found := elfMap[i]; found {
			if num == 0 {
				fmt.Printf("(%d) ", recipe)
			} else if num == 1 {
				fmt.Printf("[%d] ", recipe)
			} else {
				panic("bad elf num")
			}
		} else {
			fmt.Printf(" %d  ", recipe)
		}
	}
	fmt.Println()
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	recipes := [1000000]int{}
	recipes[0] = 3
	recipes[1] = 7

	numRecipes := 2

	elves := [2]int{0, 1}

	dumpRecipes(recipes[0:numRecipes], elves[:])

	for numRecipes <= *maxNumRecipes+10 {
		newSum := recipes[elves[0]] + recipes[elves[1]]
		newRecipes := []int{}
		if newSum == 0 {
			newRecipes = []int{0}
		} else {
			for newSum > 0 {
				newRecipes = append(newRecipes, newSum%10)
				newSum /= 10
			}
		}

		for i := len(newRecipes) - 1; i >= 0; i-- {
			recipes[numRecipes] = newRecipes[i]
			numRecipes++
		}

		for i := range elves {
			off := recipes[elves[i]] + 1
			pos := (elves[i] + off) % numRecipes
			elves[i] = pos
		}

		if *verbose {
			dumpRecipes(recipes[0:numRecipes], elves[:])
		}
	}

	// for i := 0; i < numRecipes; i++ {
	// 	fmt.Print(recipes[i])
	// }
	// fmt.Println()

	fmt.Println(numRecipes)
	fmt.Println(recipes[*maxNumRecipes : *maxNumRecipes+10])
}
