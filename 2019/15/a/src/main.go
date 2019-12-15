package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/simmonmt/aoc/2019/15/a/src/puzzle"
	"github.com/simmonmt/aoc/2019/common/logger"
	"github.com/simmonmt/aoc/2019/common/pos"
	"github.com/simmonmt/aoc/2019/common/vm"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	ramPath = flag.String("ram", "", "path to file containing ram values")
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

func moveTo(b *puzzle.Board, async *vm.Async, pos pos.P2, dir puzzle.Dir) (newPos pos.P2, t puzzle.Tile) {
	async.In <- &vm.ChanIOMessage{Val: int64(dir)}
	resp := <-async.Out

	if resp.Err != nil {
		panic(fmt.Sprintf("error from vm: %v", resp.Err))
	}

	switch resp.Val {
	case 0:
		return pos, puzzle.TILE_WALL
	case 1:
		return dir.From(pos), puzzle.TILE_OPEN
	case 2:
		return dir.From(pos), puzzle.TILE_GOAL
	default:
		panic(fmt.Sprintf("bad resp %d", resp.Val))
	}

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

	b := puzzle.NewBoard()
	start := pos.P2{0, 0}
	b.Set(start, puzzle.TILE_OPEN)

	async := vm.RunAsync("vm", ram)

	puzzle.Explore(b, start, func(pos pos.P2, dir puzzle.Dir) (newPos pos.P2, t puzzle.Tile) {
		return moveTo(b, async, pos, dir)
	})

	puzzle.PrintBoard(b, pos.P2{0, 0})
}
