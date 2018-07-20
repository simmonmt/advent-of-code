package main

import (
	"flag"
	"fmt"
	"log"

	"maze"
)

var (
	passcode = flag.String("passcode", "", "passcode")
	width    = flag.Int("width", 4, "width")
	height   = flag.Int("height", 4, "height")
)

func main() {
	flag.Parse()

	if *passcode == "" {
		log.Fatal("--passcode is required")
	}

	found, path := maze.RunMaze(*width, *height, *passcode)
	if !found {
		fmt.Println("no path found")
	} else {
		fmt.Printf("path: %d steps %v\n", len(path), path)
	}
}
