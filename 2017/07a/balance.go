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
	"regexp"
	"strings"
)

var (
	pattern = regexp.MustCompile(`^([a-z]+) \(([0-9]+)\)(?: -> (.*))?$`)
)

type Tree struct {
	elems map[string]string
}

func NewTree() *Tree {
	return &Tree{elems: map[string]string{}}
}

func (t *Tree) Insert(bot, top string) {
	t.elems[top] = bot
}

func (t *Tree) Bottom() string {
	var bot string
	for _, bot = range t.elems {
		break
	}

	for {
		next, found := t.elems[bot]
		if !found {
			return bot
		}

		bot = next
	}
}

func main() {
	tree := NewTree()

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		matches := pattern.FindStringSubmatch(line)
		if len(matches) == 0 {
			log.Fatalf("failed to parse %v", line)
		}
		bot := matches[1]
		// weight, err := strconv.Atoi(matches[2])
		// if err != nil {
		// 	log.Fatalf("failed to parse weight %v in %v", matches[2], line)
		// }

		topsStr := matches[3]
		if topsStr == "" {
			continue
		}
		tops := strings.Split(topsStr, ", ")

		//fmt.Printf("%v %v\n", bot, tops)

		for _, top := range tops {
			tree.Insert(bot, top)
		}
	}

	fmt.Println(tree.Bottom())
}
