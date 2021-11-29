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
	"os"
	"strconv"

	"sleigh"
)

func calcQE(values []int) int {
	out := 1
	for _, val := range values {
		out *= val
	}
	return out
}

func main() {
	values := []int{}
	for i := 1; i < len(os.Args); i++ {
		if val, err := strconv.ParseInt(os.Args[i], 10, 32); err != nil {
			log.Fatalf("failed to parse %v: %v", os.Args[i], err)
		} else {
			values = append(values, int(val))
		}
	}

	for groupOneSize := 1; groupOneSize < len(values)-2; groupOneSize++ {
		fmt.Printf("== size %v\n", groupOneSize)
		groups := sleigh.FindWithGroupOneSize(values, groupOneSize)
		if len(groups) == 0 {
			fmt.Println("no groups found")
			continue
		}

		minQE := -1
		for _, group := range groups {
			if qe := calcQE(group); minQE == -1 || qe < minQE {
				minQE = qe
			}
		}

		fmt.Printf("found min qe %v\n", minQE)
		break
	}
}
