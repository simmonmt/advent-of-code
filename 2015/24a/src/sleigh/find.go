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

package sleigh

import "fmt"

// for groupOneSize = 1 to len-2
//   found = []
//
//   for each group(groupOneSize) groupOne
//     if avg(rest) / 2 != sum(groupOne)
//       continue
//     sideCap = avg(rest) / 2
//     if there's a packing p whose size == sidecap
//     append (groupOne, p (groupTwo), rest - p (groupThree)) to found
//
//   if found not empty
//     emit one with smallest entanglement

// object that emits all subgroups (order doesn't matter) of given
// size

// function that does bin-packing

func sumArr(vals []int) int {
	out := 0
	for _, val := range vals {
		out += val
	}
	return out
}

func FindWithGroupOneSize(values []int, groupOneSize int) [][]int {
	found := [][]int{}

	sg := NewSubgrouper(values, groupOneSize)
	for {
		groupOne, rest, ok := sg.Next()
		if !ok {
			break
		}

		cap := sumArr(groupOne)
		if cap*2 != sumArr(rest) {
			//fmt.Printf("ignoring groupOne %v rest %v; no cap %v match\n", groupOne, rest, cap)
			continue
		}

		if pack := BinPack(rest, cap); pack != nil {
			fmt.Printf("found packing for %v\n", groupOne)
			found = append(found, groupOne)
		}
	}

	return found
}
