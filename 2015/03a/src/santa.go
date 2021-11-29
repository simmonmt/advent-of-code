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
	"os"
	"strings"
)

type Position struct {
	X, Y int
}

func readCommands(in io.Reader) string {
	reader := bufio.NewReader(in)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func main() {
	commands := readCommands(os.Stdin)

	pos := Position{0, 0}
	presents := map[Position]int{}
	presents[pos]++

	for _, command := range commands {
		//fmt.Sprintf("command '%c'\n", command)
		switch command {
		case '<':
			pos = Position{pos.X - 1, pos.Y}
			break
		case '>':
			pos = Position{pos.X + 1, pos.Y}
			break
		case '^':
			pos = Position{pos.X, pos.Y - 1}
			break
		case 'v':
			pos = Position{pos.X, pos.Y + 1}
			break
		default:
			panic(fmt.Sprintf("unknown commmand '%c'", command))
		}

		presents[pos]++
	}

	fmt.Println(presents)
	fmt.Println(len(presents))
}
