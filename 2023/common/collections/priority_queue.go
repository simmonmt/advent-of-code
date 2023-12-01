package collections

import (
	"container/heap"
	"container/list"

	"golang.org/x/exp/constraints"
)

type ValueType interface {
	constraints.Integer | ~string
}

type PriorityQueue[T ValueType] interface {
	// Insert adds a value to the queue or updates its priority if the value
	// is already present in the queue. Returns true if a new value was
	// inserted.
	Insert(value T, priority int) bool

	Next() (value T, priority int)
	IsEmpty() bool
}

// Use to create a queue that returns the lowest priority first
func LessThan(a, b int) bool { return a < b }

// Use to create a queue that returns the highest priority first
func GreaterThan(a, b int) bool { return a > b }

// A simple priority queue implementation using a linked list. Shouldn't be used
// by anyone, but it was helpful as a brain-dead simple implementation to write
// a test against.

// A simple
type pqNaiveElem[T ValueType] struct {
	value    T
	priority int
}

type pqNaive[T ValueType] struct {
	elems      *list.List
	betterThan func(a, b int) bool
}

// NewNaivePriorityQueue creates a new low-performance priority queue that
// really shouldn't be used by anyone. Use NewPriorityQueue instead. betterThan
// determines the order in which values are returned by Next. Better values are
// returned first.
func NewNaivePriorityQueue[T ValueType](betterThan func(a, b int) bool) PriorityQueue[T] {
	return &pqNaive[T]{
		elems:      list.New(),
		betterThan: betterThan,
	}
}

func (q *pqNaive[T]) Insert(value T, priority int) bool {
	for elem := q.elems.Front(); elem != nil; elem = elem.Next() {
		pqElem := elem.Value.(*pqNaiveElem[T])
		if pqElem.value == value {
			pqElem.priority = priority
			return false
		}
	}

	q.elems.PushBack(&pqNaiveElem[T]{value: value, priority: priority})
	return true
}

func (q *pqNaive[T]) Next() (value T, priority int) {
	if q.IsEmpty() {
		panic("empty list")
	}

	var best *list.Element
	bestPriority := 0
	for elem := q.elems.Front(); elem != nil; elem = elem.Next() {
		elemPriority := elem.Value.(*pqNaiveElem[T]).priority

		if best == nil || q.betterThan(elemPriority, bestPriority) {
			best = elem
			bestPriority = elemPriority
		}
	}

	q.elems.Remove(best)
	return best.Value.(*pqNaiveElem[T]).value, bestPriority
}

func (q *pqNaive[T]) IsEmpty() bool {
	return q.elems.Front() == nil
}

// A heap-backed priority queue based on the container/heap PriorityQueue
// example in the Go documentation (pkg.go.dev/container/heap). An additional
// layer (the PriorityQueue interface, implemented by pqHeap) hides the
// interaction complexity (the need to use heap.* methods) and API breadth used
// by that example.
//
// Interface PriorityQueue, by pqHeap, uses pqHeapImpl to implement sort.Sort
// and heap.Interface methods to store pqHeapElem elements.

type pqHeapElem[T any] struct {
	value    T
	priority int
	index    int
}

type pqHeapImpl[T any] struct {
	arr        []*pqHeapElem[T]
	betterThan func(a, b int) bool
}

func (pqi *pqHeapImpl[T]) Len() int { return len(pqi.arr) }

func (pqi *pqHeapImpl[T]) Less(i, j int) bool {
	return pqi.betterThan(pqi.arr[i].priority, pqi.arr[j].priority)
}

func (pqi *pqHeapImpl[T]) Swap(i, j int) {
	pqi.arr[i], pqi.arr[j] = pqi.arr[j], pqi.arr[i]
	pqi.arr[i].index = i
	pqi.arr[j].index = j
}

func (pqi *pqHeapImpl[T]) Push(x any) {
	n := len(pqi.arr)
	elem := x.(*pqHeapElem[T])
	elem.index = n
	pqi.arr = append(pqi.arr, elem)
}

func (pqi *pqHeapImpl[T]) Pop() any {
	old := pqi.arr
	n := len(old)
	elem := old[n-1]
	old[n-1] = nil  // avoid memory leak
	elem.index = -1 // for safety
	pqi.arr = old[0 : n-1]
	return elem
}

func (pqi *pqHeapImpl[T]) update(elem *pqHeapElem[T], value T, priority int) {
	elem.value = value
	elem.priority = priority
	heap.Fix(pqi, elem.index)
}

type pqHeap[T ValueType] struct {
	impl  *pqHeapImpl[T]
	elems map[T]*pqHeapElem[T]
}

// NewPriorityQueue creates a new heap-backed priority queue.  betterThan
// determines the order in which values are returned by Next. Better values are
// returned first.
func NewPriorityQueue[T ValueType](betterThan func(a, b int) bool) PriorityQueue[T] {
	pq := &pqHeap[T]{
		impl: &pqHeapImpl[T]{
			arr:        []*pqHeapElem[T]{},
			betterThan: betterThan,
		},
		elems: map[T]*pqHeapElem[T]{},
	}

	heap.Init(pq.impl)
	return pq
}

func (pq *pqHeap[T]) Insert(value T, priority int) bool {
	if elem, found := pq.elems[value]; found {
		elem.priority = priority
		heap.Fix(pq.impl, elem.index)
		return false
	}

	elem := &pqHeapElem[T]{value: value, priority: priority}
	pq.elems[value] = elem
	heap.Push(pq.impl, elem)
	return true
}

func (pq *pqHeap[T]) Next() (value T, priority int) {
	elem := heap.Pop(pq.impl).(*pqHeapElem[T])
	delete(pq.elems, elem.value)
	return elem.value, elem.priority
}

func (pq *pqHeap[T]) IsEmpty() bool {
	return len(pq.impl.arr) == 0
}
