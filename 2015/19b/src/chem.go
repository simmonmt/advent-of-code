package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"

	"chem"
)

type NextPolicy interface {
	PickNext(cands []*chem.Mapping, dict *chem.Dict) int
}

type LongestNextPolicy struct{}

func (p *LongestNextPolicy) PickNext(cands []*chem.Mapping, dict *chem.Dict) int {
	longestIdx := -1
	var longestMapping *chem.Mapping

	for i, cand := range cands {
		fmt.Printf("cand %d: %v\n", i, cand.ToString(dict))

		if longestMapping == nil || len(cand.From) > len(longestMapping.From) {
			longestIdx = i
			longestMapping = cand
		}
	}

	fmt.Printf("#cands %v res %v\n", len(cands), longestIdx)

	return longestIdx
}

type RandNextPolicy struct{}

func (p *RandNextPolicy) PickNext(cands []*chem.Mapping, dict *chem.Dict) int {
	return rand.Intn(len(cands))
}

func readInput(r io.Reader, d *chem.Dict) (*chem.Mappings, []byte, error) {
	reader := bufio.NewReader(r)

	mappings := chem.NewMappings()
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		pairs := strings.SplitN(line, " => ", 2)
		from := chem.ParseMolecule(pairs[0], d)
		to := chem.ParseMolecule(pairs[1], d)
		mappings.Add(to, from)
	}

	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, nil, fmt.Errorf("no input molecule")
	}
	line = strings.TrimSpace(line)

	molecule := chem.ParseMolecule(line, d)

	return mappings, molecule, nil
}

func replace(molecule []byte, start int, mapping *chem.Mapping, dict *chem.Dict) []byte {
	newLen := len(molecule) - len(mapping.From) + len(mapping.To)
	out := make([]byte, newLen)

	moleculeCur := 0
	outCur := 0
	if start > 0 {
		copy(out, molecule[0:start])
		moleculeCur = start
		outCur = start
	}

	copy(out[outCur:], mapping.To)
	moleculeCur += len(mapping.From)
	outCur += len(mapping.To)

	if moleculeCur < len(molecule) {
		copy(out[outCur:], molecule[moleculeCur:])
	}

	// fmt.Printf("applied %v at %d\n", mapping.ToString(dict), start)
	// fmt.Printf("was: %v\n", chem.MoleculeToString(molecule, dict))
	// fmt.Printf("is : %v\n", chem.MoleculeToString(out, dict))

	return out
}

func reduce(molecule []byte, mappings *chem.Mappings, dict *chem.Dict, nextPolicy NextPolicy) []byte {
	allFoundMappings := []*chem.Mapping{}
	allFoundMappingIdxes := []int{}

	for i := 0; i < len(molecule); i++ {
		foundMappings := mappings.Find(molecule[i:], dict)
		if len(foundMappings) == 0 {
			i++
			continue
		} else if len(foundMappings) != 1 {
			panic(fmt.Sprintf("found >1 mappings at %v",
				chem.MoleculeToString(molecule[i:], dict)))
		}
		allFoundMappings = append(allFoundMappings, foundMappings[0])
		allFoundMappingIdxes = append(allFoundMappingIdxes, i)
	}

	if len(allFoundMappings) == 0 {
		return nil
	}

	chosenIdx := nextPolicy.PickNext(allFoundMappings, dict)
	return replace(molecule, allFoundMappingIdxes[chosenIdx], allFoundMappings[chosenIdx], dict)
}

func doSearch(molecule []byte, mappings *chem.Mappings, dict *chem.Dict, nextPolicy NextPolicy, finishByte byte) (int, bool) {
	for i := 1; ; i++ {
		replacement := reduce(molecule, mappings, dict, nextPolicy)
		if replacement == nil {
			return i, false
		}

		if len(replacement) == 1 && replacement[0] == finishByte {
			return i, true
		}

		//fmt.Printf("round %d: len=%d\n", i, len(replacement))
		molecule = replacement
	}
}

func main() {
	dict := chem.NewDict()

	// reduce finds all candidate mappings in a given molecule. Choose the
	// next one at random. Choosing longest doesn't work for whatever
	// reason, as it doesn't converge. Oddly if we choose at random we tend
	// to get the shortest-length convergence almost immediately.
	nextPolicy := &RandNextPolicy{}

	mappings, initial, err := readInput(os.Stdin, dict)
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	// mappings.Dump(dict)
	// fmt.Println(initial)

	eByte := dict.StrToByte("e")

	// We're picking candidates at random (see above), so we need to keep
	// trying again and again.
	shortest := -1
	for i := 0; ; i++ {
		if i != 0 && i%1000 == 0 {
			fmt.Printf("iter %v\n", i)
		}

		numRounds, found := doSearch(initial, mappings, dict, nextPolicy, eByte)
		if !found {
			continue
		}

		if shortest == -1 {
			if numRounds < shortest {
				shortest = numRounds
				fmt.Printf("new shortest %v at iter %v\n", shortest, i)
			} else {
				fmt.Printf("converged at %v at iter %v\n", numRounds, i)
			}
		}
	}
}
