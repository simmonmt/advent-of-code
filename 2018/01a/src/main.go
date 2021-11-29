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
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var freq int64

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		val, err := strconv.ParseInt(strings.TrimPrefix(line, "+"), 0, 32)
		if err != nil {
			log.Fatalf("failed to parse %v: %v", line, err)
		}

		freq += val
		fmt.Printf("%v %v %v\n", line, val, freq)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("read failed: %v", err)
	}
}
