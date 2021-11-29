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

package strutil

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2020/common/intmath"
)

func StringToInt64s(str string) ([]int64, error) {
	out := []int64{}
	for _, s := range strings.Split(str, ",") {
		v, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("bad value %v: %v", s, err)
		}

		out = append(out, v)
	}
	return out, nil
}

func StringDiff(a, b string) (bool, string) {
	minLen := intmath.IntMin(len(a), len(b))
	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			return true, fmt.Sprintf("mismatch at char %v: %v != %v",
				i, strconv.Quote(string(a[i])),
				strconv.Quote(string(b[i])))
		}
	}

	if len(a) != len(b) {
		return true, fmt.Sprintf("len %v != len %v", len(a), len(b))
	}

	return false, ""
}
