package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "verbose")
)

type Matrix struct {
	dim  int
	grid []bool
}

func NewMatrix(dim int) *Matrix {
	m := &Matrix{
		dim:  dim,
		grid: make([]bool, dim*dim),
	}

	m.Set(0, 0, true)
	m.Set(dim-1, 0, true)
	m.Set(0, dim-1, true)
	m.Set(dim-1, dim-1, true)

	return m
}

func (m *Matrix) Dim() int {
	return m.dim
}

func (m *Matrix) Clone(tgt *Matrix) {
	copy(tgt.grid, m.grid)
}

func (m *Matrix) Get(x, y int) bool {
	if x < 0 || y < 0 || x >= m.dim || y >= m.dim {
		return false
	}

	off := y*m.dim + x
	return m.grid[off]
}

func (m *Matrix) Set(x, y int, state bool) {
	shouldIgnore := !state && ((x == 0 && y == 0) || (x == m.dim-1 && y == 0) ||
		(x == 0 && y == m.dim-1) || (x == m.dim-1 && y == m.dim-1))
	if shouldIgnore {
		return
	}

	off := y*m.dim + x
	m.grid[off] = state
}

func (m *Matrix) Dump(w io.Writer) {
	writer := bufio.NewWriter(w)
	defer writer.Flush()

	off := 0
	for y := 0; y < m.dim; y++ {
		for x := 0; x < m.dim; x++ {
			if m.grid[off] {
				writer.WriteByte('#')
			} else {
				writer.WriteByte('.')
			}
			off++
		}
		writer.WriteByte('\n')
	}
}

func (m *Matrix) Count() int {
	numOn := 0
	for _, val := range m.grid {
		if val {
			numOn++
		}
	}
	return numOn
}

func readMatrix(r io.Reader) *Matrix {
	var m *Matrix

	reader := bufio.NewReader(r)
	for y := 0; ; y++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if y == 0 {
			dim := len(line)
			m = NewMatrix(dim)
		}

		for x, val := range line {
			if val == '#' {
				m.Set(x, y, true)
			}
		}
	}

	return m
}

var (
	stepOffsets = [][2]int{
		[2]int{-1, -1}, [2]int{0, -1}, [2]int{1, -1},
		[2]int{-1, 0}, [2]int{1, 0},
		[2]int{-1, 1}, [2]int{0, 1}, [2]int{1, 1},
	}
)

func nextState(x, y int, cur *Matrix) bool {
	numOn := 0
	for _, offset := range stepOffsets {
		if cur.Get(x+offset[0], y+offset[1]) {
			numOn++
		}
	}

	if cur.Get(x, y) {
		return numOn == 2 || numOn == 3
	} else {
		return numOn == 3
	}
}

func advance(cur, next *Matrix) {
	for y := 0; y < cur.Dim(); y++ {
		for x := 0; x < cur.Dim(); x++ {
			next.Set(x, y, nextState(x, y, cur))
		}
	}
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatalf("usage: %v niters\n", os.Args[0])
	}
	nIters, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		log.Fatalf("failed to parse niters: %v", err)
	}

	matrices := [2]*Matrix{}
	matrices[0] = readMatrix(os.Stdin)
	matrices[1] = NewMatrix(matrices[0].Dim())
	cur := 0

	if *verbose {
		matrices[cur].Dump(os.Stdout)
	}

	for i := 0; i < nIters; i++ {
		next := (cur + 1) % 2
		advance(matrices[cur], matrices[next])
		cur = next

		if *verbose {
			fmt.Println()
			matrices[cur].Dump(os.Stdout)
		}
	}

	fmt.Printf("num on = %d\n", matrices[cur].Count())
}
