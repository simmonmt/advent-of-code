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

package pad

import (
	"reflect"
	"testing"
)

func TestQueueEmpty(t *testing.T) {
	q := NewQueue()
	q.ExpireTo(999)

	if q.elems.Front() != nil {
		t.Errorf("expected empty; not empty")
	}
}

func TestQueueExpireAll(t *testing.T) {
	q := NewQueue()
	q.Add(1, 'a', 100)
	q.Add(2, 'b', 200)
	q.Add(3, 'c', 300)
	q.ExpireTo(999)

	if q.elems.Front() != nil {
		t.Errorf("expected empty; not empty")
	}
}

func TestQueueExpireSome(t *testing.T) {
	q := NewQueue()
	q.Add(1, 'a', 100)
	q.Add(2, 'b', 200)
	q.Add(3, 'c', 300)
	q.ExpireTo(200)

	if q.elems.Len() != 1 {
		t.Errorf("expected 1 elem, found %v", q.elems.Len())
	}
	expected := &Element{3, 'c', 300}
	if elem := q.elems.Front().Value.(*Element); !reflect.DeepEqual(expected, elem) {
		t.Errorf("found front elem %+v, want %+v", elem, expected)
	}
}

func TestQueueActiveBefore(t *testing.T) {
	q := NewQueue()
	q.Add(1, 'a', 100)
	q.Add(2, 'b', 200)
	q.Add(3, 'c', 300)

	expectedIndexes := []int{1, 2}
	activeElems := q.ActiveBefore(3)
	activeIndexes := []int{}
	for _, e := range activeElems {
		activeIndexes = append(activeIndexes, e.Index)
	}

	if !reflect.DeepEqual(expectedIndexes, activeIndexes) {
		t.Errorf("ActiveBefore() = %v, want %v", activeIndexes, expectedIndexes)
	}
}
