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
	Read() int64
	Write(int64)
}

type SaverIO struct {
	input     []int64
	inputAddr int
	output    []int64
}

func NewSaverIO(input ...int64) *SaverIO {
	return &SaverIO{
		input:     input,
		inputAddr: 0,
		output:    []int64{},
	}
}

func (io *SaverIO) Read() int64 {
	if io.inputAddr >= len(io.input) {
		panic("out of input")
	}

	in := io.input[io.inputAddr]
	io.inputAddr++
	return in
}

func (io *SaverIO) Write(val int64) {
	io.output = append(io.output, val)
	logger.LogF("output: %v", val)
}

func (io *SaverIO) Written() []int64 {
	return io.output
}

func (io *SaverIO) String() string {
	return fmt.Sprintf("in=[%v]", io.input[io.inputAddr:])
}

type ChanIOMessage struct {
	Val int64
	Err error
}

type chanIO struct {
	in, out chan *ChanIOMessage
}

func NewChanIO(in, out chan *ChanIOMessage) IO {
	return &chanIO{
		in:  in,
		out: out,
	}
}

func (io *chanIO) Read() int64 {
	msg, ok := <-io.in
	if !ok {
		panic("sender closed unexpectedly")
	}

	return msg.Val
}

func (io *chanIO) Write(val int64) {
	io.out <- &ChanIOMessage{Val: val}
}
