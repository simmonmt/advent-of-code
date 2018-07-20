package pad

import (
	"crypto/md5"
	"fmt"
)

func hexItoa(b byte) rune {
	if b < 10 {
		return rune('0' + b)
	} else {
		return rune('a' + (b - 10))
	}
}

func doHash(str string) string {
	h := md5.Sum([]byte(str))

	out := make([]rune, len(h)*2)
	for i, b := range h {
		out[i*2] = hexItoa(b >> 4)
		out[i*2+1] = hexItoa(b & 0xf)
	}

	return string(out)
}

func makeNormalHash(salt string, index int) string {
	return doHash(fmt.Sprintf("%v%d", salt, index))
}

type Hasher interface {
	MakeHash(salt string, index int) string
}

type NormalHasher struct{}

func (h *NormalHasher) MakeHash(salt string, index int) string {
	return makeNormalHash(salt, index)
}

type StretchedHasher struct{}

func (h *StretchedHasher) MakeHash(salt string, index int) string {
	hStr := makeNormalHash(salt, index)
	for i := 0; i < 2016; i++ {
		hStr = doHash(hStr)
	}
	return hStr
}
