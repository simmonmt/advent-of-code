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

package main

import (
	"flag"
	"fmt"

	"discs"
)

var (
	smallInput = flag.Bool("small_input", false, "use small input")
	addBDisc   = flag.Bool("add_b_disc", false, "add the B disc")
)

func main() {
	flag.Parse()

	var descs []discs.DiscDesc
	if *smallInput {
		descs = []discs.DiscDesc{
			discs.DiscDesc{5, 4},
			discs.DiscDesc{2, 1},
		}
	} else {
		descs = []discs.DiscDesc{
			discs.DiscDesc{7, 0},
			discs.DiscDesc{13, 0},
			discs.DiscDesc{3, 2},
			discs.DiscDesc{5, 2},
			discs.DiscDesc{17, 0},
			discs.DiscDesc{19, 7},
		}
	}

	if *addBDisc {
		descs = append(descs, discs.DiscDesc{11, 0})
	}

	posns := make([]int, len(descs))
	for i := range posns {
		posns[i] = descs[i].Start
	}

	for t := 1; ; t++ {
		discs.Advance(posns)
		if discs.Success(descs, posns) {
			fmt.Println(t)
			break
		}
	}
}
