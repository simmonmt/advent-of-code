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
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/simmonmt/aoc/2019/common/logger"
	"github.com/simmonmt/aoc/2019/common/vm"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
	ramPath = flag.String("ram", "", "path to file containing ram values")
)

type nicCmd struct {
	x, y int64
}

type nicIO struct {
	addr int

	queue     []int64
	queueLock sync.Mutex

	valStash   [2]int64
	numStashed int

	dispatch func(dest int, cmd *nicCmd)
}

func NewNicIO(addr int, dispatch func(dest int, cmd *nicCmd)) *nicIO {
	return &nicIO{
		addr:       addr,
		queue:      []int64{int64(addr)},
		dispatch:   dispatch,
		numStashed: 0,
	}
}

func (io *nicIO) Recv(cmd *nicCmd) {
	io.queueLock.Lock()
	io.queue = append(io.queue, cmd.x, cmd.y)
	io.queueLock.Unlock()
}

func (io *nicIO) Read() int64 {
	io.queueLock.Lock()
	defer io.queueLock.Unlock()

	if len(io.queue) == 0 {
		// Block for 1ms so the other goroutines get a chance to run. If
		// we don't sleep, the Go runtime will never preempt us on its
		// own.
		time.Sleep(time.Duration(1) * time.Millisecond)
		return -1
	}

	v := io.queue[0]
	io.queue = io.queue[1:]
	fmt.Printf("vm %d io read => %v\n", io.addr, v)
	return v
}

func (io *nicIO) Write(val int64) {
	if io.numStashed < 2 {
		io.valStash[io.numStashed] = val
		io.numStashed++
		return
	}

	dest := int(io.valStash[0])
	cmd := &nicCmd{x: io.valStash[1], y: val}
	io.numStashed = 0
	io.dispatch(dest, cmd)
}

type nicState struct {
	addr   int
	ram    vm.Ram
	nicMap []*nicState
	io     *nicIO
}

func NewNic(addr int, ram vm.Ram, nicMap []*nicState) *nicState {
	s := &nicState{
		addr:   addr,
		ram:    ram,
		nicMap: nicMap,
	}

	s.io = NewNicIO(addr, func(addr int, cmd *nicCmd) { s.dispatch(addr, cmd) })

	return s
}

func (s *nicState) dispatch(dest int, cmd *nicCmd) {
	fmt.Printf("vm %d => %d cmd %v\n", s.addr, dest, *cmd)

	if dest == 255 {
		fmt.Printf("sending to 255: %+v", cmd)
		return
	}

	nic := s.nicMap[dest]
	nic.Recv(cmd)
}

func (s *nicState) Start() {
	go func() {
		if err := vm.Run(s.ram, s.io); err != nil {
			panic(fmt.Sprintf("nic %d err: %v", s.addr, err))
		}
		fmt.Printf("nic %d terminating\n", s.addr)
	}()
}

func (s *nicState) Recv(cmd *nicCmd) {
	s.io.Recv(cmd)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *ramPath == "" {
		log.Fatalf("--ram is required")
	}

	ram, err := vm.NewRamFromFile(*ramPath)
	if err != nil {
		log.Fatal(err)
	}

	nicMap := make([]*nicState, 50)

	for i := range nicMap {
		nicMap[i] = NewNic(i, ram.Clone(), nicMap)
	}

	for _, nic := range nicMap {
		nic.Start()
	}

	for {
	}
}
