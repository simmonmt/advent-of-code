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
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"instr"
	"logger"
)

var (
	unscrambled = flag.String("unscrambled", "", "unscrambled input")
	verbose     = flag.Bool("verbose", false, "verbose")
)

func readInput(r io.Reader) ([]instr.Instr, error) {
	insts := []instr.Instr{}

	reader := bufio.NewReader(r)
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		inst, err := instr.Parse(line)
		if err != nil {
			return nil, fmt.Errorf("%d: %v", lineNum, err)
		}

		insts = append(insts, inst)
	}

	return insts, nil
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *unscrambled == "" {
		log.Fatal("--unscrambled is required")
	}

	insts, err := readInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	state := make([]byte, len(*unscrambled))
	copy(state, []byte(*unscrambled))
	for i := len(insts) - 1; i >= 0; i-- {
		//for i, inst := range insts {
		inst := insts[i]
		if !inst.Exec(state) {
			log.Fatalf("%d: exec failed; state %v", i, string(state))
		}

		logger.LogF("ran \"%v\", state now %v\n", inst.String(), string(state))
	}

	fmt.Println(string(state))
}
