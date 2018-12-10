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
	Points     map[Point]bool
	MinX, MinY int
	MaxX, MaxY int
	W, H       int
}

func (s *Sky) Dump() {
	// for x := s.MinX; x <= s.MaxX; x++ {
	// 	fmt.Printf("%3d ", x)
	// }
	// fmt.Println()

	for y := s.MinY; y <= s.MaxY; y++ {
		// fmt.Printf("%3d ", y)

		for x := s.MinX; x <= s.MaxX; x++ {
			p := Point{x, y}
			if _, found := s.Points[p]; found {
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
		sky := Sky{Points: map[Point]bool{}}
		for _, p := range points {
			sky.Points[p] = true
			if p.PosX < sky.MinX {
				sky.MinX = p.PosX
			}
			if p.PosX > sky.MaxX {
				sky.MaxX = p.PosX
			}
			if p.PosY < sky.MinY {
				sky.MinY = p.PosY
			}
			if p.PosY > sky.MaxY {
				sky.MaxY = p.PosY
			}
		}
		sky.W = sky.MaxX - sky.MinX
		sky.H = sky.MaxY - sky.MinY

		fmt.Printf("rep %d w %d h %d\n", rep, sky.W, sky.H)

		if lastSky != nil {
			if sky.W > lastSky.W {
				break
			}
		}

		//sky.Dump()

		updated := make([]Point, len(points))
		for i := range points {
			np := Point{PosX: points[i].PosX + vels[i].VelX,
				PosY: points[i].PosY + vels[i].VelY,
			}
			updated[i] = np
		}
		points = updated
		lastSky = &sky
	}
	lastSky.Dump()
}
