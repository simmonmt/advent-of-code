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
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
)

func readPrereqs() (map[string][]string, error) {
	prereqs := map[string][]string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		prereq := parts[1]
		name := parts[7]

		if _, ok := prereqs[name]; !ok {
			prereqs[name] = []string{}
		}
		prereqs[name] = append(prereqs[name], prereq)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return prereqs, nil
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	prereqs, err := readPrereqs()
	if err != nil {
		log.Fatal(err)
	}

	todo := map[string]bool{}
	for name, ps := range prereqs {
		todo[name] = true
		for _, p := range ps {
			todo[p] = true
		}
	}

	order := ""

	done := map[string]bool{}
	for len(todo) > 0 {
		cands := []string{}

		for step, _ := range todo {
			ps := prereqs[step]
			if ps == nil {
				ps = []string{}
			}

			alldone := true
			for _, p := range ps {
				if _, ok := done[p]; !ok {
					alldone = false
					break
				}
			}
			if !alldone {
				continue
			}

			cands = append(cands, step)
		}

		if len(cands) == 0 {
			panic("empty")
		}

		sort.Strings(cands)
		chosen := cands[0]
		done[chosen] = true
		delete(todo, chosen)
		order += chosen
	}

	fmt.Println(order)
}
