// Copyright 2022 Google LLC
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

package intmath

import (
	"fmt"
	"strconv"
)

func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func IntMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Abs(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func Uint64Max(a, b uint64) uint64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func AtoiOrDie(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse %v: %v", s, err))
	}
	return val
}

var (
	kPrimes = []int{
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97,
	}
)

func GCD(vs ...int) int {
	d := 1

	for _, p := range kPrimes {
		for {
			all := true
			for _, v := range vs {
				if v < p || v%p != 0 {
					all = false
				}
			}

			if !all {
				break
			}

			d *= p
			for i := range vs {
				vs[i] /= p
			}
		}
	}

	for _, v := range vs {
		if v > kPrimes[len(kPrimes)-1] {
			panic(fmt.Sprintf("%d too big", v))
		}
	}

	return d
}
