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
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	pattern = regexp.MustCompile(`^p=<(-?[0-9]+),(-?[0-9]+),(-?[0-9]+)>, *v=<(-?[0-9]+),(-?[0-9]+),(-?[0-9]+)>, *a=<(-?[0-9]+),(-?[0-9]+),(-?[0-9]+)>`)
)

type Pos struct {
	X, Y, Z int
}

type Particle struct {
	Position, Velocity, Acceleration Pos
}

func newParticle(position, velocity, acceleration Pos) *Particle {
	return &Particle{position, velocity, acceleration}
}

func (p *Particle) Advance() {
	p.Velocity.X += p.Acceleration.X
	p.Velocity.Y += p.Acceleration.Y
	p.Velocity.Z += p.Acceleration.Z

	p.Position.X += p.Velocity.X
	p.Position.Y += p.Velocity.Y
	p.Position.Z += p.Velocity.Z
}

func posFromStrings(xStr, yStr, zStr string) (*Pos, error) {
	x, err := strconv.Atoi(xStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse x %v: %v", err)
	}
	y, err := strconv.Atoi(yStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse y %v: %v", err)
	}
	z, err := strconv.Atoi(zStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse z %v: %v", err)
	}

	return &Pos{x, y, z}, nil
}

// p=<1609,-863,-779>, v=<-15,54,-69>, a=<-10,0,14>

func readParticles(in io.Reader) ([]*Particle, error) {
	particles := []*Particle{}

	reader := bufio.NewReader(in)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		matches := pattern.FindStringSubmatch(strings.TrimSpace(line))
		if matches == nil {
			return nil, fmt.Errorf("failed to parse particle %v", line)
		}

		position, err := posFromStrings(matches[1], matches[2], matches[3])
		if err != nil {
			return nil, fmt.Errorf("failed to parse position from %v: %v", line, err)
		}

		velocity, err := posFromStrings(matches[4], matches[5], matches[6])
		if err != nil {
			return nil, fmt.Errorf("failed to parse velocity from %v: %v", line, err)
		}

		acceleration, err := posFromStrings(matches[7], matches[8], matches[9])
		if err != nil {
			return nil, fmt.Errorf("failed to parse acceleration from %v: %v", line, err)
		}

		particles = append(particles, newParticle(*position, *velocity, *acceleration))
	}

	return particles, nil
}

func hasTurned(a, b int) bool {
	if a < 0 {
		return b <= 0
	} else if a > 0 {
		return b >= 0
	} else {
		return true
	}
}

func posHasTurned(a, b *Pos) bool {
	return hasTurned(a.X, b.X) &&
		hasTurned(a.Y, b.Y) &&
		hasTurned(a.Z, b.Z)
}

func numParticlesHaveTurned(particles []*Particle) (int, int) {
	numTurned := 0
	unturnedExample := -1

	for i, part := range particles {
		if posHasTurned(&part.Velocity, &part.Acceleration) {
			numTurned++
		} else if unturnedExample == -1 {
			unturnedExample = i
		}
	}
	return numTurned, unturnedExample
}

func numParticlesOnRightSide(particles []*Particle) (int, int) {
	numDone := 0
	notDoneExample := -1

	for i, part := range particles {
		if posHasTurned(&part.Position, &part.Acceleration) {
			numDone++
		} else if notDoneExample == -1 {
			notDoneExample = i
		}
	}

	return numDone, notDoneExample
}

func abs(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func advanceUntilDone(particles []*Particle, pred func([]*Particle) (int, int)) {
	for i := 0; ; i++ {
		numDone, notDoneExample := pred(particles)
		fmt.Printf("iter %d; %d done, %d left", i, numDone, len(particles)-numDone)
		if notDoneExample >= 0 {
			fmt.Printf(" ex %d: %v", notDoneExample, particles[notDoneExample])
		}
		fmt.Println()
		if numDone == len(particles) {
			return
		}

		for _, part := range particles {
			part.Advance()
		}
	}
}

func main() {
	particles, err := readParticles(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read particles: %v", err)
	}

	// Step 1: Advance the particles until the
	// sign(velocity)=sign(acceleration). This means nobody else
	// will turn around.

	// for _, part := range particles {
	// 	fmt.Println(*part)
	// }

	fmt.Println("Step 1: wait for turning")
	advanceUntilDone(particles, numParticlesHaveTurned)

	fmt.Println("Step 2: wait for right side")
	advanceUntilDone(particles, numParticlesOnRightSide)

	// for _, part := range particles {
	// 	fmt.Println(*part)
	// }

	// Step 3: Find the one that's closest
	minDist := 0
	minDistNum := -1
	for i, part := range particles {
		partDist := abs(part.Position.X) + abs(part.Position.Y) + abs(part.Position.Z)
		if minDistNum == -1 || partDist < minDist {
			minDist = partDist
			minDistNum = i
		}
	}

	fmt.Printf("closest particle: %d\n", minDistNum)
}
