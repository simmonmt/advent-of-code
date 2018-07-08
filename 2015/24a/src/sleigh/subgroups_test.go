package sleigh

import (
	"reflect"
	"testing"
)

type result struct {
	group []int
	rest  []int
}

func testSubgroup(t *testing.T, in []int, size int, expected []result) {
	sg := NewSubgrouper(in, size)

	actual := []result{}
	for {
		if group, rest, ok := sg.Next(); ok {
			actual = append(actual, result{group, rest})
		} else {
			break
		}
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestOne(t *testing.T) {
	expected := []result{
		result{group: []int{1}, rest: []int{2, 3, 4, 5}},
		result{group: []int{2}, rest: []int{1, 3, 4, 5}},
		result{group: []int{3}, rest: []int{1, 2, 4, 5}},
		result{group: []int{4}, rest: []int{1, 2, 3, 5}},
		result{group: []int{5}, rest: []int{1, 2, 3, 4}},
	}

	testSubgroup(t, []int{1, 2, 3, 4, 5}, 1, expected)
}

func TestTwo(t *testing.T) {
	in := []int{1, 2, 3, 4}

	expected := []result{
		result{group: []int{1, 2}, rest: []int{3, 4}},
		result{group: []int{1, 3}, rest: []int{2, 4}},
		result{group: []int{1, 4}, rest: []int{2, 3}},
		result{group: []int{2, 3}, rest: []int{1, 4}},
		result{group: []int{2, 4}, rest: []int{1, 3}},
		result{group: []int{3, 4}, rest: []int{1, 2}},
	}

	testSubgroup(t, in, 2, expected)
}

func TestThree(t *testing.T) {
	in := []int{1, 2, 3, 4, 5}

	expected := []result{
		result{group: []int{1, 2, 3}, rest: []int{4, 5}},
		result{group: []int{1, 2, 4}, rest: []int{3, 5}},
		result{group: []int{1, 2, 5}, rest: []int{3, 4}},
		result{group: []int{1, 3, 4}, rest: []int{2, 5}},
		result{group: []int{1, 3, 5}, rest: []int{2, 4}},
		result{group: []int{1, 4, 5}, rest: []int{2, 3}},
		result{group: []int{2, 3, 4}, rest: []int{1, 5}},
		result{group: []int{2, 3, 5}, rest: []int{1, 4}},
		result{group: []int{2, 4, 5}, rest: []int{1, 3}},
		result{group: []int{3, 4, 5}, rest: []int{1, 2}},
	}

	testSubgroup(t, in, 3, expected)
}

func testSzEqualsGroupLen(t *testing.T) {
	expected := []result{
		result{group: []int{1, 2, 3, 4, 5}, rest: []int{}},
	}

	testSubgroup(t, []int{1, 2, 3, 4, 5}, 5, expected)
}
