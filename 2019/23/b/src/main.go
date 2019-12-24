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

var (
	lastInput, lastOutput, lastNat time.Time
	lastLock                       sync.Mutex

	natThresh = time.Duration(100) * time.Millisecond
)

type nicCmd struct {
	x, y int64
}

type natState struct {
	nicMap []*nicState

	cmd     *nicCmd
	cmdLock sync.Mutex
}

func NewNat(nicMap []*nicState) *natState {
	return &natState{
		nicMap: nicMap,
	}
}

func (s *natState) Recv(cmd *nicCmd) {
	s.cmdLock.Lock()
	s.cmd = cmd
	s.cmdLock.Unlock()
}

func (s *natState) Start() {
	go func() {
		for {
			var toSend *nicCmd
			now := time.Now()

			//fmt.Printf("nat looking\n")

			lastLock.Lock()
			if now.Sub(lastInput) > natThresh && now.Sub(lastOutput) > natThresh && now.Sub(lastNat) > natThresh {
				s.cmdLock.Lock()
				toSend = s.cmd
				s.cmd = nil
				s.cmdLock.Unlock()

				lastNat = now
			}
			lastLock.Unlock()

			if toSend != nil {
				fmt.Printf("nat sending %+v\n", toSend)
				s.nicMap[0].Recv(toSend)
			}

			//fmt.Printf("nat sleeping\n")

			time.Sleep(time.Duration(50) * time.Millisecond)
			//fmt.Printf("nat awake\n")
		}

		//fmt.Printf("nat done")
	}()
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

	lastLock.Lock()
	lastInput = time.Now()
	lastLock.Unlock()

	v := io.queue[0]
	io.queue = io.queue[1:]
	//fmt.Printf("vm %d io read => %v\n", io.addr, v)
	return v
}

func (io *nicIO) Write(val int64) {
	lastLock.Lock()
	lastOutput = time.Now()
	lastLock.Unlock()

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
	nat    *natState
	io     *nicIO
}

func NewNic(addr int, ram vm.Ram, nat *natState, nicMap []*nicState) *nicState {
	s := &nicState{
		addr:   addr,
		ram:    ram,
		nicMap: nicMap,
		nat:    nat,
	}

	s.io = NewNicIO(addr, func(addr int, cmd *nicCmd) { s.dispatch(addr, cmd) })

	return s
}

func (s *nicState) dispatch(dest int, cmd *nicCmd) {
	//fmt.Printf("vm %d => %d cmd %v\n", s.addr, dest, *cmd)

	if dest == 255 {
		s.nat.Recv(cmd)
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

	lastLock.Lock()
	lastInput = time.Now()
	lastOutput = time.Now()
	lastNat = time.Now()
	lastLock.Unlock()

	nicMap := make([]*nicState, 50)

	nat := NewNat(nicMap)

	for i := range nicMap {
		nicMap[i] = NewNic(i, ram.Clone(), nat, nicMap)
	}

	nat.Start()
	for _, nic := range nicMap {
		nic.Start()
	}

	for {
		time.Sleep(time.Duration(1) * time.Second)
	}
}
