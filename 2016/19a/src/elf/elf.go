package elf

import (
	"container/list"
	"fmt"

	"logger"
)

type Elf struct {
	Name     uint
	Presents uint
}

func InitElves(num int) *list.List {
	elves := list.New()
	for i := 1; i <= num; i++ {
		elves.PushBack(&Elf{Name: uint(i), Presents: uint(1)})
	}
	return elves
}

func Print(elves *list.List) {
	for elem := elves.Front(); elem != nil; elem = elem.Next() {
		e := elem.Value.(*Elf)
		fmt.Printf("  %d %d\n", e.Name, e.Presents)
	}
}

func Play(elves *list.List) uint {
	for round := 1; elves.Len() > 1; round++ {
		logger.LogF("round %v, elves: %v\n", round, elves.Len())

		elem := elves.Front()
		for elem != nil {
			nElem := elem.Next()
			if nElem == nil {
				nElem = elves.Front()
				if nElem == nil {
					break
				}
			}

			e := elem.Value.(*Elf)
			ne := nElem.Value.(*Elf)

			e.Presents += ne.Presents
			elves.Remove(nElem)
			elem = elem.Next()
		}

		if logger.Enabled() {
			Print(elves)
		}
	}

	return elves.Front().Value.(*Elf).Name
}
