package area

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2022/common/mtsmath"
	"github.com/simmonmt/aoc/2022/common/pos"
)

type Area1D struct {
	From, To int
}

// ParseArea1D parses strings like "3-4" into Area1D instances. It
// requires the left value to be less than or equal to the right
// value.
func ParseArea1D(s string) (Area1D, error) {
	left, right, ok := strings.Cut(s, "-")
	if !ok {
		return Area1D{}, fmt.Errorf("bad range cut")
	}

	parseNum := func(s string) (int, error) {
		num, err := strconv.ParseInt(s, 0, 32)
		if err != nil {
			return 0, err
		}
		if num <= 0 {
			return 0, fmt.Errorf("num out of range")
		}
		return int(num), nil
	}

	var a Area1D
	var err error
	a.From, err = parseNum(left)
	if err != nil {
		return Area1D{}, err
	}
	a.To, err = parseNum(right)
	if err != nil {
		return Area1D{}, err
	}

	if a.From > a.To {
		return Area1D{}, fmt.Errorf("from > to")
	}

	return a, nil
}

func (a Area1D) Contains(o Area1D) bool {
	return a.From <= o.From && a.To >= o.To
}

func (a Area1D) ContainsVal(v int) bool {
	return v >= a.From && v <= a.To
}

func (a Area1D) Overlaps(o Area1D) bool {
	if a.From <= o.From {
		return a.To >= o.From
	} else {
		return a.From <= o.To
	}
}

func (a Area1D) Size() int {
	return a.To - a.From + 1
}

// Merge joins two ranges that *must* overlap. The return value is
// undefined if they don't.
func (a Area1D) Merge(o Area1D) Area1D {
	return Area1D{
		From: mtsmath.Min(a.From, o.From),
		To:   mtsmath.Max(a.To, o.To),
	}
}

func (a Area1D) String() string {
	return fmt.Sprintf("%d-%d", a.From, a.To)
}

// TODO: More efficient algorithm
func Merge1DRanges(ranges []Area1D) []Area1D {
	for i := 0; i < len(ranges); i++ {
		for {
			changed := false
			out := ranges[0 : i+1]
			for j := i + 1; j < len(ranges); j++ {
				if out[i].Overlaps(ranges[j]) {
					out[i] = out[i].Merge(ranges[j])
					changed = true
				} else if out[i].To+1 == ranges[j].From {
					out[i].To = ranges[j].To
					changed = true
				} else if ranges[j].To+1 == out[i].From {
					out[i].From = ranges[j].From
					changed = true
				} else {
					out = append(out, ranges[j])
				}
			}
			ranges = out
			if !changed {
				break
			}
		}
	}

	return ranges
}

type Area2D struct {
	From, To pos.P2
}

func (a Area2D) Contains(o Area2D) bool {
	return a.From.X <= o.From.X && a.To.X >= o.To.X &&
		a.From.Y <= o.From.Y && a.To.Y >= o.To.Y
}

func (a Area2D) String() string {
	return fmt.Sprintf("(%s)-(%s)", a.From, a.To)
}

type Area3D struct {
	From, To pos.P3
}

func (a Area3D) Contains(o Area3D) bool {
	return a.From.X <= o.From.X && a.To.X >= o.To.X &&
		a.From.Y <= o.From.Y && a.To.Y >= o.To.Y &&
		a.From.Z <= o.From.Z && a.To.Z >= o.To.Z
}

func (a Area3D) Overlaps(o Area3D) bool {
	axisOverlaps := func(aLo, aHi, oLo, oHi int) bool {
		return oLo <= aHi && oHi >= aLo
	}

	if !axisOverlaps(a.From.X, a.To.X, o.From.X, o.To.X) {
		return false
	}
	if !axisOverlaps(a.From.Y, a.To.Y, o.From.Y, o.To.Y) {
		return false
	}
	if !axisOverlaps(a.From.Z, a.To.Z, o.From.Z, o.To.Z) {
		return false
	}
	return true
}

func (a Area3D) String() string {
	return fmt.Sprintf("(%s)-(%s)", a.From, a.To)
}
