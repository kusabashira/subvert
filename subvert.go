package main

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	branches = regexp.MustCompile(`(?:[^,\\]|\\.)*`)
)

func newMatcher(pat string) (m *regexp.Regexp, err error) {
	pat = strings.Replace(pat, `\,`, `\\,`, -1)
	pat = `"` + pat + `"`
	pat, err = strconv.Unquote(pat)
	if err != nil {
		return nil, err
	}
	sp := branches.FindAllString(pat, -1)
	for i := 0; i < len(sp); i++ {
		sp[i] = strings.Replace(sp[i], `\,`, `,`, -1)
		sp[i] = regexp.QuoteMeta(sp[i])
	}
	pat = "(" + strings.Join(sp, "|") + ")"
	return regexp.Compile(pat)
}
