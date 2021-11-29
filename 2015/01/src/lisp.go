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
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)

		firstBasement := -1
		floor := 0
		for i, c := range line {
			switch c {
			case '(':
				floor++
				break
			case ')':
				floor--
				break
			default:
				log.Fatalf("unknown char %v\n", c)
			}

			if floor == -1 && firstBasement == -1 {
				firstBasement = i + 1
			}
		}

		fmt.Printf("%d first basement: %v\n", lineNum, firstBasement)
		fmt.Printf("%d end floor: %v\n", lineNum, floor)
	}
}
