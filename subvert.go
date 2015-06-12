package main

import (
	"regexp"
	"strings"
)

var (
	branches = regexp.MustCompile(`(?:[^,\\]|\\.)*`)
)

func newMatcher(pat string) (*regexp.Regexp, error) {
	sp := branches.FindAllString(pat, -1)
	for i := 0; i < len(sp); i++ {
		sp[i] = strings.Replace(sp[i], `\,`, `,`, -1)
		sp[i] = regexp.QuoteMeta(sp[i])
	}
	pat = "(" + strings.Join(sp, "|") + ")"
	return regexp.Compile(pat)
}
