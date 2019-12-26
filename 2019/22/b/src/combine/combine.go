package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"reflect"

	"github.com/simmonmt/aoc/2019/22/b/src/puzzle"
	"github.com/simmonmt/aoc/2019/common/logger"
)

var (
	verbose   = flag.Bool("verbose", false, "verbose")
	indexFlag = flag.Int64("index", 2020, "index to use for eval")
	input     = flag.String("input", "", "input file")
	mod       = flag.Int("mod", 119315717514047, "mod")
	numRuns   = flag.Int("num_runs", -1, "num combining runs")
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

func readCommands(path string) ([]*puzzle.Command, error) {
	lines, err := readInput(path)
	if err != nil {
		return nil, err
	}

	cmds, err := puzzle.ParseCommands(lines)
	if err != nil {
		return nil, err
	}

	return cmds, nil
}

func runCommands(cmds []*puzzle.Command, index int64, numRuns int) []int64 {
	out := []int64{index}

	for i := 1; i <= numRuns; i++ {
		index = puzzle.ForwardCommandsForIndex(cmds, int64(*mod), index)
		out = append(out, index)
	}

	return out
}

func verifyCmds(cmds []*puzzle.Command, refVals []int64) bool {
	vals := runCommands(cmds, *indexFlag, len(refVals)-1)
	if !reflect.DeepEqual(vals, refVals) {
		fmt.Println("difference found")
		fmt.Printf("ref : %v\n", refVals)
		fmt.Printf("eval: %v\n", vals)
		return false
	}

	return true
}

func isCut(cmd *puzzle.Command) bool {
	return cmd.Verb == puzzle.VERB_CUT_LEFT || cmd.Verb == puzzle.VERB_CUT_RIGHT
}

func isInc(cmd *puzzle.Command) bool {
	return cmd.Verb == puzzle.VERB_DEAL_WITH_INCREMENT
}

func isNew(cmd *puzzle.Command) bool {
	return cmd.Verb == puzzle.VERB_DEAL_INTO_NEW_STACK
}

func combineCuts(a, b *puzzle.Command) []*puzzle.Command {
	aVal := a.Val
	if a.Verb == puzzle.VERB_CUT_RIGHT {
		aVal = -aVal
	}

	bVal := b.Val
	if b.Verb == puzzle.VERB_CUT_RIGHT {
		bVal = -bVal
	}

	newVal := big.NewInt(int64(aVal))
	newVal.Add(newVal, big.NewInt(int64(bVal)))
	if newVal.Cmp(big.NewInt(0)) < 0 {
		newVal.Mod(newVal, big.NewInt(int64(*mod)))
		newVal.Sub(newVal, big.NewInt(int64(*mod)))
	} else {
		newVal.Mod(newVal, big.NewInt(int64(*mod)))
	}
	val := int(newVal.Int64())

	//val := aVal + bVal
	if val > 0 {
		return []*puzzle.Command{
			&puzzle.Command{Verb: puzzle.VERB_CUT_LEFT, Val: val},
		}
	} else {
		return []*puzzle.Command{
			&puzzle.Command{Verb: puzzle.VERB_CUT_RIGHT, Val: -val},
		}
	}
}

func combineIncs(a, b *puzzle.Command) []*puzzle.Command {
	newVal := big.NewInt(int64(a.Val))
	newVal.Mul(newVal, big.NewInt(int64(b.Val)))
	if newVal.Cmp(big.NewInt(0)) < 0 {
		newVal.Mod(newVal, big.NewInt(int64(*mod)))
		newVal.Sub(newVal, big.NewInt(int64(*mod)))
	} else {
		newVal.Mod(newVal, big.NewInt(int64(*mod)))
	}

	// newVal := big.NewInt(int64(a.Val))
	// newVal.Mul(newVal, big.NewInt(int64(b.Val)))
	//newVal.Mod(newVal, big.NewInt(int64(*mod)))

	return []*puzzle.Command{
		&puzzle.Command{
			Verb: puzzle.VERB_DEAL_WITH_INCREMENT,
			Val:  int(newVal.Int64()),
		},
	}
}

func combineCutInc(cut, inc *puzzle.Command) []*puzzle.Command {
	newVal := big.NewInt(int64(cut.Val))
	newVal.Mul(newVal, big.NewInt(int64(inc.Val)))
	newVal.Mod(newVal, big.NewInt(int64(*mod)))

	newCut := &puzzle.Command{
		Verb: cut.Verb,
		Val:  int(newVal.Int64()),
	}

	if newCut.Val < 0 {
		if newCut.Verb == puzzle.VERB_CUT_LEFT {
			newCut.Verb = puzzle.VERB_CUT_RIGHT
		} else {
			newCut.Verb = puzzle.VERB_CUT_LEFT
		}

		newCut.Val = -newCut.Val
	}

	return []*puzzle.Command{inc, newCut}
}

func combineCutNew(cut, deal *puzzle.Command) []*puzzle.Command {
	newCutVerb := puzzle.VERB_CUT_LEFT
	if cut.Verb == puzzle.VERB_CUT_LEFT {
		newCutVerb = puzzle.VERB_CUT_RIGHT
	}

	newCut := &puzzle.Command{Verb: newCutVerb, Val: cut.Val}

	return []*puzzle.Command{deal, newCut}
}

func combineNewInc(deal, inc *puzzle.Command) []*puzzle.Command {
	return []*puzzle.Command{
		&puzzle.Command{Verb: puzzle.VERB_CUT_RIGHT, Val: 1},
		&puzzle.Command{Verb: puzzle.VERB_DEAL_INTO_NEW_STACK, Val: -deal.Val},
	}
}

func combineNewNew(a, b *puzzle.Command) []*puzzle.Command {
	return []*puzzle.Command{}
}

func eliminateNew() []*puzzle.Command {
	return []*puzzle.Command{
		&puzzle.Command{Verb: puzzle.VERB_DEAL_WITH_INCREMENT, Val: -1},
		&puzzle.Command{Verb: puzzle.VERB_CUT_LEFT, Val: 1},
	}
}

type checkDesc struct {
	check1, check2 func(*puzzle.Command) bool
	rewrite        func(a, b *puzzle.Command) []*puzzle.Command
}

var (
	checkDescs = []*checkDesc{
		&checkDesc{isCut, isCut, combineCuts},
		&checkDesc{isInc, isInc, combineIncs},
		&checkDesc{isCut, isInc, combineCutInc},
		&checkDesc{isCut, isNew, combineCutNew},
		&checkDesc{isNew, isInc, combineNewInc},
		&checkDesc{isNew, isNew, combineNewNew},
	}
)

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

	refVals := runCommands(cmds, *indexFlag, 10)

	changed := true
	for num := 0; *numRuns == -1 || num < *numRuns; num++ {
		changed = false
		newCmds := []*puzzle.Command{}
		for i := 0; i < len(cmds); {
			numMatched := 0
			if i < len(cmds)-1 && !changed {
				for _, desc := range checkDescs {
					if desc.check1(cmds[i]) && desc.check2(cmds[i+1]) {
						rewritten := desc.rewrite(cmds[i], cmds[i+1])
						// fmt.Printf("%v %v to ", cmds[i], cmds[i+1])
						// for _, c := range rewritten {
						// 	fmt.Printf(" %v", *c)
						// }
						// fmt.Println()

						if len(rewritten) > 0 {
							newCmds = append(newCmds, rewritten...)
						}
						numMatched = 2
						break
					}
				}
			}

			if numMatched == 0 && !changed {
				if isNew(cmds[i]) {
					newCmds = append(newCmds, eliminateNew()...)
					numMatched = 1
				}
			}

			if numMatched > 0 {
				changed = true
				i += numMatched
			} else {
				newCmds = append(newCmds, cmds[i])
				i++
			}
		}

		cmds = newCmds

		if !verifyCmds(cmds, refVals) {
			log.Fatal("verification failure")
		}

		if !changed {
			break
		}
	}

	for _, cmd := range cmds {
		switch cmd.Verb {
		case puzzle.VERB_DEAL_INTO_NEW_STACK:
			fmt.Println("deal into new stack")
			break
		case puzzle.VERB_CUT_LEFT:
			fmt.Printf("cut %d\n", cmd.Val)
			break
		case puzzle.VERB_CUT_RIGHT:
			fmt.Printf("cut -%d\n", cmd.Val)
			break
		case puzzle.VERB_DEAL_WITH_INCREMENT:
			fmt.Printf("deal with increment %d\n", cmd.Val)
			break
		}
	}
}
