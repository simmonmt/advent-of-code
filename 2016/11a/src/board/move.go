package board

import (
	"fmt"
	"object"
	"strings"
)

type Move struct {
	dest uint8
	objs []object.Object
}

func (m Move) String() string {
	out := fmt.Sprintf("{Dest: %d Objs:[", m.dest)

	os := []string{}
	for _, obj := range m.objs {
		os = append(os, obj.String())
	}
	out += strings.Join(os, ",")
	out += "]}"
	return out
}

func newMove(dest uint8, objs ...object.Object) *Move {
	m := &Move{
		dest: dest,
		objs: make([]object.Object, len(objs)),
	}
	copy(m.objs, objs)
	return m
}
