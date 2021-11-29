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
)

var (
	seq = flag.Int("seq", -1, "seq")
)

func main() {
	flag.Parse()

	objs := []string{
		"pointer",    // 1
		"hypercube",  // 2
		"cake",       // 4
		"tambourine", // 8
		"mouse",      // 16
		"coin",       // 32
		"mug",        // 64
		"monolith",   // 128
	}

	i := *seq
	for j := 0; j < 8; j++ {
		if (i & (1 << j)) != 0 {
			fmt.Printf("drop %s\n", objs[j])
		}
	}
}
