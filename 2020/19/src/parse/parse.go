package parse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	rulePattern = regexp.MustCompile(`^([0-9]+): (.*)$`)
)

func makeRuleMap(rules []string) (map[int]string, error) {
	ruleMap := map[int]string{}

	for lineNum, rule := range rules {
		parts := rulePattern.FindStringSubmatch(rule)
		if parts == nil {
			return nil, fmt.Errorf("%d: bad rule %v", lineNum, rule)
		}

		ruleNum, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("%d: bad rule num %v",
				lineNum, parts[1])
		}

		ruleMap[ruleNum] = parts[2]
	}

	return ruleMap, nil
}

func resolveSeq(ruleMap map[int]string, body string) string {
	body = strings.TrimSpace(body)
	logger.LogF("resolving seq '%v'", body)

	out := strings.Builder{}
	for _, part := range strings.Split(body, " ") {
		num, err := strconv.Atoi(part)
		if err != nil {
			panic(fmt.Sprintf("bad num %v: %v", part, err))
		}

		out.WriteString(resolve(ruleMap, num))
	}
	return out.String()
}

func resolve(ruleMap map[int]string, num int) string {
	body, found := ruleMap[num]
	if !found {
		panic("bad rule ref")
	}

	parts := strings.Split(body, "|")
	if len(parts) > 1 {
		a := resolveSeq(ruleMap, parts[0])
		b := resolveSeq(ruleMap, parts[1])
		return fmt.Sprintf("(?:%s|%s)", a, b)
	}

	if strings.HasPrefix(body, `"`) {
		return strings.Trim(body, `"`)
	}

	return resolveSeq(ruleMap, body)
}

func Parse(rules []string, num int) (string, error) {
	ruleMap, err := makeRuleMap(rules)
	if err != nil {
		return "", nil
	}

	return resolve(ruleMap, num), nil
}
