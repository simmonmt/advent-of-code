package vm

import (
	"fmt"

	"github.com/simmonmt/aoc/2019/common/logger"
)

type IO interface {
	Read() int
	Write(int)
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

func (io *ioImpl) String() string {
	return fmt.Sprintf("in=[%v]", io.input[io.inputAddr:])
}
