package util

func StrToBoolArray(in string) []bool {
	out := make([]bool, len(in))
	for i, c := range in {
		out[i] = c == '1'
	}
	return out
}

func BoolArrayToStr(in []bool) string {
	out := make([]rune, len(in))
	for i, b := range in {
		if b {
			out[i] = '1'
		} else {
			out[i] = '0'
		}
	}
	return string(out)
}
