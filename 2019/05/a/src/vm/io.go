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

package vm

import (
	"fmt"

	"github.com/simmonmt/aoc/2019/common/logger"
)

type IO interface {
	Read() int
	Write(int)
	Written() []int
	String() string
}

type ioImpl struct {
	input     []int
	inputAddr int
	output    []int
}

func NewIO(input ...int) IO {
	return &ioImpl{
		input:     input,
		inputAddr: 0,
		output:    []int{},
	}
}

func (io *ioImpl) Read() int {
	if io.inputAddr >= len(io.input) {
		panic("out of input")
	}

	in := io.input[io.inputAddr]
	io.inputAddr++
	return in
}

func (io *ioImpl) Write(val int) {
	io.output = append(io.output, val)
	logger.LogF("output: %v", val)
}

func (io *ioImpl) Written() []int {
	return io.output
}

func (io *ioImpl) String() string {
	return fmt.Sprintf("in=[%v]", io.input[io.inputAddr:])
}
