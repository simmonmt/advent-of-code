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

import "container/list"

type Element struct {
	Index       int
	Repeater    rune
	GoodThrough int
}

type Queue struct {
	elems *list.List
}

func NewQueue() *Queue {
	return &Queue{elems: list.New()}
}

func (q *Queue) Add(index int, repeater rune, goodThrough int) {
	elem := &Element{
		Index:       index,
		Repeater:    repeater,
		GoodThrough: goodThrough,
	}

	q.elems.PushBack(elem)
}

func (q *Queue) ExpireTo(num int) {
	for {
		listElem := q.elems.Front()
		if listElem == nil {
			return
		}

		elem := listElem.Value.(*Element)
		if elem.GoodThrough <= num {
			q.elems.Remove(listElem)
		} else {
			return
		}
	}

	panic("unreachable")
}

func (q *Queue) ActiveBefore(index int) []*Element {
	out := []*Element{}
	for le := q.elems.Front(); le != nil; le = le.Next() {
		elem := le.Value.(*Element)
		if elem.Index < index {
			out = append(out, elem)
		}
	}
	return out
}

func (q *Queue) Delete(match *Element) {
	for listElem := q.elems.Front(); listElem != nil; listElem = listElem.Next() {
		elem := listElem.Value.(*Element)
		if elem.Index == match.Index && elem.Repeater == match.Repeater {
			q.elems.Remove(listElem)
			return
		}
	}
}
