package command

import (
	"fmt"
	"regexp"
	"strconv"

	"screen"
)

var (
	rectPattern         = regexp.MustCompile(`^rect ([0-9]+)x([0-9]+)$`)
	rotateRowPattern    = regexp.MustCompile(`^rotate row y=([0-9]+) by ([0-9]+)$`)
	rotateColumnPattern = regexp.MustCompile(`^rotate column x=([0-9]+) by ([0-9]+)$`)
)

type Command interface {
	Execute(s *screen.Screen)
}

type rect struct {
	a, b int
}

func NewRect(a, b int) Command {
	return &rect{a, b}
}

func ParseRect(str string) (Command, error) {
	matches := rectPattern.FindStringSubmatch(str)
	if matches == nil {
		return nil, fmt.Errorf("invalid rect command")
	}

	a, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse a in rect command")
	}

	b, err := strconv.ParseUint(matches[2], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse a in rect command")
	}

	return NewRect(int(a), int(b)), nil
}

func (c *rect) Execute(s *screen.Screen) {
	s.Rect(c.a, c.b)
}

type rotateRow struct {
	y, by int
}

func NewRotateRow(y, by int) Command {
	return &rotateRow{y, by}
}

func ParseRotateRow(str string) (Command, error) {
	matches := rotateRowPattern.FindStringSubmatch(str)
	if matches == nil {
		return nil, fmt.Errorf("invalid rotate row command")
	}

	y, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse y in rotate row command")
	}

	by, err := strconv.ParseUint(matches[2], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse by in rotate row command")
	}

	return NewRotateRow(int(y), int(by)), nil
}

func (c *rotateRow) Execute(s *screen.Screen) {
	s.RotateRow(c.y, c.by)
}

type rotateColumn struct {
	x, by int
}

func NewRotateColumn(x, by int) Command {
	return &rotateColumn{x, by}
}

func ParseRotateColumn(str string) (Command, error) {
	matches := rotateColumnPattern.FindStringSubmatch(str)
	if matches == nil {
		return nil, fmt.Errorf("invalid rotate column command")
	}

	x, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse y in rotate column command")
	}

	by, err := strconv.ParseUint(matches[2], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse by in rotate column command")
	}

	return NewRotateColumn(int(x), int(by)), nil
}

func (c *rotateColumn) Execute(s *screen.Screen) {
	s.RotateColumn(c.x, c.by)
}
