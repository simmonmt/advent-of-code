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

// func doReplacement(in []string, from string, to []string) [][]string {
// 	ret := [][]string{}

// 	for i, a := range in {
// 		if a == from {
// 			// fmt.Printf("matched %v to %v at %v\n", from, to, i)
// 			out := make([]string, len(in)+len(to)-1)
// 			off := 0
// 			if i > 0 {
// 				copy(out[off:i], in[0:i])
// 				off += i
// 				// fmt.Printf("first %v\n", out)
// 			}
// 			copy(out[off:off+len(to)], to)
// 			off += len(to)
// 			// fmt.Printf("mid %v\n", out)

// 			if i < len(in) {
// 				copy(out[off:], in[i+1:])
// 				// fmt.Printf("end %v\n", out)
// 			}

// 			// fmt.Printf("out %v\n", out)
// 			ret = append(ret, out)
// 		}
// 	}

// 	return ret
// }

// func doReplacements(in []string, repls map[string][][]string) map[string][]string {
// 	//fmt.Printf("doReplacements in %v\n", in)

// 	allResults := map[string][]string{}
// 	for from, tos := range repls {
// 		for _, to := range tos {
// 			// fmt.Printf("%v %v\n", from, to)
// 			results := doReplacement(in, from, to)
// 			for _, result := range results {
// 				key := strings.Join(result, " ")
// 				// fmt.Println(key)
// 				allResults[key] = result
// 			}
// 		}
// 	}

// 	return allResults
// }

// func same(a, b []string) bool {
// 	if len(a) != len(b) {
// 		return false
// 	}
// 	for i, v := range a {
// 		if v != b[i] {
// 			return false
// 		}
// 	}
// 	return true
// }

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

func reduce(molecule []byte, mappings *chem.Mappings, d *chem.Dict, results *chem.Results) {
	for i := 0; i < len(molecule); i++ {
		foundMappings := mappings.Find(molecule[i:], d)
		// for _, mapping := range foundMappings {
		// 	fmt.Printf("found mapping %v\n", mapping.ToString(d))
		// }

		for _, mapping := range foundMappings {
			results.Add(replace(molecule, i, mapping, d))
		}
	}
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

		for _, molecule := range molecules {
			reduce(molecule, mappings, dict, results)
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
