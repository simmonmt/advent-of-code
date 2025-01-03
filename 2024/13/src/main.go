package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime/pprof"
	"strconv"

	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/logger"
	"github.com/simmonmt/aoc/2024/common/pos"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

	buttonPattern = regexp.MustCompile(`^Button ([AB]): X\+([0-9]+), Y\+([0-9]+)$`)
	prizePattern  = regexp.MustCompile(`^Prize: X=([0-9]+), Y=([0-9]+)$`)
)

type Machine struct {
	A, B  pos.P2
	Prize pos.P2
}

func parseButton(line string, want string) (pos.P2, error) {
	parts := buttonPattern.FindStringSubmatch(line)
	if len(parts) != 4 {
		return pos.P2{}, fmt.Errorf("bad match")
	}

	if parts[1] != want {
		return pos.P2{}, fmt.Errorf("bad button %v", parts[1])
	}

	var p pos.P2
	var err error
	if p.X, err = strconv.Atoi(parts[2]); err != nil {
		return pos.P2{}, fmt.Errorf("bad X: %v", err)
	}
	if p.Y, err = strconv.Atoi(parts[3]); err != nil {
		return pos.P2{}, fmt.Errorf("bad Y: %v", err)
	}

	return p, nil
}

func parsePrize(line string) (pos.P2, error) {
	parts := prizePattern.FindStringSubmatch(line)
	if len(parts) != 3 {
		return pos.P2{}, fmt.Errorf("bad match")
	}

	var p pos.P2
	var err error
	if p.X, err = strconv.Atoi(parts[1]); err != nil {
		return pos.P2{}, fmt.Errorf("bad X: %v", err)
	}
	if p.Y, err = strconv.Atoi(parts[2]); err != nil {
		return pos.P2{}, fmt.Errorf("bad Y: %v", err)
	}

	return p, nil
}

func parseInput(lines []string) ([]*Machine, error) {
	groups := filereader.BlankSeparatedGroupsFromLines(lines)
	machines := make([]*Machine, len(groups))

	for i, group := range groups {
		if len(group) != 3 {
			return nil, fmt.Errorf("bad group %d - not 3 lines", i+1)
		}

		var err error
		machine := &Machine{}
		if machine.A, err = parseButton(group[0], "A"); err != nil {
			return nil, fmt.Errorf("group %d: bad A: %v", i+1, err)
		}
		if machine.B, err = parseButton(group[1], "B"); err != nil {
			return nil, fmt.Errorf("group %d: bad B: %v", i+1, err)
		}
		if machine.Prize, err = parsePrize(group[2]); err != nil {
			return nil, fmt.Errorf("group %d: bad prize: %v", i+1, err)
		}
		machines[i] = machine
	}

	return machines, nil
}

func expensiveSolveMachine(machine *Machine) int {
	var steps int
	var p pos.P2
	for steps = 1; ; steps++ {
		if steps > 10000 {
			panic("overflow")
		}

		p.Add(machine.B)
		if p.X >= machine.Prize.X || p.Y >= machine.Prize.Y {
			break
		}
	}

	minCost := -1
	for numA := 100; numA >= 0; numA-- {
		for numB := 100; numB >= 0; numB-- {
			p := pos.P2{
				X: machine.A.X*numA + machine.B.X*numB,
				Y: machine.A.Y*numA + machine.B.Y*numB,
			}
			if !p.Equals(machine.Prize) {
				continue
			}
			cost := numA*3 + numB
			if minCost == -1 || cost < minCost {
				minCost = cost
			}
		}
	}
	return minCost
}

func isSolution(machine *Machine, numA, numB, lim int) bool {
	if lim > 0 && (numA > lim || numB > lim) {
		return false
	}

	p := pos.P2{
		X: machine.A.X*numA + machine.B.X*numB,
		Y: machine.A.Y*numA + machine.B.Y*numB,
	}
	return p.Equals(machine.Prize)
}

func findDimSeq(start, want int, calc func(n int) int) []int {
	out := []int{}

	seen := map[int]bool{}
	for i := start; len(out) < want; i++ {
		mod := calc(i)
		if mod == 0 {
			out = append(out, i)
			seen = map[int]bool{}
		} else {
			if _, found := seen[mod]; found {
				return out
			}
			seen[mod] = true
		}
	}

	return out
}

func cheapSolveMachine(machine *Machine, lim int) int {
	maxB := machine.Prize.X / machine.B.X
	maxB = min(maxB, machine.Prize.Y/machine.B.Y)
	if lim > 0 {
		maxB = min(maxB, lim)
	}

	// Find the first two solutions for X and Y -- b values that have
	// corresponding integer a values.
	logger.Infof("find bXSeq, bYSeq")
	bXSeq := findDimSeq(0, 2, func(n int) int {
		return (machine.Prize.X - machine.B.X*n) % machine.A.X
	})
	bYSeq := findDimSeq(0, 2, func(n int) int {
		return (machine.Prize.Y - machine.B.Y*n) % machine.A.Y
	})

	// Handle cases where no solutions exist
	if len(bXSeq) == 0 || len(bYSeq) == 0 {
		return -1 // Couldn't build sequence for one of them
	} else if len(bXSeq) == 1 || len(bYSeq) == 1 {
		panic("not infinite")
	}

	// Both sequences have 2 elements
	//
	// Find values of b that have integer a solutions for both X and Y by
	// walking the sequences from the previous step until we find elements
	// that are present in both sequences.
	//
	// First we make sure there will be such a solution. We know that the
	// solution will look like
	//
	//   x0 + xd * I = y0 + yd * J
	//
	// Where x0 is the first element in bXSeq and xd is the delta between
	// the two elements in bXSeq. I and J are unknown integers. We can
	// detect sequences that don't have common values by checking to see if
	// there are integer I and J. That is, by seeing if there are solutions
	// to
	//
	//   (x0 - y0 + xd * I) % yd = 0
	//
	// We can do this by iterating through I=[0,yd) but I'm not sure why.
	// https://libraryguides.centennialcollege.ca/c.php?g=717548&p=5121955
	// says I can. If none of those values of I satisfy the equation,
	// there's no solution, so we can deem the machine unsolveable.
	logger.Infof("find common from x %v y %v", bXSeq, bYSeq)

	deltaBX, deltaBY := bXSeq[1]-bXSeq[0], bYSeq[1]-bYSeq[0]

	found := false
	for i := range deltaBY {
		if (bXSeq[0]-bYSeq[0]+deltaBX*i)%deltaBY == 0 {
			found = true
			break
		}
	}
	if !found {
		return -1 // no solution
	}

	// We know overlaps exist, so go find the first two.
	bSeq := []int{}
	for bx, by := bXSeq[0], bYSeq[0]; len(bSeq) < 2; {
		if bx == by {
			bSeq = append(bSeq, bx)
			bx += deltaBX
			by += deltaBY
		} else if bx < by {
			bx += deltaBX
		} else {
			by += deltaBY
		}
	}
	logger.Infof("common %v", bSeq)

	// We know bSeq has two elements
	if lim > 0 && bSeq[0] > lim {
		return -1 // First common is beyond limit
	}

	// bSeq contains the first two elements of a sequence of b values that
	// have corresponding a values in X and Y. Unfortunately we haven't
	// checked to see if the a values are the same. That is, for a given b
	// value, the a value that solves X may not be the same as the a value
	// that solves Y. We therefore need to iterate through the bSeq sequence
	// to find b values whose X and Y a values are the same. We end up with
	// the first two b values of the a sequence whose values are solutions
	// to the machine.
	//
	// It may be possible to combine this search with the previous one, but
	// keeping them separate (I think) makes the logic simpler.

	solSeq := []int{}
	lastAX, lastAY := -1, -1
	b := bSeq[0]
	inc := bSeq[1] - bSeq[0]
	for len(solSeq) < 2 {
		ax := (machine.Prize.X - machine.B.X*b) / machine.A.X
		ay := (machine.Prize.Y - machine.B.Y*b) / machine.A.Y
		if ax < 0 || ay < 0 {
			break // b overshot the Prize coordinates, so stop looking
		}
		if ax == ay {
			// We found a solution
			logger.Infof("equal at %v", b)
			solSeq = append(solSeq, b)
			b += inc
		} else if lastAX != lastAY {
			// We've seen two non-solutions in a row. They should be
			// converging on a common point - a value of b that has
			// a single a value that solves both X and Y. The rate
			// of convergence is constant, so we can fast-forward to
			// the convergence point (which could be quite a ways
			// away). If the convergence point is non-integral,
			// we'll never reach it (we'll skip over it), so we give
			// up searching. We'll also never reach it if the number
			// of steps is negative, so we give up then too.
			axDelta, ayDelta := ax-lastAX, ay-lastAY

			n := ay - ax
			d := axDelta - ayDelta
			if n%d != 0 {
				break
			}
			steps := n / d

			logger.Infof("b %v last %v,%v now %v,%v delta %v,%v steps %v",
				b, lastAX, lastAY, ax, ay, axDelta, ayDelta, steps)
			if steps < 0 {
				break
			}
			b += inc * steps
		} else {
			b += inc
		}

		lastAX, lastAY = ax, ay
	}

	logger.Infof("valid %v", solSeq)

	// Calculate the final result from the various possible forms of
	numA, numB := 0, 0
	if sz := len(solSeq); sz == 0 {
		return -1 // no sequence
	} else if lim > 0 && solSeq[0] > lim {
		return -1 // First common is beyond limit
	} else if sz == 1 {
		numB = solSeq[0]
		numA = (machine.Prize.X - numB*machine.B.X) / machine.A.X
	} else {
		// Find the largest b value in the sequence that's less than the
		// Prize's X coordinate, and calculate the a value from
		// that. Note that we can do our math entirely with the X
		// coordinates since the previous step made sure that solSeq
		// only contains values that work for X and Y.
		v := machine.Prize.X - (solSeq[0] * machine.B.X)
		db := solSeq[1] - solSeq[0]
		d := db * machine.B.X

		if lim > 0 && v+d > lim {
			numB = solSeq[0]
			numA = v / machine.A.X
		} else {
			numB = solSeq[0] + (v/d)*db
			numA = (machine.Prize.X - numB*machine.B.X) / machine.A.X
		}
	}

	// Sanity check
	if !isSolution(machine, numA, numB, lim) {
		panic(fmt.Sprintf("mismatch machine %v A %v B %v x %v y %v solSeq %v",
			*machine, numA, numB, bXSeq, bYSeq, solSeq))
	}
	return numA*3 + numB
}

func solveA(machines []*Machine) int {
	sum := 0
	for _, machine := range machines {
		exp := expensiveSolveMachine(machine)
		chp := cheapSolveMachine(machine, 100)

		if exp != chp {
			logger.Fatalf("disagreement machine %v exp %v chp %v", machine, exp, chp)
		}
		if chp > 0 {
			sum += chp
		}
	}

	return sum
}

func solveB(machines []*Machine) int {
	sum := 0
	for i, mp := range machines {
		machine := *mp
		machine.Prize.X += 10000000000000
		machine.Prize.Y += 10000000000000
		logger.Infof("B %d %v", i, machine)
		if cost := cheapSolveMachine(&machine, -1); cost > 0 {
			logger.Infof("solution")
			sum += cost
		}
	}
	return sum
	//return -1
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *input == "" {
		logger.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		logger.Fatalf("failed to read input: %v", err)
	}

	input, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
