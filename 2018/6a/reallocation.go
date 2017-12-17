package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

type Tracker struct {
	configs map[string]bool
}

func NewTracker() *Tracker {
	return &Tracker{configs: map[string]bool{}}
}

func (t *Tracker) Insert(config []int) bool {
	str := ""
	for _, c := range config {
		str += strconv.Itoa(c) + ","
	}

	if _, found := t.configs[str]; found {
		return true
	}

	t.configs[str] = true
	return false
}

func (t *Tracker) Dump() {
	all := []string{}

	for config, _ := range t.configs {
		all = append(all, config)
	}
	sort.Strings(all)

	for _, config := range all {
		fmt.Println(config)
	}
}

func findMaxIndex(vals []int) int {
	maxVal := vals[0]
	maxIndex := 0

	for i, val := range vals {
		if val > maxVal {
			maxVal = val
			maxIndex = i
		}
	}

	return maxIndex
}

func main() {
	if len(os.Args) == 1 {
		log.Fatalf("Usage: %v bank [bank...]", os.Args[0])
	}

	banks := []int{}
	for _, str := range os.Args[1:] {
		bank, err := strconv.Atoi(str)
		if err != nil {
			log.Fatalf("failed to parse bank %v", str)
		}

		banks = append(banks, bank)
	}

	tracker := NewTracker()
	for cycle := 1; ; cycle++ {
		//fmt.Printf("start: %v\n", banks)

		maxIndex := findMaxIndex(banks)
		toDist := banks[maxIndex]
		//fmt.Printf("max in %d, restributing %d\n", maxIndex, toDist)

		banks[maxIndex] = 0
		distIndex := maxIndex + 1
		for toDist > 0 {
			banks[distIndex%len(banks)]++
			distIndex++
			toDist--
		}

		//fmt.Printf("finish: %v\n\n", banks)

		if tracker.Insert(banks) {
			fmt.Println(cycle)
			break
		}
	}
}
