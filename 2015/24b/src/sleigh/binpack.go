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

package sleigh

func OneBinPack(values []int, cap int) []int {
	// The easy way, with lots of allocations

	//fmt.Printf("BinPack(%v, %v)\n", values, cap)

	if cap == 0 {
		//fmt.Printf("  easy out; cap==0")
		return []int{}
	} else if len(values) == 0 {
		//fmt.Printf("  easy fail; len==0, cap>0 %v\n", cap)
		return nil
	}

	for i := 0; i < len(values); i++ {
		cand := values[i]
		//fmt.Printf("  trying cand %v\n", cand)
		if cand > cap {
			//fmt.Printf("    doesn't fit\n")
			continue
		}

		rest := values[i+1 : len(values)]
		if found := OneBinPack(rest, cap-cand); found != nil {
			out := make([]int, len(found)+1)
			out[0] = cand
			for j, val := range found {
				out[j+1] = val
			}
			//fmt.Printf("  worked; returning %v\n", out)
			return out
		}
	}

	//fmt.Printf("  end of list; return nil\n")
	return nil
}

func AllBinPacks(values []int, cap int) [][]int {
	// The easy way, with lots of allocations

	//fmt.Printf("BinPack(%v, %v)\n", values, cap)

	out := [][]int{}

	if cap == 0 {
		//fmt.Printf("  easy out; cap==0")
		out = append(out, []int{})
		return out
	} else if len(values) == 0 {
		//fmt.Printf("  easy fail; len==0, cap>0 %v\n", cap)
		return out
	}

	for i := 0; i < len(values); i++ {
		cand := values[i]
		//fmt.Printf("  trying cand %v\n", cand)
		if cand > cap {
			//fmt.Printf("    doesn't fit\n")
			continue
		}

		rest := values[i+1 : len(values)]
		foundPacks := AllBinPacks(rest, cap-cand)
		for _, foundPack := range foundPacks {
			withCand := make([]int, len(foundPack)+1)
			withCand[0] = cand
			for j, val := range foundPack {
				withCand[j+1] = val
			}
			out = append(out, withCand)
		}
	}

	//fmt.Printf("  end of list; return nil\n")
	return out
}
