package main

import (
	"reflect"
	"regexp"
	"testing"
)

var genMatcherTests = []struct {
	src string
	dst *regexp.Regexp
}{
	{`abc`, regexp.MustCompile(`(abc)`)},
	{`abcdef`, regexp.MustCompile(`(abcdef)`)},

	{`a,b`, regexp.MustCompile(`(a|b)`)},
	{`a,bc,def`, regexp.MustCompile(`(a|bc|def)`)},

	{`a\,b`, regexp.MustCompile(`(a,b)`)},
	{`a\,bc\,def`, regexp.MustCompile(`(a,bc,def)`)},

	{`a\,b,c`, regexp.MustCompile(`(a,b|c)`)},
	{`a,bc\,def`, regexp.MustCompile(`(a|bc,def)`)},
}

func TestGenMatcher(t *testing.T) {
	for _, test := range genMatcherTests {
		expect := test.dst
		actual, err := newMatcher(test.src)
		if err != nil {
			t.Errorf("NewSubvert(%q) returns %q, want nil",
				test.src, err)
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("%q: got %q, want %q",
				test.src, actual, expect)
		}
	}
}
