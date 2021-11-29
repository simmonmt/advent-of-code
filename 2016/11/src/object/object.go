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

package object

import "fmt"

type Object int8

func Microchip(num int8) Object {
	return Object(num)
}

func Generator(num int8) Object {
	return Object(-num)
}

func (o Object) IsMicrochip() bool {
	return o > 0
}

func (o Object) Num() int {
	if o < 0 {
		return int(-o)
	} else {
		return int(o)
	}
}

func (o Object) String() string {
	if o < 0 {
		return fmt.Sprintf("%dG", -o)
	} else {
		return fmt.Sprintf("%dM", o)
	}
}

func (o Object) Serialize() byte {
	if o < 0 {
		return byte('A' + -o - 1)
	} else {
		return byte('a' + o - 1)
	}
}

func Deserialize(ser byte) Object {
	if ser > 'Z' {
		return Microchip(int8(ser - 'a' + 1))
	} else {
		return Generator(int8(ser - 'A' + 1))
	}
}

type Objects []Object

func (o Objects) Len() int {
	return len(o)
}

func abs(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func (o Objects) Less(i, j int) bool {
	absI := abs(int(o[i]))
	absJ := abs(int(o[j]))

	if absI < absJ {
		return true
	} else if absI > absJ {
		return false
	} else {
		return o[i] < 0
	}
}

func (o Objects) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}
