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
	pat2 := "(" + strings.Join(sp, "|") + ")"
	pat3 := strings.Replace(pat2, `\,`, `,`, -1)
	return regexp.Compile(pat3)
}
