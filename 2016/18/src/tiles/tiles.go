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

package tiles

import "fmt"

type Row []bool

func MakeRow(str string) (Row, error) {
	out := make([]bool, len(str))
	for i, t := range str {
		if t == '^' {
			out[i] = true
		} else if t == '.' {
			out[i] = false
		} else {
			return nil, fmt.Errorf(`invalid row "%v"`, str)
		}
	}
	return out, nil
}

func (r Row) String() string {
	out := make([]rune, len(r))
	for i, t := range r {
		if t {
			out[i] = '^'
		} else {
			out[i] = '.'
		}
	}
	return string(out)
}

func (r Row) IsTrap(i int) bool {
	if i < 0 || i >= len(r) {
		return false
	}
	return r[i]
}

func (r Row) NumSafe() int {
	n := 0
	for _, t := range r {
		if !t {
			n++
		}
	}
	return n
}

func (r Row) Next() Row {
	var next Row = make([]bool, len(r))

	for i := range r {
		left := r.IsTrap(i - 1)
		center := r.IsTrap(i)
		right := r.IsTrap(i + 1)

		a := left && center && !right
		b := !left && center && right
		c := left && !center && !right
		d := !left && !center && right

		isTrap := a || b || c || d
		next[i] = isTrap
	}

	return next
}
