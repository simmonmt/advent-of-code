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

	"logger"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	numPlayers = flag.Int("num_players", -1, "num players")
	lastMarble = flag.Int("last_marble", -1, "last marble")
)

type Marbles struct {
	Elems *list.List
	Cur   *list.Element
}

func NewMarbles() *Marbles {
	m := &Marbles{
		Elems: list.New(),
	}

	m.Cur = m.Elems.PushBack(0)

	return m
}

func (m *Marbles) Dump(step int) {
	fmt.Printf("%3d: ", step)

	for e := m.Elems.Front(); e != nil; e = e.Next() {
		if e == m.Cur {
			fmt.Printf("(%3d) ", e.Value)
		} else {
			fmt.Printf(" %3d  ", e.Value)
		}
	}
	fmt.Println()
}

func (m *Marbles) Insert(num int) int {
	scoreDelta := 0

	if num%23 != 0 {
		e := m.cw(m.Cur)
		m.Cur = m.Elems.InsertAfter(num, e)
	} else {
		scoreDelta = num
		e := m.ccw(m.Cur, 7)
		m.Cur = m.cw(e)
		scoreDelta += e.Value.(int)
		m.Elems.Remove(e)
	}

	return scoreDelta
}

func (m *Marbles) cw(e *list.Element) *list.Element {
	e = e.Next()
	if e == nil {
		e = m.Elems.Front()
	}
	return e
}

func (m *Marbles) ccw(e *list.Element, num int) *list.Element {
	for i := 0; i < num; i++ {
		e = e.Prev()
		if e == nil {
			e = m.Elems.Back()
		}
	}
	return e
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *numPlayers == -1 {
		log.Fatal("--num_players is required")
	}
	if *lastMarble == -1 {
		log.Fatal("--last_marble is required")
	}

	marbles := NewMarbles()
	//marbles.Dump(0)

	marbleNum := 1
	playerNum := 1
	scores := map[int]int{}
	for {
		scoreDelta := marbles.Insert(marbleNum)
		scores[playerNum] += scoreDelta

		//marbles.Dump(playerNum)

		marbleNum++
		if marbleNum > *lastMarble {
			break
		}

		playerNum++
		if playerNum > *numPlayers {
			playerNum = 1
		}
	}

	//fmt.Println(scores)

	highScore := -1
	highPlayer := -1
	for p, s := range scores {
		if s > highScore {
			highScore = s
			highPlayer = p
		}
	}

	fmt.Printf("p %v s %v\n", highPlayer, highScore)
}
