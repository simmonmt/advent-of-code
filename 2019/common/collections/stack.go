package collections

type Stack struct {
	elems []interface{}
	last  int
}

func NewStack() *Stack {
	return &Stack{
		elems: []interface{}{},
		last:  -1,
	}
}

func (s *Stack) Push(elem interface{}) {
	s.last++
	if s.last == len(s.elems) {
		newElems := make([]interface{}, len(s.elems)+50)
		copy(newElems, s.elems)
		s.elems = newElems
	}

	s.elems[s.last] = elem
}

func (s *Stack) Pop() interface{} {
	if s.last < 0 {
		panic("pop empty stack")
	}

	ret := s.elems[s.last]
	s.last--
	return ret
}

func (s *Stack) Peek() interface{} {
	if s.last < 0 {
		panic("peek empty stack")
	}

	return s.elems[s.last]
}

func (s *Stack) Empty() bool {
	return s.last < 0
}

func (s *Stack) All() []interface{} {
	return s.elems[0:(s.last + 1)]
}
