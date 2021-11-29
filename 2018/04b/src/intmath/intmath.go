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

package intmath

import (
	"fmt"
	"strconv"
)

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
