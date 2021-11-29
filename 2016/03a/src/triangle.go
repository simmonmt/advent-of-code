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
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	inputPattern = regexp.MustCompile(`\w+`)
)

func readInput(r io.Reader) ([][3]int, error) {
	out := [][3]int{}

	reader := bufio.NewReader(r)
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		matches := inputPattern.FindAllString(line, -1)
		if matches != nil && len(matches) != 3 {
			return nil, fmt.Errorf("%d: expected 3 words, found %v", lineNum, matches)
		}

		dims := [3]int{}
		for i, match := range matches {
			match = strings.TrimSpace(match)

			dim, err := strconv.ParseUint(match, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("failed to parse dim %v: %v", match, err)
			}
			dims[i] = int(dim)
		}

		out = append(out, dims)
	}

	return out, nil
}

func main() {
	triangles, err := readInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	numPossible := 0
	for _, tri := range triangles {
		possible := tri[0]+tri[1] > tri[2] && tri[0]+tri[2] > tri[1] && tri[1]+tri[2] > tri[0]
		if possible {
			numPossible++
		}
	}

	fmt.Printf("num possible = %v\n", numPossible)
}
