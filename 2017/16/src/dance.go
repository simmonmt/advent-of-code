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

type Sequence struct {
	posns []byte
}

func NewSequence() *Sequence {
	posns := make([]byte, 16)
	for i := range posns {
		posns[i] = byte(int('a') + i)
	}

	return &Sequence{posns}
}

func (s *Sequence) Spin(num int) {
	arr := s.posns[len(s.posns)-num:]
	arr = append(arr, s.posns[0:len(s.posns)-num]...)
	s.posns = arr
}

func (s *Sequence) Exchange(a, b int) {
	s.posns[a], s.posns[b] = s.posns[b], s.posns[a]
}

func (s *Sequence) Partner(a, b byte) {
	var aPos, bPos int
	for i, s := range s.posns {
		if s == a {
			aPos = i
		}
		if s == b {
			bPos = i
		}
	}

	s.posns[aPos], s.posns[bPos] = s.posns[bPos], s.posns[aPos]
}

func (s *Sequence) Dump(out io.Writer) {
	fmt.Fprintf(out, "%v\n", string(s.posns))
}

func readSteps(in io.Reader) ([]string, error) {
	reader := bufio.NewReader(in)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)
	return strings.Split(line, ","), nil
}

const (
	CMDSpin = iota
	CMDExchange
	CMDPartner
)

type Command struct {
	Cmd  int
	A, B int
}

func main() {
	steps, err := readSteps(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	cmds := []Command{}
	for _, step := range steps {
		var cmd Command

		switch step[0] {
		case 's':
			num, _ := strconv.Atoi(step[1:])
			cmd = Command{CMDSpin, num, -1}
			break
		case 'x':
			parts := strings.Split(step[1:], "/")
			a, _ := strconv.Atoi(parts[0])
			b, _ := strconv.Atoi(parts[1])

			cmd = Command{CMDExchange, a, b}
			break
		case 'p':
			parts := strings.Split(step[1:], "/")
			cmd = Command{CMDPartner, int(parts[0][0]),
				int(parts[1][0])}
			break
		}

		cmds = append(cmds, cmd)
	}

	seen := map[string]int{}

	seq := NewSequence()
	for i := 0; i < 100; i++ {
		if i%100000 == 0 {
			fmt.Println(i)
		}
		for _, cmd := range cmds {
			switch cmd.Cmd {
			case CMDSpin:
				seq.Spin(cmd.A)
				break
			case CMDExchange:
				seq.Exchange(cmd.A, cmd.B)
				break
			case CMDPartner:
				seq.Partner(byte(cmd.A), byte(cmd.B))
				break
			}

		}

		if was, found := seen[string(seq.posns)]; found {
			fmt.Printf("found %d was %d\n", i, was)
		} else {
			seen[string(seq.posns)] = i
		}
	}

	// seq.Dump(os.Stdout)
	// seq.Partner("c", "h")
	seq.Dump(os.Stdout)
}
