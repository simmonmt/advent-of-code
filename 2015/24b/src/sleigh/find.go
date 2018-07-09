package sleigh

import "fmt"

// for groupOneSize = 1 to len-2
//   found = []
//
//   for each group(groupOneSize) groupOne
//     if avg(rest) / 2 != sum(groupOne)
//       continue
//     sideCap = avg(rest) / 2
//     if there's a packing p whose size == sidecap
//     append (groupOne, p (groupTwo), rest - p (groupThree)) to found
//
//   if found not empty
//     emit one with smallest entanglement

// object that emits all subgroups (order doesn't matter) of given
// size

// function that does bin-packing

func sumArr(vals []int) int {
	out := 0
	for _, val := range vals {
		out += val
	}
	return out
}

func removeElems(allVals, toRemove []int) []int {
	toRemoveMap := map[int]bool{}
	for _, v := range toRemove {
		toRemoveMap[v] = true
	}

	out := []int{}
	for _, v := range allVals {
		if toRemoveMap[v] {
			toRemoveMap[v] = false
			continue
		}
		out = append(out, v)
	}

	return out
}

func FindWithGroupOneSize(values []int, groupOneSize int) [][]int {
	found := [][]int{}

	sg := NewSubgrouper(values, groupOneSize)
	for {
		groupOne, nonGroupOne, ok := sg.Next()
		if !ok {
			break
		}

		cap := sumArr(groupOne)
		if cap*3 != sumArr(nonGroupOne) {
			continue
		}

		foundGroupTwos := AllBinPacks(nonGroupOne, cap)
		fmt.Printf("found %v possible group two packings for %v\n", len(foundGroupTwos), groupOne)
		for _, foundGroupTwo := range foundGroupTwos {
			rest := removeElems(nonGroupOne, foundGroupTwo) // not group 1 or 2
			if OneBinPack(rest, cap) != nil {
				// TODO: de-dup in found
				found = append(found, groupOne)
			}
		}
	}

	return found
}
