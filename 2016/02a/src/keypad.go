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

type Keypad struct {
	digits     [][]int
	x, y       int
	dimX, dimY int
}

func NewKeypad() *Keypad {
	return &Keypad{
		dimX: 3,
		dimY: 3,
		digits: [][]int{
			[]int{1, 2, 3},
			[]int{4, 5, 6},
			[]int{7, 8, 9},
		},
		x: 1,
		y: 1,
	}
}

func (k *Keypad) Up() {
	if k.y > 0 {
		k.y--
	}
}

func (k *Keypad) Down() {
	if k.y+1 < k.dimY {
		k.y++
	}
}

func (k *Keypad) Left() {
	if k.x > 0 {
		k.x--
	}
}

func (k *Keypad) Right() {
	if k.x+1 < k.dimX {
		k.x++
	}
}

func (k *Keypad) Cur() int {
	return k.digits[k.y][k.x]
}

func main() {
	keypad := NewKeypad()

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)

		for _, cmd := range []byte(line) {
			switch cmd {
			case 'U':
				keypad.Up()
				break
			case 'D':
				keypad.Down()
				break
			case 'L':
				keypad.Left()
				break
			case 'R':
				keypad.Right()
				break
			default:
				log.Fatalf("unknown cmd %v", string(cmd))
			}
		}

		fmt.Println(keypad.Cur())
	}
}
