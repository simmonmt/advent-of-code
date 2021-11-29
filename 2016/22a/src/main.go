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
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"grid"
	"logger"
	"node"
	"solver"
)

var (
	sizePattern = regexp.MustCompile(`^/dev/grid/node-x([0-9]+)-y([0-9]+) +([0-9]+)T +([0-9]+)T +[0-9]+T +[0-9]+%$`)

	verbose = flag.Bool("verbose", false, "verbose")
)

func readInput(r io.Reader) (width uint8, height uint8, nodes []node.Node, err error) {
	yxnodes := []node.Node{}

	var maxX, maxY uint64

	reader := bufio.NewReader(r)
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if !strings.HasPrefix(line, "/") {
			continue
		}

		line = strings.TrimSpace(line)
		matches := sizePattern.FindStringSubmatch(line)
		if matches == nil {
			return 0, 0, nil, fmt.Errorf("%d: failed to parse", lineNum)
		}

		x, err := strconv.ParseUint(matches[1], 10, 8)
		if err != nil {
			return 0, 0, nil, fmt.Errorf("%d: failed to parse x: %v", lineNum, err)
		}
		y, err := strconv.ParseUint(matches[2], 10, 8)
		if err != nil {
			return 0, 0, nil, fmt.Errorf("%d: failed to parse y: %v", lineNum, err)
		}
		size, err := strconv.ParseUint(matches[3], 10, 16)
		if err != nil {
			return 0, 0, nil, fmt.Errorf("%d: failed to parse size: %v", lineNum, err)
		}
		used, err := strconv.ParseUint(matches[4], 10, 16)
		if err != nil {
			return 0, 0, nil, fmt.Errorf("%d: failed to parse used: %v", lineNum, err)
		}

		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}

		yxnodes = append(yxnodes, *node.New(uint16(size), uint16(used)))
	}

	width = uint8(maxX + 1)
	height = uint8(maxY + 1)

	// They gave us nodes sorted by x then y. We needed y then
	// x. Rearrange the nodes array.
	nodes = make([]node.Node, len(yxnodes))
	for i, n := range yxnodes {
		x := i / int(height)
		y := i % int(height)
		nodes[y*int(width)+x] = n
	}

	return width, height, nodes, nil
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	width, height, nodes, err := readInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(nodes)

	// numViable := 0
	// for i := 0; i < len(nodes); i++ {
	// 	for j := 0; j < len(nodes); j++ {
	// 		if i == j {
	// 			continue
	// 		}

	// 		a, b := nodes[i], nodes[j]
	// 		if a.Used == 0 {
	// 			continue
	// 		}
	// 		if a.Used > (b.Size - b.Used) {
	// 			continue
	// 		}

	// 		numViable++
	// 	}
	// }

	// fmt.Printf("#viable = %v\n", numViable)

	g, err := grid.New(width, height, uint8(width-1), 0, nodes)
	if err != nil {
		log.Fatal(err)
	}

	// g2, err := grid.Deserialize(width, height, g.Serialize())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(g2.Get(35, 21))

	// return

	found, numSteps := solver.Solve(width, height, g)
	fmt.Printf("found %v numSteps %v\n", found, numSteps)
}
