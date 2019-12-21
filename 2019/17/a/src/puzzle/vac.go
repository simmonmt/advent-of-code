package puzzle

import (
	"github.com/simmonmt/aoc/2019/common/dir"
	"github.com/simmonmt/aoc/2019/common/pos"
)

type VacDir dir.Dir

func (vd VacDir) String() string {
	switch vd {
	case VacDir(dir.DIR_NORTH):
		return "^"
	case VacDir(dir.DIR_SOUTH):
		return "v"
	case VacDir(dir.DIR_WEST):
		return "<"
	case VacDir(dir.DIR_EAST):
		return ">"
	default:
		panic("bad vac dir")
	}
}

func ParseVacDir(r rune) VacDir {
	switch r {
	case '^':
		return VacDir(dir.DIR_NORTH)
	case 'v':
		return VacDir(dir.DIR_SOUTH)
	case '<':
		return VacDir(dir.DIR_WEST)
	case '>':
		return VacDir(dir.DIR_EAST)
	default:
		panic("bad vac dir")
	}
}

type Vac struct {
	Pos pos.P2
	Dir VacDir
}
