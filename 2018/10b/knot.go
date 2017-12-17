package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func readLengths(in io.Reader) ([]int, error) {
	reader := bufio.NewReader(in)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)
	lengths := []int{}
	for _, c := range line {
		lengths = append(lengths, int(c))
	}
	return lengths, nil
}

type WrappingList []int

func NewWrappingList(size int) WrappingList {
	l := make([]int, size)
	for i := range l {
		l[i] = i
	}
	return l
}

func (l *WrappingList) Sub(start, end int) []int {
	sub := []int{}
	for i := start; i <= end; i++ {
		sub = append(sub, (*l)[i%len(*l)])
	}
	return sub
}

func (l *WrappingList) Write(arr []int, start int) {
	for i, val := range arr {
		(*l)[(start+i)%len(*l)] = val
	}
}

func reverseArray(arr []int) []int {
	rev := make([]int, len(arr))
	for i, val := range arr {
		rev[len(arr)-i-1] = val
	}
	return rev
}

func runHash(list WrappingList, lengths []int, curPos, skipSize int) (int, int) {
	for _, length := range lengths {
		start := curPos
		end := curPos + length - 1

		src := list.Sub(start, end)
		rev := reverseArray(src)
		//fmt.Printf("replacing %d-%d %v with %v\n", start, end%len(list), src, rev)
		list.Write(rev, start)
		//fmt.Printf("list: %v\n", list)

		curPos = (curPos + length + skipSize) % len(list)
		skipSize++
	}

	return curPos, skipSize
}

func xorArray(arr []int) int {
	out := 0
	for _, val := range arr {
		out = out ^ val
	}
	fmt.Printf("xor of %d-value %v is %v\n", len(arr), arr, out)
	return out
}

func main() {
	listLen := 256
	if len(os.Args) != 1 {
		log.Fatalf("Usage: %v [listLen]\n", os.Args[0])
	}

	lengths, err := readLengths(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read lengths: %v\n", err)
	}
	lengths = append(lengths, 17, 31, 73, 47, 23)
	fmt.Printf("lengths: %v\n", lengths)

	list := NewWrappingList(listLen)
	fmt.Printf("list: %v\n", list)

	curPos := 0
	skipSize := 0
	for i := 0; i < 64; i++ {
		curPos, skipSize = runHash(list, lengths, curPos, skipSize)
	}

	dense := []int{}
	for i := 0; i < len(list); i += 16 {
		fmt.Printf("sending %d:%d\n", i, i+15)
		dense = append(dense, xorArray(list[i:i+16]))
	}

	out := ""
	for _, val := range dense {
		out += fmt.Sprintf("%02x", val)
	}

	fmt.Printf("%s\n", out)
}
