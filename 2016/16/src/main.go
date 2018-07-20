package main

import (
	"flag"
	"fmt"
	"log"

	"data"
	"util"
)

var (
	initialState = flag.String("initial_state", "", "initial state")
	goalLength   = flag.Int("goal_length", -1, "goal length")
)

func main() {
	flag.Parse()

	if *initialState == "" {
		log.Fatal("--initial_state is required")
	}
	if *goalLength == -1 {
		log.Fatal("--goal_length is required")
	}

	d := util.StrToBoolArray(*initialState)
	for len(d) < *goalLength {
		d = data.Grow(d)
	}

	d = d[0:*goalLength]
	cksum := data.Checksum(d)

	fmt.Println(util.BoolArrayToStr(cksum))
}
