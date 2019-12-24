package puzzle

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
		out[i*inc%len(in)] = in[i]
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
