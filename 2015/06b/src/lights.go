package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var (
	commandPattern = regexp.MustCompile(`(turn on|turn off|toggle) ([0-9]+),([0-9]+) through ([0-9]+),([0-9]+)`)
)

type Pos struct {
	X, Y int
}

type Matrix struct {
	mat  [][]int
	w, h int
}

func NewMatrix(w, h int) *Matrix {
	mat := make([][]int, h)
	for y := 0; y < h; y++ {
		mat[y] = make([]int, w)
	}

	return &Matrix{
		mat: mat,
		w:   w,
		h:   h,
	}
}

func (m *Matrix) Set(a, b Pos, delta int) {
	for y := a.Y; y <= b.Y; y++ {
		for x := a.X; x <= b.X; x++ {
			m.mat[y][x] += delta
			if m.mat[y][x] < 0 {
				m.mat[y][x] = 0
			}
		}
	}
}

func (m *Matrix) Toggle(a, b Pos) {
	for y := a.Y; y <= b.Y; y++ {
		for x := a.X; x <= b.X; x++ {
			m.mat[y][x] += 2
		}
	}
}

func (m *Matrix) Count() int {
	total := 0

	for _, row := range m.mat {
		for _, val := range row {
			total += val
		}
	}

	return total
}

func (m *Matrix) Dump() {
	for _, row := range m.mat {
		for _, val := range row {
			fmt.Printf("%02d ", val)
		}
		fmt.Printf("\n")
	}
}

func mkPos(xStr, yStr string) (Pos, error) {
	x, err := strconv.Atoi(xStr)
	if err != nil {
		return Pos{}, err
	}
	y, err := strconv.Atoi(yStr)
	if err != nil {
		return Pos{}, err
	}

	return Pos{X: x, Y: y}, nil
}

func main() {
	mat := NewMatrix(1000, 1000)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		match := commandPattern.FindStringSubmatch(line)
		if match == nil {
			log.Fatalf("failed to parse %v", match)
		}
		fmt.Println(match)

		command := match[1]
		a, err := mkPos(match[2], match[3])
		if err != nil {
			log.Fatalf("failed to parse a in %v: %v", line, err)
		}
		b, err := mkPos(match[4], match[5])
		if err != nil {
			log.Fatalf("failed to parse b in %v: %v", line, err)
		}

		switch command {
		case "turn on":
			mat.Set(a, b, 1)
			break
		case "turn off":
			mat.Set(a, b, -1)
			break
		case "toggle":
			mat.Toggle(a, b)
			break
		default:
			log.Fatalf("unknown command %v in %v", command, line)
		}
		// mat.Dump()
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("reading stdin: %v", err)
	}

	fmt.Println(mat.Count())
}
