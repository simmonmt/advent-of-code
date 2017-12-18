package knot

import (
	"fmt"
)

type hashArray struct {
	vals []int
}

func newHashArray(size int) *hashArray {
	arr := &hashArray{}
	arr.vals = make([]int, size)
	for i := range arr.vals {
		arr.vals[i] = i
	}
	return arr
}

func (h *hashArray) SubRange(base, size int) []int {
	out := make([]int, size)

	for i := 0; i < size; i++ {
		off := (base + i) % len(h.vals)
		out[i] = h.vals[off]
	}

	return out
}

func (h *hashArray) ReplRange(base int, vals []int) {
	for i := 0; i < len(vals); i++ {
		off := (base + i) % len(h.vals)
		h.vals[off] = vals[i]
	}
}

func reverse(arr []int) []int {
	out := make([]int, len(arr))
	for i := 0; i < len(arr); i++ {
		out[len(out)-i-1] = arr[i]
	}
	return out
}

func hashRound(arr *hashArray, lengths []int, pos, skip int) (int, int) {
	for _, length := range lengths {
		sub := arr.SubRange(pos, length)
		sub = reverse(sub)
		arr.ReplRange(pos, sub)

		pos += length + skip
		skip++
	}

	return pos, skip
}

func xorChunk(vals []int) int {
	out := 0
	for _, val := range vals {
		out ^= val
	}
	return out
}

func Hash(in string) string {
	lengths := make([]int, len(in))
	for i, r := range in {
		lengths[i] = int(r)
	}
	lengths = append(lengths, 17, 31, 73, 47, 23)
	//fmt.Printf("lengths: %v\n", lengths)

	size := 256
	arr := newHashArray(size)
	pos := 0
	skip := 0

	for i := 0; i < 64; i++ {
		pos, skip = hashRound(arr, lengths, pos, skip)
	}

	xord := []int{}
	for i := 0; i < size; i += 16 {
		xord = append(xord, xorChunk(arr.vals[i:i+16]))
	}

	out := ""
	for _, val := range xord {
		out += fmt.Sprintf("%02x", val)
	}
	return out
}
