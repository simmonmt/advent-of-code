package data

func Grow(a []bool) []bool {
	out := make([]bool, len(a)*2+1)
	copy(out, a)
	for i := range a {
		outPos := len(a)*2 - i
		outVal := !a[i]
		out[outPos] = outVal
	}
	return out
}

func checksumRound(in []bool) []bool {
	out := make([]bool, len(in)/2)
	for i := range out {
		out[i] = in[i*2] == in[i*2+1]
	}
	return out
}

func Checksum(in []bool) []bool {
	if len(in)%2 == 1 {
		panic("odd checksum input")
	}

	for {
		sum := checksumRound(in)
		if len(sum)%2 == 1 {
			return sum
		}
		in = sum
	}
}
