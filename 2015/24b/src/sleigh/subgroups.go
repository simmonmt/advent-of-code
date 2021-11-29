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

type Subgrouper struct {
	group   []int
	sz      int
	indexes []int
	done    bool
}

func NewSubgrouper(group []int, sz int) *Subgrouper {
	if sz > len(group) {
		panic(fmt.Sprintf("sz %d > len(group) %d", sz, len(group)))
	}

	indexes := make([]int, sz)
	for i := 0; i < sz; i++ {
		indexes[i] = i
	}

	return &Subgrouper{
		group:   group,
		sz:      sz,
		indexes: indexes,
		done:    false,
	}
}

func (sg *Subgrouper) Next() (subgroup, rest []int, ok bool) {
	if sg.done {
		ok = false
		return
	}

	subgroup = make([]int, sg.sz)
	rest = make([]int, len(sg.group)-sg.sz)
	curIdx := 0
	for i, val := range sg.group {
		if curIdx < sg.sz && sg.indexes[curIdx] == i {
			subgroup[curIdx] = val
			curIdx++
		} else {
			rest[i-curIdx] = val
		}
	}

	for i := 0; i < sg.sz; i++ {
		subgroup[i] = sg.group[sg.indexes[i]]
	}

	i := len(sg.indexes) - 1
	carry := false
	for {
		if i == -1 {
			// We carried off the beginning
			sg.done = true
			break
		}

		// if sz = 3,
		// index[2] max is len(group)-1
		// index[1] max is len(group)-2
		// index[0] max is len(group)-3

		indexMax := len(sg.group) - (len(sg.indexes) - i)
		if sg.indexes[i]+1 > indexMax {
			carry = true
			i--
			continue
		}

		sg.indexes[i]++
		if carry {
			for j := i + 1; j < sg.sz; j++ {
				sg.indexes[j] = sg.indexes[j-1] + 1
			}
		}
		break
	}

	ok = true
	return
}
