package main

import (
	"flag"
	"fmt"

	"board"
	"game"
	"logger"
	"object"
)

var (
	audit      = flag.Bool("audit", false, "audit results")
	logging    = flag.Bool("verbose", false, "enable logging")
	smallInput = flag.Bool("small_input", false, "use small input")
	dumpSeen   = flag.Bool("dump_seen", false, "dump seen after completion")
)

func main() {
	flag.Parse()
	logger.Init(*logging)

	var b *board.Board

	// Small input.
	if *smallInput {
		b = board.New(map[object.Object]uint8{
			object.Microchip(1): 1, // hydrogen
			object.Generator(1): 2, // hydrogen
			object.Microchip(2): 1, // lithium
			object.Generator(2): 3, // lithium
		})
	} else {
		// Contest input.
		b = board.New(map[object.Object]uint8{
			object.Microchip(1): 1, // promethium
			object.Generator(1): 1, // promethium
			object.Generator(2): 2, // cobalt
			object.Generator(3): 2, // curium
			object.Generator(4): 2, // ruthenium
			object.Generator(5): 2, // plutonium
			object.Microchip(2): 3, // cobalt
			object.Microchip(3): 3, // curium
			object.Microchip(4): 3, // ruthenium
			object.Microchip(5): 3, // plutonium
		})
	}

	minMoves, seen := game.Play(b)

	if minMoves == nil {
		fmt.Println("no solutions found")
	} else {
		fmt.Printf("minMoves %d\n", len(minMoves))
	}

	if *audit {
		game.Audit(b, minMoves)
	}

	if *dumpSeen {
		for s, v := range seen {
			fmt.Printf("seen %v = %+v\n", s, v)
		}
	}
}
