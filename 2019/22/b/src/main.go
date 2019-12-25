package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/simmonmt/aoc/2019/22/b/src/puzzle"
	"github.com/simmonmt/aoc/2019/common/logger"
)

var (
	verbose  = flag.Bool("verbose", false, "verbose")
	input    = flag.String("input", "", "input file")
	numCards = flag.Int64("num_cards", 119315717514047, "number of cards")
	maxRuns  = flag.Int("max_runs", -1, "max runs")
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

func findValue(cards []int, want int) int {
	for i, card := range cards {
		if card == want {
			return i
		}
	}
	return -1
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	cmds, err := puzzle.ParseCommands(lines)
	if err != nil {
		log.Fatal(err)
	}

	cards := make([]int, 10007)
	for i := range cards {
		cards[i] = i
	}

	cards = puzzle.RunCommands(cards, cmds)

	fmt.Printf("card 2019 location: %v\n", findValue(cards, 2019))

	// Found via experimentation. I took a bunch of primes (10007
	// is prime) and looked for when they started to repeat. Note
	// that in some cases they repeat sooner, but there's always a
	// restart at numCards-2. This of course defeats any attempts
	// to find the cycle start by brute force for large numCards.
	// cycleStart := int64(*numCards - 1)

	// wantShuffles := int64(101741582076661)
	// needShuffles := wantShuffles % cycleStart

	// fmt.Printf("cycle start   %v\n", cycleStart)
	// fmt.Printf("want shuffles %v\n", wantShuffles)
	// fmt.Printf("need shuffles %v\n", needShuffles)

	// return

	//numCards := int64(119315717514047)

	// cards := make([]int, *numCards)
	// for i := range cards {
	// 	cards[i] = i
	// }

	// fmt.Println("forward")
	// n := int64(2020)
	// for i := 1; *maxRuns < 0 || i <= *maxRuns; i++ {
	// 	// cards = puzzle.RunCommands(cards, cmds)
	// 	// fmt.Printf("%v\n", cards[0:10])

	// 	n = puzzle.ForwardCommandsForIndex(cmds, int64(*numCards), n)

	// 	fmt.Println(n)
	// 	// fmt.Printf("now i=%d %d value %d\n", i, n, cards[n])

	// 	if n == 0 {
	// 		fmt.Printf("repeat at %d (%d)\n", i, int(*numCards)/i)
	// 		return
	// 	}
	// 	if i%100000 == 0 {
	// 		fmt.Printf("i=%d n=%v\n", i, n)
	// 	}
	// 	// if when, found := fwdCache[n]; found {
	// 	// 	fmt.Printf("fwd %d repeat from %d\n", i, when)
	// 	// 	break
	// 	// } else {
	// 	// 	fwdCache[n] = i
	// 	// }
	// }

	// fmt.Println("reverse")
	// n := 2020
	// for i := 1; *maxRuns == -1 || i < *maxRuns; i++ {
	// 	n = puzzle.ReverseCommandsForIndex(cmds, int(*numCards), n)
	// }
	// fmt.Printf("n=%d\n", n)

	// fmt.Println("done")

	// fmt.Printf("at pos 2020 is %v\n", cards[2020])
}
