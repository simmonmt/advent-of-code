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
