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

// 4855969 too high

package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	regZero = flag.Int("reg_zero", 0, "register 0 initial value")
)

func main() {
	flag.Parse()

	reg := [6]int64{}
	reg[0] = int64(*regZero)

	vals := map[int64]int{}
	iter := 0
six:
	iter++

	reg[3] = reg[4] | 0x10000
	reg[4] = 12670166

eight:
	reg[2] = reg[3] & 0xff
	reg[4] = reg[4] + reg[2]
	reg[4] = reg[4] & 0xffffff
	reg[4] = reg[4] * 65899
	reg[4] = reg[4] & 0xffffff

	if 256 > reg[3] {
		reg[2] = 1
	} else {
		reg[2] = 0
	}

	if reg[2] == 1 {
		goto twentyeight
	}

	for reg[2] = 0; ; reg[2]++ {
		reg[5] = reg[2] + 1
		reg[5] = reg[5] * 256

		if reg[5] > reg[3] {
			break
		}
	}

	reg[3] = reg[2]
	goto eight

twentyeight:
	fmt.Printf("0=%v 4=%v iter %v\n", reg[0], reg[4], iter)
	if was, found := vals[reg[4]]; found {
		fmt.Printf("repeat at %v was %v delta %v\n",
			iter, was, iter-was)
	}
	vals[reg[4]] = iter

	if reg[4] == reg[0] {
		fmt.Println(reg)
		os.Exit(0)
	} else {
		goto six
	}
}
