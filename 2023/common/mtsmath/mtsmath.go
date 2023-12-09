// Copyright 2023 Google LLC
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

package mtsmath

import (
	"fmt"
	"strconv"

	"golang.org/x/exp/constraints"
)

func Abs[V constraints.Integer | constraints.Float](val V) V {
	if val > 0 {
		return val
	} else {
		return -val
	}
}

func AtoiOrDie(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse %v: %v", s, err))
	}
	return val
}

func gcd(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func GCD(nums ...int64) int64 {
	out := nums[0]
	for i := 1; i < len(nums); i++ {
		out = gcd(out, nums[i])
	}
	return out
}

func LCM(nums ...int64) int64 {
	out := int64(nums[0])
	for i := 1; i < len(nums); i++ {
		out = out * (int64(nums[i]) / gcd(out, int64(nums[i])))
	}
	return out
}
