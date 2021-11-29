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
	"strconv"
	"strings"
)

type Deck struct {
	name        string
	cap         int
	num         int
	cards       []int
	front, back int
}

func newDeck(name string, cards []int, cap int) *Deck {
	d := &Deck{
		name:  name,
		cap:   cap,
		front: 0,
		back:  len(cards) - 1,
		num:   len(cards),
	}

	d.cards = make([]int, d.cap)
	copy(d.cards, cards)

	return d
}

func (d *Deck) Name() string {
	return d.name
}

func (d *Deck) Num() int {
	return d.num
}

func (d *Deck) Empty() bool {
	return d.num == 0
}

func (d *Deck) Pop() (val int) {
	if d.Empty() {
		panic("pop empty")
	}

	val = d.cards[d.front]
	d.front = (d.front + 1) % d.cap
	d.num--
	return
}

func (d *Deck) Push(val int) {
	if d.num == d.cap {
		panic("push full")
	}

	d.back = (d.back + 1) % d.cap
	d.cards[d.back] = val
	d.num++
}

func (d *Deck) Cards() []int {
	out := []int{}
	for i := 0; i < d.num; i++ {
		out = append(out, d.cards[(d.front+i)%d.cap])
	}
	return out
}

func (d *Deck) Clone() *Deck {
	return d.CloneFirstN(d.num)
}

func (d *Deck) CloneFirstN(n int) *Deck {
	if n > d.num {
		panic("n too many")
	}

	nd := &Deck{
		name:  d.name,
		cap:   d.cap,
		front: 0,
		back:  n - 1,
		num:   n,
	}

	nd.cards = make([]int, nd.cap)
	copy(nd.cards, d.Cards()[0:n])

	return nd
}

func (d *Deck) String() string {
	out := make([]string, d.num)
	for i := 0; i < d.num; i++ {
		out[i] = strconv.Itoa(d.cards[(d.front+i)%d.cap])
	}
	return strings.Join(out, ", ")
}
