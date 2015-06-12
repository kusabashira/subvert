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
	// one branch
	{`abc`, regexp.MustCompile(`(abc)`)},
	{`abcdef`, regexp.MustCompile(`(abcdef)`)},

	// multiple branches
	{`a,b`, regexp.MustCompile(`(a|b)`)},
	{`a,bc,def`, regexp.MustCompile(`(a|bc|def)`)},

	// use escape
	{`a\,b`, regexp.MustCompile(`(a,b)`)},
	{`a\,bc\,def`, regexp.MustCompile(`(a,bc,def)`)},

	// multiple branches with escape
	{`a\,b,c`, regexp.MustCompile(`(a,b|c)`)},
	{`a,bc\,def`, regexp.MustCompile(`(a|bc,def)`)},

	// regexp quote
	{`a+b`, regexp.MustCompile(`(a\+b)`)},
	{`(a|bc)*def`, regexp.MustCompile(`(\(a\|bc\)\*def)`)},

	// unquote special values
	{`a\\bc`, regexp.MustCompile("(a\\\\bc)")},
	{`a\tb\,c`, regexp.MustCompile("(a\tb,c)")},
	{`a\tbc\n\ndef`, regexp.MustCompile("(a\tbc\n\ndef)")},

	// multiple sequenses
	{`a~b`, regexp.MustCompile("(a)(b)")},
	{`a~bc~def`, regexp.MustCompile("(a)(bc)(def)")},
	{`a,b~c`, regexp.MustCompile("(a|b)(c)")},
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
