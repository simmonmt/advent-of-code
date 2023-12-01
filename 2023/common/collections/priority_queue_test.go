package collections

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func testPriorityQueue(t *testing.T, q PriorityQueue[int]) {
	in := make([]int, 50)
	for i := 0; i < 50; i++ {
		in[i] = i
	}
	rand.Shuffle(50, func(i, j int) {
		in[i], in[j] = in[j], in[i]
	})

	if got, want := q.IsEmpty(), true; !got {
		t.Fatalf("IsEmpty() = %v, want %v", got, want)
	}

	for _, elem := range in {
		if got, want := q.Insert(elem*10, elem), true; got != want {
			t.Errorf("Insert(%v, %v) = %v, want %v",
				elem*10, elem, got, want)
		}
		if got, want := q.IsEmpty(), false; got {
			t.Fatalf("IsEmpty() = %v, want %v", got, want)
		}
	}

	out := []int{}
	for !q.IsEmpty() {
		val, pri := q.Next()
		if val != pri*10 {
			t.Errorf("got val %v with pri %v, val=pri*10",
				val, pri)
		}
		out = append(out, pri)
	}

	want := in[:]
	sort.Ints(want)

	if !reflect.DeepEqual(want, out) {
		t.Errorf("wanted priorities %v, got %v", want, out)
	}

	// Now test update

	in = []int{4, 2, 1, 5, 6} // missing 3
	for _, elem := range in {
		if got, want := q.Insert(elem*10, elem), true; got != want {
			t.Errorf("Insert(%v, %v) = %v, want %v",
				elem*10, elem, got, want)
		}
	}

	if got, want := q.Insert(50, 3), false; got != want {
		t.Errorf("Update(50,3) = %v, want %v", got, want)
	}

	want = []int{10, 20, 50, 40, 60}

	out = []int{}
	for !q.IsEmpty() {
		val, _ := q.Next()
		out = append(out, val)
	}

	if !reflect.DeepEqual(want, out) {
		t.Errorf("wanted priorities %v, got %v", want, out)
	}
}

func TestPriorityQueue(t *testing.T) {
	type Constructor struct {
		name string
		f    func() PriorityQueue[int]
	}

	ctors := []Constructor{
		Constructor{
			"naive",
			func() PriorityQueue[int] {
				return NewNaivePriorityQueue[int](LessThan)
			},
		},
		Constructor{
			"heap",
			func() PriorityQueue[int] {
				return NewPriorityQueue[int](LessThan)
			},
		},
	}

	for _, ctor := range ctors {
		t.Run(ctor.name, func(t *testing.T) {
			testPriorityQueue(t, ctor.f())
		})
	}
}
