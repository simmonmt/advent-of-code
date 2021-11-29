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
	"log"

	"data"
	"util"
)

var (
	initialState = flag.String("initial_state", "", "initial state")
	goalLength   = flag.Int("goal_length", -1, "goal length")
)

func main() {
	flag.Parse()

	if *initialState == "" {
		log.Fatal("--initial_state is required")
	}
	if *goalLength == -1 {
		log.Fatal("--goal_length is required")
	}

	d := util.StrToBoolArray(*initialState)
	for len(d) < *goalLength {
		d = data.Grow(d)
	}

	d = d[0:*goalLength]
	cksum := data.Checksum(d)

	fmt.Println(util.BoolArrayToStr(cksum))
}
