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
	"knot"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %v key", os.Args[0])
	}
	key := os.Args[1]

	grid := make([][]bool, 128)
	for i := range grid {
		row := make([]bool, 128)

		rowHashKey := fmt.Sprintf("%s-%d", key, i)
		rowHash := knot.Hash(rowHashKey)

		// rowHash = "a0c2017" + rowHash[7:]

		for i, r := range rowHash {
			num, err := strconv.ParseUint(string(r), 16, 64)
			if err != nil {
				log.Fatalf("failed to parse %v\n", r)
			}

			// fmt.Printf("row[i*4]=(num&8) row[%d]=%d&8=%d=%v\n",
			// 	i*4, num, (num & 0x8), (num&0x8) != 0)

			row[i*4] = (num & 0x8) != 0
			row[i*4+1] = (num & 0x4) != 0
			row[i*4+2] = (num & 0x2) != 0
			row[i*4+3] = (num & 0x1) != 0
		}

		grid[i] = row
	}

	for i := 0; i < 8; i++ {
		row := grid[i]
		for j := 0; j < 8; j++ {
			if row[j] {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println("")
	}

	numFilled := 0
	for _, row := range grid {
		for _, cell := range row {
			if cell {
				numFilled++
			}
		}
	}
	fmt.Println(numFilled)
}
