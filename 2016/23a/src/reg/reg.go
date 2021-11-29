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

package reg

import (
	"fmt"
	"strings"
)

type Reg int

const (
	A Reg = iota
	B
	C
	D
)

func (r Reg) String() string {
	switch r {
	case A:
		return "a"
	case B:
		return "b"
	case C:
		return "c"
	case D:
		return "d"
	default:
		return "UNKNOWN"
	}
}

func FromString(name string) (Reg, error) {
	switch strings.ToLower(name) {
	case "a":
		return A, nil
	case "b":
		return B, nil
	case "c":
		return C, nil
	case "d":
		return D, nil
	default:
		return A, fmt.Errorf("unknown reg '%v'", name)
	}
}
