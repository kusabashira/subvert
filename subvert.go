package main

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	branches  = regexp.MustCompile(`(?:[^,\\]|\\.)*`)
	sequenses = regexp.MustCompile(`(?:[^~\\]|\\.)*`)
)

func newMatcher(pat string) (m *regexp.Regexp, err error) {
	pat = strings.Replace(pat, `\,`, `\\,`, -1)
	pat = strings.Replace(pat, `\~`, `\\~`, -1)
	pat = `"` + pat + `"`
	pat, err = strconv.Unquote(pat)
	if err != nil {
		return nil, err
	}
	sls := sequenses.FindAllString(pat, -1)
	for si := 0; si < len(sls); si++ {
		bls := branches.FindAllString(sls[si], -1)
		for bi := 0; bi < len(bls); bi++ {
			bls[bi] = strings.Replace(bls[bi], `\,`, `,`, -1)
			bls[bi] = strings.Replace(bls[bi], `\~`, `~`, -1)
			bls[bi] = regexp.QuoteMeta(bls[bi])
		}
		sls[si] = "(" + strings.Join(bls, "|") + ")"
	}
	pat = strings.Join(sls, "")
	return regexp.Compile(pat)
}
