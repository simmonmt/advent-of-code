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
	"strings"
)

// (0,0) (1,0) (2,0)
//    (0,1) (1,1) (2,1)
// (0,2) (1,2) (2,2)

func readDirections(in io.Reader) ([]string, error) {
	reader := bufio.NewReader(in)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)
	return strings.Split(line, ","), nil
}

func advance(x, y int, dir string) (int, int) {
	switch dir {
	case "n":
		x += 1
		break
	case "s":
		x -= 1
		break
	case "se":
		if y%2 == 0 {
			x -= 1
		}
		y += 1
		break
	case "ne":
		if y%2 != 0 {
			x += 1
		}
		y += 1
		break
	case "sw":
		if y%2 == 0 {
			x -= 1
		}
		y -= 1
		break
	case "nw":
		if y%2 != 0 {
			x += 1
		}
		y -= 1
		break
	default:
		panic(fmt.Sprintf("unknown dir %v", dir))
	}
	return x, y
}

func distance(destX, destY int) int {
	x, y := 0, 0
	dist := 0
	for x != destX || y != destY {
		if y == destY {
			if x < destX {
				x++
			} else {
				x--
			}
		} else if x == destX {
			if y > destY {
				if y%2 == 0 {
					x, y = advance(x, y, "nw")
				} else {
					x, y = advance(x, y, "sw")
				}
			} else {
				if y%2 == 0 {
					x, y = advance(x, y, "ne")
				} else {
					x, y = advance(x, y, "se")
				}
			}
		} else if x < destX && y < destY {
			x, y = advance(x, y, "ne")
		} else if x < destX && y > destY {
			x, y = advance(x, y, "nw")
		} else if x > destX && y < destY {
			x, y = advance(x, y, "se")
		} else if x > destX && y > destY {
			x, y = advance(x, y, "sw")
		} else {
			panic("unknown path")
		}

		dist++
		//fmt.Printf("x %v y %v dist %v\n", x, y, dist)
	}
	return dist
}

func main() {
	dirs, err := readDirections(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read directions: %v\n", err)
	}

	maxDist := 0
	x, y := 0, 0
	for _, dir := range dirs {
		x, y = advance(x, y, dir)

		dist := distance(x, y)
		if dist > maxDist {
			maxDist = dist
			fmt.Printf("maxDist now %d\n", maxDist)
		}
		//fmt.Printf("dir %v now x %v y %v\n", dir, x, y)
	}
	//fmt.Printf("x %v y %v\n", x, y)

	dist := distance(x, y)

	fmt.Printf("dist %d\n", dist)
	fmt.Printf("max dist %d\n", maxDist)
}
