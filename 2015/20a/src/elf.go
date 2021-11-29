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
	"math"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %v goal", os.Args[0])
	}
	goal, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("failed to parse goal %v", os.Args[1])
	}

	maxPresents := 0
	for houseNum := 1; ; houseNum++ {
		numPresents := 0 //houseNum * 10
		// fmt.Printf("house %v:", houseNum)

		// TODO(simmonmt): count up to sqrt, use div & rem as
		// elf numbers.
		lim := int(math.Sqrt(float64(houseNum)))
		for elfNum := 1; elfNum <= lim; elfNum++ {
			factors := map[int]bool{}

			if houseNum%elfNum == 0 {
				// fmt.Printf(" %v", elfNum)
				factors[elfNum] = true

				other := houseNum / elfNum
				if other != 1 {
					// fmt.Printf(" %v", other)
					factors[other] = true
				}
			}

			for factor := range factors {
				numPresents += factor * 10
			}

		}
		// fmt.Println()

		// fmt.Printf("house %v numPresents %v\n", houseNum, numPresents)

		if numPresents > maxPresents {
			maxPresents = numPresents
		}
		if houseNum != 0 && houseNum%1000 == 0 {
			fmt.Printf("house %v max %v\n", houseNum, maxPresents)
		}

		//fmt.Printf("house %v presents %v\n", houseNum, numPresents)
		if numPresents >= goal {
			fmt.Printf("house %v presents %v\n", houseNum, numPresents)
			break
		}
	}
}
