package pad

func HasRepeats(in string, wantLen int) []rune {
	reps := []rune{}
	var last rune
	var streak int
	for i, r := range in {
		if i == 0 {
			last = r
			streak = 1
		} else {
			if r == last {
				streak++
				if streak == wantLen {
					reps = append(reps, last)
				}
			} else {
				last = r
				streak = 1
			}
		}
	}

	return reps
}
