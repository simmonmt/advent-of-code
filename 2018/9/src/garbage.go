package main

import (
	"bufio"
	"fmt"
	lexer "go-lexer-master"
	"os"
	"strings"
)

const (
	EOFToken = iota
	GarbageToken
	GroupStartToken
	GroupEndToken
	ThingSepToken
)

func GarbageSize(g string) int {
	size := 0

	g = g[1 : len(g)-1]
	cancelling := false
	for _, c := range g {
		if cancelling {
			cancelling = false
		} else if c == '!' {
			cancelling = true
		} else {
			size++
		}
	}
	return size
}

func GarbageState(l *lexer.L) lexer.StateFunc {
	for {
		switch l.Peek() {
		case '!':
			l.Next()
			l.Next()
			continue
		case '>':
			l.Next()
			l.Emit(GarbageToken)
			return InitialState
		default:
			l.Next()
		}
	}
}

func InitialState(l *lexer.L) lexer.StateFunc {
	switch l.Peek() {
	case '<':
		return GarbageState
	case '{':
		l.Next()
		l.Emit(GroupStartToken)
		return InitialState
	case '}':
		l.Next()
		l.Emit(GroupEndToken)
		return InitialState
	case ',':
		l.Next()
		l.Emit(ThingSepToken)
		return InitialState
	case lexer.EOFRune:
		return nil
	}

	l.Error(fmt.Sprintf("unknown char %v", l.Peek()))
	return InitialState
}

type Lexer struct {
	lexer.L
}

func NewLexer(src string) *Lexer {
	return &Lexer{*lexer.New(src, InitialState)}
}

func (l *Lexer) Lex() (lexer.TokenType, string) {
	tok, done := l.NextToken()
	if done {
		return EOFToken, ""
	} else {
		return tok.Type, tok.Value
	}
}

type Group struct {
	Subs []Group
}

func NewGroup() Group {
	return Group{Subs: []Group{}}
}

func ParseGroup(lexer *Lexer) (*Group, int) {
	subs := []Group{}
	garbageSize := 0

	for {
		tok, val := lexer.Lex()
		switch tok {
		case GarbageToken:
			// fmt.Printf("garbage val: %v\n", val)
			// fmt.Printf("garbage size: %v\n", GarbageSize(val))
			garbageSize += GarbageSize(val)
			continue
		case GroupStartToken:
			group, subGarbageSize := ParseGroup(lexer)
			subs = append(subs, *group)
			garbageSize += subGarbageSize
			continue
		case GroupEndToken:
			return &Group{Subs: subs}, garbageSize
		case ThingSepToken:
			continue
		case EOFToken:
			return &Group{Subs: subs}, garbageSize
		}
	}
}

func ScoreGroup(group *Group, base int) int {
	score := base
	for _, sub := range group.Subs {
		score += ScoreGroup(&sub, base+1)
	}
	return score
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		fmt.Println(line)

		lexer := NewLexer(line)
		lexer.Start()

		group, garbageSize := ParseGroup(lexer)
		fmt.Printf("total score for all groups: %v\n", ScoreGroup(group, 0))
		fmt.Printf("total garbage size: %v\n", garbageSize)
	}
}
