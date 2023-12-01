// Copyright 2023 Google LLC
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

package collections

type Stack struct {
	elems []any
	last  int
}

func NewStack() *Stack {
	return &Stack{
		elems: []any{},
		last:  -1,
	}
}

func (s *Stack) Push(elem any) {
	s.last++
	if s.last == len(s.elems) {
		newElems := make([]any, len(s.elems)+50)
		copy(newElems, s.elems)
		s.elems = newElems
	}

	s.elems[s.last] = elem
}

func (s *Stack) Pop() any {
	if s.last < 0 {
		panic("pop empty stack")
	}

	ret := s.elems[s.last]
	s.last--
	return ret
}

func (s *Stack) Peek() any {
	if s.last < 0 {
		panic("peek empty stack")
	}

	return s.elems[s.last]
}

func (s *Stack) Empty() bool {
	return s.last < 0
}

func (s *Stack) All() []any {
	return s.elems[0:(s.last + 1)]
}
