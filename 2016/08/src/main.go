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
	"strings"

	"command"
	"screen"
)

var (
	screenWidth  = flag.Int("width", 50, "screen width")
	screenHeight = flag.Int("height", 6, "screen height")
)

func readInput(r io.Reader) ([]command.Command, error) {
	cmds := []command.Command{}

	reader := bufio.NewReader(r)
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)

		var cmd command.Command
		switch {
		case strings.HasPrefix(line, "rect "):
			cmd, err = command.ParseRect(line)
			break
		case strings.HasPrefix(line, "rotate row "):
			cmd, err = command.ParseRotateRow(line)
			break
		case strings.HasPrefix(line, "rotate column "):
			cmd, err = command.ParseRotateColumn(line)
			break
		default:
			return nil, fmt.Errorf("unknown command %v", line)
		}

		if err != nil {
			return nil, fmt.Errorf("%d: %v", lineNum, err)
		}

		cmds = append(cmds, cmd)
	}

	return cmds, nil
}

func main() {
	flag.Parse()

	cmds, err := readInput(os.Stdin)
	if err != nil {
		log.Fatal(err.Error())
	}

	s := screen.NewScreen(*screenWidth, *screenHeight)
	for _, cmd := range cmds {
		cmd.Execute(s)
		s.Print()
	}

	fmt.Printf("num on = %d\n", s.CountOn())
}
