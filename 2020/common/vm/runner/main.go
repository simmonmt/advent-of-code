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
