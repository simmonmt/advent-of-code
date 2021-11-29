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
)

var (
	row = flag.Int("row", -1, "row number")
	col = flag.Int("col", -1, "col number")
)

// calculate the number of the code in the first column for each row

func main() {
	flag.Parse()

	if *row == -1 || *col == -1 {
		log.Fatalf("--row and --col are required")
	}

	curCodeNum := 1
	curRow := 1

	for curRow < *row {
		curRow++
		curCodeNum += (curRow - 1)
	}
	fmt.Printf("row %v codeNum %v\n", curRow, curCodeNum)

	curCol := 1
	colInc := curRow + 1

	for curCol < *col {
		curCodeNum += colInc
		colInc++
		curCol++
	}

	fmt.Printf("row %v col %v codeNum %v\n", curRow, curCol, curCodeNum)
}
