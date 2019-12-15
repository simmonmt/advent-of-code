package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/simmonmt/aoc/2019/common/intmath"
	"github.com/simmonmt/aoc/2019/common/logger"
	"github.com/simmonmt/aoc/2019/common/pos"
)

var (
	verbose  = flag.Bool("verbose", false, "verbose")
	input    = flag.String("input", "", "input file")
	maxSteps = flag.Int("maxsteps", -1, "number of steps (-1 = infinite)")
)

func readInput(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := []string{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lines, nil
}

type Planet struct {
	Name     string
	Pos, Vel [3]int
}

func (p *Planet) String() string {
	return fmt.Sprintf("id=%d pos=%s vel=%s", p.Name, p.Pos, p.Vel)
}

func (p *Planet) PotentialEnergy() int {
	return intmath.Abs(p.Pos[0]) + intmath.Abs(p.Pos[1]) + intmath.Abs(p.Pos[2])
}

func (p *Planet) KineticEnergy() int {
	return intmath.Abs(p.Vel[0]) + intmath.Abs(p.Vel[1]) + intmath.Abs(p.Vel[2])
}

func printPlanets(planets []Planet) {
	for _, p := range planets {
		fmt.Println(p)
	}
}

func applyGravity(planets []Planet) {
	adjust := func(pa, pb int, va, vb *int) {
		if pa > pb {
			*va -= 1
			*vb += 1
		} else if pa < pb {
			*va += 1
			*vb -= 1
		}
	}

	for i := range planets {
		for j := 0; j < i; j++ {
			if i == j {
				continue
			}

			adjust(planets[i].Pos[0], planets[j].Pos[0],
				&planets[i].Vel[0], &planets[j].Vel[0])
			adjust(planets[i].Pos[1], planets[j].Pos[1],
				&planets[i].Vel[1], &planets[j].Vel[1])
			adjust(planets[i].Pos[2], planets[j].Pos[2],
				&planets[i].Vel[2], &planets[j].Vel[2])
		}
	}
}

func adjustPositions(planets []Planet) {
	for i := range planets {
		planets[i].Pos[0] += planets[i].Vel[0]
		planets[i].Pos[1] += planets[i].Vel[1]
		planets[i].Pos[2] += planets[i].Vel[2]
	}
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	var planets [4]Planet
	for i, line := range lines {
		planet := &planets[i]
		planet.Name = strconv.Itoa(i)

		p, err := pos.P3FromString(line)
		if err != nil {
			log.Fatalf("bad coord %v: %v", line, err)
		}

		planet.Pos[0], planet.Pos[1], planet.Pos[2] = p.X, p.Y, p.Z
	}

	type CacheKey struct {
		pos, vel [4]int
	}

	// Each axis is computed entirely independently, so we can treat them
	// separately. Look for a cycle in each axis (again,
	// independently). Each axis cycle will have a different length, so then
	// look for the time when those cycles line up.

	cache := [3]map[CacheKey]int64{}
	for i := range cache {
		cache[i] = map[CacheKey]int64{}
	}

	var cycleStart [3]int64 = [3]int64{-1, -1, -1}
	var cycleLen [3]int64

	var numSteps int64
	for numSteps = 0; *maxSteps < 0 || numSteps < int64(*maxSteps); numSteps++ {
		if logger.Enabled() {
			fmt.Printf("%10d: ", numSteps)
			p := &planets[0]

			fmt.Printf("p=%4d v=%4d  pe=%4d ke=%4d  ",
				p.Pos[0], p.Vel[0],
				p.PotentialEnergy(), p.KineticEnergy())
			fmt.Println()
		}

		keepGoing := false
		for i := 0; i < 3; i++ {
			if cycleStart[i] != -1 {
				continue
			} else {
				keepGoing = true
			}

			key := CacheKey{
				pos: [4]int{planets[0].Pos[i], planets[1].Pos[i], planets[2].Pos[i], planets[3].Pos[i]},
				vel: [4]int{planets[0].Vel[i], planets[1].Vel[i], planets[2].Vel[i], planets[3].Vel[i]},
			}

			if val, ok := cache[i][key]; ok {
				fmt.Printf("found %d at %d\n", i, numSteps)
				cycleStart[i] = val
				cycleLen[i] = numSteps - val
			}
			cache[i][key] = numSteps
		}
		if !keepGoing {
			break
		}

		applyGravity(planets[:])
		adjustPositions(planets[:])
	}

	fmt.Printf("cycle start: %v\n", cycleStart)
	fmt.Printf("cycle len: %v\n", cycleLen)

	longest := int64(-1)
	longestIdx := -1
	for i, l := range cycleLen {
		if longest == -1 || l > longest {
			longest = l
			longestIdx = i
		}
	}

	matchNum := 0
	for t := cycleStart[longestIdx]; ; t += longest {
		matches := 0
		for i := 0; i < 3; i++ {
			if (t-cycleStart[i])%cycleLen[i] == 0 {
				matches++
			}
		}
		if matches == 3 {
			matchNum++
			if matchNum == 2 {
				fmt.Println(t)
				break
			}
		}
	}
}
