package intmath

func Abs(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func Uint64Max(a, b uint64) uint64 {
	if a > b {
		return a
	} else {
		return b
	}
}
