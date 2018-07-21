package extent

import "strings"

type Extents []*Extent

func (e Extents) String() string {
	out := make([]string, len(e))
	for i := range e {
		out[i] = e[i].String()
	}
	return strings.Join(out, ",")
}

func (e Extents) Len() int { return len(e) }

func (e Extents) Swap(i, j int) {
	e[i].Start, e[j].Start = e[j].Start, e[i].Start
	e[i].End, e[j].End = e[j].End, e[i].End
}

func (e Extents) Less(i, j int) bool {
	if e[i].Start < e[j].Start {
		return true
	} else if e[i].Start > e[j].Start {
		return false
	} else {
		return e[i].End < e[j].End
	}
}

// Extents must be sorted
func (e Extents) Merge() Extents {
	out := make([]*Extent, len(e))

	outIdx := 0
	for i := 0; i < len(e); {
		// Try to merge as much into cur as possible
		cur := e[i]

		numMerged := 0
		for j := i + 1; j < len(e); j++ {
			if merged := cur.Merge(e[j]); merged != nil {
				cur = merged
				numMerged++
			} else {
				break
			}
		}

		out[outIdx] = cur
		outIdx++
		i += numMerged + 1
	}

	// out was sized pessimistically, assuming no merges. There may be
	// considerable wasted space now that we've merged, so resize the
	// array. If we don't resize, we'd have to return a slice of the
	// oversized array, which would keep it around longer than necessary.
	return append([]*Extent(nil), out[:outIdx]...)
}
