package collections

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
