// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	Pos, Vel pos.P3
}

func (p *Planet) String() string {
	return fmt.Sprintf("id=%d pos=%s vel=%s", p.Name, p.Pos, p.Vel)
}

func (p *Planet) PotentialEnergy() int {
	return intmath.Abs(p.Pos.X) + intmath.Abs(p.Pos.Y) + intmath.Abs(p.Pos.Z)
}

func (p *Planet) KineticEnergy() int {
	return intmath.Abs(p.Vel.X) + intmath.Abs(p.Vel.Y) + intmath.Abs(p.Vel.Z)
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

			adjust(planets[i].Pos.X, planets[j].Pos.X,
				&planets[i].Vel.X, &planets[j].Vel.X)
			adjust(planets[i].Pos.Y, planets[j].Pos.Y,
				&planets[i].Vel.Y, &planets[j].Vel.Y)
			adjust(planets[i].Pos.Z, planets[j].Pos.Z,
				&planets[i].Vel.Z, &planets[j].Vel.Z)
		}
	}
}

func adjustPositions(planets []Planet) {
	for i := range planets {
		planets[i].Pos.X += planets[i].Vel.X
		planets[i].Pos.Y += planets[i].Vel.Y
		planets[i].Pos.Z += planets[i].Vel.Z
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

	planets := make([]Planet, len(lines))
	for i, line := range lines {
		planet := &planets[i]
		planet.Name = strconv.Itoa(i)

		var err error
		if planet.Pos, err = pos.P3FromString(line); err != nil {
			log.Fatalf("bad coord %v: %v", line, err)
		}
	}

	var numSteps int
	for numSteps = 0; *maxSteps < 0 || numSteps < *maxSteps; numSteps++ {
		if logger.Enabled() {
			logger.LogF("numsteps = %d start", numSteps)
			printPlanets(planets)
		}

		applyGravity(planets)
		adjustPositions(planets)
	}

	fmt.Printf("final (numsteps = %d)\n", numSteps)
	printPlanets(planets)

	total := 0
	for _, p := range planets {
		potential := p.PotentialEnergy()
		kinetic := p.KineticEnergy()
		subtotal := potential * kinetic
		fmt.Printf("pot=%d kin=%d, tot=%d\n", potential, kinetic, subtotal)
		total += subtotal
	}

	fmt.Printf("system energy=%d\n", total)
}
