package strutil

import (
	"fmt"
	"strconv"
	"strings"
)

func StringToInt64s(str string) ([]int64, error) {
	out := []int64{}
	for _, s := range strings.Split(str, ",") {
		v, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("bad value %v: %v", s, err)
		}

		out = append(out, v)
	}
	return out, nil
}
