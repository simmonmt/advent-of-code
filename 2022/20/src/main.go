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

// -9668 wrong

package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func parseInput(lines []string) ([]int, error) {
	out := []int{}
	for _, line := range lines {
		num, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("bad num %v: %v", line, err)
		}
		out = append(out, num)
	}
	return out, nil
}

type TreeNode struct {
	up, before, after *TreeNode
	size              int
	val               int
}

func (n *TreeNode) fillList(l []int) {
	if len(l) != n.size {
		panic(fmt.Sprintf("size mismatch; want %v got %v", len(l), n.size))
	}

	i := 0
	if n.before != nil {
		n.before.fillList(l[i : i+n.before.size])
		i += n.before.size
	}
	l[i] = n.val
	i++
	if n.after != nil {
		n.after.fillList(l[i:])
	}
}

func (n *TreeNode) Check() {
	wantSz := 1
	if n.before != nil {
		n.before.Check()
		if n.before.up != n {
			panic("bad n.before")
		}
		wantSz += n.before.size
	}
	if n.after != nil {
		n.after.Check()
		if n.after.up != n {
			panic("bad n.after")
		}
		wantSz += n.after.size
	}

	if wantSz != n.size {
		panic("bad size")
	}
}

func (n *TreeNode) AsList() []int {
	n.Check()

	l := make([]int, n.size)
	n.fillList(l)
	return l
}

func (n *TreeNode) Forward() *TreeNode {
	if n.size == 1 {
		var p, t *TreeNode = n.up, n
		for p != nil {
			if p.before == t {
				return p
			}

			t, p = p, p.up
		}
		panic("off the top")
	} else if n.after != nil {
		t := n.after
		for t.before != nil {
			t = t.before
		}
		return t
	} else {
		var p, t *TreeNode = n.up, n
		for p != nil {
			if p.after == t {
				return p
			}
			t, p = p, p.up
		}
		panic("off the top")
	}
}

func (n *TreeNode) Index() int {
	off := 0
	if n.before != nil {
		off = n.before.size
	}

	t, p := n, n.up
	for {
		if p == nil {
			return off
		}
		if p.before == t {
			t, p = p, p.up
			continue
		}

		return p.Index() + 1 + off
	}
}

func (n *TreeNode) FindIndex(i int) *TreeNode {
	t := n

	for {
		beforeSz := 0
		if t.before != nil {
			beforeSz = t.before.size
			if i < beforeSz {
				t = t.before
				continue
			}
		}

		i -= beforeSz
		if i == 0 {
			return t
		}

		i--
		t = t.after
		if t == nil {
			panic("bad sizes")
		}
	}
}

func (n *TreeNode) Remove(head **TreeNode) {
	if n.before == nil && n.after == nil {
		// fmt.Println("no children case")

		if n.up == nil {
			panic("can't remove last")
		}

		p := n.up
		n.up = nil
		if p.before == n {
			p.before = nil
		} else if p.after == n {
			p.after = nil
		}

		for ; p != nil; p = p.up {
			p.size--
		}
		return
	}

	if (n.before == nil) != (n.after == nil) {
		//fmt.Println("one child case")

		// has before or after but not both, so replace this
		// node with the one that exists.
		rep := n.before
		if rep == nil {
			rep = n.after
		}

		if n.up == nil {
			*head = rep
			return
		}
		p := n.up

		rep.up = p
		if p.before == n {
			p.before = rep
		} else {
			p.after = rep
		}

		for ; p != nil; p = p.up {
			p.size--
		}

		n.size = 1
		n.before, n.after = nil, nil
		n.up = nil
		return
	}

	//fmt.Println("two children case")

	// has both before and after
	var rep *TreeNode
	if n.before.size < n.after.size {
		// find least from bigger after subtree
		//fmt.Println("from after")
		for rep = n.after; rep.before != nil; rep = rep.before {
		}
	} else {
		// find greatest from bigger before subtree
		//fmt.Println("from before")
		for rep = n.before; rep.after != nil; rep = rep.after {
		}
	}

	//fmt.Printf("rep is %+v\n", rep)

	rep.Remove(head)
	rep.size = 1
	rep.after = n.after
	if n.after != nil {
		n.after.up = rep
		rep.size += n.after.size
	}
	rep.before = n.before
	if n.before != nil {
		n.before.up = rep
		rep.size += n.before.size
	}

	if n.up == nil {
		*head = rep
		rep.up = nil
	} else {
		p := n.up
		rep.up = p
		if p.before == n {
			p.before = rep
		} else {
			p.after = rep
		}
	}

	n.size = 1
	n.before, n.after = nil, nil
	n.up = nil
	return
}

func (n *TreeNode) InsertAfter(dest *TreeNode) {
	if dest.after == nil {
		dest.after = n
		n.up = dest

	} else {
		dest = dest.after
		for dest.before != nil {
			dest = dest.before
		}

		dest.before = n
		n.up = dest
	}

	for p := dest; p != nil; p = p.up {
		p.size += n.size
	}
}

func makeTree(nums []int, ptrs []*TreeNode) *TreeNode {
	if len(nums) == 1 {
		head := &TreeNode{size: 1, val: nums[0]}
		ptrs[0] = head
		return head
	}

	beforeSz := (len(nums) - 1) / 2
	afterSz := len(nums) - beforeSz - 1

	head := &TreeNode{size: len(nums), val: nums[beforeSz]}
	ptrs[beforeSz] = head

	if beforeSz > 0 {
		head.before = makeTree(nums[0:beforeSz], ptrs[0:beforeSz])
		head.before.up = head
	}
	if afterSz > 0 {
		head.after = makeTree(nums[beforeSz+1:], ptrs[beforeSz+1:])
		head.after.up = head
	}

	return head
}

func mixIteration(ptr *TreeNode, head *TreeNode, ptrs []*TreeNode) *TreeNode {
	if ptr.val == 0 {
		return head
	}

	cur := ptr.Index()

	// one we remove cur there will be len(ptrs)-1 elements in the
	// list, so we have to do our math using len-1.

	// the index before cur's
	beforeCur := (cur - 1) % len(ptrs)

	// the index before where cur will go
	beforeNew := beforeCur + ptr.val
	for beforeNew < len(ptrs)-1 {
		beforeNew += len(ptrs) - 1
	}
	beforeNew %= len(ptrs) - 1

	ptr.Remove(&head)
	ptr.InsertAfter(head.FindIndex(beforeNew))
	return head
}

func mix(head *TreeNode, ptrs []*TreeNode) *TreeNode {
	for _, ptr := range ptrs {
		//fmt.Println(i)
		head.Check()
		head = mixIteration(ptr, head, ptrs)
	}

	return head
}

func lookup(head *TreeNode, i int) int {
	return head.FindIndex(i % head.size).val
}

func solveA(nums []int) int {
	fmt.Println(len(nums), "nums")

	ptrs := make([]*TreeNode, len(nums))
	head := makeTree(nums, ptrs)

	head = mix(head, ptrs)

	var zero *TreeNode
	for _, ptr := range ptrs {
		if ptr.val == 0 {
			zero = ptr
			break
		}
	}

	fmt.Printf("zero at %v\n", zero.Index())

	sum := 0
	for _, n := range []int{1000, 2000, 3000} {
		sum += lookup(head, zero.Index()+n)
	}

	return sum
}

func solveB(nums []int) int {
	return -1
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		log.Fatal(err)
	}

	input, err := parseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
