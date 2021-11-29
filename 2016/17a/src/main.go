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

	"maze"
)

var (
	passcode = flag.String("passcode", "", "passcode")
	width    = flag.Int("width", 4, "width")
	height   = flag.Int("height", 4, "height")
)

func main() {
	flag.Parse()

	if *passcode == "" {
		log.Fatal("--passcode is required")
	}

	found, path := maze.RunMaze(*width, *height, *passcode)
	if !found {
		fmt.Println("no path found")
	} else {
		fmt.Printf("path: %v\n", path)
	}
}
