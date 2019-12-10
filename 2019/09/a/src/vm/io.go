package vm

import (
	"fmt"

	"github.com/simmonmt/aoc/2019/common/logger"
)

type IO interface {
	Read() int
	Write(int)
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

type ChanIOMessage struct {
	Val int
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

func (io *chanIO) Read() int {
	msg, ok := <-io.in
	if !ok {
		panic("sender closed unexpectedly")
	}

	return msg.Val
}

func (io *chanIO) Write(val int) {
	io.out <- &ChanIOMessage{Val: val}
}
