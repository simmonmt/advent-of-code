package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/simmonmt/aoc/2019/17/a/src/puzzle"
	"github.com/simmonmt/aoc/2019/common/dir"
	"github.com/simmonmt/aoc/2019/common/logger"
	"github.com/simmonmt/aoc/2019/common/pos"
	"github.com/simmonmt/aoc/2019/common/vm"
)

var (
	verbose     = flag.Bool("verbose", false, "verbose")
	ramPath     = flag.String("ram", "", "path to file containing ram values")
	programPath = flag.String("program", "", "path to optional file containing program")
)

type Program struct {
	Main, A, B, C, Continuous string
}

func readProgram(path string) (*Program, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := []string{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	if len(lines) != 5 {
		return nil, fmt.Errorf("bad program")
	}

	return &Program{
		Main:       lines[0],
		A:          lines[1],
		B:          lines[2],
		C:          lines[3],
		Continuous: lines[4],
	}, nil
}

func sendCommand(ch chan *vm.ChanIOMessage, str string) {
	fmt.Printf("sending %v\n", str)
	for _, r := range str {
		ch <- &vm.ChanIOMessage{Val: int64(r)}
	}
	ch <- &vm.ChanIOMessage{Val: 10}
}

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

	var program *Program
	if *programPath != "" {
		program, err = readProgram(*programPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Program: %+v\n", program)

		ram.Write(0, 2)
	}

	async := vm.RunAsync("vm", ram)

	if program != nil {
		go func() {
			sendCommand(async.In, program.Main)
			sendCommand(async.In, program.A)
			sendCommand(async.In, program.B)
			sendCommand(async.In, program.C)
			sendCommand(async.In, program.Continuous)
		}()
	}

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

		if msg.Val > 255 {
			fmt.Printf("large value %v\n", msg.Val)
			continue
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
			//fmt.Printf("read unexpected %v %c\n", msg.Val, r)
		}
	}

	intPs := puzzle.FindIntersections(b)
	align := 0
	for intP := range intPs {
		align += intP.X * intP.Y
	}
	fmt.Printf("sum alignment %d\n", align)

	puzzle.DumpBoard(b, &vac, nil)
	fmt.Printf("vac %v\n", vac)

}
