package main

import (
	"container/list"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/logger"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func parseInput(lines []string) (string, error) {
	return lines[0], nil
}

type Elem struct {
	ID          int
	Start, Size int
}

func buildFS(input string) (fs []int, files, frees *list.List) {
	sz := 0
	for _, r := range input {
		sz += int(r - '0')
	}

	fs = make([]int, sz)
	files, frees = list.New(), list.New()
	id := 0
	isFile := true
	cur := 0

	fill := func(start, sz, v int) {
		for i := 0; i < sz; i++ {
			fs[start+i] = v
		}
	}

	for _, r := range input {
		n := int(r - '0')

		if isFile {
			fill(cur, n, id)
			files.PushBack(&Elem{ID: id, Start: cur, Size: n})
			id++
		} else if n != 0 {
			fill(cur, n, -1)
			frees.PushBack(&Elem{ID: id, Start: cur, Size: n})
		}

		cur += n
		isFile = !isFile
	}

	return
}

func score(fs []int) int64 {
	sum := int64(0)
	for i := 0; i < len(fs); i++ {
		if fs[i] == -1 {
			break
		}
		sum += int64(i * fs[i])
	}

	return sum
}

func doSolveA(fs []int, files, frees *list.List) {
	for file := files.Back(); file != nil; {
		fileElem := file.Value.(*Elem)
		freeElem := frees.Front().Value.(*Elem)

		for fileElem.Size > 0 && freeElem.Start < fileElem.Start {
			fs[freeElem.Start] = fileElem.ID
			fs[fileElem.Start+fileElem.Size-1] = -1

			freeElem.Start++
			freeElem.Size--
			if freeElem.Size == 0 {
				frees.Remove(frees.Front())
				front := frees.Front()
				if front == nil {
					// Ran out of free space
					return
				}
				freeElem = front.Value.(*Elem)
			}

			fileElem.Size--
		}

		if fileElem.Size != 0 {
			// Ran out of free space to the left of file
			return
		}

		next := file.Prev()
		files.Remove(file)
		file = next
	}
}

func solveA(input string) int64 {
	fs, files, frees := buildFS(input)
	doSolveA(fs, files, frees)
	return score(fs)
}

func solveB(input string) int64 {
	return -1
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *input == "" {
		logger.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		logger.Fatalf("failed to read input: %v", err)
	}

	input, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
