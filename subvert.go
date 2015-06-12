package main

import (
	"regexp"
)

func newMatcher(pat string) (*regexp.Regexp, error) {
	return regexp.Compile(pat)
}
