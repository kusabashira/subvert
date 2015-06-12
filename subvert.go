package main

import (
	"regexp"
	"strings"
)

func newMatcher(pat string) (*regexp.Regexp, error) {
	sp := strings.Split(pat, ",")
	pat2 := "(" + strings.Join(sp, "|") + ")"
	return regexp.Compile(pat2)
}
