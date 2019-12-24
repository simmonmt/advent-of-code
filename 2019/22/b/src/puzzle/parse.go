package puzzle

import (
	"fmt"
	"strconv"
	"strings"
)

type Verb int

const (
	VERB_DEAL_INTO_NEW_STACK Verb = iota
	VERB_DEAL_WITH_INCREMENT
	VERB_CUT_LEFT
	VERB_CUT_RIGHT
)

type Command struct {
	Verb Verb
	Val  int
}

func parseCommand(str string) (*Command, error) {
	switch {
	case str == "deal into new stack":
		return &Command{VERB_DEAL_INTO_NEW_STACK, 0}, nil
	case strings.HasPrefix(str, "cut "):
		parts := strings.Split(str, " ")
		val, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		if val < 0 {
			return &Command{VERB_CUT_RIGHT, -val}, nil
		} else {
			return &Command{VERB_CUT_LEFT, val}, nil
		}
	case strings.HasPrefix(str, "deal with increment"):
		parts := strings.Split(str, " ")
		val, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, err
		}

		return &Command{VERB_DEAL_WITH_INCREMENT, val}, nil
	default:
		return nil, fmt.Errorf("bad command")
	}
}

func ParseCommands(lines []string) ([]*Command, error) {
	cmds := []*Command{}
	for _, line := range lines {
		cmd, err := parseCommand(line)
		if err != nil {
			return nil, fmt.Errorf(`error parsing "%v": %v`, line, err)
		}
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}
