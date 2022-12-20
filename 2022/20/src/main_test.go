// Copyright 2022 Google LLC
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
	_ "embed"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	//go:embed sample.txt
	rawSample   string
	sampleLines []string
)

func TestParseInput(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	want := []int{1, 2, -3, 3, -2, 0, 4}
	if !reflect.DeepEqual(input, want) {
		t.Errorf("parseInput(sampleLines) = %v, want %v",
			input, want)
	}
}

func TestTreeAsList(t *testing.T) {
	head := &TreeNode{
		size: 4,
		before: &TreeNode{
			size: 2,
			before: &TreeNode{
				size: 1,
				val:  9,
			},
			val: 8,
		},
		val: 7,
		after: &TreeNode{
			size: 1,
			val:  6,
		},
	}

	head.before.before.up = head.before
	head.before.up = head
	head.after.up = head

	want := []int{9, 8, 7, 6}
	if got := head.AsList(); !reflect.DeepEqual(got, want) {
		t.Errorf("AsList() = %v, want %v", got, want)
	}
}

func TestMakeTree(t *testing.T) {
	nums, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	ptrs := make([]*TreeNode, len(nums))
	head := makeTree(nums, ptrs)

	want := []int{1, 2, -3, 3, -2, 0, 4}
	if got := head.AsList(); !reflect.DeepEqual(got, want) {
		t.Errorf("result = %v, want %v", got, want)
	}

	got := make([]int, len(nums))
	for i, ptr := range ptrs {
		got[i] = ptr.val
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %v, want %v", got, want)
	}
}

func TestForward(t *testing.T) {
	nums, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	ptrs := make([]*TreeNode, len(nums))
	makeTree(nums, ptrs)

	for i := 0; i < len(ptrs)-1; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got, want := ptrs[i].Forward().val, ptrs[i+1].val; got != want {
				t.Errorf("forward(idx %v val %v) => %v, want %v",
					i, ptrs[i].val, got, want)
			}
		})
	}
}

func TestIndex(t *testing.T) {
	nums, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	ptrs := make([]*TreeNode, len(nums))
	head := makeTree(nums, ptrs)

	for i := 0; i < len(ptrs)-1; i++ {
		t.Run("Index"+strconv.Itoa(i), func(t *testing.T) {
			if got := ptrs[i].Index(); got != i {
				t.Errorf("index(idx %v val %v) => %v, want %v",
					i, ptrs[i].val, got, i)
			}
		})

		t.Run("FindIndex"+strconv.Itoa(i), func(t *testing.T) {
			if got, want := head.FindIndex(i).val, ptrs[i].val; got != want {
				t.Errorf("FindIndex(%d).val = %v, want %v",
					i, got, want)
			}
		})
	}
}

func TestRemoveSample(t *testing.T) {
	nums, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < len(nums); i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ptrs := make([]*TreeNode, len(nums))
			head := makeTree(nums, ptrs)

			wantElems := []int{}
			if i > 0 {
				wantElems = append(wantElems, nums[0:i]...)
			}
			if i < len(nums)-1 {
				wantElems = append(wantElems, nums[i+1:]...)
			}

			wantNode := &TreeNode{size: 1, val: ptrs[i].val}

			ptrs[i].Remove(&head)

			if got := head.AsList(); !reflect.DeepEqual(got, wantElems) {
				t.Errorf("got %v, want %v", got, wantElems)
			}

			if got := ptrs[i]; !reflect.DeepEqual(got, wantNode) {
				t.Errorf("got %+v, want %+v", got, wantNode)
			}

		})
	}
}

func TestRemove(t *testing.T) {
	type TestCase struct {
		removals []int
		want     []int
	}

	testCases := []TestCase{
		TestCase{[]int{0, 1, 2, 3}, []int{4, 5, 6}},
	}

	nums := []int{0, 1, 2, 3, 4, 5, 6}
	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ptrs := make([]*TreeNode, len(nums))
			head := makeTree(nums, ptrs)

			for _, removal := range tc.removals {
				ptrs[removal].Remove(&head)
				head.Check()

				wantNode := &TreeNode{size: 1, val: ptrs[removal].val}
				if got := ptrs[removal]; !reflect.DeepEqual(got, wantNode) {
					t.Errorf("removal %v: got %+v, want %+v",
						removal, got, wantNode)
				}
			}

			if got := head.AsList(); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	nums, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < len(nums)-1; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ptrs := make([]*TreeNode, len(nums))
			head := makeTree(nums, ptrs)

			ptrs[0].Remove(&head)
			head.Check()
			victim := ptrs[0]

			afterIdx := i
			after := head.FindIndex(afterIdx)

			victim.InsertAfter(after)
			head.Check()

			want := []int{}
			for j := 1; j < len(nums); j++ {
				want = append(want, nums[j])
				if j-1 == i {
					want = append(want, victim.val)
				}
			}

			if got := head.AsList(); !reflect.DeepEqual(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}

}

func TestMixIteration(t *testing.T) {
	type TestCase struct {
		nums      []int
		iterateOn int
		want      []int
	}

	testCases := []TestCase{
		TestCase{[]int{98, 97, 01, 95, 94}, 2, []int{98, 97, 95, 01, 94}}, //0
		TestCase{[]int{98, 97, 02, 95, 94}, 2, []int{98, 97, 95, 94, 02}}, //1
		TestCase{[]int{98, 97, 03, 95, 94}, 2, []int{98, 03, 97, 95, 94}}, //2
		TestCase{[]int{98, 97, 04, 95, 94}, 2, []int{98, 97, 04, 95, 94}}, //3
		TestCase{[]int{98, 97, 05, 95, 94}, 2, []int{98, 97, 95, 05, 94}}, //4

		TestCase{ //5
			[]int{3, 2, 1, 4, 5},
			2,
			[]int{3, 2, 4, 1, 5},
		},
		TestCase{ //6
			[]int{3, 2, -1, 4, 5},
			2,
			[]int{3, -1, 2, 4, 5},
		},
		TestCase{ // 7: negative past the start
			[]int{-1, 2, 3, 4, 5},
			0,
			[]int{2, 3, 4, -1, 5},
		},
		TestCase{ // 8: negative more than past the start
			[]int{-2, 1, 3, 4, 5},
			0,
			[]int{1, 3, -2, 4, 5},
		},
		TestCase{ // 9: negative that's >len
			[]int{-10, 2, 3, 4, 5},
			0,
			[]int{2, 3, -10, 4, 5},
		},
		TestCase{ // 10: positive just past the start
			[]int{5, 4, 3, 2, 1},
			4,
			[]int{5, 1, 4, 3, 2},
		},
		TestCase{ // 11: positive more than past the start
			[]int{5, 4, 3, 1, 2},
			4,
			[]int{5, 4, 2, 3, 1},
		},
		TestCase{ // 12: positive that's >len
			[]int{5, 4, 3, 2, 10},
			4,
			[]int{5, 4, 10, 3, 2},
		},
		TestCase{ // 13: zeros don't move
			[]int{5, 4, 0, 2, 1},
			2,
			[]int{5, 4, 0, 2, 1},
		},
		TestCase{ // 14:
			[]int{5, 4, 3, 2, 1},
			0,
			[]int{4, 5, 3, 2, 1},
		},
		TestCase{ // 15:
			[]int{6, 4, 3, 2, 1},
			0,
			[]int{4, 3, 6, 2, 1},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ptrs := make([]*TreeNode, len(tc.nums))
			head := makeTree(tc.nums, ptrs)

			was := head.AsList()

			head = mixIteration(ptrs[tc.iterateOn], head, ptrs)
			if got := head.AsList(); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("in %v got %v, want %v", was, got, tc.want)
			}
		})
	}
}

func TestMix(t *testing.T) {
	nums, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	ptrs := make([]*TreeNode, len(nums))
	head := makeTree(nums, ptrs)

	wants := [][]int{
		[]int{2, 1, -3, 3, -2, 0, 4},
		[]int{1, -3, 2, 3, -2, 0, 4},
		[]int{1, 2, 3, -2, -3, 0, 4},
		[]int{1, 2, -2, -3, 0, 3, 4},
		[]int{1, 2, -3, 0, 3, 4, -2},
		[]int{1, 2, -3, 0, 3, 4, -2},
		[]int{1, 2, -3, 4, 0, 3, -2},
	}

	for i, ptr := range ptrs {
		head = mixIteration(ptr, head, ptrs)
		//fmt.Printf("after moving %v, list %v\n", ptr.val, head.AsList())

		if got := head.AsList(); !reflect.DeepEqual(got, wants[i]) {
			t.Fatalf("step %d: got %v, want %v", i, got, wants[i])
		}
	}
}

func TestSolveA(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(input), 3; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(input), 1623178306; got != want {
		t.Errorf("solveB(sample) = %v, want %v", got, want)
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)

	sampleLines = strings.Split(rawSample, "\n")
	if len(sampleLines) > 0 && sampleLines[len(sampleLines)-1] == "" {
		sampleLines = sampleLines[0 : len(sampleLines)-1]
	}

	os.Exit(m.Run())
}
