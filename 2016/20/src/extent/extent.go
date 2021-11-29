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

package extent

import (
	"fmt"
	"strconv"
	"strings"

	"intmath"
)

type Extent struct {
	Start, End uint64
}

func Parse(str string) (*Extent, error) {
	parts := strings.SplitN(str, "-", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid extent %v", str)
	}

	start, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf(`invalid start "%v": %v`, err)
	}
	end, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf(`invalid end "%v": %v`, err)
	}

	if start > end {
		return nil, fmt.Errorf("start %v > end %v", start, end)
	}

	return &Extent{Start: start, End: end}, nil
}

func (e Extent) String() string {
	return fmt.Sprintf("%d-%d", e.Start, e.End)
}

func (e *Extent) Remove(remove *Extent) []*Extent {
	if remove.Start <= e.Start {
		if remove.End >= e.End {
			return nil
		} else {
			return []*Extent{&Extent{Start: remove.End + 1, End: e.End}}
		}
	} else if remove.Start >= e.End {
		return nil
	} else if remove.End >= e.End {
		// remove.start > e.start
		return []*Extent{&Extent{Start: e.Start, End: remove.Start - 1}}
	} else {
		// remove.start > e.start && remove.end < e.end
		return []*Extent{
			&Extent{Start: e.Start, End: remove.Start - 1},
			&Extent{Start: remove.End + 1, End: e.End},
		}
	}
}

func (e *Extent) Merge(other *Extent) *Extent {
	if other.Start <= e.Start {
		if other.End >= e.Start-1 {
			return &Extent{
				Start: other.Start,
				End:   intmath.Uint64Max(other.End, e.End),
			}
		} else {
			return nil // disjoint
		}
	} else if other.Start <= e.End+1 {
		return &Extent{
			Start: e.Start,
			End:   intmath.Uint64Max(other.End, e.End),
		}
	} else {
		return nil // disjoint
	}
}
