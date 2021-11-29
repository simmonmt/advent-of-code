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

package code

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	repeatPattern = regexp.MustCompile(`([0-9]+)x([0-9]+)\)`)
)

func Decode(in string) (string, error) {
	chars := []byte(in)
	out := []byte{}
	for i := 0; i < len(chars); {
		c := chars[i]

		if c == '(' {
			if i+1 == len(chars) {
				return "", fmt.Errorf("start of repeat pattern at end")
			}
			matches := repeatPattern.FindSubmatch(chars[i+1:])
			if matches == nil {
				return "", fmt.Errorf("invalid repeat pattern at \"%v\"", chars[i+1:])
			}

			numCharsStr := string(matches[1])
			numChars, err := strconv.ParseUint(numCharsStr, 10, 32)
			if err != nil {
				return "", fmt.Errorf("failed to parse num chars in \"%v\"", string(numCharsStr))
			}

			numRepsStr := string(matches[2])
			numReps, err := strconv.ParseUint(numRepsStr, 10, 32)
			if err != nil {
				return "", fmt.Errorf("failed to parse num reps in \"%v\"", string(numRepsStr))
			}

			nextIdx := i + len(matches[0]) + 1
			nextEnd := nextIdx + int(numChars)
			if nextEnd > len(chars) {
				return "", fmt.Errorf("repeat %vx%v reaches past end (next %v, end %v)",
					numChars, numReps, nextIdx, nextEnd)
			}

			for j := 0; j < int(numReps); j++ {
				out = append(out, chars[nextIdx:nextEnd]...)
			}
			i = int(nextEnd)
		} else {
			out = append(out, c)
			i++
		}
	}
	return string(out), nil
}
