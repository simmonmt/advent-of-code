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
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

func parseMolecule(in string) []string {
	out := []string{}
	for _, c := range in {
		if unicode.IsUpper(c) {
			out = append(out, string(c))
		} else {
			out[len(out)-1] += string(c)
		}
	}
	return out
}

func readInput(r io.Reader) (map[string][][]string, []string, error) {
	reader := bufio.NewReader(r)

	repls := map[string][][]string{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		pairs := strings.SplitN(line, " => ", 2)
		from, to := pairs[0], pairs[1]

		if _, found := repls[from]; !found {
			repls[from] = [][]string{}
		}
		repls[from] = append(repls[from], parseMolecule(to))
	}

	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, nil, fmt.Errorf("no input molecule")
	}
	line = strings.TrimSpace(line)

	molecule := parseMolecule(line)

	return repls, molecule, nil
}

func doReplacement(in []string, from string, to []string) [][]string {
	ret := [][]string{}

	for i, a := range in {
		if a == from {
			// fmt.Printf("matched %v to %v at %v\n", from, to, i)
			out := make([]string, len(in)+len(to)-1)
			off := 0
			if i > 0 {
				copy(out[off:i], in[0:i])
				off += i
				// fmt.Printf("first %v\n", out)
			}
			copy(out[off:off+len(to)], to)
			off += len(to)
			// fmt.Printf("mid %v\n", out)

			if i < len(in) {
				copy(out[off:], in[i+1:])
				// fmt.Printf("end %v\n", out)
			}

			// fmt.Printf("out %v\n", out)
			ret = append(ret, out)
		}
	}

	return ret
}

func main() {
	repls, molecule, err := readInput(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	fmt.Println(repls)
	fmt.Println(molecule)

	allResults := map[string]bool{}
	for from, tos := range repls {
		for _, to := range tos {
			// fmt.Printf("%v %v\n", from, to)
			results := doReplacement(molecule, from, to)
			for _, result := range results {
				key := strings.Join(result, " ")
				// fmt.Println(key)
				allResults[key] = true
			}
		}
	}

	// for r := range allResults {
	// 	fmt.Println(r)
	// }

	fmt.Printf("num distinct = %d\n", len(allResults))

	//fmt.Println(doReplacement(molecule, "H", []string{"H", "O"}))
}
