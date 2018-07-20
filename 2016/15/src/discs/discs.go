package discs

type DiscDesc struct {
	NPos  int
	Start int
}

func Advance(posns []int) {
	for i := range posns {
		posns[i]++
	}
}

func Success(descs []DiscDesc, posns []int) bool {
	for i := range posns {
		if (posns[i]+i+1)%descs[i].NPos != 0 {
			return false
		}
	}

	return true
}
