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
	"log"
	"math"
	"os"
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %v goal", os.Args[0])
	}
	goal, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("failed to parse goal %v", os.Args[1])
	}

	printer := message.NewPrinter(language.English)

	maxPresents := 0
	for houseNum := 1; ; houseNum++ {
		numPresents := 0 //houseNum * 10
		// fmt.Printf("house %v:", houseNum)

		// TODO(simmonmt): count up to sqrt, use div & rem as
		// elf numbers.
		lim := int(math.Sqrt(float64(houseNum)))
		for i := 1; i <= lim; i++ {
			elves := map[int]bool{}

			if houseNum%i == 0 {
				// fmt.Printf(" %v", elfNum)
				elves[i] = true

				other := houseNum / i
				if other != 1 {
					// fmt.Printf(" %v", other)
					elves[other] = true
				}
			}

			for elf := range elves {
				if houseNum > elf*50 {
					// fmt.Printf("house %v: skipping %v\n", houseNum, elf)
					continue
				}

				numPresents += elf * 11
			}

		}
		// fmt.Println()

		// fmt.Printf("house %v numPresents %v\n", houseNum, numPresents)

		if numPresents > maxPresents {
			maxPresents = numPresents
		}
		if houseNum != 0 && houseNum%1000 == 0 {
			printer.Printf("house %d max %d\n", houseNum, maxPresents)
		}

		//fmt.Printf("house %v presents %v\n", houseNum, numPresents)
		if numPresents >= goal {
			printer.Printf("house %d presents %d\n", houseNum, numPresents)
			break
		}
	}
}
