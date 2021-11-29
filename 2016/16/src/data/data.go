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

package data

func Grow(a []bool) []bool {
	out := make([]bool, len(a)*2+1)
	copy(out, a)
	for i := range a {
		outPos := len(a)*2 - i
		outVal := !a[i]
		out[outPos] = outVal
	}
	return out
}

func checksumRound(in []bool) []bool {
	out := make([]bool, len(in)/2)
	for i := range out {
		out[i] = in[i*2] == in[i*2+1]
	}
	return out
}

func Checksum(in []bool) []bool {
	if len(in)%2 == 1 {
		panic("odd checksum input")
	}

	for {
		sum := checksumRound(in)
		if len(sum)%2 == 1 {
			return sum
		}
		in = sum
	}
}
