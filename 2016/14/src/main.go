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
	"math"
	"sort"

	"logger"
	"pad"
)

var (
	salt      = flag.String("salt", "", "salt")
	stretched = flag.Bool("stretched", false, "use stretched hasher")
	verbose   = flag.Bool("verbose", false, "verbose")
)

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *salt == "" {
		log.Fatal("--salt is required")
	}

	var hasher pad.Hasher
	if *stretched {
		hasher = &pad.StretchedHasher{}
	} else {
		hasher = &pad.NormalHasher{}
	}

	hashQueue := pad.NewQueue()

	finishIdx := math.MaxInt32
	keyIndexes := []int{}
	for i := 0; i < finishIdx; i++ {
		if i != 0 && i%1000 == 0 {
			fmt.Printf("%d: #keys: %d, finish %v\n", i, len(keyIndexes), finishIdx)
		}

		h := hasher.MakeHash(*salt, i)
		logger.LogF("%d: hash %v\n", i, h)

		if reps := pad.HasRepeats(h, 3); len(reps) > 0 {
			// We only consider the first one
			logger.LogF("%d: adding 3-rep to hash queue: %x, exp %v\n", i, reps[0], i+1000)
			hashQueue.Add(i, reps[0], i+1000)
		}

		// Off in the future, we want to verify whether the 3-reps have
		// corresponding 5-reps. We use ActiveBefore because we want to
		// exclude any added this iteration.
		if reps := pad.HasRepeats(h[:], 5); len(reps) > 0 {
			logger.LogF("%d: found 5-reps %v\n", i, reps)
			activeElems := hashQueue.ActiveBefore(i)
			logger.LogF("%d: found 5-rep %x active %v\n", i, h, activeElems)
			for _, activeElem := range activeElems {
				for _, rep := range reps {
					if activeElem.Repeater == rep {
						keyIndexes = append(keyIndexes, activeElem.Index)
						logger.LogF("adding key index %v for %v, now %d keys\n",
							activeElem.Index, rep, len(keyIndexes))
						hashQueue.Delete(activeElem)
					}
				}
			}
		}

		if finishIdx == math.MaxInt32 && len(keyIndexes) >= 64 {
			finishIdx = i + 1000
			fmt.Printf("resetting finishidx to %v\n", finishIdx)
		}

		hashQueue.ExpireTo(i)
	}

	sort.Ints(keyIndexes)

	for i, keyIndex := range keyIndexes {
		fmt.Printf("%3d: %v\n", i+1, keyIndex)
	}
}
