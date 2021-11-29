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

// This is the Go version of the assembly program in ../input.txt

package main

import (
	"flag"
	"fmt"
)

var (
	runlong = flag.Bool("runlong", false, "")
)

func main() {
	flag.Parse()

	//reg := [6]int{} //0, 0, 1, 0, 987, 151}

	var c int
	if *runlong {
		c = 10551387
		//d = 10550400
	} else {
		c = 987
		//d = 151
	}

	out := 0
	for a := 1; a <= c; a++ {
		for b := 1; b <= c; b++ {
			if a*b == c {
				out += a
			}
			// This conditional added after analysis of
			// transliterated assembly. Without this, the
			// loop is way too expensive to run with
			// --runlong
			if a*b > c {
				break
			}
		}
	}

	fmt.Println(out)
}
