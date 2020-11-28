package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/simmonmt/aoc/2020/common/logger"
	"github.com/simmonmt/aoc/2020/common/strutil"
	"github.com/simmonmt/aoc/2020/common/vm"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	ramPath = flag.String("ram", "", "path to file containing ram values")
	input   = flag.String("input", "", "input values")
)

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *ramPath == "" {
		log.Fatalf("--ram is required")
	}

	if *input == "" {
		log.Fatalf("--input is required")
	}

	ram, err := vm.NewRamFromFile(*ramPath)
	if err != nil {
		log.Fatal(err)
	}

	inputValues, err := strutil.StringToInt64s(*input)
	if err != nil {
		log.Fatal(err)
	}

	async := vm.RunAsync("vm", ram)

	go func() {
		for _, v := range inputValues {
			async.In <- &vm.ChanIOMessage{Val: v}
		}
		logger.LogLn("input done")
	}()

	for {
		msg, ok := <-async.Out
		if !ok {
			break
		}
		fmt.Printf("got %+v\n", msg)
	}
	fmt.Print("program terminated")
}
