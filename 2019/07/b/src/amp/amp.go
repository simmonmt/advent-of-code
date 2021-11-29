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
