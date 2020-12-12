package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/simmonmt/aoc/2020/common/dir"
	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/intmath"
	"github.com/simmonmt/aoc/2020/common/logger"
	"github.com/simmonmt/aoc/2020/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func posToCompass(p pos.P2) string {
	ew := "E"
	if p.X < 0 {
		ew = "W"
	}

	ns := "S"
	if p.Y < 0 {
		ns = "N"
	}

	return fmt.Sprintf("%s%d%s%d",
		ew, intmath.Abs(p.X), ns, intmath.Abs(p.Y))
}

func solveA(lines []string) {
	startPos := pos.P2{X: 0, Y: 0}

	p := startPos
	d := dir.DIR_EAST

	for _, line := range lines {
		cmd := line[0]
		num := intmath.AtoiOrDie(line[1:])

		switch string(cmd) {
		case "N":
			fallthrough
		case "S":
			fallthrough
		case "E":
			fallthrough
		case "W":
			p = dir.Parse(string(cmd)).StepsFrom(p, num)
		case "L":
			for ; num > 0; num -= 90 {
				d = d.Left()
			}
		case "R":
			for ; num > 0; num -= 90 {
				d = d.Right()
			}
		case "F":
			p = d.StepsFrom(p, num)
		}

		logger.LogF("line %v result %s %s", line, d, p)
	}

	fmt.Printf("A: %v\n", p.ManhattanDistance(startPos))
}

func rotateWaypointLeft(shipPos, wpRelPos pos.P2, deg int) pos.P2 {
	for ; deg > 0; deg -= 90 {
		wpRelPos = pos.P2{X: wpRelPos.Y, Y: -wpRelPos.X}
	}
	return wpRelPos
}

func rotateWaypointRight(shipPos, wpRelPos pos.P2, deg int) pos.P2 {
	for ; deg > 0; deg -= 90 {
		wpRelPos = pos.P2{X: -wpRelPos.Y, Y: wpRelPos.X}
	}
	return wpRelPos
}

func solveB(lines []string) {
	shipStartPos := pos.P2{X: 0, Y: 0}

	shipPos := shipStartPos
	shipDir := dir.DIR_EAST
	wpRelPos := pos.P2{X: shipStartPos.X + 10, Y: shipStartPos.Y - 1}

	for _, line := range lines {
		cmd := line[0]
		num := intmath.AtoiOrDie(line[1:])

		switch string(cmd) {
		case "N":
			fallthrough
		case "S":
			fallthrough
		case "E":
			fallthrough
		case "W":
			wpRelPos = dir.Parse(string(cmd)).StepsFrom(wpRelPos, num)
		case "L":
			wpRelPos = rotateWaypointLeft(shipPos, wpRelPos, num)
		case "R":
			wpRelPos = rotateWaypointRight(shipPos, wpRelPos, num)
		case "F":
			for i := 0; i < num; i++ {
				shipPos.Add(wpRelPos)
			}
		}

		logger.LogF("line %v result ship %s %s wp %v",
			line, shipDir, posToCompass(shipPos),
			posToCompass(wpRelPos))
	}

	fmt.Printf("B: %v\n", shipPos.ManhattanDistance(shipStartPos))
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		log.Fatal(err)
	}

	solveA(lines)
	solveB(lines)
}
