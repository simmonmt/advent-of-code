package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"chem"
)

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

func reduce(molecule []byte, mappings *chem.Mappings, dict *chem.Dict, results *chem.Results) bool {
	foundReplacement := false

	i := 0
	for i < len(molecule) {
		foundMappings := mappings.Find(molecule[i:], dict)
		if len(foundMappings) == 0 {
			i++
			continue
		} else if len(foundMappings) != 1 {
			panic(fmt.Sprintf("found >1 mappings at %v",
				chem.MoleculeToString(molecule[i:], dict)))
		}
		foundMapping := foundMappings[0]
		foundReplacement = true

		replaced := replace(molecule, i, foundMapping, dict)

		fmt.Printf("%v => %v @ %v; was %v now %v\n",
			chem.MoleculeToString(foundMapping.From, dict),
			chem.MoleculeToString(foundMapping.To, dict), i,
			chem.MoleculeToString(molecule, dict),
			chem.MoleculeToString(replaced, dict))

		i += len(foundMapping.To)
		molecule = replaced
	}

	results.Add(molecule)
	return foundReplacement
}

func main() {
	dict := chem.NewDict()

	mappings, initial, err := readInput(os.Stdin, dict)
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	// mappings.Dump(dict)
	// fmt.Println(initial)

	eByte := dict.StrToByte("e")

	molecules := [][]byte{initial}
	for i := 1; ; i++ {
		results := chem.NewResults()

		foundReplacement := false
		for _, molecule := range molecules {
			if reduce(molecule, mappings, dict, results) {
				foundReplacement = true
			}
		}

		if !foundReplacement {
			fmt.Printf("round %d: no replacement; molecules[0] %v\n",
				i, chem.MoleculeToString(molecules[0], dict))
			break
		}

		molecules = results.Get()
		var shortest int
		for j, molecule := range molecules {
			if j == 0 || len(molecule) < shortest {
				shortest = len(molecule)
			}

			if len(molecule) == 1 && molecule[0] == eByte {
				log.Fatalf("found match")
			}
		}

		fmt.Printf("round %d: num=%d, shortest=%d\n", i, len(molecules), shortest)
	}
}
