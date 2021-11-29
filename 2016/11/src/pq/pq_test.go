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

package pq

import (
	"container/heap"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	priQ := &PriorityQueue{}
	heap.Init(priQ)

	for _, pri := range []uint{10, 4, 2, 1, 5, 12, 3} {
		item := &Item{Value: strconv.Itoa(int(pri)), Priority: pri}
		heap.Push(priQ, item)
	}

	fmt.Println("contents:")
	for _, item := range *priQ {
		fmt.Printf("  %+v\n", item)
	}

	out := []uint{}
	for i := 0; i < 5; i++ {
		out = append(out, heap.Pop(priQ).(*Item).Priority)
	}

	expected := []uint{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(expected, out) {
		t.Errorf("expected %v, got %v\n", expected, out)
	}
}
