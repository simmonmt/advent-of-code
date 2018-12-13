package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"logger"
)

var (
	verbose  = flag.Bool("verbose", false, "verbose")
	numTicks = flag.Int("num_ticks", -1, "num ticks")
)

type TrackOrientation int
type CarDirection int

const (
	TO_NONE TrackOrientation = iota
	TO_LEFTRIGHT
	TO_UPDOWN
	TO_UPRIGHT
	TO_UPLEFT
	TO_INTERSECTION

	CD_UP CarDirection = iota
	CD_DOWN
	CD_LEFT
	CD_RIGHT
)

func (to TrackOrientation) String() string {
	switch to {
	case TO_NONE:
		return " "
	case TO_LEFTRIGHT:
		return "-"
	case TO_UPDOWN:
		return "|"
	case TO_UPRIGHT:
		return "/"
	case TO_UPLEFT:
		return "\\"
	case TO_INTERSECTION:
		return "+"
	default:
		panic("bad to")
	}
}

func (cd CarDirection) String() string {
	switch cd {
	case CD_UP:
		return "^"
	case CD_DOWN:
		return "v"
	case CD_LEFT:
		return "<"
	case CD_RIGHT:
		return ">"
	default:
		panic("unknown dir")
	}
}

type Track [][]TrackOrientation

func (t *Track) Dump(cars map[Loc]Car) {
	for y := 0; y < len(*t); y++ {
		row := (*t)[y]
		for x := 0; x < len(row); x++ {
			loc := Loc{x, y}
			if car, found := cars[loc]; found {
				if car.CrashedWith == nil {
					fmt.Print(car.D)
				} else {
					fmt.Print("X")
				}
			} else {
				fmt.Print(row[x])
			}
		}
		fmt.Println()
	}
}

type Loc struct {
	X, Y int
}

func (l Loc) String() string {
	return fmt.Sprintf("(%d,%d)", l.X, l.Y)
}

type Car struct {
	Num         int
	NumTurns    int
	CrashedWith []int
	L           Loc
	D           CarDirection
}

func NewCar(num int, l Loc, d CarDirection) Car {
	return Car{
		Num:         num,
		NumTurns:    0,
		CrashedWith: nil,
		L:           l,
		D:           d,
	}
}

func (c Car) String() string {
	crashed := "ok"
	if c.CrashedWith != nil {
		with := []string{}
		for _, w := range c.CrashedWith {
			with = append(with, strconv.Itoa(w))
		}
		crashed = strings.Join(with, ":")
	}
	return fmt.Sprintf("#%d@(%d,%d),%s,%s,nt%d", c.Num, c.L.X, c.L.Y, crashed, c.D, c.NumTurns)
}

func readInput() (*Track, map[Loc]Car, error) {
	track := Track{}
	cars := map[Loc]Car{}
	scanner := bufio.NewScanner(os.Stdin)

	addCar := func(x, y int, cd CarDirection) {
		carNum := len(cars)
		loc := Loc{x, y}
		cars[loc] = NewCar(carNum, loc, cd)
	}

	for y := 0; scanner.Scan(); y++ {
		line := scanner.Text()

		row := []TrackOrientation{}
		for x, c := range line {
			dir := TO_NONE
			switch c {
			case '-':
				dir = TO_LEFTRIGHT
			case '|':
				dir = TO_UPDOWN
			case '/':
				dir = TO_UPRIGHT
			case '\\':
				dir = TO_UPLEFT
			case '+':
				dir = TO_INTERSECTION
			case ' ':
				dir = TO_NONE
			case '^':
				dir = TO_UPDOWN
				addCar(x, y, CD_UP)
			case 'v':
				dir = TO_UPDOWN
				addCar(x, y, CD_DOWN)
			case '<':
				dir = TO_LEFTRIGHT
				addCar(x, y, CD_LEFT)
			case '>':
				dir = TO_LEFTRIGHT
				addCar(x, y, CD_RIGHT)
			default:
				panic(fmt.Sprintf("unknown c %c at %v %v\n", c, x, y))
			}

			row = append(row, dir)
		}

		track = append(track, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("read failed: %v", err)
	}

	return &track, cars, nil
}

type ByLocation []*Car

func (a ByLocation) Len() int      { return len(a) }
func (a ByLocation) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByLocation) Less(i, j int) bool {
	if a[i].L.Y != a[j].L.Y {
		return a[i].L.Y < a[j].L.Y
	}
	return a[i].L.X < a[j].L.X
}

func turnLeft(cd CarDirection, loc Loc) (CarDirection, Loc) {
	switch cd {
	case CD_LEFT:
		return CD_DOWN, Loc{loc.X, loc.Y + 1}
	case CD_RIGHT:
		return CD_UP, Loc{loc.X, loc.Y - 1}
	case CD_UP:
		return CD_LEFT, Loc{loc.X - 1, loc.Y}
	case CD_DOWN:
		return CD_RIGHT, Loc{loc.X + 1, loc.Y}
	default:
		panic("bad cd")
	}
}

func turnRight(cd CarDirection, loc Loc) (CarDirection, Loc) {
	switch cd {
	case CD_LEFT:
		return CD_UP, Loc{loc.X, loc.Y - 1}
	case CD_RIGHT:
		return CD_DOWN, Loc{loc.X, loc.Y + 1}
	case CD_UP:
		return CD_RIGHT, Loc{loc.X + 1, loc.Y}
	case CD_DOWN:
		return CD_LEFT, Loc{loc.X - 1, loc.Y}
	default:
		panic("bad cd")
	}
}

func straightLoc(cd CarDirection, loc Loc) Loc {
	switch cd {
	case CD_LEFT:
		return Loc{loc.X - 1, loc.Y}
	case CD_RIGHT:
		return Loc{loc.X + 1, loc.Y}
	case CD_UP:
		return Loc{loc.X, loc.Y - 1}
	case CD_DOWN:
		return Loc{loc.X, loc.Y + 1}
	default:
		panic("bad cd")
	}
}

func advanceCar(track *Track, car Car) (CarDirection, bool, Loc) {
	to := (*track)[car.L.Y][car.L.X]
	var newDir CarDirection
	var newLoc Loc
	switch to {
	case TO_LEFTRIGHT:
		if car.D == CD_RIGHT {
			return car.D, false, Loc{car.L.X + 1, car.L.Y}
		} else if car.D == CD_LEFT {
			return car.D, false, Loc{car.L.X - 1, car.L.Y}
		} else {
			panic("bad cd in leftright")
		}
	case TO_UPDOWN:
		if car.D == CD_UP {
			return car.D, false, Loc{car.L.X, car.L.Y - 1}
		} else if car.D == CD_DOWN {
			return car.D, false, Loc{car.L.X, car.L.Y + 1}
		} else {
			panic("bad cd in updown")
		}
	case TO_INTERSECTION:
		var newDir CarDirection
		var newLoc Loc
		switch (car.NumTurns) % 3 {
		case 0:
			newDir, newLoc = turnLeft(car.D, car.L)
		case 1:
			newDir = car.D
			newLoc = straightLoc(car.D, car.L)
		case 2:
			newDir = CD_RIGHT
			newDir, newLoc = turnRight(car.D, car.L)
		}

		return newDir, true, newLoc

	case TO_UPLEFT:
		switch car.D {
		case CD_RIGHT:
			fallthrough
		case CD_LEFT:
			newDir, newLoc = turnRight(car.D, car.L)
		case CD_UP:
			fallthrough
		case CD_DOWN:
			newDir, newLoc = turnLeft(car.D, car.L)
		default:
			panic("bad cd in upleft")
		}
		return newDir, false, newLoc

	case TO_UPRIGHT:
		switch car.D {
		case CD_RIGHT:
			fallthrough
		case CD_LEFT:
			newDir, newLoc = turnLeft(car.D, car.L)
		case CD_UP:
			fallthrough
		case CD_DOWN:
			newDir, newLoc = turnRight(car.D, car.L)
		default:
			panic("bad cd in upright")
		}
		return newDir, false, newLoc

	default:
		panic(fmt.Sprintf("bad to %s", to))
	}

}

func advanceCars(track *Track, cars map[Loc]Car) (map[Loc]Car, bool) {
	newCars := map[Loc]Car{}
	for l, c := range cars {
		newCars[l] = c
	}

	newCrash := false
	for _, car := range cars {
		if car.CrashedWith != nil {
			continue
		}
		newDir, turned, newLoc := advanceCar(track, car)
		//fmt.Printf("nd %s nl %s\n", newDir, newLoc)

		delete(newCars, car.L)

		if hittee, found := newCars[newLoc]; found && hittee.Num != car.Num {
			// car hit another
			if hittee.CrashedWith == nil {
				hittee.CrashedWith = []int{}
			}
			hittee.CrashedWith = append(hittee.CrashedWith, car.Num)
			newCars[newLoc] = hittee
			newCrash = true
		} else {
			car.D = newDir
			car.L = newLoc
			if turned {
				car.NumTurns++
			}
			newCars[newLoc] = car
		}
	}

	return newCars, newCrash
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	track, cars, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	//track.Dump(cars)
	for t := 1; ; t++ {
		if *numTicks != -1 && t > *numTicks {
			fmt.Printf("out of ticks")
			break
		}

		newCars, didCrash := advanceCars(track, cars)
		//fmt.Printf("\ntick = %d\n", t)
		//track.Dump(newCars)
		cars = newCars

		if didCrash {
			fmt.Println("crash")
			break
		}
	}

	for _, c := range cars {
		if c.CrashedWith == nil {
			continue
		}
		fmt.Println(c.L)
	}
}
