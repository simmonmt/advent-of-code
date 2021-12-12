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
	"container/list"
	"flag"
	"fmt"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Graph struct {
	edges map[string][]string
}

func NewGraph() *Graph {
	return &Graph{
		edges: map[string][]string{},
	}
}

func (g *Graph) hasEdge(from, to string) bool {
	foundTos, found := g.edges[from]
	if !found {
		return false
	}

	for _, foundTo := range foundTos {
		if to == foundTo {
			return true
		}
	}

	return false
}

func (g *Graph) addDirectedEdge(from, to string) {
	if g.hasEdge(from, to) {
		panic("double add")
	}

	if g.edges[from] == nil {
		g.edges[from] = []string{to}
	} else {
		g.edges[from] = append(g.edges[from], to)
	}
}

func (g *Graph) AddEdge(from, to string) {
	g.addDirectedEdge(from, to)
	g.addDirectedEdge(to, from)
}

func (g *Graph) NodesFrom(from string) (tos []string, found bool) {
	tos, found = g.edges[from]
	return
}

func parseGraph(lines []string) (*Graph, error) {
	g := NewGraph()

	for i, line := range lines {
		parts := strings.SplitN(line, "-", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("%d: bad line %v", i, line)
		}

		from, to := parts[0], parts[1]
		g.AddEdge(from, to)
	}

	return g, nil
}

func readInput(path string) ([]string, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, err
	}

	return lines, err
}

func isSmallNode(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsLower(r)
}

func dfs(g *Graph, curPath *list.List, seen map[string]bool, cb func(path *list.List, cur string)) {
	cur := curPath.Back().Value.(string)
	cb(curPath, cur)

	tos, found := g.NodesFrom(cur)
	if !found {
		panic("can't get out")
	}

	for _, to := range tos {
		isSmall := isSmallNode(to)

		if isSmall {
			if found := seen[to]; found {
				continue // can't revisit small nodes
			}
			seen[to] = true
		}

		curPath.PushBack(to)
		dfs(g, curPath, seen, cb)
		curPath.Remove(curPath.Back())

		if isSmall {
			seen[to] = false
		}
	}
}

func pathToString(l *list.List) string {
	out := []string{}
	for e := l.Front(); e != nil; e = e.Next() {
		out = append(out, e.Value.(string))
	}
	return strings.Join(out, ",")
}

func allPaths(g *Graph) []string {
	curPath := list.New()
	curPath.PushBack("start")

	seen := map[string]bool{"start": true}

	paths := []string{}

	dfs(g, curPath, seen, func(path *list.List, cur string) {
		if cur == "end" {
			paths = append(paths, pathToString(path))
		}
	})

	return paths
}

func solveA(g *Graph) {
	paths := allPaths(g)
	fmt.Println("A", len(paths))
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

	g, err := parseGraph(lines)
	if err != nil {
		log.Fatal(err)
	}

	solveA(g)
}
