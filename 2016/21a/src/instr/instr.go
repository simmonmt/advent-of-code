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

package instr

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

var (
	swapPosPattern     = regexp.MustCompile(`^swap position ([0-9]+) with position ([0-9]+)$`)
	swapCharPattern    = regexp.MustCompile(`^swap letter ([a-z]) with letter ([a-z])$`)
	rotatePattern      = regexp.MustCompile(`^rotate (left|right) ([0-9]+) steps?$`)
	rotateMagicPattern = regexp.MustCompile(`^rotate based on position of letter ([a-z])$`)
	reversePattern     = regexp.MustCompile(`^reverse positions ([0-9]+) through ([0-9]+)$`)
	movePattern        = regexp.MustCompile(`^move position ([0-9]+) to position ([0-9]+)$`)
)

type Instr interface {
	Exec(str []byte) bool
	String() string
}

type swap struct {
	pos1, pos2 int
	ch1, ch2   rune
}

func newSwapPos(pos1, pos2 int) Instr {
	return &swap{pos1: pos1, pos2: pos2}
}

func newSwapChar(ch1, ch2 rune) Instr {
	return &swap{pos1: -1, pos2: -1, ch1: ch1, ch2: ch2}
}

func (inst *swap) Exec(str []byte) bool {
	var pos1, pos2 int
	if inst.pos1 >= 0 {
		pos1, pos2 = inst.pos1, inst.pos2
	} else {
		pos1 = bytes.IndexRune(str, inst.ch1)
		pos2 = bytes.IndexRune(str, inst.ch2)
		if pos1 == -1 || pos2 == -1 {
			return false
		}
	}

	str[pos1], str[pos2] = str[pos2], str[pos1]
	return true
}

func (inst *swap) String() string {
	if inst.pos1 < 0 {
		return fmt.Sprintf("swap %c with %c", inst.ch1, inst.ch2)
	} else {
		return fmt.Sprintf("swap [%v] with [%v]", inst.pos1, inst.pos2)
	}
}

type rotate struct {
	left bool
	num  int
}

func newRotate(left bool, num int) Instr {
	return &rotate{left: left, num: num}
}

func (inst *rotate) String() string {
	dir := "right"
	if inst.left {
		dir = "left"
	}

	return fmt.Sprintf("rotate %v %d", dir, inst.num)
}

func (inst *rotate) Exec(str []byte) bool {
	out := make([]byte, len(str))

	for i := range str {
		var newPos int
		if inst.left {
			// I'm sure there's a prettier way to do this
			newPos = (len(str)*2 + (i - (inst.num % len(str)))) % len(str)
		} else {
			newPos = (i + inst.num) % len(str)
		}
		out[newPos] = str[i]
	}

	copy(str, out)
	return true
}

type rotateMagic struct {
	ch rune
}

func newRotateMagic(ch rune) Instr {
	return &rotateMagic{ch: ch}
}

func (inst *rotateMagic) String() string {
	return fmt.Sprintf("rotate magic %c", inst.ch)
}

func (inst *rotateMagic) Exec(str []byte) bool {
	pos := bytes.IndexRune(str, inst.ch)
	if pos == -1 {
		return false
	}

	numRots := 1 + pos
	if pos >= 4 {
		numRots++
	}

	out := make([]byte, len(str))

	for i := range str {
		newPos := (i + numRots) % len(str)
		out[newPos] = str[i]
	}

	copy(str, out)
	return true
}

type reverse struct {
	pos1, pos2 int
}

func newReverse(pos1, pos2 int) Instr {
	return &reverse{pos1: pos1, pos2: pos2}
}

func (inst *reverse) String() string {
	return fmt.Sprintf("reverse %v %v", inst.pos1, inst.pos2)
}

func (inst *reverse) Exec(str []byte) bool {
	out := make([]byte, len(str))
	for i := range str {
		if i < inst.pos1 || i > inst.pos2 {
			out[i] = str[i]
		} else {
			out[inst.pos2-(i-inst.pos1)] = str[i]
		}
	}
	copy(str, out)
	return true
}

type move struct {
	pos1, pos2 int
}

func newMove(pos1, pos2 int) Instr {
	return &move{pos1: pos1, pos2: pos2}
}

func (inst *move) String() string {
	return fmt.Sprintf("move %v %v", inst.pos1, inst.pos2)
}

func (inst *move) Exec(str []byte) bool {
	tmp := make([]byte, len(str))
	for i := range str {
		if i < inst.pos1 {
			tmp[i] = str[i]
		} else if i > inst.pos1 {
			tmp[i-1] = str[i]
		}
	}

	toInsert := str[inst.pos1]
	for i := range str {
		if i < inst.pos2 {
			str[i] = tmp[i]
		} else if i == inst.pos2 {
			str[i] = toInsert
		} else {
			str[i] = tmp[i-1]
		}
	}
	return true
}

func parsePositions(str1, str2 string) (int, int, error) {
	pos1, err := strconv.ParseUint(str1, 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("bad 1st position: %v", err)
	}
	pos2, err := strconv.ParseUint(str2, 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("bad 2nd position: %v", err)
	}
	return int(pos1), int(pos2), nil
}

func Parse(str string) (Instr, error) {
	if matches := swapPosPattern.FindStringSubmatch(str); matches != nil {
		pos1, pos2, err := parsePositions(matches[1], matches[2])
		if err != nil {
			return nil, err
		}
		return newSwapPos(int(pos1), int(pos2)), nil

	} else if matches := swapCharPattern.FindStringSubmatch(str); matches != nil {
		return newSwapChar(rune(matches[1][0]), rune(matches[2][0])), nil

	} else if matches := rotatePattern.FindStringSubmatch(str); matches != nil {
		left := matches[1] == "left"
		num, err := strconv.ParseUint(matches[2], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("bad num: %v", err)
		}

		return newRotate(left, int(num)), nil

	} else if matches := rotateMagicPattern.FindStringSubmatch(str); matches != nil {
		return newRotateMagic(rune(matches[1][0])), nil

	} else if matches := reversePattern.FindStringSubmatch(str); matches != nil {
		pos1, pos2, err := parsePositions(matches[1], matches[2])
		if err != nil {
			return nil, err
		}
		return newReverse(int(pos1), int(pos2)), nil

	} else if matches := movePattern.FindStringSubmatch(str); matches != nil {
		pos1, pos2, err := parsePositions(matches[1], matches[2])
		if err != nil {
			return nil, err
		}
		return newMove(int(pos1), int(pos2)), nil

	} else {
		return nil, fmt.Errorf("unknown instruction")
	}

	panic("unreachable")
}
