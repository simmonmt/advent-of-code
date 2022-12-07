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
	"container/list"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func readInput(path string) ([]string, error) {
	lines, err := filereader.Lines(path)
	if err != nil {
		return nil, err
	}

	return lines, nil
}

type INode struct {
	Name     string
	Parent   *INode
	Size     int
	Children []*INode
}

func MakeDirINode(name string, parent *INode) *INode {
	return &INode{
		Name:     name,
		Parent:   parent,
		Size:     -1,
		Children: []*INode{},
	}
}

func MakeFileINode(name string, size int) *INode {
	return &INode{
		Name: name,
		Size: size,
	}
}

type Mode int

const (
	MODE_CMD = 1
	MODE_LS  = 2
)

func buildFilesystem(lines []string) (*INode, error) {
	root := MakeDirINode("/", nil)
	cur := root
	mode := MODE_CMD

	makeError := func(i int, msg string) error {
		return fmt.Errorf("line %d: %s", i+1, msg)
	}

	for i, line := range lines {
		if mode == MODE_LS {
			if line[0] == '$' {
				mode = MODE_CMD
			} else {
				a, b, ok := strings.Cut(line, " ")
				if !ok {
					return nil, makeError(i, "failed split")
				}

				if strings.HasPrefix(a, "dir") {
					cur.Children = append(cur.Children,
						MakeDirINode(b, cur))
				} else {
					sz, err := strconv.Atoi(a)
					if err != nil {
						return nil, makeError(i, "bad size")
					}
					cur.Children = append(cur.Children,
						MakeFileINode(b, sz))
				}
			}
		}

		// We use a separate if statement for this comparison because
		// the mode may have changed during the previous one (that is we
		// thought we were doing a MODE_LS but saw a line beginning with
		// "$" which means no wait we're back in command mode).
		if mode == MODE_CMD {
			parts := strings.Split(line, " ")
			if len(parts) == 2 && parts[1] == "ls" {
				mode = MODE_LS
			} else if len(parts) == 3 && parts[1] == "cd" {
				if parts[2] == "/" {
					cur = root
				} else if parts[2] == ".." {
					cur = cur.Parent
					if cur == nil {
						return nil, makeError(i, "cd .. past /")
					}
				} else {
					var dest *INode
					for _, child := range cur.Children {
						if child.Name == parts[2] {
							dest = child
							break
						}
					}
					if dest == nil {
						return nil, makeError(i, "cd to nonexistent dir")
					}
					cur = dest
				}
			} else {
				return nil, makeError(i, "bad command")
			}
		}
	}

	return root, nil
}

func dumpINodeHeader(in *INode, indent int) {
	if in.Size == -1 {
		fmt.Printf("%*s- %s (dir)\n", indent, "", in.Name)
	} else {
		fmt.Printf("%*s- %s (file, size=%d)\n",
			indent, "", in.Name, in.Size)
	}
}

func dumpFilesystemLevel(in *INode, indent int) {
	dumpINodeHeader(in, indent)
	for _, child := range in.Children {
		dumpFilesystemLevel(child, indent+2)
	}
}

func dumpFilesystem(root *INode) {
	dumpFilesystemLevel(root, 0)
}

func pathToString(path *list.List) string {
	out := []string{}
	for elem := path.Front(); elem != nil; elem = elem.Next() {
		out = append(out, elem.Value.(string))
	}
	if len(out) == 0 {
		return "/"
	} else {
		return strings.Join(out, "/")
	}
}

func dirSizes(path *list.List, in *INode) (map[string]int, int) {
	subs := map[string]int{}
	dirSize := 0
	for _, child := range in.Children {
		if child.Size >= 0 {
			dirSize += child.Size
		} else {
			path.PushBack(child.Name)
			childSubs, childSize := dirSizes(path, child)
			path.Remove(path.Back())

			dirSize += childSize
			for path, sz := range childSubs {
				subs[path] = sz
			}
		}
	}

	subs[pathToString(path)] = dirSize
	return subs, dirSize
}

func solveA(root *INode) int {
	sizes, _ := dirSizes(list.New(), root)
	sum := 0
	for _, sz := range sizes {
		if sz > 100000 {
			continue
		}
		sum += sz
	}

	return sum
}

func solveB(root *INode) int {
	sizes, _ := dirSizes(list.New(), root)

	totalDisk := 70000000
	usedDisk := sizes["/"]
	freeDisk := totalDisk - usedDisk
	needForUpdate := 30000000
	toFreeForUpdate := needForUpdate - freeDisk

	closestSize := -1

	for _, sz := range sizes {
		if sz < toFreeForUpdate {
			continue
		}

		if closestSize == -1 || sz < closestSize {
			closestSize = sz
		}
	}

	return closestSize
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	fs, err := buildFilesystem(lines)
	if err != nil {
		log.Fatalf("failed to build filesystem: %v", err)
	}

	if logger.Enabled() {
		dumpFilesystem(fs)
	}

	fmt.Println("A", solveA(fs))
	fmt.Println("B", solveB(fs))
}
