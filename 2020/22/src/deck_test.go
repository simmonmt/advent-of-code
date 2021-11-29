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
	"reflect"
	"testing"
)

func seq(min, max int) []int {
	out := make([]int, max-min+1)
	for i := min; i <= max; i++ {
		out[i-min] = i
	}
	return out
}

func TestDeckInit(t *testing.T) {
	deck := newDeck("test", seq(0, 9), 10)
	got := deck.Cards()
	want := seq(0, 9)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestDeckLoopingAround(t *testing.T) {
	deck := newDeck("test", seq(0, 9), 10)
	got := []int{}
	for i := 0; i < 100; i++ {
		val := deck.Pop()
		got = append(got, val)
		deck.Push(val + 10)
	}

	want := seq(0, 99)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestDeckDrain(t *testing.T) {
	deck := newDeck("test", seq(0, 9), 10)

	got := []int{}
	for !deck.Empty() {
		got = append(got, deck.Pop())
	}

	want := seq(0, 9)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
		return
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("pop empty didn't panic")
		}
	}()

	deck.Pop()
}

func TestDeckPush(t *testing.T) {
	deck := newDeck("test", []int{}, 10)

	for i := 0; i < 10; i++ {
		deck.Push(i)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("push full didn't panic")
		}
	}()

	deck.Push(10)
}
