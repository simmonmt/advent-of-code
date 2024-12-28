package collections

func MapKeys[T comparable, U any](in map[T]U) []T {
	out := []T{}
	for k := range in {
		out = append(out, k)
	}
	return out
}
