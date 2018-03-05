package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

type Dict struct {
	from map[string]byte
	to map[byte]string
	last uint16
}

func NewDict() *Dict {
	return &Dict{
		from: map[string]byte[},
		to :map[byte]string{},
		last: 0,
	}
}

func (d *Dict) Add(s string) byte {
	d.last++
	if d.last > 0xff {
		panic("too many")
	}
	return byte(d.last)
}

func parseMolecule(in string) []string {
	out := []string{}
	for _, c := range in {
		if unicode.IsUpper(c) {
			out = append(out, string(c))
		} else {
			out[len(out)-1] += string(c)
		}
	}
	return out
}

func readInput(r io.Reader) (map[string][][]string, []string, error) {
	reader := bufio.NewReader(r)

	repls := map[string][][]string{}
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
		from, to := pairs[0], pairs[1]

		if _, found := repls[from]; !found {
			repls[from] = [][]string{}
		}
		repls[from] = append(repls[from], parseMolecule(to))
	}

	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, nil, fmt.Errorf("no input molecule")
	}
	line = strings.TrimSpace(line)

	molecule := parseMolecule(line)

	return repls, molecule, nil
}

func doReplacement(in []string, from string, to []string) [][]string {
	ret := [][]string{}

	for i, a := range in {
		if a == from {
			// fmt.Printf("matched %v to %v at %v\n", from, to, i)
			out := make([]string, len(in)+len(to)-1)
			off := 0
			if i > 0 {
				copy(out[off:i], in[0:i])
				off += i
				// fmt.Printf("first %v\n", out)
			}
			copy(out[off:off+len(to)], to)
			off += len(to)
			// fmt.Printf("mid %v\n", out)

			if i < len(in) {
				copy(out[off:], in[i+1:])
				// fmt.Printf("end %v\n", out)
			}

			// fmt.Printf("out %v\n", out)
			ret = append(ret, out)
		}
	}

	return ret
}

func doReplacements(in []string, repls map[string][][]string) map[string][]string {
	//fmt.Printf("doReplacements in %v\n", in)

	allResults := map[string][]string{}
	for from, tos := range repls {
		for _, to := range tos {
			// fmt.Printf("%v %v\n", from, to)
			results := doReplacement(in, from, to)
			for _, result := range results {
				key := strings.Join(result, " ")
				// fmt.Println(key)
				allResults[key] = result
			}
		}
	}

	return allResults
}

func same(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func main() {
	repls, goal, err := readInput(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	// fmt.Println(repls)
	fmt.Printf("%d: %v\n", len(goal), goal)

	// for r := range allResults {
	// 	fmt.Println(r)
	// }

	molecules := map[string][]string{"e": []string{"e"}}
	for i := 0; ; i++ {
		allResults := map[string][]string{}

		maxLen := -1
		droppedTooBig := 0

		for _, molecule := range molecules {
			results := doReplacements(molecule, repls)
			//fmt.Printf("got results %v\n", results)

			for k, v := range results {
				if maxLen == -1 || len(v) > maxLen {
					maxLen = len(v)
				}

				if len(v) > len(goal) {
					droppedTooBig++
					continue
				} else if len(v) == len(goal) {
					if same(v, goal) {
						log.Fatalf("found it at i=%d", i)
					}
				}

				allResults[k] = v
			}
		}

		fmt.Printf("i=%d, num=%d, maxLen=%d, tooBig=%d\n",
			i, len(allResults), maxLen, droppedTooBig)
		molecules = allResults
	}
}
