package main

import (
	"flag"
	"fmt"
	"log"

	"board"
	"game"
	"logger"
	"object"
)

var (
	logging  = flag.Bool("verbose", false, "enable logging")
	inputSet = flag.String("input_set", "", "input set to use -- small, a, or b")
)

func main() {
	flag.Parse()
	logger.Init(*logging)

	var b *board.Board

	// Small input.
	switch *inputSet {
	case "small":
		b = board.New(map[object.Object]uint8{
			object.Microchip(1): 1, // hydrogen
			object.Generator(1): 2, // hydrogen
			object.Microchip(2): 1, // lithium
			object.Generator(2): 3, // lithium
		})
		break

	case "a":
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
		break

	case "b":
		// Contest input.
		b = board.New(map[object.Object]uint8{
			object.Microchip(1): 1, // promethium
			object.Generator(1): 1, // promethium
			object.Microchip(6): 1, // elerium
			object.Generator(6): 1, // elerium
			object.Microchip(7): 1, // dilithium
			object.Generator(7): 1, // dilithium
			object.Generator(2): 2, // cobalt
			object.Generator(3): 2, // curium
			object.Generator(4): 2, // ruthenium
			object.Generator(5): 2, // plutonium
			object.Microchip(2): 3, // cobalt
			object.Microchip(3): 3, // curium
			object.Microchip(4): 3, // ruthenium
			object.Microchip(5): 3, // plutonium
		})
		break

	default:
		log.Fatalf("unknown input set \"%v\"", *inputSet)
	}

	minMoves := game.Play(b)

	if minMoves == nil {
		fmt.Println("no solutions found")
	} else {
		fmt.Printf("minMoves %d\n", len(minMoves)-1)
	}
}
