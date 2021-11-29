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
	"regexp"

	"intmath"
	"logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")

	inputPat = regexp.MustCompile(`^position=< *(-?\d+), *(-?\d+)> velocity=< *(-?\d+), *(-?\d+)>$`)
)

type Point struct {
	PosX, PosY int
}

type Vel struct {
	VelX, VelY int
}

type Sky struct {
	Rep        int
	Points     []Point
	MinX, MinY int
	MaxX, MaxY int
	W, H       int
}

func NewSky(rep int, points []Point) *Sky {
	s := &Sky{
		Rep:    rep,
		Points: points,
		MinX:   -1,
		MaxX:   -1,
		MinY:   -1,
		MaxY:   -1,
	}

	for _, p := range points {
		if s.MinX == -1 || p.PosX < s.MinX {
			s.MinX = p.PosX
		}
		if s.MaxX == -1 || p.PosX > s.MaxX {
			s.MaxX = p.PosX
		}
		if s.MinY == -1 || p.PosY < s.MinY {
			s.MinY = p.PosY
		}
		if s.MaxY == -1 || p.PosY > s.MaxY {
			s.MaxY = p.PosY
		}
	}

	s.W = s.MaxX - s.MinX
	s.H = s.MaxY - s.MinY
	return s
}

func (s *Sky) Dump() {
	fmt.Printf("rep %d w %d h %d minx %v maxx %v miny %v maxy %v\n",
		s.Rep, s.W, s.H, s.MinX, s.MaxX, s.MinY, s.MaxY)

	coords := map[Point]bool{}
	for _, p := range s.Points {
		coords[p] = true
	}

	for y := s.MinY; y <= s.MaxY; y++ {
		for x := s.MinX; x <= s.MaxX; x++ {
			p := Point{x, y}
			if _, found := coords[p]; found {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func readInput() ([]Point, []Vel, error) {
	points := []Point{}
	vels := []Vel{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		parts := inputPat.FindStringSubmatch(line)
		if parts == nil {
			return nil, nil, fmt.Errorf("failed to parse %v", line)
		}

		posX := intmath.AtoiOrDie(parts[1])
		posY := intmath.AtoiOrDie(parts[2])
		velX := intmath.AtoiOrDie(parts[3])
		velY := intmath.AtoiOrDie(parts[4])

		points = append(points, Point{posX, posY})
		vels = append(vels, Vel{velX, velY})
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("read failed: %v", err)
	}

	return points, vels, nil
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	points, vels, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	var lastSky *Sky

	for rep := 0; ; rep++ {
		sky := NewSky(rep, points)

		if lastSky != nil {
			shrunk := sky.W <= lastSky.W || sky.H <= lastSky.H
			if !shrunk {
				break
			}
		}

		updated := make([]Point, len(points))
		for i := range points {
			np := Point{PosX: points[i].PosX + vels[i].VelX,
				PosY: points[i].PosY + vels[i].VelY,
			}
			updated[i] = np
		}
		points = updated
		lastSky = sky
	}
	lastSky.Dump()
}
