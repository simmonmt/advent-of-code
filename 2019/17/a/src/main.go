package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/simmonmt/aoc/2019/17/a/src/puzzle"
	"github.com/simmonmt/aoc/2019/common/dir"
	"github.com/simmonmt/aoc/2019/common/logger"
	"github.com/simmonmt/aoc/2019/common/pos"
	"github.com/simmonmt/aoc/2019/common/vm"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	ramPath = flag.String("ram", "", "path to file containing ram values")
)

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *ramPath == "" {
		log.Fatalf("--ram is required")
	}

	ram, err := vm.NewRamFromFile(*ramPath)
	if err != nil {
		log.Fatal(err)
	}

	async := vm.RunAsync("vm", ram)

	vac := puzzle.Vac{Pos: pos.P2{0, 0}, Dir: puzzle.VacDir(dir.DIR_NORTH)}
	p := pos.P2{0, 0}
	b := puzzle.NewBoard()
	for {
		msg, ok := <-async.Out
		if !ok {
			break
		}
		if msg.Err != nil {
			panic(fmt.Sprintf("async out err: %v", msg.Err))
		}

		fmt.Printf("%c", rune(msg.Val))

		switch r := rune(msg.Val); r {
		case '#':
			fallthrough
		case '.':
			b.Set(p, r)
			p.X++
			break
		case '^':
			fallthrough
		case 'v':
			fallthrough
		case '<':
			fallthrough
		case '>':
			b.Set(p, '#')
			vac.Pos = p
			vac.Dir = puzzle.ParseVacDir(r)
			p.X++
			break

		case 10:
			p.X = 0
			p.Y++
			break
		default:
			panic(fmt.Sprintf("bad val %d %c\n", r, r))
		}
	}

	intPs := puzzle.FindIntersections(b)
	align := 0
	for intP := range intPs {
		align += intP.X * intP.Y
	}
	fmt.Printf("sum alignment %d\n", align)

	puzzle.DumpBoard(b, &vac, intPs)
	fmt.Printf("vac %v\n", vac)

}
