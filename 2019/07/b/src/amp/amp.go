package amp

import (
	"fmt"

	vm "github.com/simmonmt/aoc/2019/07/b/src/vm"
)

type Amp struct {
	In    chan *vm.ChanIOMessage
	Out   chan *vm.ChanIOMessage
	phase int
	ram   vm.Ram
}

func Start(phase int, ram vm.Ram) *Amp {
	amp := &Amp{
		In:    make(chan *vm.ChanIOMessage, 2),
		Out:   make(chan *vm.ChanIOMessage, 2),
		phase: phase,
		ram:   ram,
	}

	amp.In <- &vm.ChanIOMessage{Val: phase}

	go run(amp)

	return amp
}

func run(amp *Amp) {
	io := vm.NewChanIO(amp.In, amp.Out)
	if err := vm.Run(amp.ram, io, 0); err != nil {
		amp.Out <- &vm.ChanIOMessage{
			Err: fmt.Errorf("phase %d amp failed: %v", amp.phase, err),
		}
	}

	close(amp.Out)
	fmt.Sprintf("phase %d amp terminating\n", amp.phase)
}
