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

package puzzle

import (
	"fmt"

	"github.com/simmonmt/aoc/2019/common/intmath"
)

var (
	kPrimes = []int{
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97,
	}
)

func Factor(n, d int) (int, int) {
	if intmath.Abs(n) > kPrimes[len(kPrimes)-1] || intmath.Abs(d) > kPrimes[len(kPrimes)-1] {
		panic(fmt.Sprintf("too large n=%d d=%d", n, d))
	}

	for {
		changed := false
		for _, p := range kPrimes {
			if p > intmath.Abs(n) && p > intmath.Abs(d) {
				break
			}

			if n%p == 0 && d%p == 0 {
				n /= p
				d /= p
				changed = true
			}
		}
		if !changed {
			break
		}
	}

	return n, d
}
