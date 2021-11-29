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
	"log"
	"os"
	"sort"
	"strings"

	"logger"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	numWorkers = flag.Int("num_workers", 2, "num workers")
	stepBase   = flag.Int("step_base", 60, "step base")
)

type Worker struct {
	Cur  string
	Left int
}

func readPrereqs() (map[string][]string, error) {
	prereqs := map[string][]string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		prereq := parts[1]
		name := parts[7]

		if _, ok := prereqs[name]; !ok {
			prereqs[name] = []string{}
		}
		prereqs[name] = append(prereqs[name], prereq)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return prereqs, nil
}

func stepLen(name string) int {
	return *stepBase + int([]byte(name)[0]-'A'+1)
}

func availWorkers(workers []Worker) []int {
	avail := []int{}

	for i, w := range workers {
		if w.Cur == "" {
			avail = append(avail, i)
		}
	}

	return avail
}

func assignWorker(worker *Worker, work string) {
	if worker.Cur != "" {
		panic("already assigned")
	}
	worker.Cur = work
	worker.Left = stepLen(work)
}

func dump(t int, workers []Worker) {
	fmt.Printf("%5d ", t)
	for _, w := range workers {
		fmt.Printf("%5s(%3d) ", w.Cur, w.Left)
	}
	fmt.Println()
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	prereqs, err := readPrereqs()
	if err != nil {
		log.Fatal(err)
	}

	todo := map[string]bool{}
	for name, ps := range prereqs {
		todo[name] = true
		for _, p := range ps {
			todo[p] = true
		}
	}

	order := ""

	t := 0
	workers := make([]Worker, *numWorkers)

	working := map[string]bool{}
	done := map[string]bool{}
	for len(todo) > 0 || len(working) > 0 {
		avail := availWorkers(workers)
		if len(avail) > 0 {
			// find work to do for available workers
			cands := []string{}
			for step, _ := range todo {
				ps := prereqs[step]
				if ps == nil {
					ps = []string{}
				}

				alldone := true
				for _, p := range ps {
					if _, ok := done[p]; !ok {
						alldone = false
						break
					}
				}
				if !alldone {
					continue
				}

				cands = append(cands, step)
			}

			// assign work (if any) to available workers
			if len(cands) > 0 {
				sort.Strings(cands)

				for _, cand := range cands {
					if len(avail) == 0 {
						break
					}

					availIdx := avail[0]
					avail = avail[1:]

					assignWorker(&workers[availIdx], cand)
					working[cand] = true
					delete(todo, cand)
				}
			}
		}

		dump(t, workers)

		// advance time and consume work
		t++
		for i := range workers {
			w := &workers[i]
			if w.Cur == "" {
				continue
			}

			w.Left--
			if w.Left == 0 {
				delete(working, w.Cur)
				done[w.Cur] = true
				order += w.Cur
				w.Cur = ""
			}
		}
	}

	fmt.Println(t)
	fmt.Println(order)

}
