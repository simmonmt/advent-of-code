package collections

import (
	"cmp"
	"iter"
	"slices"
)

func MapKeys[T comparable, U any](in map[T]U) []T {
	out := []T{}
	for k := range in {
		out = append(out, k)
	}
	return out
}

func SortedMapIter[T cmp.Ordered, U any](in map[T]U) iter.Seq2[T, U] {
	keys := MapKeys(in)
	slices.Sort(keys)

	return func(yield func(T, U) bool) {
		for _, k := range keys {
			if !yield(k, in[k]) {
				return
			}
		}
	}
}
