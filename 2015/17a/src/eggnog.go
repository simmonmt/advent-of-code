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
	"strconv"
	"strings"
)

func foo(amt, goal int, containers []int) int {
	//fmt.Printf("foo(amt=%d,goal=%d,containers=%d)\n", amt, goal, containers)

	if amt == goal {
		return 1
	} else if amt > goal || len(containers) == 0 {
		return 0
	}

	found := foo(amt+containers[0], goal, containers[1:])
	found += foo(amt, goal, containers[1:])
	return found
}

func readInput(r io.Reader) ([]int, error) {
	vals := []int{}

	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		val, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			return nil, fmt.Errorf("failed to parse %v: %v", line, err)
		}

		vals = append(vals, val)
	}

	return vals, nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %v goal", os.Args[0])
	}
	goal, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("invalid goal %v: %v", os.Args[1], err)
	}

	containers, err := readInput(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read containers: %v", err)
	}
	fmt.Println(containers)

	fmt.Println(foo(0, goal, containers))
}
