package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	linePattern = regexp.MustCompile(`^([^\(]+) \(contains ([^\)]+)\)$`)
)

type Ing string
type Gen string

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		log.Fatal(err)
	}

	allIngs := map[Ing]bool{}
	genContainers := map[Gen][][]Ing{}
	numFoodsByIng := map[Ing]int{}
	for _, line := range lines {
		parts := linePattern.FindStringSubmatch(line)
		if parts == nil {
			log.Fatalf("bad line: %v", line)
		}

		ings := []Ing{}
		for _, ing := range strings.Split(parts[1], " ") {
			ings = append(ings, Ing(ing))
		}

		gens := []Gen{}
		for _, gen := range strings.Split(parts[2], ", ") {
			gens = append(gens, Gen(gen))
		}

		for _, gen := range gens {
			if _, found := genContainers[gen]; !found {
				genContainers[gen] = [][]Ing{}
			}
			genContainers[gen] = append(genContainers[gen], ings)
		}

		for _, ing := range ings {
			numFoodsByIng[ing]++
			allIngs[ing] = true
		}
	}

	genCands := map[Gen][]Ing{}
	for gen, ingLists := range genContainers {
		var commonIngs []Ing
		if len(ingLists) == 1 {
			commonIngs = ingLists[0]
		} else {
			ingCounts := map[Ing]int{}
			for _, ingList := range ingLists {
				for _, ing := range ingList {
					ingCounts[ing]++
				}
			}

			commonIngs = []Ing{}
			for ing, num := range ingCounts {
				if num == len(ingLists) {
					commonIngs = append(commonIngs, ing)
				}
			}
		}

		logger.LogF("gen %v common %v of %v", gen, len(commonIngs), len(allIngs))

		genCands[gen] = commonIngs
	}
	logger.LogF("genCands %v", genCands)

	maybeGenIngs := map[Ing]bool{}
	for _, ings := range genCands {
		for _, ing := range ings {
			maybeGenIngs[ing] = true
		}
	}

	nonGenIngs := []Ing{}
	for ing := range allIngs {
		if _, found := maybeGenIngs[ing]; !found {
			nonGenIngs = append(nonGenIngs, ing)
		}
	}
	logger.LogF("non gen ings: %v", nonGenIngs)

	numFoods := 0
	for _, ing := range nonGenIngs {
		if num, found := numFoodsByIng[ing]; found {
			numFoods += num
		}
	}

	fmt.Printf("A: %d\n", numFoods)

	results := map[Gen]Ing{}
	for {
		var toRemove Ing
		for gen, ings := range genCands {
			if len(ings) == 1 {
				results[gen] = ings[0]
				toRemove = ings[0]
				break
			}
		}

		if toRemove == "" {
			break
		}

		for gen, ings := range genCands {
			newIngs := []Ing{}
			for _, ing := range ings {
				if ing != toRemove {
					newIngs = append(newIngs, ing)
				}
			}
			genCands[gen] = newIngs
		}
	}

	allGens := []Gen{}
	for gen := range results {
		allGens = append(allGens, gen)
	}
	sort.Slice(allGens, func(i, j int) bool { return allGens[i] < allGens[j] })

	out := []string{}
	for _, gen := range allGens {
		out = append(out, string(results[gen]))
	}
	fmt.Printf("B: %v\n", strings.Join(out, ","))
}
