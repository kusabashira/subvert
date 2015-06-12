package main

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	sequenses = regexp.MustCompile(`(?:[^~\\]|\\.)*`)
	branches  = regexp.MustCompile(`(?:[^,\\]|\\.)*`)
)

func newMatcher(expr string) (m *regexp.Regexp, err error) {
	expr = strings.Replace(expr, `\,`, `\\,`, -1)
	expr = strings.Replace(expr, `\~`, `\\~`, -1)
	expr, err = strconv.Unquote(`"` + expr + `"`)
	if err != nil {
		return nil, err
	}

	sls := sequenses.FindAllString(expr, -1)
	for si := 0; si < len(sls); si++ {
		bls := branches.FindAllString(sls[si], -1)
		for bi := 0; bi < len(bls); bi++ {
			bls[bi] = strings.Replace(bls[bi], `\,`, `,`, -1)
			bls[bi] = strings.Replace(bls[bi], `\~`, `~`, -1)
			bls[bi] = regexp.QuoteMeta(bls[bi])
		}
		sls[si] = "(" + strings.Join(bls, "|") + ")"
	}
	return regexp.Compile(strings.Join(sls, ""))
}
