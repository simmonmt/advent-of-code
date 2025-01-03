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

func MapValues[T comparable, U any](in map[T]U) []U {
	out := []U{}
	for _, v := range in {
		out = append(out, v)
	}
	return out
}

func CloneMap[T comparable, U any](in map[T]U) map[T]U {
	out := make(map[T]U, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func FilterMap[T comparable, U any](in map[T]U, filter func(key T, val U) bool) map[T]U {
	out := map[T]U{}
	for k, v := range in {
		if filter(k, v) {
			out[k] = v
		}
	}
	return out
}

func OneMapKey[T comparable, U any](in map[T]U, def T) T {
	for k := range in {
		return k
	}
	return def
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
