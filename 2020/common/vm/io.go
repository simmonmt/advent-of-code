package vm

import (
	"fmt"

	"github.com/simmonmt/aoc/2020/common/logger"
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
