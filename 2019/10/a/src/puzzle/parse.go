package puzzle

func ParseMap(lines []string) map[Pos]bool {
	m := map[Pos]bool{}
	for y, line := range lines {
		for x, r := range line {
			if r == '#' {
				m[Pos{x, y}] = true
			}
		}
	}
	return m
}
