package puzzle

import (
	"math/big"
)

var (
	cmdDispatch = map[Verb]func([]int, int) []int{
		VERB_DEAL_INTO_NEW_STACK: DealIntoNewStack,
		VERB_DEAL_WITH_INCREMENT: DealWithIncrement,
		VERB_CUT_LEFT:            CutLeft,
		VERB_CUT_RIGHT:           CutRight,
	}
)

func RunCommands(cards []int, cmds []*Command) []int {
	for _, cmd := range cmds {
		f, found := cmdDispatch[cmd.Verb]
		if !found {
			panic("no cmd")
		}
		cards = f(cards, cmd.Val)
	}
	return cards
}

func DealIntoNewStack(in []int, _ int) []int {
	out := make([]int, len(in))
	for i := range out {
		out[i] = in[len(in)-i-1]
	}
	return out
}

func DealWithIncrement(in []int, inc int) []int {
	out := make([]int, len(in))
	for i := range out {
		newI := i * inc
		for newI < 0 {
			newI += len(in)
		}
		newI = newI % len(in)
		out[newI] = in[i]
	}
	return out
}

func CutLeft(in []int, l int) []int {
	out := make([]int, len(in))
	copy(out, in[l:])
	copy(out[len(out)-l:], in)
	return out
}

func CutRight(in []int, l int) []int {
	out := make([]int, len(in))
	copy(out, in[len(in)-l:])
	copy(out[l:], in)
	return out
}

func extGCD(a, b int) (s, t int) {
	x, y, u, v := 0, 1, 1, 0
	for a != 0 {
		q, r := b/a, b%a
		m, n := x-u*q, y-v*q
		b, a, x, y, u, v = a, r, u, v, m, n
	}

	return x, y
}

type sCacheEnt struct {
	val int
	sz  int
}

var (
	sCache = map[sCacheEnt]int{}
)

func cachedExtGCD(val, sz int) (s int) {
	ent := sCacheEnt{val, sz}
	if s, found := sCache[ent]; found {
		return s
	}
	s, _ = extGCD(val, sz)
	sCache[ent] = s
	return s
}

func ReverseCommandsForIndex(cmds []*Command, sz int, index int) int {
	for i := len(cmds) - 1; i >= 0; i-- {
		switch cmd := cmds[i]; cmd.Verb {
		case VERB_DEAL_INTO_NEW_STACK:
			index = sz - index - 1
			break
		case VERB_DEAL_WITH_INCREMENT:
			s := cachedExtGCD(cmd.Val, sz)
			if s < 0 {
				s += sz
			}

			accum := big.NewInt(int64(index))
			accum.Mul(accum, big.NewInt(int64(s)))
			accum.Mod(accum, big.NewInt(int64(sz)))

			index = int(accum.Int64())
			break
		case VERB_CUT_LEFT:
			if index >= sz-cmd.Val {
				index -= sz - cmd.Val
			} else {
				index += cmd.Val
			}
			break
		case VERB_CUT_RIGHT:
			if index < cmd.Val {
				index += sz - cmd.Val
			} else {
				index -= cmd.Val
			}
			break

		default:
			panic("unimplemented")
		}
	}
	return index
}

func ForwardCommandsForIndex(cmds []*Command, sz int64, index int64) int64 {
	for _, cmd := range cmds {
		switch cmd.Verb {
		case VERB_DEAL_INTO_NEW_STACK:
			index = sz - index - 1
			break
		case VERB_DEAL_WITH_INCREMENT:
			newIndex := big.NewInt(int64(index))
			newIndex.Mul(newIndex, big.NewInt(int64(cmd.Val)))
			newIndex.Mod(newIndex, big.NewInt(sz))
			index = newIndex.Int64()
			break
		case VERB_CUT_LEFT:
			if index < int64(cmd.Val) {
				index += sz - int64(cmd.Val)
			} else {
				index -= int64(cmd.Val)
			}
			break
		case VERB_CUT_RIGHT:
			if index >= sz-int64(cmd.Val) {
				index -= sz - int64(cmd.Val)
			} else {
				index += int64(cmd.Val)
			}
			break
		default:
			panic("unimplemented")
		}
	}
	return index
}
