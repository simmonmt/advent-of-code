package utils

func Pluralize(s bool) string {
	if s {
		return "s"
	} else {
		return ""
	}
}
