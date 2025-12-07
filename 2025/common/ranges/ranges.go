package ranges

type IncRange struct {
	From, To int
}

func (f IncRange) Contains(n int) bool {
	return n >= f.From && n <= f.To
}

func (r IncRange) Overlaps(other IncRange) bool {
	a, b := &r, &other
	if other.From <= r.From {
		a, b = &other, &r
	}

	if a.To < b.From {
		return false // a ends before b begins
	} else if a.To > b.To {
		return true // a completely encloses b
	}
	return a.To <= b.To
}

func (r IncRange) Merge(other IncRange) (IncRange, bool) {
	if !r.Overlaps(other) {
		return IncRange{}, false
	}

	return IncRange{
		From: min(r.From, other.From),
		To:   max(r.To, other.To),
	}, true
}
