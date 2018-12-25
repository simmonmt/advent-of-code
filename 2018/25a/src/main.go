package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"intmath"
	"logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
)

type Pos struct {
	X, Y, Z, T int
}

func (p Pos) Dist(o Pos) int {
	return intmath.Abs(o.X-p.X) + intmath.Abs(o.Y-p.Y) +
		intmath.Abs(o.Z-p.Z) + intmath.Abs(o.T-p.T)
}

func readInput() ([]Pos, error) {
	posns := []Pos{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, ",")
		x := intmath.AtoiOrDie(parts[0])
		y := intmath.AtoiOrDie(parts[1])
		z := intmath.AtoiOrDie(parts[2])
		t := intmath.AtoiOrDie(parts[3])

		posns = append(posns, Pos{x, y, z, t})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return posns, nil
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	posns, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	consts := make([][]Pos, len(posns))

	for i, p := range posns {
		consts[i] = []Pos{p}
	}

	for {
		moved := false
		for i, c := range consts {
			for _, pos := range c {
				move := false
				for j, o := range consts {
					if i == j || o == nil {
						continue
					}

					for _, op := range o {
						if pos.Dist(op) <= 3 {
							move = true
							break
						}
					}

					if move == true {
						consts[i] = append(consts[i], consts[j]...)
						consts[j] = nil
						moved = true
						goto done
					}
				}
			}
		}

	done:
		if !moved {
			break
		}
	}

	num := 0
	for _, c := range consts {
		if c == nil {
			continue
		}
		fmt.Println(c)
		num++
	}
	fmt.Println(num)
}
