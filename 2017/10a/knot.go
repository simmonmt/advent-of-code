package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func readLengths(in io.Reader) ([]int, error) {
	reader := bufio.NewReader(in)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)
	parts := strings.Split(line, ",")

	lengths := []int{}
	for _, part := range parts {
		length, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("failed to parse part %v: %v", part, err)
		}
		lengths = append(lengths, length)
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

func main() {
	listLen := 256
	if len(os.Args) == 2 {
		var err error
		listLen, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalf("failed to parse list len %v\n", listLen)
		}
	} else if len(os.Args) != 1 {
		log.Fatalf("Usage: %v [listLen]\n", os.Args[0])
	}

	lengths, err := readLengths(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read lengths: %v\n", err)
	}
	fmt.Printf("lengths: %v\n", lengths)

	list := NewWrappingList(listLen)
	fmt.Printf("list: %v\n", list)

	curPos := 0
	skipSize := 0
	for _, length := range lengths {
		start := curPos
		end := curPos + length - 1

		src := list.Sub(start, end)
		rev := reverseArray(src)
		fmt.Printf("replacing %d-%d %v with %v\n", start, end%listLen, src, rev)
		list.Write(rev, start)
		fmt.Printf("list: %v\n", list)

		curPos = (curPos + length + skipSize) % len(list)
		skipSize++
	}

	fmt.Printf("out: %d\n", list[0]*list[1])
}
